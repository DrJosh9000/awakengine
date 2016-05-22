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
	ImageKey() string
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
	Semiobject
	deltaZ int
	txt    string
	chars  []oneChar
	next   int
}

func (s *Text) parts() drawList {
	l := make(drawList, len(s.chars))
	for i := range s.chars {
		l[i] = &s.chars[i]
	}
	return l
}

func (s *Text) Z() int { return s.Semiobject.Z() + s.deltaZ }

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

// NewText computes a new Text. You *must* specify width.
func NewText(txt string, width int, pos vec.I2, font Font, parent Semiobject, deltaZ int) (*Text, error) {
	text := &Text{
		Pos:        pos,
		Font:       font,
		Semiobject: parent,
		deltaZ:     deltaZ,
		txt:        txt,
		next:       0,
	}
	text.layout(width)
	return text, nil
}

func (s *Text) layout(width int) {
	chars := make([]oneChar, 0, len(s.txt))
	cm := s.Metrics()
	x, y := 0, 0
	wordStartC, wordStartI := 0, 0 // chars index, txt index
	wrapIt := func(end int) {
		if x < width {
			return
		}
		x = 0
		y += s.LineHeight()
		// Fix previous word.
		for i, j := wordStartC, wordStartI; j < end; i, j = i+1, j+1 {
			c := s.txt[j]
			ci := cm[c]
			chars[i].pos = vec.I2{x, y}
			x += ci.XAdvance
		}
	}
	for i := range s.txt {
		if s.txt[i] == '\n' {
			x = 0
			y += s.LineHeight()
			wordStartC = len(chars)
			wordStartI = i + 1
			continue
		}
		c := s.txt[i]
		ci := cm[c]
		if s.txt[i] == ' ' {
			wrapIt(i)
			wordStartC = len(chars)
			wordStartI = i + 1
			x += ci.XAdvance
			continue
		}
		chars = append(chars, oneChar{
			Text: s,
			pos:  vec.I2{x, y},
			c:    c,
		})
		x += ci.XAdvance
	}
	wrapIt(len(s.txt))
	s.chars = chars
	s.Size = vec.I2{width, y + s.LineHeight()}
}
