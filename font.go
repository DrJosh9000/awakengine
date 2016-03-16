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

type CharMetrics map[byte]CharInfo

type Font interface {
	Source() string
	Metrics() CharMetrics
	LineHeight() int
	YOffset() int
}

type CharInfo struct {
	X, Y, Width, Height, XOffset, YOffset, XAdvance int
}

// AdvancingText renders text, animated typewriter-style.
type AdvancingText struct {
	Pos, size vec.I2
	InWorld   bool

	*oneCharacter
	flat *ebiten.Image
}

// oneCharacter lays out an entire string, then renders one character.
type oneCharacter struct {
	Font
	txt string
	idx int
	adv []vec.I2
}

// Len implements ImageParts.
func (s *oneCharacter) Len() int { return 1 }

// Src implements ImageParts.
func (s *oneCharacter) Src(int) (x0, y0, x1, y1 int) {
	m := s.Font.Metrics()
	ci := m[s.txt[s.idx]]
	return ci.X, ci.Y, ci.X + ci.Width, ci.Y + ci.Height
}

// Dst implements ImageParts.
func (s *oneCharacter) Dst(int) (x0, y0, x1, y1 int) {
	m := s.Font.Metrics()
	ci := m[s.txt[s.idx]]
	t := s.adv[s.idx].Add(vec.I2{ci.XOffset, ci.YOffset + s.Font.YOffset()})
	return t.X, t.Y, t.X + ci.Width, t.Y + ci.Height
}

// Advance draws another character to the flat image.
func (s *AdvancingText) Advance() error {
	if s.idx < len(s.txt) {
		if err := s.flat.DrawImage(allImages[s.Font.Source()], &ebiten.DrawImageOptions{ImageParts: s.oneCharacter}); err != nil {
			return err
		}
	}
	s.idx++
	return nil
}

// Draw draws the string to the screen.
func (s *AdvancingText) Draw(screen *ebiten.Image) error {
	return screen.DrawImage(s.flat, &ebiten.DrawImageOptions{ImageParts: s})
}

// NewText computes a new AdvancingText. You *must* specify width.
func NewText(txt string, width int, pos vec.I2, font Font, inWorld bool) (*AdvancingText, error) {
	adv := make([]vec.I2, len(txt))
	cm := font.Metrics()
	x, y := 0, 0
	wordStart := 0
	wrapIt := func(end int) {
		if x >= width {
			x = 0
			y += font.LineHeight()
			// Fix previous word.
			for j := wordStart; j < end; j++ {
				adv[j] = vec.I2{x, y}
				ci := cm[txt[j]]
				x += ci.XAdvance
			}
		}
	}
	for i := range txt {
		if txt[i] == '\n' {
			x = 0
			y += font.LineHeight()
			wordStart = i + 1
			continue
		}
		if txt[i] == ' ' {
			wrapIt(i)
			wordStart = i + 1
		}
		adv[i] = vec.I2{x, y}
		ci := cm[txt[i]]
		x += ci.XAdvance
	}
	wrapIt(len(txt))
	f, err := ebiten.NewImage(width, y+font.LineHeight(), ebiten.FilterNearest)
	if err != nil {
		return nil, err
	}
	return &AdvancingText{
		Pos:     pos,
		InWorld: inWorld,
		flat:    f,
		size:    vec.I2{width, y + font.LineHeight()},
		oneCharacter: &oneCharacter{
			txt: txt,
			adv: adv,
			idx: 0,
		},
	}, nil
}

// Len implements ImageParts.
func (s *AdvancingText) Len() int { return 1 }

// Src implements ImageParts.
func (s *AdvancingText) Src(int) (x0, y0, x1, y1 int) {
	x1, y1 = s.size.C()
	return
}

// Dst implements ImageParts.
func (s *AdvancingText) Dst(int) (x0, y0, x1, y1 int) {
	t := s.Pos
	if s.InWorld {
		t = t.Sub(camPos)
	}
	x0, y0 = t.C()
	x1, y1 = t.Add(s.size).C()
	return
}
