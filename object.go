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

type Drawable interface {
	ImageKey() string
	Dst() (x0, y0, x1, y1 int)
	Src() (x0, y0, x1, y1 int) // relative to the image referred to by ImageKey()
}

type Rect struct{ X0, Y0, X1, Y1 int }

func (r *Rect) C() (x0, y0, x1, y1 int) { return r.X0, r.Y0, r.X1, r.Y1 }

// StaticDrawable implements Drawable with struct fields.
type StaticDrawable struct {
	Key  string
	S, D Rect
}

func (s *StaticDrawable) ImageKey() string          { return s.Key }
func (s *StaticDrawable) Dst() (x0, y0, x1, y1 int) { return s.D.C() }
func (s *StaticDrawable) Src() (x0, y0, x1, y1 int) { return s.S.C() }

// Semiobject is some import stuff for logical grouping.
type Semiobject interface {
	Fixed() bool   // true if the object never moves - X, Y, or Z (but can still be world-relative).
	InWorld() bool // true if the object exists in world-coordinates, false if screen coordinates
	Retire() bool  // true if the object will never draw again and can be removed from the draw list
	Visible() bool
	Z() int
}

// StaticSemiobject implements Semiobjects with struct fields.
type StaticSemiobject struct {
	F, IW, R, V bool
	Zed         int
}

func (s *StaticSemiobject) Fixed() bool   { return s.F }
func (s *StaticSemiobject) InWorld() bool { return s.IW }
func (s *StaticSemiobject) Retire() bool  { return s.R }
func (s *StaticSemiobject) Visible() bool { return s.V }
func (s *StaticSemiobject) Z() int        { return s.Zed }

// Object is everything, everything is an object.
type Object interface {
	Drawable
	Semiobject
}

type Parent struct{ Semiobject }

func (p Parent) Z() int { return p.Semiobject.Z() + 1 }

type drawList []Object

func (d drawList) Len() int           { return len(d) }
func (d drawList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d drawList) Less(i, j int) bool { return d[i].Z() < d[j].Z() }

func (d drawList) Sort() { sort.Sort(d) }

// subslice reorders d so the non-kept objects are at the end,
// and provides the subslice of remaining elements.
func (d drawList) subslice(keep func(int) bool) drawList {
	b := len(d)
	for i := 0; i < b; {
		if keep(i) {
			i++
			continue
		}
		b--
		d[i], d[b] = d[b], d[i]
	}
	return d[:b]
}

// cull removes invisible objects.
func (d drawList) cull() drawList {
	return d.subslice(func(i int) bool {
		if !d[i].Visible() {
			return false
		}
		if x0, y0, x1, y1 := d.Dst(i); x1 <= 0 || y1 <= 0 || x0 > camSize.X || y0 > camSize.Y {
			return false
		}
		return true
	})
}

// gc removes retired objects.
func (d drawList) gc() drawList {
	return d.subslice(func(i int) bool {
		return !d[i].Retire()
	})
}

func (d drawList) Dst(i int) (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = d[i].Dst()
	if !d[i].InWorld() {
		return
	}
	x0 -= camPos.X
	y0 -= camPos.Y
	x1 -= camPos.X
	y1 -= camPos.Y
	return
}

func (d drawList) Src(i int) (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = d[i].Src()
	o := compositeOffset[d[i].ImageKey()]
	x0 += o.X
	y0 += o.Y
	x1 += o.X
	y1 += o.Y
	return
}

func (d drawList) draw(screen *ebiten.Image) error {
	return screen.DrawImage(composite, &ebiten.DrawImageOptions{
		ImageParts: d,
	})
}
