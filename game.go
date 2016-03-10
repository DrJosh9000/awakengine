package awakengine

import (
	"image/color"
	"log"
	"os"
	"sort"

	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
    "github.com/hajimehoshi/ebiten/ebitenutil"
)

const windowTitle = "A walk in the park"

var (
	// Debug controls display of debug graphics.
	Debug bool
	// LevelPreview enables showing the whole level, with no triggers.
	LevelPreview bool

	gameFrame int

	mouseDown bool

	pixelSize    = 3
	camSize      = vec.I2{267, 150}
	camSizeTiles = camSize.Add(vec.I2{tileSize - 1, tileSize - 1}).Div(tileSize).Add(vec.I2{1, 1})
	camPos       = vec.I2{0, 0} // top left corner, pixels.

	goalAckMarker = Transient{
		Anim: &Anim{
			Key:       "mark",
			Offset:    vec.I2{15, 15},
			Frames:    4,
			FrameSize: vec.I2{32, 32},
			Mode:      AnimOneShot,
		},
		birth: -999,
	}

	obstacles, paths *vec.Graph

	dialogueStack   []DialogueLine
	currentDialogue *DialogueDisplay

	playerLayer = []Sprite{
		player,
		&goalAckMarker,
	}
)

// Load prepares assets for use by the game.
func Load() error {
	if LevelPreview {
		pixelSize = 1
		camSize = vec.I2{1024, 1024}
		triggers = nil
	}
	if Debug {
		log.Printf("camSize: %v, camSizeTiles: %v", camSize, camSizeTiles)
	}

	if err := LoadAllImages(); err != nil {
		return err
	}

	b, err := NewBubble(vec.I2{10, camSize.Y - 80}, vec.I2{camSize.X - 20, 70})
	if err != nil {
		return err
	}
	dialogueBubble = b

	if err := LoadTerrain(); err != nil {
		return err
	}

	obstacles, paths = terrain.ObstaclesAndPaths(playerFatUL, playerFatDR)
	return nil
}

// Run runs the game (ebiten.Run) in addition to setting up any necessary GIF recording.
func Run(rf string, frameCount int) error {
	up := update
	if rf != "" {
		f, err := os.Create(rf)
		if err != nil {
			return err
		}
		defer f.Close()
		up = ebitenutil.RecordScreenAsGIF(up, f, frameCount)
	}
	return ebiten.Run(up, camSize.X, camSize.Y, pixelSize, windowTitle)
}

// drawDebug draws debugging graphics onto the screen if Debug is true.
func drawDebug(screen *ebiten.Image) error {
	if !Debug {
		return nil
	}
	obsView := GraphView{
		edges:        obstacles.Edges(),
		edgeColour:   color.RGBA{0xff, 0, 0, 0xff},
		normalColour: color.RGBA{0, 0xff, 0, 0xff},
	}
	if err := screen.DrawLines(obsView); err != nil {
		return err
	}
	pathsView := GraphView{
		edges:        paths.Edges(),
		edgeColour:   color.RGBA{0, 0, 0xff, 0xff},
		normalColour: color.Transparent,
	}
	if err := screen.DrawLines(pathsView); err != nil {
		return err
	}
	if len(player.path) > 0 {
		u := player.Pos().Sub(camPos)
		for _, v := range player.path {
			v = v.Sub(camPos)
			if err := screen.DrawLine(u.X, u.Y, v.X, v.Y, color.RGBA{0, 0xff, 0xff, 0xff}); err != nil {
				return err
			}
			u = v
		}
	}
	return nil
}

// update is the main update function.
func update(screen *ebiten.Image) error {
	// Read inputs
	md := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	e := Event{Pos: vec.NewI2(ebiten.CursorPosition()).Add(camPos)}
	switch {
	case md && !mouseDown:
		mouseDown = true
		e.Type = EventMouseDown
	case !md && mouseDown:
		mouseDown = false
		e.Type = EventMouseUp
	}

	// Do we proceed with the game, or with the dialogue display?
	if currentDialogue == nil {
		// Got any triggers?
		for k, trig := range triggers {
			if !trig.fired && trig.active() {
				// All dependencies fired?
				for _, dep := range trig.depends {
					if !triggers[dep].fired {
						continue
					}
				}
				if Debug {
					log.Printf("firing %s with %d dialogues", k, len(trig.dialogues))
				}
				if trig.fire != nil {
					trig.fire()
				}
				dialogueStack = trig.dialogues
				currentDialogue = nil
				player.state.a = playerIdle
				player.path = nil
				if len(dialogueStack) > 0 {
					d, err := DialogueFromLine(dialogueStack[0])
					if err != nil {
						return err
					}
					currentDialogue = d
				}
				trig.fired = true
				break
			}
		}
		if currentDialogue == nil {
			gameFrame++
			player.Update(gameFrame, e)
		}
	} else if currentDialogue.Update(e) {
		// Play
		dialogueStack = dialogueStack[1:]
		currentDialogue = nil
		if len(dialogueStack) > 0 {
			d, err := DialogueFromLine(dialogueStack[0])
			if err != nil {
				return err
			}
			currentDialogue = d
		}
	}

	// Update camera to focus on player.
	camPos = player.pos.I2().Sub(camSize.Div(2)).ClampLo(vec.I2{0, 0}).ClampHi(terrain.size.Mul(tileSize).Sub(camSize))

	// Draw all the things.
	if err := terrain.Draw(screen); err != nil {
		return err
	}

	// Tiny sort.
	sort.Sort(ByYPos(playerLayer))
	for _, s := range playerLayer {
		if err := (SpriteParts{s, true}.Draw(screen)); err != nil {
			return err
		}
	}

	// Any doodads overlapping the player?
	pp := player.Pos()
	pu := pp.Sub(player.Anim().offset)
	pd := pu.Add(player.Anim().frameSize)
	for _, dd := range terrain.doodads {
		if pp.Y >= dd.pos.Y {
			continue
		}
		tu := dd.pos.Sub(dd.Anim().offset)
		td := tu.Add(dd.Anim().frameSize)
		if tu.Y > pd.Y || td.Y < pu.Y {
			// td.Y < pu.Y is essentially given, but consistency.
			continue
		}
		if tu.X > pd.X || td.X < pu.X {
			continue
		}
		if err := (SpriteParts{dd, true}.Draw(screen)); err != nil {
			return err
		}
	}

	// The W is special. All hail the W!
	wu := theW.pos.Sub(theW.Anim().offset)
	wd := wu.Add(theW.Anim().frameSize)
	cd := camPos.Add(camSize)
	if (wu.X < cd.X || wd.X >= camPos.X) && (wu.Y < cd.Y || wd.Y >= camPos.X) {
		if err := (SpriteParts{theW, true}.Draw(screen)); err != nil {
			return err
		}
	}

	if currentDialogue != nil {
		if err := currentDialogue.Draw(screen); err != nil {
			return err
		}
	}
	return drawDebug(screen)
	//return nil
}
