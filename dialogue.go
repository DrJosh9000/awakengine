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

type ButtonSpec struct {
	Label  string
	Action func()
}

// DialogueLine is information for displaying a singe line of dialogue in a display.
type DialogueLine struct {
	Avatar   *SheetFrame
	Text     string
	Buttons  []ButtonSpec
	AutoNext bool
	Slowness int
}

// Dialogue is all the things needed for displaying blocking dialogue text.
type DialogueDisplay struct {
	*View
	bubble   *Bubble
	buttons  []*Button
	frame    int // frame number for this dialogue.
	text     *Text
	complete bool

	line *DialogueLine
}

// DialogueFromLine creates a new DialogueDisplay.
func DialogueFromLine(line *DialogueLine, scene *Scene) *DialogueDisplay {
	camSize := scene.View.Size()
	basePos := vec.I2{10, camSize.Y - 74}
	baseSize := vec.I2{camSize.X - 20, 64}
	textPos := vec.I2{10, 10}
	if line.Avatar != nil {
		// Provide space for the avatar.
		textPos.X += line.Avatar.Sheet.FrameSize.X + 5
	}
	bk, _ := game.BubbleKey()
	d := &DialogueDisplay{
		View:  &View{},
		line:  line,
		frame: 0,
		text: &Text{
			View: &View{},
			Text: line.Text,
			Font: game.Font(),
		},
		bubble: &Bubble{
			View: &View{},
			Key:  bk,
		},
	}
	d.View.SetParent(scene.HUD)
	d.bubble.View.SetParent(d.View)
	d.bubble.View.SetPosition(basePos)
	d.bubble.View.SetSize(baseSize)
	d.bubble.View.SetOffset(bubblePartSize)
	d.text.View.SetParent(d.bubble.View)
	d.text.View.SetPosition(textPos)
	d.text.View.SetSize(vec.I2{camSize.X - textPos.X - 35, 0})
	d.text.Layout(false) // Rolls out the text for each Advance.
	p := vec.I2{textPos.X + 15, baseSize.Y - 20}
	for _, s := range line.Buttons {
		d.buttons = append(d.buttons, NewButton(
			s.Label,
			s.Action,
			vec.Rect{p, p.Add(vec.I2{40, 11})},
			d.bubble.View))
		p.X += 65
	}
	return d
}

func (d *DialogueDisplay) AddToScene(s *Scene) {
	d.bubble.AddToScene(s)
	d.text.AddToScene(s)
	for _, b := range d.buttons {
		b.AddToScene(s)
	}
	if d.line.Avatar == nil {
		return
	}
	s.AddPart(&struct {
		*SheetFrame
		*View
	}{d.line.Avatar, d.bubble.View})
}

func (d *DialogueDisplay) finish() {
	d.complete = true
	for d.text.next < len(d.text.chars) {
		d.text.Advance()
	}
}

// Update updates things in the dialogue, based on user input or passage of time.
// Returns true if the event is handled.
func (d *DialogueDisplay) Handle(event Event) bool {
	for _, b := range d.buttons {
		if b.Handle(event) {
			d.retire = true
			return true
		}
	}
	if d.complete && d.line.AutoNext {
		d.retire = true
		return true
	}
	if event.Type == EventMouseUp {
		if d.complete && len(d.buttons) == 0 {
			d.retire = true
			return true
		}
		if !d.line.AutoNext {
			d.finish()
		}
	}
	if !d.complete {
		if d.line.Slowness < 0 {
			d.finish()
		}
		if d.line.Slowness == 0 || d.frame%d.line.Slowness == 0 {
			d.text.Advance()
			if d.text.next >= len(d.text.chars) {
				d.complete = true
			}
		}
	}
	d.frame++
	return false
}
