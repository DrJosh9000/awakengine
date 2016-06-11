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

var bubblePartSize = vec.I2{5, 5}

// Bubble renders a bubble at any size larger than bubblePartSize.Mul(3)
type Bubble struct {
	*View
	Key   string
	added bool
}

func (b *Bubble) AddToScene(s *Scene) {
	if b.added {
		return
	}
	b.added = true
	s.AddPart(
		&bubblePart{b, 0}, &bubblePart{b, 1}, &bubblePart{b, 2},
		&bubblePart{b, 3}, &bubblePart{b, 4}, &bubblePart{b, 5},
		&bubblePart{b, 6}, &bubblePart{b, 7}, &bubblePart{b, 8},
	)
}

type bubblePart struct {
	*Bubble
	i int
}

func (b *bubblePart) ImageKey() string { return b.Bubble.Key }

func (b *bubblePart) Src() (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(b.i, 3).EMul(bubblePartSize).C()
	x1, y1 = x0+bubblePartSize.X, y0+bubblePartSize.X
	return
}

func (b *bubblePart) Dst() (x0, y0, x1, y1 int) {
	j, k := vec.Div(b.i, 3).C()
	x0, y0, x1, y1 = b.View.LogicalBounds().C()
	switch j {
	case 0:
		x1 = x0 + bubblePartSize.X
	case 1:
		x0 += bubblePartSize.X
		x1 -= bubblePartSize.X
	case 2:
		x0 = x1 - bubblePartSize.X
	}
	switch k {
	case 0:
		y1 = y0 + bubblePartSize.X
	case 1:
		y0 += bubblePartSize.X
		y1 -= bubblePartSize.X
	case 2:
		y0 = y1 - bubblePartSize.X
	}
	return
}
