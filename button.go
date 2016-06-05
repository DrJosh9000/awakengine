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

type Button struct {
	*Bubble
	*Text
	Action func()
}

func NewButton(text string, action func(), ul, dr vec.I2, par ChildOf) *Button {
	sz := dr.Sub(ul)
	bk, _ := game.BubbleKey()
	b := &Button{
		Action: action,
		Bubble: &Bubble{
			UL:      ul,
			DR:      dr,
			Key:     bk,
			ChildOf: par,
		},
		Text: &Text{
			Text: text,
			Size: sz.Sub(vec.I2{10, 10}),
			Font: game.Font(),
		},
	}
	b.Text.ChildOf = ChildOf{b.Bubble}
	b.Text.Layout(true)
	b.Text.Pos = ul.Add(sz.Sub(b.Text.Size).Div(2)) // Centre text within button.
	return b
}

func (b *Button) Handle(e Event) (handled bool) {
	k1, k2 := game.BubbleKey()
	if e.Pos.InRect(b.Bubble.UL, b.Bubble.DR) {
		switch {
		case e.MouseDown:
			b.Text.Invert = true
			b.Bubble.Key = k2
		case e.Type == EventMouseUp:
			b.Text.Invert = false
			b.Bubble.Key = k1
			b.Action()
			return true
		}
		return false
	}
	b.Text.Invert = false
	b.Bubble.Key = k1
	return false
}

func (b *Button) AddToScene(s *Scene) {
	b.Bubble.AddToScene(s)
	b.Text.AddToScene(s)
}
