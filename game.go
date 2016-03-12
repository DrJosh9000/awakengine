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

var (
	// Debug controls display of debug graphics.
	Debug bool

	game      Game
	gameFrame int

	mouseDown bool

	pixelSize = 3
	camSize   = vec.I2{267, 150}
	camPos    = vec.I2{0, 0}
	title     = "AwakEngine"

	terrain          *Terrain
	obstacles, paths *vec.Graph

	triggers        map[string]*Trigger
	dialogueStack   []DialogueLine
	currentDialogue *DialogueDisplay

	units []Unit
)

// Unit can be told to update and provide information for drawing.
// Examples of units include the player character, NPCs, etc.
type Unit interface {
	Sprite                         // for drawing
	Update(frame int, event Event) // time moves on, so compute new state
}

// Level abstracts things needed for a base terrain/level.
type Level interface {
	// Doodads provides objects above the base, that can be flattened onto the terrain
	// most of the time.
	Doodads() []*Doodad

	// Source is the paletted PNG to use as the base terrain layer - pixel at (x,y) becomes
	// the tile at (x,y).
	Source() string

	// TileInfos maps indexes to information about the terrain.
	TileInfos() map[uint8]TileInfo

	// Tiles is an image containing square tiles.
	Tiles() (key string, tileSize int)
}

// Game abstracts the non-engine parts of the game: the story, art, level design, etc.
type Game interface {
	// Terrain provides the base level.
	Level() Level

	// Triggers provide some dynamic behaviour.
	Triggers() map[string]*Trigger

	// Units provides all units in the level.
	Units() []Unit

	// Viewport is the size of the window and the pixels in the window.
	Viewport() (camSize vec.I2, pixelSize int, title string)
}

// Load prepares assets for use by the game.
func Load(g Game, debug bool) error {
	game = g
	Debug = debug
	camSize, pixelSize, title = game.Viewport()

	if err := loadAllImages(); err != nil {
		return err
	}

	b, err := NewBubble(vec.I2{10, camSize.Y - 80}, vec.I2{camSize.X - 20, 70})
	if err != nil {
		return err
	}
	dialogueBubble = b

	t, err := loadTerrain(game.Level())
	if err != nil {
		return err
	}
	terrain = t
	obstacles, paths = t.ObstaclesAndPaths(playerFatUL, playerFatDR)

	triggers = game.Triggers()
	units = game.Units()
	return nil
}

// Run runs the game (ebiten.Run) i n addition to setting up any necessary GIF recording.
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
	return ebiten.Run(up, camSize.X, camSize.Y, pixelSize, title)
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
			if !trig.Fired && trig.Active() {
				// All dependencies fired?
				for _, dep := range trig.Depends {
					if !triggers[dep].Fired {
						continue
					}
				}
				if Debug {
					log.Printf("firing %s with %d dialogues", k, len(trig.Dialogues))
				}
				if trig.Fire != nil {
					trig.Fire()
				}
				dialogueStack = trig.Dialogues
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
