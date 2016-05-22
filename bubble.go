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

// Bubble is an ImageParts that renders a bubble at any size larger than 15x15.
type Bubble struct {
	pos, sz vec.I2
	imgkey  string
	deltaZ  int
	//flat    *ebiten.Image
	Semiobject // parent
}

/*
// NewBubble prepares a bubble of the correct size.
func NewBubble(pos, size vec.I2, imgkey string) (*Bubble, error) {
	/
		f, err := ebiten.NewImage(size.X, size.Y, ebiten.FilterNearest)
		if err != nil {
			return nil, fmt.Errorf("creating image: %v", err)
		}
		if err := Draw(f, imgkey, AllBubbleParts(size)); err != nil {
			return nil, fmt.Errorf("drawing bubble: %v", err)
		}
	return &Bubble{
		pos:    pos,
		sz:     size,
		imgkey: imgkey,
		//flat: f,
	}, nil
}
*/

func (b *Bubble) ImageKey() string { return b.imgkey }
func (b *Bubble) Z() int           { return b.Semiobject.Z() + b.deltaZ }

func (b *Bubble) parts() drawList {
	l := make(drawList, 9)
	for i := 0; i < 9; i++ {
		l[i] = bubblePart{b, i}
	}
	return l
}

/*

// Src implements ImageParts.
func (b *Bubble) Src() (x0, y0, x1, y1 int) {
	x1, y1 = b.dr.Sub(b.ul).C()
	return
}

// Dst implements ImageParts.
func (b *Bubble) Dst() (x0, y0, x1, y1 int) {
	return b.ul.X, b.ul.Y, b.dr.X, b.dr.Y
}

// Draw draws the bubble to the screen.
func (b *Bubble) Draw(screen *ebiten.Image) error {
	return screen.DrawImage(b.flat, &ebiten.DrawImageOptions{ImageParts: b})
}
*/

type bubblePart struct {
	*Bubble
	i int
}

// Src implements ImageParts.
func (b bubblePart) Src() (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(b.i, 3).Mul(bubblePartSize).C()
	return x0, y0, x0 + bubblePartSize, y0 + bubblePartSize
}

// Dst implements ImageParts.
func (b bubblePart) Dst() (x0, y0, x1, y1 int) {
	j, k := vec.Div(b.i, 3).C()
	switch j {
	case 0:
		x1 = bubblePartSize
	case 1:
		x0 = bubblePartSize
		x1 = b.sz.X - bubblePartSize
	case 2:
		x0 = b.sz.X - bubblePartSize
		x1 = b.sz.X
	}
	switch k {
	case 0:
		y1 = bubblePartSize
	case 1:
		y0 = bubblePartSize
		y1 = b.sz.Y - bubblePartSize
	case 2:
		y0 = b.sz.Y - bubblePartSize
		y1 = b.sz.Y
	}
	return x0 + b.pos.X, y0 + b.pos.Y, x1 + b.pos.X, y1 + b.pos.Y
}
