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

type CharMetrics map[byte]CharInfo

type Font interface {
	ImageKey(invert bool) string
	Metrics() CharMetrics
	LineHeight() int
	YOffset() int
}

type CharInfo struct {
	X, Y, Width, Height, XOffset, YOffset, XAdvance int
}

type Text struct {
	Pos, Size vec.I2
	Font
	Parent
	Text   string
	Invert bool
	chars  []oneChar
	next   int
}

func (s *Text) ImageKey() string { return s.Font.ImageKey(s.Invert) }

func (s *Text) parts() drawList {
	l := make(drawList, len(s.chars))
	for i := range s.chars {
		l[i] = drawPosition{&s.chars[i]}
	}
	return l
}

func (s *Text) Z() int { return s.Semiobject.Z() + 1 }

type oneChar struct {
	*Text
	pos     vec.I2
	c       byte
	visible bool
}

// Src implements ImageParts.
func (s *oneChar) Src() (x0, y0, x1, y1 int) {
	m := s.Metrics()
	ci := m[s.c]
	return ci.X, ci.Y, ci.X + ci.Width, ci.Y + ci.Height
}

// Dst implements ImageParts.
func (s *oneChar) Dst() (x0, y0, x1, y1 int) {
	m := s.Metrics()
	ci := m[s.c]
	x0, y0 = s.Text.Pos.X+s.pos.X+ci.XOffset, s.Text.Pos.Y+s.pos.Y+ci.YOffset+s.YOffset()
	return x0, y0, x0 + ci.Width, y0 + ci.Height
}

func (s *oneChar) Visible() bool { return s.visible && s.Text.Visible() }

// Advance makes the next character visible.
func (s *Text) Advance() error {
	if s.next < len(s.chars) {
		s.chars[s.next].visible = true
	}
	s.next++
	return nil
}

// Layout causes the text to lay out all the characters, and update
// the size to exactly contain the text. Text will be wrapped to the
// existing Size.X as a width.
func (s *Text) Layout(visible bool) {
	width := s.Size.X
	maxW := 0
	chars := make([]oneChar, 0, len(s.Text))
	cm := s.Metrics()
	x, y := 0, 0
	wordStartC, wordStartI := 0, 0 // chars index, Text index
	wrapIt := func(end int) {
		if x < width {
			return
		}
		if x > maxW {
			maxW = x
		}
		x = 0
		y += s.LineHeight()
		// Fix previous word.
		for i, j := wordStartC, wordStartI; j < end; i, j = i+1, j+1 {
			c := s.Text[j]
			ci := cm[c]
			chars[i].pos = vec.I2{x, y}
			x += ci.XAdvance
		}
	}
	for i := range s.Text {
		if s.Text[i] == '\n' {
			x = 0
			y += s.LineHeight()
			wordStartC = len(chars)
			wordStartI = i + 1
			continue
		}
		c := s.Text[i]
		ci := cm[c]
		if s.Text[i] == ' ' {
			wrapIt(i)
			wordStartC = len(chars)
			wordStartI = i + 1
			x += ci.XAdvance
			continue
		}
		chars = append(chars, oneChar{
			Text:    s,
			pos:     vec.I2{x, y},
			c:       c,
			visible: visible,
		})
		x += ci.XAdvance
	}
	wrapIt(len(s.Text))
	if x > maxW {
		maxW = x
	}
	s.chars = chars
	s.Size = vec.I2{maxW, y + s.LineHeight()}
}
