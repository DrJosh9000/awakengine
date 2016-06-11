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

func NewButton(text string, action func(), bounds vec.Rect, parent *View) *Button {
	sz := bounds.Size()
	bk, _ := game.BubbleKey()
	b := &Button{
		Action: action,
		Bubble: &Bubble{
			View: &View{},
			Key:  bk,
		},
		Text: &Text{
			View: &View{},
			Text: text,
			Font: game.Font(),
		},
	}
	b.Bubble.SetParent(parent)
	b.Bubble.SetZ(1)
	b.Text.SetParent(b.Bubble.View)
	b.Text.SetZ(1)
	// Initial size is inset from the bubble by the bubble part size on all sides.
	sz = sz.Sub(bubblePartSize.Mul(2))
	b.Text.SetSize(sz)
	b.Text.Layout(true)
	// Text should now have the minimal size. Centre the text within button, but offset slightly.
	b.Text.SetPosition(sz.Sub(b.Text.View.Size()).Div(2).Sub(vec.I2{1, 1}))
	return b
}

func (b *Button) Handle(e Event) (handled bool) {
	k1, k2 := game.BubbleKey()
	if b.Bubble.View.Bounds().Contains(e.ScreenPos) {
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
