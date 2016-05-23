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

// Anim describes an animated sprite that might play.
type Anim struct {
	Key       string
	Offset    vec.I2
	Frames    int
	FrameSize vec.I2
}

func (a *Anim) ImageKey() string { return a.Key }

type Sprite interface {
	Anim() *Anim
	Pos() vec.I2
	Frame() int
	Update(t int)
}

// TODO: implement particles

type SpriteObject struct {
	Sprite
	Semiobject
}

func (s SpriteObject) ImageKey() string { return s.Anim().ImageKey() }

func (s SpriteObject) Dst() (x0, y0, x1, y1 int) {
	a := s.Anim()
	b := s.Pos().Sub(a.Offset)
	c := b.Add(a.FrameSize)
	return b.X, b.Y, c.X, c.Y
}

func (s SpriteObject) Src() (x0, y0, x1, y1 int) {
	a, f := s.Anim(), s.Frame()
	f %= a.Frames
	x0 = f * a.FrameSize.X
	return x0, 0, x0 + a.FrameSize.X, a.FrameSize.Y
}

// StaticSprite just displays whatever frame number it is given.
type StaticSprite struct {
	A *Anim
	F int
	P vec.I2
}

func (s *StaticSprite) Anim() *Anim { return s.A }
func (s *StaticSprite) Frame() int  { return s.F }
func (s *StaticSprite) Pos() vec.I2 { return s.P }
func (s *StaticSprite) Update(int)  {}
