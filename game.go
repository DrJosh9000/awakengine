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
	"fmt"
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

	// One frame of animation for every (animPeriod) frames rendered.
	// So animation FPS = 60 / animPeriod.
	animPeriod = 3
	pixelSize  = 3
	camSize    = vec.I2{267, 150}
	camPos     = vec.I2{0, 0}
	title      = "AwakEngine"

	terrain          *Terrain
	obstacles, paths *vec.Graph

	triggers        map[string]*Trigger
	dialogueStack   []DialogueLine
	currentDialogue *DialogueDisplay

	player  Unit
	sprites []Sprite
)

// Unit can be told to update and provide information for drawing.
// Examples of units include the player character, NPCs, etc. Or it
// could be a unit in an RTS.
type Unit interface {
	// GoIdle asks the unit to stop whatever it's doing ("at ease").
	GoIdle()

	// Footprint is the rectangle relative to the sprite position with the ground area of the unit.
	Footprint() (ul, dr vec.I2)

	// Path is the path the unit is currently following. The current position is
	// implied as the first point.
	Path() []vec.I2

	// Sprite is here for drawing purposes.
	Sprite

	// Update asks the unit to update its own state, including the current event.
	Update(frame int, event Event)
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
	TileInfos() []TileInfo

	// Tiles is an image containing square tiles.
	Tiles() (key string, tileSize int)
}

// Game abstracts the non-engine parts of the game: the story, art, level design, etc.
type Game interface {
	// Font is the general/default typeface to use.
	Font() Font

	// Terrain provides the base level.
	Level() Level

	// Player provides the player unit.
	Player() Unit

	// Sprites provides all sprites in the level (include the player).
	Sprites() []Sprite

	// Triggers provide some dynamic behaviour.
	Triggers() map[string]*Trigger

	// Viewport is the size of the window and the pixels in the window.
	Viewport() (camSize vec.I2, pixelSize, animPeriod int, title string)
}

// load prepares assets for use by the game.
func load(g Game, debug bool) error {
	game = g
	Debug = debug
	camSize, pixelSize, animPeriod, title = game.Viewport()

	if err := loadAllImages(); err != nil {
		return fmt.Errorf("loading images: %v", err)
	}

	player = game.Player()
	sprites = game.Sprites()
	triggers = game.Triggers()

	t, err := loadTerrain(game.Level())
	if err != nil {
		return fmt.Errorf("loading terrain: %v", err)
	}
	terrain = t

	b, err := NewBubble(vec.I2{10, camSize.Y - 80}, vec.I2{camSize.X - 20, 70})
	if err != nil {
		return fmt.Errorf("loading bubble: %v", err)
	}
	dialogueBubble = b

	// TODO: compute unfattened static obstacles and fully dynamic paths to support
	// multiple units.
	// Invert the footprint to fatten the obstacles with.
	ul, dr := player.Footprint()
	ul = ul.Mul(-1)
	dr = dr.Mul(-1)
	obstacles, paths = t.ObstaclesAndPaths(dr, ul)
	return nil
}

// Run runs the game (ebiten.Run) in addition to setting up any necessary GIF recording.
func Run(g Game, debug bool, rf string, frameCount int) error {
	if err := load(g, debug); err != nil {
		return err
	}
	up := update
	if rf != "" {
		f, err := os.Create(rf)
		if err != nil {
			return fmt.Errorf("creating recording file: %v", err)
		}
		defer f.Close()
		up = ebitenutil.RecordScreenAsGIF(up, f, frameCount)
	}
	return ebiten.Run(up, camSize.X, camSize.Y, pixelSize, title)
}

/*
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
	if len(player.Path()) > 0 {
		u := player.Pos().Sub(camPos)
		for _, v := range player.Path() {
			v = v.Sub(camPos)
			if err := screen.DrawLine(u.X, u.Y, v.X, v.Y, color.RGBA{0, 0xff, 0xff, 0xff}); err != nil {
				return err
			}
			u = v
		}
	}
	return nil
}
*/

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
			if !trig.Fired && trig.Active(gameFrame) {
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
					trig.Fire(gameFrame)
				}
				dialogueStack = trig.Dialogues
				currentDialogue = nil
				player.GoIdle()
				if len(dialogueStack) > 0 {
					d, err := DialogueFromLine(dialogueStack[0])
					if err != nil {
						return err
					}
					currentDialogue = d
				}
				trig.Fired = true
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
	camPos = player.Pos().Sub(camSize.Div(2)).ClampLo(vec.I2{}).ClampHi(terrain.Size().Sub(camSize))

	// Draw all the things.
	if err := terrain.Draw(screen); err != nil {
		return err
	}

	// Tiny sort.
	sort.Sort(SpritesByYPos(sprites))
	for _, s := range sprites {
		if err := (SpriteParts{s, true}.Draw(screen)); err != nil {
			return err
		}
	}

	// Any doodads overlapping the player?
	pp := player.Pos()
	pu := pp.Sub(player.Anim().Offset)
	pd := pu.Add(player.Anim().FrameSize)
	for _, dd := range terrain.doodads {
		if pp.Y >= dd.P.Y {
			continue
		}
		tu := dd.P.Sub(dd.Anim().Offset)
		td := tu.Add(dd.Anim().FrameSize)
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

	/*
		// The W is special. All hail the W!
		wu := theW.pos.Sub(theW.Anim().Offset)
		wd := wu.Add(theW.Anim().FrameSize)
		cd := camPos.Add(camSize)
		if (wu.X < cd.X || wd.X >= camPos.X) && (wu.Y < cd.Y || wd.Y >= camPos.X) {
			if err := (SpriteParts{theW, true}.Draw(screen)); err != nil {
				return err
			}
		}
	*/

	if currentDialogue != nil {
		if err := currentDialogue.Draw(screen); err != nil {
			return err
		}
	}
	//return drawDebug(screen)
	return nil
}

// Navigate attempts to construct a path within the terrain.
func Navigate(from, to vec.I2) []vec.I2 {
	path, err := vec.FindPath(obstacles, paths, from, to, camPos, camPos.Add(camSize))
	if err != nil {
		// Go near to the cursor position.
		e, q := obstacles.NearestPoint(to)
		if Debug {
			log.Printf("nearest edge: %v point: %v", e, q)
		}
		q = q.Add(e.V.Sub(e.U).Normal().Sgn()) // Adjust it slightly...
		path2, err2 := vec.FindPath(obstacles, paths, from, q, camPos, camPos.Add(camSize))
		if err2 != nil {
			// Ok... Go as far as we can go.
			p2, y := obstacles.NearestBlock(from, to)
			if y {
				to = p2.Sub(p2.Sub(from).Sgn())
			}
			path2 = []vec.I2{to}
		}
		path = path2
	}
	if Debug {
		log.Printf("path: %v", path)
	}
	return path
}
