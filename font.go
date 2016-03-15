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

const (
	munroHeight  = 11
	munroYOffset = -2
)

var (
	munroMap = map[byte]charInfo{
		' ':  {width: 0, x: 1, y: 9, xoffset: 0, height: 0, yoffset: 11, xadvance: 3},
		'!':  {width: 1, x: 2, y: 2, xoffset: 0, height: 7, yoffset: 4, xadvance: 2},
		'"':  {width: 3, x: 4, y: 1, xoffset: 0, height: 2, yoffset: 3, xadvance: 4},
		'#':  {width: 5, x: 8, y: 3, xoffset: 0, height: 6, yoffset: 5, xadvance: 6},
		'$':  {width: 3, x: 14, y: 1, xoffset: 0, height: 9, yoffset: 3, xadvance: 4},
		'%':  {width: 7, x: 18, y: 2, xoffset: 0, height: 7, yoffset: 4, xadvance: 8},
		'&':  {width: 6, x: 26, y: 2, xoffset: 0, height: 7, yoffset: 4, xadvance: 7},
		'\'': {width: 1, x: 33, y: 1, xoffset: 0, height: 2, yoffset: 3, xadvance: 2},
		'(':  {width: 2, x: 35, y: 1, xoffset: 0, height: 9, yoffset: 3, xadvance: 3},
		')':  {width: 2, x: 38, y: 1, xoffset: 0, height: 9, yoffset: 3, xadvance: 3},
		'*':  {width: 3, x: 41, y: 1, xoffset: 0, height: 3, yoffset: 3, xadvance: 4},
		'+':  {width: 5, x: 45, y: 3, xoffset: 0, height: 5, yoffset: 5, xadvance: 6},
		',':  {width: 1, x: 51, y: 8, xoffset: 0, height: 2, yoffset: 10, xadvance: 2},
		'-':  {width: 3, x: 53, y: 6, xoffset: 0, height: 1, yoffset: 8, xadvance: 4},
		'.':  {width: 1, x: 57, y: 8, xoffset: 0, height: 1, yoffset: 10, xadvance: 2},
		'/':  {width: 3, x: 59, y: 1, xoffset: 0, height: 9, yoffset: 3, xadvance: 4},
		'0':  {width: 4, x: 63, y: 2, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'1':  {width: 2, x: 1, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 3},
		'2':  {width: 3, x: 4, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'3':  {width: 3, x: 8, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'4':  {width: 4, x: 12, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'5':  {width: 3, x: 17, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'6':  {width: 4, x: 21, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'7':  {width: 3, x: 26, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'8':  {width: 4, x: 30, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'9':  {width: 4, x: 35, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		':':  {width: 1, x: 40, y: 14, xoffset: 0, height: 4, yoffset: 7, xadvance: 2},
		';':  {width: 1, x: 42, y: 14, xoffset: 0, height: 5, yoffset: 7, xadvance: 2},
		'<':  {width: 5, x: 44, y: 12, xoffset: 0, height: 5, yoffset: 5, xadvance: 5},
		'=':  {width: 5, x: 50, y: 13, xoffset: 0, height: 3, yoffset: 6, xadvance: 5},
		'>':  {width: 5, x: 56, y: 12, xoffset: 0, height: 5, yoffset: 5, xadvance: 5},
		'?':  {width: 3, x: 62, y: 11, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'@':  {width: 7, x: 1, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 8},
		'A':  {width: 4, x: 9, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'B':  {width: 4, x: 14, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'C':  {width: 3, x: 19, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'D':  {width: 4, x: 23, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'E':  {width: 3, x: 28, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'F':  {width: 3, x: 32, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'G':  {width: 4, x: 36, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'H':  {width: 4, x: 41, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'I':  {width: 1, x: 46, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 2},
		'J':  {width: 2, x: 48, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 3},
		'K':  {width: 4, x: 51, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'L':  {width: 3, x: 56, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'M':  {width: 5, x: 60, y: 20, xoffset: 0, height: 7, yoffset: 4, xadvance: 6},
		'N':  {width: 4, x: 1, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'O':  {width: 4, x: 6, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'P':  {width: 4, x: 11, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'Q':  {width: 4, x: 16, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'R':  {width: 4, x: 21, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'S':  {width: 3, x: 26, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'T':  {width: 3, x: 30, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'U':  {width: 4, x: 34, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'V':  {width: 5, x: 39, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 6},
		'W':  {width: 5, x: 45, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 6},
		'X':  {width: 5, x: 51, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 6},
		'Y':  {width: 5, x: 57, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 6},
		'Z':  {width: 3, x: 63, y: 28, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'[':  {width: 2, x: 1, y: 36, xoffset: 0, height: 11, yoffset: 2, xadvance: 3},
		'\\': {width: 3, x: 4, y: 37, xoffset: 0, height: 9, yoffset: 3, xadvance: 4},
		']':  {width: 2, x: 8, y: 36, xoffset: 0, height: 11, yoffset: 2, xadvance: 3},
		'^':  {width: 3, x: 11, y: 37, xoffset: 0, height: 2, yoffset: 3, xadvance: 4},
		'_':  {width: 5, x: 15, y: 45, xoffset: 0, height: 1, yoffset: 11, xadvance: 6},
		'a':  {width: 4, x: 21, y: 40, xoffset: 0, height: 5, yoffset: 6, xadvance: 5},
		'b':  {width: 4, x: 26, y: 38, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'c':  {width: 3, x: 31, y: 40, xoffset: 0, height: 5, yoffset: 6, xadvance: 4},
		'd':  {width: 4, x: 35, y: 38, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'e':  {width: 4, x: 40, y: 40, xoffset: 0, height: 5, yoffset: 6, xadvance: 5},
		'f':  {width: 3, x: 45, y: 38, xoffset: 0, height: 7, yoffset: 4, xadvance: 3},
		'g':  {width: 4, x: 49, y: 40, xoffset: 0, height: 7, yoffset: 6, xadvance: 5},
		'h':  {width: 4, x: 54, y: 38, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'i':  {width: 1, x: 59, y: 38, xoffset: 0, height: 7, yoffset: 4, xadvance: 2},
		'j':  {width: 2, x: 61, y: 38, xoffset: -1, height: 9, yoffset: 4, xadvance: 2},
		'k':  {width: 4, x: 64, y: 38, xoffset: 0, height: 7, yoffset: 4, xadvance: 5},
		'l':  {width: 1, x: 1, y: 48, xoffset: 0, height: 7, yoffset: 4, xadvance: 2},
		'm':  {width: 7, x: 3, y: 50, xoffset: 0, height: 5, yoffset: 6, xadvance: 8},
		'n':  {width: 4, x: 11, y: 50, xoffset: 0, height: 5, yoffset: 6, xadvance: 5},
		'o':  {width: 4, x: 16, y: 50, xoffset: 0, height: 5, yoffset: 6, xadvance: 5},
		'p':  {width: 4, x: 21, y: 50, xoffset: 0, height: 7, yoffset: 6, xadvance: 5},
		'q':  {width: 5, x: 26, y: 50, xoffset: 0, height: 7, yoffset: 6, xadvance: 5},
		'r':  {width: 3, x: 32, y: 50, xoffset: 0, height: 5, yoffset: 6, xadvance: 4},
		's':  {width: 3, x: 36, y: 50, xoffset: 0, height: 5, yoffset: 6, xadvance: 4},
		't':  {width: 3, x: 40, y: 48, xoffset: 0, height: 7, yoffset: 4, xadvance: 4},
		'u':  {width: 4, x: 44, y: 50, xoffset: 0, height: 5, yoffset: 6, xadvance: 5},
		'v':  {width: 5, x: 49, y: 50, xoffset: 0, height: 5, yoffset: 6, xadvance: 6},
		'w':  {width: 7, x: 55, y: 50, xoffset: 0, height: 5, yoffset: 6, xadvance: 8},
		'x':  {width: 5, x: 63, y: 50, xoffset: 0, height: 5, yoffset: 6, xadvance: 6},
		'y':  {width: 4, x: 1, y: 62, xoffset: 0, height: 7, yoffset: 6, xadvance: 5},
		'z':  {width: 3, x: 6, y: 62, xoffset: 0, height: 5, yoffset: 6, xadvance: 4},
		'{':  {width: 3, x: 10, y: 59, xoffset: 0, height: 9, yoffset: 3, xadvance: 4},
		'|':  {width: 1, x: 14, y: 58, xoffset: 1, height: 11, yoffset: 2, xadvance: 4},
		'}':  {width: 3, x: 16, y: 59, xoffset: 0, height: 9, yoffset: 3, xadvance: 4},
		'~':  {width: 5, x: 20, y: 63, xoffset: 0, height: 3, yoffset: 7, xadvance: 6},
	}
)

type charInfo struct {
	x, y, width, height, xoffset, yoffset, xadvance int
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
	txt string
	idx int
	adv []vec.I2
}

// Len implements ImageParts.
func (s *oneCharacter) Len() int { return 1 }

// Src implements ImageParts.
func (s *oneCharacter) Src(int) (x0, y0, x1, y1 int) {
	ci := munroMap[s.txt[s.idx]]
	return ci.x, ci.y, ci.x + ci.width, ci.y + ci.height
}

// Dst implements ImageParts.
func (s *oneCharacter) Dst(int) (x0, y0, x1, y1 int) {
	ci := munroMap[s.txt[s.idx]]
	t := s.adv[s.idx].Add(vec.I2{ci.xoffset, ci.yoffset + munroYOffset})
	return t.X, t.Y, t.X + ci.width, t.Y + ci.height
}

// Advance draws another character to the flat image.
func (s *AdvancingText) Advance() error {
	if s.idx < len(s.txt) {
		if err := s.flat.DrawImage(allImages["munro"], &ebiten.DrawImageOptions{ImageParts: s.oneCharacter}); err != nil {
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
func NewText(txt string, width int, pos vec.I2, inWorld bool) (*AdvancingText, error) {
	adv := make([]vec.I2, len(txt))
	x, y := 0, 0
	wordStart := 0
	wrapIt := func(end int) {
		if x >= width {
			x = 0
			y += munroHeight
			// Fix previous word.
			for j := wordStart; j < end; j++ {
				adv[j] = vec.I2{x, y}
				ci := munroMap[txt[j]]
				x += ci.xadvance
			}
		}
	}
	for i := range txt {
		if txt[i] == '\n' {
			x = 0
			y += munroHeight
			wordStart = i + 1
			continue
		}
		if txt[i] == ' ' {
			wrapIt(i)
			wordStart = i + 1
		}
		adv[i] = vec.I2{x, y}
		ci := munroMap[txt[i]]
		x += ci.xadvance
	}
	wrapIt(len(txt))
	f, err := ebiten.NewImage(width, y+munroHeight, ebiten.FilterNearest)
	if err != nil {
		return nil, err
	}
	return &AdvancingText{
		Pos:     pos,
		InWorld: inWorld,
		flat:    f,
		size:    vec.I2{width, y + munroHeight},
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
