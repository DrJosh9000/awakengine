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
}

func NewButton(text string, ul, dr vec.I2, parent Parent) *Button {
	sz := dr.Sub(ul)
	b := &Button{
		Bubble: &Bubble{
			ul:     ul,
			dr:     dr,
			imgkey: game.BubbleKey(),
			Parent: parent,
		},
		Text: &Text{
			Text: text,
			Size: sz,
			Font: game.Font(),
		},
	}
	b.Text.Parent = Parent{b.Bubble}
	b.Text.Layout(true)
	b.Text.Pos = ul.Add(sz.Sub(b.Text.Size).Div(2)) // Centre text within button.
	return b
}

func (b *Button) parts() drawList {
	return append(b.Bubble.parts(), b.Text.parts()...)
}
