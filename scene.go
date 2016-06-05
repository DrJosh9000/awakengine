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
	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
)

// Scene manages drawing one scene.
//
// Using a Scene as a parent will subtract the camera position. To use screen coordinates,
// use a parent of nil.
type Scene struct {
	CameraPos  vec.I2
	CameraSize vec.I2

	fixed          drawList
	fixedNeedsSort bool
	loose          drawList
	dispFixed      drawList
	dispLoose      drawList
	dispMerged     drawList
}

// AddObject adds objects to the pipeline.
func (s *Scene) AddObject(objs ...Object) {
	for _, o := range objs {
		dp, ok := o.(drawPosition)
		if !ok {
			dp = drawPosition{o}
		}
		if o.Fixed() {
			s.fixedNeedsSort = true
			s.fixed = append(s.fixed, dp)
		} else {
			s.loose = append(s.loose, dp)
		}
	}
}

func (s *Scene) sortFixedIfNeeded() {
	if !s.fixedNeedsSort {
		return
	}
	s.fixed.Sort()
	s.fixedNeedsSort = false
}

func (s *Scene) CameraFocus(p vec.I2) {
	s.CameraPos = p.Sub(s.CameraSize.Div(2)).ClampLo(vec.I2{}).ClampHi(terrain.Size().Sub(s.CameraSize))
}
func (s *Scene) Dst() (x0, y0, x1, y1 int) { return -s.CameraPos.X, -s.CameraPos.Y, 0, 0 }

func (s *Scene) Draw(screen *ebiten.Image) error { return s.dispMerged.draw(screen) }

func (s *Scene) Update() {
	// Reorganise objects to display.
	s.fixed = s.fixed.gc(s.fixed[:0])
	s.loose = s.loose.gc(s.loose[:0])
	s.sortFixedIfNeeded()
	s.dispFixed = s.fixed.cull(s.dispFixed[:0], s)
	s.dispLoose = s.loose.cull(s.dispLoose[:0], s)
	s.dispLoose.Sort()
	s.dispMerged = merge(s.dispMerged[:0], s.dispFixed, s.dispLoose)
}

func (s *Scene) Fixed() bool        { return true }
func (s *Scene) Parent() Semiobject { return nil }
func (s *Scene) Retire() bool       { return false }
func (s *Scene) Visible() bool      { return true }
func (s *Scene) Z() int             { return 0 }
