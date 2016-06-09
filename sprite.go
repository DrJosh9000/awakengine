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

type StaticOffset vec.I2

func (s StaticOffset) Offset(int) vec.I2 { return vec.I2(s) }

type StaticPlayback int

func (s StaticPlayback) Frame() int { return int(s) }

// Playback describes playing an animation according to different frame durations and looping.
type Playback struct {
	SF, DF        int
	FrameDuration []int // model frames to spend in each animation frame, assumes 1 for each frame otherwise
	LoopTo        int   // return to this frame number when complete
}

func (p *Playback) Reset()     { p.SF, p.DF = 0, -1 }
func (p *Playback) Frame() int { return p.SF }
func (p *Playback) Update(int) {
	p.DF++
	if p.DF < p.FrameDuration[p.SF] {
		return
	}
	p.DF = 0
	p.SF++
	if p.SF >= len(p.FrameDuration) {
		p.SF = p.LoopTo
	}
}

type Sprite interface {
	// Templatey things.
	ImageKey() string
	Dst() (x0, y0, x1, y1 int)
	FrameSrc(frame int) (x0, y0, x1, y1 int)
	Offset(frame int) vec.I2

	// Instancey things.
	Frame() int
	Pos() vec.I2
	Update(t int)
}

type SpriteObject struct {
	Sprite
}

func (s SpriteObject) Dst() (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = s.Sprite.Dst()
	p := s.Pos().Sub(s.Offset(s.Frame()))
	return x0 + p.X, y0 + p.Y, x1 + p.X, y1 + p.Y
}
func (s SpriteObject) Src() (x0, y0, x1, y1 int) { return s.FrameSrc(s.Frame()) }

// StaticSprite just displays whatever frame number it is given, forever.
type StaticSprite struct {
	*Sheet
	StaticOffset
	StaticPlayback
	P vec.I2
}

func (s *StaticSprite) Pos() vec.I2 { return s.P }
func (s *StaticSprite) Update(int)  {}
