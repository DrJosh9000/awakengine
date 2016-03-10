package awakengine

import (
	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
)

var dialogueBubble *Bubble

// DialogueLine is information for displaying a singe line of dialogue in a display.
type DialogueLine struct {
	Avatars *Anim
    Frame int
	Text string
}

// Dialogue is all the things needed for displaying blocking dialogue text.
type DialogueDisplay struct {
	frame    int // frame number for this dialogue.
	text     *AdvancingText
	complete bool
	avatar   *Static
}

// NewDialogue creates a new DialogueDisplay.
func DialogueFromLine(line DialogueLine) (*DialogueDisplay, error) {
	textPos := vec.I2{20, camSize.Y - 80 + 5}
	if line.av != avatarNone {
		// Provide space for the avatar.
		textPos.X += avatarsAnim.frameSize.X + 5
	}
	t, err := NewText(line.text, camSize.X-textPos.X-20, textPos, false)
	if err != nil {
		return nil, err
	}
	return &DialogueDisplay{
		frame: 0,
		avatar: &Static{
			anim:  avatarsAnim,
			frame: int(line.av),
			pos:   vec.I2{15, camSize.Y - 80 + 2},
		},
		text: t,
	}, nil
}

// Draw draws the dialogue.
func (d *DialogueDisplay) Draw(screen *ebiten.Image) error {
	if err := dialogueBubble.Draw(screen); err != nil {
		return err
	}
	if d.avatar.frame >= 0 {
		if err := (SpriteParts{d.avatar, false}).Draw(screen); err != nil {
			return err
		}
	}
	if err := d.text.Draw(screen); err != nil {
		return err
	}
	return nil
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
		for d.text.idx < len(d.text.txt) {
			d.text.Advance()
		}
	}
	if !d.complete {
		if d.frame%2 == 0 {
			d.text.Advance()
		}
		if d.text.idx >= len(d.text.txt) {
			d.complete = true
		}
	}
	d.frame++
	return
}
