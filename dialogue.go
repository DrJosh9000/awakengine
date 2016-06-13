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

import (
	"log"

	"github.com/DrJosh9000/vec"
)

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
	avatar  *ImageView
	bubble  *Bubble
	buttons []*Button
	text    *Text

	complete bool
	frame    int // frame number for this dialogue.

	line *DialogueLine
}

// NewDialogueDisplay creates a new DialogueDisplay.
func NewDialogueDisplay(scene *Scene) *DialogueDisplay {
	bk, _ := game.BubbleKey()
	//_, bk := game.BubbleKey()
	d := &DialogueDisplay{
		View: &View{},
		avatar: &ImageView{
			View: &View{},
		},
		text: &Text{
			View: &View{},
			Font: game.Font(),
		},
		bubble: &Bubble{
			View: &View{},
			Key:  bk,
		},
	}

	camSize := scene.View.Size()
	size := vec.I2{camSize.X - 10, 84}

	d.SetParent(scene.HUD)
	d.SetPositionAndSize(vec.I2{5, camSize.Y - 89}, size)

	d.bubble.SetParent(d.View)
	d.bubble.SetSize(size)
	d.bubble.SetZ(1)

	d.avatar.SetParent(d.bubble.View)
	d.avatar.SetPosition(bubblePartSize)
	d.avatar.SetZ(1)

	d.text.SetParent(d.bubble.View)
	d.text.SetZ(1)
	return d
}

// Layout rearranges views
func (d *DialogueDisplay) Layout(line *DialogueLine) {
	// Dispose of any old buttons first.
	for _, b := range d.buttons {
		b.Dispose()
	}
	d.buttons = nil

	// Reset things...
	d.complete = false
	d.frame = 0

	// Refresh properties from the line.
	d.avatar.SheetFrame = line.Avatar
	d.frame = 0
	d.line = line
	d.text.Text = line.Text

	size := d.Size()
	textPos := vec.I2{10, 10}

	if line.Avatar != nil {
		// Provide space for the avatar.
		textPos.X += line.Avatar.Sheet.FrameSize.X + 5
		d.avatar.SetSize(line.Avatar.Sheet.FrameSize)
	}
	d.avatar.SetVisible(line.Avatar != nil)

	d.text.SetPositionAndSize(textPos, vec.I2{size.X - textPos.X - 15, 0})
	d.text.Layout(line.Slowness < 0)

	p := vec.I2{textPos.X + 15, size.Y - 40}
	for _, s := range line.Buttons {
		btn := NewButton(
			s.Label,
			s.Action,
			vec.Rect{p, p.Add(vec.I2{65, 25})},
			d.bubble.View)
		d.buttons = append(d.buttons, btn)
		p.X += 75
	}
}

func (d *DialogueDisplay) AddToScene(scene *Scene) {
	d.bubble.AddToScene(scene)
	d.text.AddToScene(scene)
	for _, b := range d.buttons {
		b.AddToScene(scene)
	}
	scene.AddPart(d.avatar)
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
			// log.Printf("dialogue: button handled event")
			return true
		}
	}
	if d.complete && d.line.AutoNext {
		// log.Printf("dialogue: complete and autonext")
		return true
	}
	if event.Type == EventMouseUp {
		if d.complete && len(d.buttons) == 0 {
			// log.Printf("dialogue: clicked, complete, and no buttons")
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

// Management of global dialogue state (dialogue, dialogueStack)

func playNextDialogue() {
	if len(dialogueStack) == 0 {
		if dialogue != nil {
			if config.Debug {
				log.Printf("disposing a dialogue")
			}
			dialogue.Dispose()
		}
		dialogue = nil
		return
	}
	if dialogue == nil {
		if config.Debug {
			log.Printf("creating a dialogue")
		}
		dialogue = NewDialogueDisplay(scene)
	}
	if config.Debug {
		log.Printf("laying out a dialogue")
	}
	dialogue.Layout(dialogueStack[0])
	dialogue.AddToScene(scene)
	dialogueStack = dialogueStack[1:]
}

// PushDialogueToBack makes some dialogue the dialogue to play after all the current dialogue is finished.
func PushDialogueToBack(dl ...*DialogueLine) {
	dialogueStack = append(dialogueStack, dl...)
}

// PushDialogue makes some dialogue the next dialogue to play.
func PushDialogue(dl ...*DialogueLine) {
	dialogueStack = append(dl, dialogueStack...)
}
