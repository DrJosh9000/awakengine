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

// Transient is a sprite that starts at a given birth.
type Transient struct {
	P     vec.I2
	Birth int
	A     *Anim
}

// Anim implements Sprite.
func (t *Transient) Anim() *Anim { return t.A }

// Frame implements Sprite.
func (t *Transient) Frame() int { return (gameFrame - t.Birth) / animPeriod }

// Pos implements Sprite.
func (t *Transient) Pos() vec.I2 { return t.P }

// Static just draws a frame.
type Static struct {
	P vec.I2
	F int
	A *Anim
}

// Anim implements Sprite.
func (s *Static) Anim() *Anim { return s.A }

// Frame implements Sprite.
func (s *Static) Frame() int { return s.F }

// Pos implements Sprite.
func (s *Static) Pos() vec.I2 { return s.P }

// Sprite is all the information required to draw an animated thingy at a point on screen.
type Sprite interface {
	Anim() *Anim
	Frame() int
	Pos() vec.I2 // logical / world position.
}

// SpritesByYPos orders Sprites by Y position (least to greatest).
type SpritesByYPos []Sprite

// Len implements sort.Interface.
func (b SpritesByYPos) Len() int { return len(b) }

// Less implements sort.Interface.
func (b SpritesByYPos) Less(i, j int) bool { return b[i].Pos().Y < b[j].Pos().Y }

// Swap implements sort.Interface.
func (b SpritesByYPos) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// SpriteParts implements ebiten.ImageParts for sprite drawing.
type SpriteParts struct {
	Sprite
	InWorld bool
}

// Draw draws the sprite to the screen.
func (s SpriteParts) Draw(screen *ebiten.Image) error {
	return screen.DrawImage(s.Anim().Image(), &ebiten.DrawImageOptions{ImageParts: s})
}

// Len implements ebiten.ImageParts.
func (s SpriteParts) Len() int { return 1 }

// Dst implements ebiten.ImageParts.
func (s SpriteParts) Dst(i int) (x0, y0, x1, y1 int) {
	a := s.Anim()
	b := s.Pos().Sub(a.Offset)
	if s.InWorld {
		b = b.Sub(camPos)
	}
	c := b.Add(a.FrameSize)
	return b.X, b.Y, c.X, c.Y
}

// Src implements ebiten.ImageParts.
func (s SpriteParts) Src(i int) (x0, y0, x1, y1 int) {
	a, f := s.Anim(), s.Frame()
	switch a.Mode {
	case AnimOneShot:
		if f >= a.Frames {
			return
		}
	case AnimLoop:
		f %= a.Frames
	}
	x0 = f * a.FrameSize.X
	return x0, 0, x0 + a.FrameSize.X, a.FrameSize.Y
}
