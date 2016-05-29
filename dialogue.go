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

// DialogueLine is information for displaying a singe line of dialogue in a display.
type DialogueLine struct {
	Avatars *Sheet
	Index   int
	Text    string
}

// Dialogue is all the things needed for displaying blocking dialogue text.
type DialogueDisplay struct {
	bubble   *Bubble
	frame    int // frame number for this dialogue.
	text     *Text
	complete bool
	retire   bool
	avatar   *SheetFrame
	visible  bool
}

// DialogueFromLine creates a new DialogueDisplay.
func DialogueFromLine(line DialogueLine) (*DialogueDisplay, error) {
	textPos := vec.I2{20, camSize.Y - 80 + 5}
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
	d := &DialogueDisplay{
		frame:   0,
		avatar:  avatar,
		visible: true,
	}
	t, err := NewText(line.Text, camSize.X-textPos.X-20, textPos, game.Font(), d, 0)
	if err != nil {
		return nil, err
	}
	d.text = t
	d.bubble = &Bubble{
		pos:        vec.I2{10, camSize.Y - 80},
		sz:         vec.I2{camSize.X - 20, 70},
		imgkey:     game.BubbleKey(),
		deltaZ:     -1,
		Semiobject: d,
	}
	return d, nil
}

func (d *DialogueDisplay) Retire() bool  { return d.retire }
func (d *DialogueDisplay) InWorld() bool { return false }
func (d *DialogueDisplay) Visible() bool { return d.visible }
func (d *DialogueDisplay) Z() int        { return dialogueZ }

func (d *DialogueDisplay) parts() drawList {
	l := d.bubble.parts()
	if d.avatar != nil {
		l = append(l, &struct {
			*SheetFrame
			*DialogueDisplay
		}{d.avatar, d})
	}
	l = append(l, d.text.parts()...)
	return l
}

// Update updates things in the dialogue, based on user input or passage of time.
func (d *DialogueDisplay) Update(event Event) (dismiss bool) {
	if event.Type == EventMouseUp {
		if d.complete {
			dismiss = true
			return
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
