// Copyright 2016 Josh Deprez
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package awakengine

import (
	"sort"

	"github.com/hajimehoshi/ebiten"
)

// drawPosition adjusts the source rectangle to refer to the texture atlas, and
// destination rectangle of in-world objects to refer to screen coordinates.
type drawPosition struct{ Object }

func (p drawPosition) Dst() (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = p.Object.Dst()
	if !p.Object.InWorld() {
		return
	}
	x0 -= camPos.X
	y0 -= camPos.Y
	x1 -= camPos.X
	y1 -= camPos.Y
	return
}

func (p drawPosition) Src() (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = p.Object.Src()
	o := compositeOffset[p.Object.ImageKey()]
	x0 += o.X
	y0 += o.Y
	x1 += o.X
	y1 += o.Y
	return
}

func (p drawPosition) InWorld() bool    { return false } // Once positioned, always in screen coordinates.
func (p drawPosition) ImageKey() string { return "" }    // Once positioned, always in texture atlas coordinates.

// drawList is a Z-sortable list of objects.
type drawList []drawPosition

// makeDrawLists makes a fixed and loose drawList out of the slices of objects being passed in.
func makeDrawLists(objs ...[]Object) (fixed, loose drawList) {
	var f, l drawList
	for _, ol := range objs {
		for _, o := range ol {
			r := &l
			if o.Fixed() {
				r = &f
			}
			// Is it a drawPosition already?
			if dp, ok := o.(drawPosition); ok {
				*r = append(*r, dp)
			} else {
				*r = append(*r, drawPosition{o})
			}
		}
	}
	return f, l
}

// Implementing sort.Interface
func (d drawList) Len() int           { return len(d) }
func (d drawList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d drawList) Less(i, j int) bool { return d[i].Z() < d[j].Z() }

// Convenience function.
func (d drawList) Sort() { sort.Sort(d) }

// Implementing ebiten.ImageParts
func (d drawList) Dst(i int) (x0, y0, x1, y1 int) { return d[i].Dst() }
func (d drawList) Src(i int) (x0, y0, x1, y1 int) { return d[i].Src() }

/*
// reorder reorders d and returns b such that the kept objects are d[:b] and non-kept objects are d[b:].
func (d drawList) reorder(keep func(int) bool) int {
	b := len(d)
	for i := 0; i < b; {
		if keep(i) {
			i++
			continue
		}
		b--
		d[i], d[b] = d[b], d[i]
	}
	return b
}
*/

// subsliceNew allocates a new slice and copies the items for which keep is true.
func (d drawList) subslice(keep func(Object) bool) drawList {
	r := make(drawList, 0, len(d))
	for _, o := range d {
		if keep(o) {
			r = append(r, o)
		}
	}
	return r
}

/*
func (d drawList) sublistIndexes(keep func(Object) bool) []int {
	r := make([]int, 0, len(d))
	for i, o := range d {
		if keep(o) {
			r = append(r, i)
		}
	}
	return r
}
*/

// cull removes invisible objects.
func (d drawList) cull() drawList {
	return d.subslice(func(o Object) bool {
		if !o.Visible() {
			return false
		}
		if x0, y0, x1, y1 := o.Dst(); x1 <= 0 || y1 <= 0 || x0 > camSize.X || y0 > camSize.Y {
			return false
		}
		return true
	})
}

// gc removes retired objects.
func (d drawList) gc() drawList {
	return d.subslice(func(o Object) bool {
		return !o.Retire()
	})
}

// merge merges two sorted drawLists into a combined list.
func merge(a, b drawList) drawList {
	r := make(drawList, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if x, y := a[i], b[j]; x.Z() < y.Z() {
			r = append(r, x)
			i++
		} else {
			r = append(r, y)
			j++
		}
	}
	for ; i < len(a); i++ {
		r = append(r, a[i])
	}
	for ; j < len(b); j++ {
		r = append(r, b[j])
	}
	return r
}

/*
// mergeDrawList draws from two draw lists which are assumed to be in order.
type mergeDrawList struct {
	a, b           drawList
	si, sj, di, dj int
}

func (m *mergeDrawList) Len() int {
	m.di, m.dj, m.si, m.sj = 0, 0, 0, 0
	return len(m.a) + len(m.b)
}

func (m *mergeDrawList) Dst(int) (x0, y0, x1, y1 int) {
	if x, y := m.a[m.di], m.b[m.dj]; x.Z() < y.Z() {
		x0, y0, x1, y1 = m.a.Dst(m.di)
		m.di++
		return
	}
	x0, y0, x1, y1 = m.b.Dst(m.dj)
	m.dj++
	return
}

func (m *mergeDrawList) Src(int) (x0, y0, x1, y1 int) {
	if x, y := m.a[m.si], m.b[m.sj]; x.Z() < y.Z() {
		x0, y0, x1, y1 = m.a.Src(m.si)
		m.si++
		return
	}
	x0, y0, x1, y1 = m.a.Src(m.sj)
	m.sj++
	return
}

func (m *mergeDrawList) draw(screen *ebiten.Image) error {
	return screen.DrawImage(composite, &ebiten.DrawImageOptions{
		ImageParts: m,
	})
}
*/

func (d drawList) draw(screen *ebiten.Image) error {
	return screen.DrawImage(composite, &ebiten.DrawImageOptions{
		ImageParts: d,
	})
}
