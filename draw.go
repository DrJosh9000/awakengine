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
	"fmt"
	"sort"

	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
)

type Dst interface {
	Dst() (x0, y0, x1, y1 int)
}

type SumDst struct {
	d1, d2 Dst
}

func (a *SumDst) Dst() (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = a.d1.Dst()
	x2, y2, x3, y3 := a.d2.Dst()
	return x0 + x2, y0 + y2, x1 + x3, y1 + y3
}

type DrawSize vec.I2

func (d *DrawSize) Dst() (x0, y0, x1, y1 int) { return 0, 0, d.X, d.Y }

type Offset vec.I2

func (o *Offset) Dst() (x0, y0, x1, y1 int) { return o.X, o.Y, o.X, o.Y }

type Parent interface {
	Parent() Semiobject
}

// ScreenDst determines the final rectangle on screen for anything, accounting
// for parent offsets.
func ScreenDst(o interface{}) (x0, y0, x1, y1 int) {
	if o == nil {
		return
	}
	if r, ok := o.(Dst); ok {
		x0, y0, x1, y1 = r.Dst()
	}
	if r, ok := o.(Parent); ok && r != nil {
		x, y, _, _ := ScreenDst(r.Parent())
		return x0 + x, y0 + y, x1 + x, y1 + y
	}
	return
}

// drawPosition adjusts the source rectangle to refer to the texture atlas, and
// destination rectangle of relative objects to refer to screen coordinates.
type drawPosition struct{ Object }

func (p drawPosition) Dst() (x0, y0, x1, y1 int) { return ScreenDst(p.Object) }

func (p drawPosition) Src() (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = p.Object.Src()
	o, ok := compositeOffset[p.Object.ImageKey()]
	if !ok {
		panic(fmt.Sprintf("unknown image key %q", p.Object.ImageKey()))
	}
	return x0 + o.X, y0 + o.Y, x1 + o.X, y1 + o.Y
}

func (p drawPosition) Parent() Semiobject { return nil } // Once positioned, always in screen coordinates.
func (p drawPosition) ImageKey() string   { return "" }  // Once positioned, always in texture atlas coordinates.

// drawList is a Z-sortable list of objects in texture atlas/screen coordinates.
type drawList []drawPosition

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
func (d drawList) cull(dst drawList, scene *Scene) drawList {
	return d.subslice(dst, func(o Object) bool {
		if !o.Visible() {
			return false
		}
		if x0, y0, x1, y1 := o.Dst(); x1 <= 0 || y1 <= 0 || x0 > scene.CameraSize.X || y0 > scene.CameraSize.Y {
			return false
		}
		return true
	})
}

// shouldRetire recursively figures out if o should retire.
func shouldRetire(o Semiobject) bool {
	if o.Retire() {
		return true
	}
	if p := o.Parent(); p != nil {
		return shouldRetire(p)
	}
	return false
}

// gc removes retired objects. dst can be d[:0].
func (d drawList) gc(dst drawList) drawList {
	return d.subslice(dst, func(o Object) bool {
		return !shouldRetire(o)
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
