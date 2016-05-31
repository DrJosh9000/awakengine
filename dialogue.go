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

const dialogueZ = 100000

type ButtonSpec struct {
	Label  string
	Action func()
}

// DialogueLine is information for displaying a singe line of dialogue in a display.
type DialogueLine struct {
	Avatars *Sheet
	Index   int
	Text    string
	Buttons []ButtonSpec
}

// Dialogue is all the things needed for displaying blocking dialogue text.
type DialogueDisplay struct {
	bubble   *Bubble
	buttons  []*Button
	frame    int // frame number for this dialogue.
	text     *Text
	complete bool
	retire   bool
	avatar   *SheetFrame
	visible  bool
}

// DialogueFromLine creates a new DialogueDisplay.
func DialogueFromLine(line *DialogueLine) (*DialogueDisplay, error) {
	basePos := vec.I2{10, camSize.Y - 80}
	baseSize := vec.I2{camSize.X - 20, 70}
	textPos := basePos.Add(vec.I2{5, 5})
	var avatar *SheetFrame
	if line.Avatars != nil && line.Index >= 0 {
		// Provide space for the avatar.
		textPos.X += line.Avatars.FrameSize.X + 5
		avatar = &SheetFrame{
			Sheet: line.Avatars,
			F:     line.Index,
			P:     vec.I2{15, camSize.Y - 80 + 2},
		}
	}
	bk, _ := game.BubbleKey()
	d := &DialogueDisplay{
		frame:   0,
		avatar:  avatar,
		visible: true,
		text: &Text{
			Text: line.Text,
			Pos:  textPos,
			Size: vec.I2{camSize.X - textPos.X - 20, 0},
			Font: game.Font(),
		},
		bubble: &Bubble{
			ul:     basePos,
			dr:     basePos.Add(baseSize),
			imgkey: bk,
		},
	}
	d.bubble.Parent = Parent{d}
	d.text.Parent = Parent{d.bubble}
	d.text.Layout(false) // Rolls out the text for each Advance.
	p := vec.I2{textPos.X + 15, basePos.Y + baseSize.Y - 30}
	for _, s := range line.Buttons {
		d.buttons = append(d.buttons, NewButton(s.Label, s.Action, p, p.Add(vec.I2{50, 18}), Parent{d.bubble}))
		p.X += 65
	}
	return d, nil
}

func (d *DialogueDisplay) Fixed() bool   { return true }
func (d *DialogueDisplay) Retire() bool  { return d.retire }
func (d *DialogueDisplay) InWorld() bool { return false }
func (d *DialogueDisplay) Visible() bool { return d.visible }
func (d *DialogueDisplay) Z() int        { return dialogueZ }

func (d *DialogueDisplay) parts() drawList {
	l := d.bubble.parts()
	l = append(l, d.text.parts()...)
	if d.avatar != nil {
		l = append(l, &struct {
			*SheetFrame
			Parent
		}{d.avatar, Parent{d.bubble}})
	}
	for _, b := range d.buttons {
		l = append(l, b.parts()...)
	}
	return l
}

// Update updates things in the dialogue, based on user input or passage of time.
func (d *DialogueDisplay) Handle(event Event) (dismiss bool) {
	for _, b := range d.buttons {
		if b.Handle(event) {
			return true
		}
	}
	if event.Type == EventMouseUp {
		if d.complete && len(d.buttons) == 0 {
			return true
		}
		// Finish.
		d.complete = true
		for d.text.next < len(d.text.chars) {
			d.text.Advance()
		}
	}
	if !d.complete {
		d.text.Advance()
		d.text.Advance()
		if d.text.next >= len(d.text.chars) {
			d.complete = true
		}
	}
	d.frame++
	return
}
