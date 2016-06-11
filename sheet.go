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

// FrameInfo describes how to play a frame in an animation.
type FrameInfo struct {
	Duration int    // in model frames. -1 means infinite, 0 means 1
	Next     int    // next model frame index, no special meaning for 0...
	Offset   vec.I2 // subtract from position to get top-left of destination
}

// BasicFrameInfos is a convenience for making FrameInfos that all have the same
// offset and duration, and all have the next frame as the next frame.
func BasicFrameInfos(n, duration int, offset vec.I2) []FrameInfo {
	fi := make([]FrameInfo, n)
	for i := range fi {
		fi[i].Duration = duration
		fi[i].Next = (i + 1) % n
		fi[i].Offset = offset
	}
	return fi
}

type Sheet struct {
	Columns    int
	Key        string
	FrameInfos []FrameInfo
	FrameSize  vec.I2
}

// Src returns the source rectangle for frame number f.
func (s *Sheet) FrameSrc(f int) (x0, y0, x1, y1 int) {
	f %= len(s.FrameInfos)
	if s.Columns == 0 {
		x0, y0 = vec.NewI2(f, 0).EMul(s.FrameSize).C()
		x1, y1 = x0+s.FrameSize.X, y0+s.FrameSize.Y
		return
	}
	x0, y0 = vec.Div(f, s.Columns).EMul(s.FrameSize).C()
	x1, y1 = x0+s.FrameSize.X, y0+s.FrameSize.Y
	return
}

// SheetFrame lets you specify a frame in addition to a sheet.
type SheetFrame struct {
	*Sheet
	FrameNo int
}

func (s *SheetFrame) ImageKey() string          { return s.Sheet.Key }
func (s *SheetFrame) Src() (x0, y0, x1, y1 int) { return s.Sheet.FrameSrc(s.FrameNo) }

// ImageView displays one frame of a sheet filling a view.
type ImageView struct {
	*View
	*SheetFrame
}

func (v *ImageView) Src() (x0, y0, x1, y1 int) { return v.SheetFrame.Src() }
func (v *ImageView) Dst() (x0, y0, x1, y1 int) { return v.View.LogicalBounds().C() }
