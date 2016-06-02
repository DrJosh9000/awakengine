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

// drawList is a Z-sortable list of objects in texture atlas/screen coordinates.
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

// subslice appends the keep items to dst one at a time, returning the final slice.
// dst can be d[:0].
func (d drawList) subslice(dst drawList, keep func(Object) bool) drawList {
	for _, o := range d {
		if keep(o) {
			dst = append(dst, o)
		}
	}
	return dst
}

// cull removes invisible objects and places visible objects in dst. dst can be d[:0].
// Visibility is determined by calling Visible() and by testing the Dst rectangle.
func (d drawList) cull(dst drawList) drawList {
	return d.subslice(dst, func(o Object) bool {
		if !o.Visible() {
			return false
		}
		if x0, y0, x1, y1 := o.Dst(); x1 <= 0 || y1 <= 0 || x0 > camSize.X || y0 > camSize.Y {
			return false
		}
		return true
	})
}

// gc removes retired objects. dst can be d[:0].
func (d drawList) gc(dst drawList) drawList {
	return d.subslice(dst, func(o Object) bool {
		return !o.Retire()
	})
}

// merge merges two sorted drawLists into a combined list.
func merge(dst, a, b drawList) drawList {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if x, y := a[i], b[j]; x.Z() < y.Z() {
			dst = append(dst, x)
			i++
		} else {
			dst = append(dst, y)
			j++
		}
	}
	dst = append(dst, a[i:]...)
	dst = append(dst, b[j:]...)
	return dst
}

func (d drawList) draw(screen *ebiten.Image) error {
	return screen.DrawImage(composite, &ebiten.DrawImageOptions{
		ImageParts: d,
	})
}
