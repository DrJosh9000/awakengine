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

type Sheet struct {
	Columns   int
	Key       string
	Frames    int
	FrameSize vec.I2
}

func (s *Sheet) ImageKey() string { return s.Key }

// Src returns the source rectangle for frame number f.
func (s *Sheet) FrameSrc(f int) (x0, y0, x1, y1 int) {
	if s.Columns == 0 {
		s.Columns = s.Frames
	}
	if s.Frames == 0 {
		s.Frames = s.Columns
	}
	f %= s.Frames
	x0, y0 = vec.Div(f, s.Columns).EMul(s.FrameSize).C()
	x1, y1 = x0+s.FrameSize.X, y0+s.FrameSize.Y
	return
}

// Dst returns the destination rectangle with the top-left corner at position p.
func (s *Sheet) PosDst(p vec.I2) (x0, y0, x1, y1 int) {
	x0, y0 = p.C()
	x1, y1 = p.Add(s.FrameSize).C()
	return
}

// SheetFrame lets you specify a frame and a position in addition to a sheet.
type SheetFrame struct {
	*Sheet
	F int
	P vec.I2
}

func (s *SheetFrame) Src() (x0, y0, x1, y1 int) { return s.FrameSrc(s.F) }
func (s *SheetFrame) Dst() (x0, y0, x1, y1 int) { return s.PosDst(s.P) }
