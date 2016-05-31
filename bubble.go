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

const bubblePartSize = 5

// Bubble renders a bubble at any size larger than 15x15.
type Bubble struct {
	ul, dr vec.I2
	imgkey string
	Parent
}

func (b *Bubble) ImageKey() string { return b.imgkey }

func (b *Bubble) parts() drawList {
	l := make(drawList, 9)
	for i := 0; i < 9; i++ {
		l[i] = drawPosition{bubblePart{b, i}}
	}
	return l
}

type bubblePart struct {
	*Bubble
	i int
}

// Src implements ImageParts.
func (b bubblePart) Src() (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(b.i, 3).Mul(bubblePartSize).C()
	x1, y1 = x0+bubblePartSize, y0+bubblePartSize
	return
}

// Dst implements ImageParts.
func (b bubblePart) Dst() (x0, y0, x1, y1 int) {
	j, k := vec.Div(b.i, 3).C()
	switch j {
	case 0:
		x0 = b.ul.X
		x1 = b.ul.X + bubblePartSize
	case 1:
		x0 = b.ul.X + bubblePartSize
		x1 = b.dr.X - bubblePartSize
	case 2:
		x0 = b.dr.X - bubblePartSize
		x1 = b.dr.X
	}
	switch k {
	case 0:
		y0 = b.ul.Y
		y1 = b.ul.Y + bubblePartSize
	case 1:
		y0 = b.ul.Y + bubblePartSize
		y1 = b.dr.Y - bubblePartSize
	case 2:
		y0 = b.dr.Y - bubblePartSize
		y1 = b.dr.Y
	}
	return
}
