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
	*View
	World, HUD *View

	fixed       drawList
	fixedSorted bool
	loose       drawList
	dispFixed   drawList
	dispLoose   drawList
	dispMerged  drawList
}

func NewScene(camSize, terrainSize vec.I2) *Scene {
	s := &Scene{
		View:  &View{},
		World: &View{},
		HUD:   &View{},
	}
	s.View.SetSize(camSize)
	s.World.SetSize(terrainSize)
	s.World.SetParent(s.View)
	s.HUD.SetSize(camSize)
	s.HUD.SetParent(s.View)
	s.HUD.SetZ(100000) // HUD over World, always
	return s
}

// AddPart adds parts to the pipeline.
func (s *Scene) AddPart(parts ...Part) {
	for _, p := range parts {
		dp, ok := p.(drawPosition)
		if !ok {
			dp = drawPosition{p}
		}
		if p.Fixed() {
			s.fixedSorted = false
			s.fixed = append(s.fixed, dp)
		} else {
			s.loose = append(s.loose, dp)
		}
	}
}

func (s *Scene) sortFixedIfNeeded() {
	if s.fixedSorted {
		return
	}
	s.fixed.Sort()
	s.fixedSorted = true
}

// CameraFocus sets the World offset such that p should be center of screen, or at least
// within the bounds of the terrain.
func (s *Scene) CameraFocus(p vec.I2) {
	sz := s.View.Size()
	p = p.Sub(sz.Div(2))
	p = p.ClampLo(vec.I2{})
	p = p.ClampHi(s.World.Size().Sub(sz))
	s.World.SetOffset(p.Mul(-1))
}
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
