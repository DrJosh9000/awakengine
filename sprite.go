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

import "github.com/DrJosh9000/vec"

type SpriteDelegate interface {
	// Instancey things
	Fixed(s *Sprite) bool
	SpriteSheet(s *Sprite) *Sheet    // Which sheet we're currently using
	Update(s *Sprite, modelTime int) // Do model updates
	Z(s *Sprite) int                 // Given our position, what's the right Z?
}

// Sprite concerns itself with displaying the right frame of an animation at the right time,
// at the right position, and at the right Z order.
type Sprite struct {
	*View        // Container for drawPosition; can be a new view.
	Pos   vec.F2 // Position
	f     int    // current frame of sprite sheet
	fd    int    // duration of current frame
	SpriteDelegate
}

func (s *Sprite) ResetAnim() { s.f, s.fd = 0, -1 }
func (s *Sprite) AdvanceAnim() {
	sheet := s.SpriteSheet(s)
	infos := sheet.FrameInfos
	if s.f < 0 || s.f >= len(infos) {
		s.ResetAnim()
	}
	s.fd++
	dur := infos[s.f].Duration
	if dur < 0 || s.fd < dur { // duration is infinite, or we aren't there yet.
		return
	}
	s.fd = 0
	s.f = infos[s.f].Next
}

func (s *Sprite) Update(t int) {
	s.AdvanceAnim()
	s.SpriteDelegate.Update(s, t)
}

func (s *Sprite) Container() *View { return s.View }
func (s *Sprite) ImageKey() string { return s.SpriteSheet(s).Key }

func (s *Sprite) Src() (x0, y0, x1, y1 int) {
	return s.SpriteSheet(s).FrameSrc(s.f)
}

func (s *Sprite) Dst() (x0, y0, x1, y1 int) {
	sheet := s.SpriteSheet(s)
	if s.f >= len(sheet.FrameInfos) {
		s.f = 0
	}
	info := &sheet.FrameInfos[s.f]
	ul := s.Pos.I2().Sub(info.Offset)
	x0, y0 = ul.C()
	x1, y1 = ul.Add(s.SpriteSheet(s).FrameSize).C()
	return
}

func (s *Sprite) Fixed() bool   { return s.SpriteDelegate.Fixed(s) }
func (s *Sprite) Retire() bool  { return s.View.Retire() }
func (s *Sprite) Visible() bool { return s.View.Visible() }
func (s *Sprite) Z() int        { return s.SpriteDelegate.Z(s) }
