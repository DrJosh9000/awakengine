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

	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const levelGeomDumpFmt = `// Copyright 2016 Josh Deprez
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

package game

import "github.com/DrJosh9000/vec"

var (
	precomputedObstacles = %#v
	precomputedPaths = %#v
)`

var (
	config *Config

	game         Game
	modelFrame   int
	displayFrame int

	mouseDown     bool
	lastCursorPos vec.I2

	pixelSize = 3
	camSize   = vec.I2{267, 150}
	camPos    = vec.I2{0, 0}
	title     = "AwakEngine"

	terrain          *Terrain
	obstacles, paths *vec.Graph

	triggers      map[string]*Trigger
	dialogueStack []DialogueLine
	dialogue      *DialogueDisplay

	player          Unit
	sprites         []Sprite
	fixedObjects    drawList
	looseObjects    drawList
	displayedFixed  drawList
	displayedLoose  drawList
	displayedMerged drawList
)

type Config struct {
	Debug           bool
	FramesPerUpdate int
	LevelGeomDump   string
	RecordingFile   string
	RecordingFrames int
}

// Handler handles events.
type Handler interface {
	Handle(e Event)
}

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

	Pos() vec.I2
}

// Level describes things needed for a base terrain/level.
type Level struct {
	Doodads                 []*Doodad // sparse terrain objects
	MapSize                 vec.I2
	TileMap, BlockMap       []uint8
	TileInfos, BlockInfos   []TileInfo // dense terrain objects
	TilesetKey, BlocksetKey string
	TileSize, BlockHeight   int

	// Obstacles and Paths are optional but speed up game start time.
	Obstacles, Paths *vec.Graph
}

// Game abstracts the non-engine parts of the game: the story, art, level design, etc.
type Game interface {
	// BubbleKey returns the key for the bubble image, and inverse.
	BubbleKey() (string, string)

	// Font is the general/default typeface to use.
	Font() Font

	// Handle handles events.
	Handle(e Event)

	// Level provides the base level.
	Level() (*Level, error)

	// Objects provides non-terrain objects.
	// Do not include Doodads.
	// Include the player object and other sprites.
	Objects() []Object

	// Player provides the player unit.
	Player() Unit

	// Triggers provide some dynamic behaviour.
	Triggers() map[string]*Trigger

	// Viewport is the size of the window and the pixels in the window.
	Viewport() (camSize vec.I2, pixelSize int, title string)
}

// load prepares assets for use by the game.
func load(g Game) error {
	game = g
	camSize, pixelSize, title = game.Viewport()

	if err := loadAllImages(); err != nil {
		return fmt.Errorf("loading images: %v", err)
	}

	player = game.Player()
	triggers = game.Triggers()

	l, err := game.Level()
	if err != nil {
		return fmt.Errorf("loading level: %v", err)
	}

	t, err := loadTerrain(l)
	if err != nil {
		return fmt.Errorf("loading terrain: %v", err)
	}
	terrain = t

	obstacles, paths = l.Obstacles, l.Paths
	if obstacles == nil || paths == nil {
		// TODO: compute unfattened static obstacles and fully dynamic paths to support
		// multiple units.
		// Invert the footprint to fatten the obstacles with.
		ul, dr := player.Footprint()
		ul = ul.Mul(-1)
		dr = dr.Mul(-1)
		obstacles, paths = t.ObstaclesAndPaths(dr, ul)
		if config.LevelGeomDump != "" {
			f, err := os.Create(config.LevelGeomDump)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = fmt.Fprintf(f, levelGeomDumpFmt, obstacles, paths)
			err = f.Close()
			if err != nil {
				return err
			}
		}
	}

	dd := make([]Object, 0, len(l.Doodads))
	for _, d := range l.Doodads {
		dd = append(dd, d)
	}

	fixedObjects, looseObjects = makeDrawLists(game.Objects(), terrain.parts(), dd)
	fixedObjects.Sort()
	return nil
}

// Run runs the game (ebiten.Run) in addition to setting up any necessary GIF recording.
func Run(g Game, cfg *Config) error {
	config = cfg
	if err := load(g); err != nil {
		return err
	}
	up := update
	if cfg.RecordingFile != "" {
		f, err := os.Create(cfg.RecordingFile)
		if err != nil {
			return fmt.Errorf("creating recording file: %v", err)
		}
		defer f.Close()
		up = ebitenutil.RecordScreenAsGIF(up, f, cfg.RecordingFrames)
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

// modelUpdate does update stuff, but no drawing. It is called once per config.FramesPerUpdate.
func modelUpdate() error {
	// Read inputs
	md := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if md {
		lastCursorPos = vec.NewI2(ebiten.CursorPosition())
	}
	tt := ebiten.Touches()
	if len(tt) > 0 {
		md = true
		lastCursorPos = vec.NewI2(tt[0].Position())
	}
	e := Event{
		Time:      modelFrame,
		Pos:       lastCursorPos.Add(camPos),
		MouseDown: md,
	}
	switch {
	case md && !mouseDown:
		e.Type = EventMouseDown
	case !md && mouseDown:
		e.Type = EventMouseUp
	}
	mouseDown = md

	// TODO: What did they just click on?

	// Do we proceed with the game, or with the dialogue display?
	if dialogue == nil {
		// Got any triggers?
		for k, trig := range triggers {
			if !trig.Fired && trig.Active(modelFrame) {
				// All dependencies fired?
				for _, dep := range trig.Depends {
					if !triggers[dep].Fired {
						continue
					}
				}
				if config.Debug {
					log.Printf("firing %s with %d dialogues", k, len(trig.Dialogues))
				}
				if trig.Fire != nil {
					trig.Fire(modelFrame)
				}
				dialogueStack = trig.Dialogues
				dialogue = nil
				player.GoIdle()
				if len(dialogueStack) > 0 {
					d, err := DialogueFromLine(&dialogueStack[0])
					if err != nil {
						return err
					}
					dialogue = d
					fixedObjects = append(fixedObjects, dialogue.parts()...)
					fixedObjects.Sort()
				}
				trig.Fired = true
				break
			}
		}
		if dialogue == nil {
			modelFrame++

			game.Handle(e)
			for _, o := range looseObjects {
				if u, ok := o.Object.(Sprite); ok {
					u.Update(modelFrame)
				}
			}
		}
	} else if dialogue.Handle(e) {
		// Play
		dialogue.retire = true
		dialogueStack = dialogueStack[1:]
		dialogue = nil
		if len(dialogueStack) > 0 {
			d, err := DialogueFromLine(&dialogueStack[0])
			if err != nil {
				return err
			}
			dialogue = d
			fixedObjects = append(fixedObjects, dialogue.parts()...)
			fixedObjects.Sort()
		}
	}

	// Update camera to focus on player.
	camPos = player.Pos().Sub(camSize.Div(2)).ClampLo(vec.I2{}).ClampHi(terrain.Size().Sub(camSize))

	fixedObjects = fixedObjects.gc(fixedObjects[:0])
	looseObjects = looseObjects.gc(looseObjects[:0])
	displayedFixed = fixedObjects.cull(displayedFixed[:0])
	displayedLoose = looseObjects.cull(displayedLoose[:0])
	displayedLoose.Sort()
	displayedMerged = merge(displayedMerged[:0], displayedFixed, displayedLoose)
	if config.Debug {
		log.Printf("{len, cap}(fixedObjects): %d, %d", len(fixedObjects), cap(fixedObjects))
		log.Printf("{len, cap}(looseObjects): %d, %d", len(looseObjects), cap(looseObjects))
		log.Printf("{len, cap}(displayedFixed): %d, %d", len(displayedFixed), cap(displayedFixed))
		log.Printf("{len, cap}(displayedLoose): %d, %d", len(displayedLoose), cap(displayedLoose))
		log.Printf("{len, cap}(displayedMerged): %d, %d", len(displayedMerged), cap(displayedMerged))
	}
	return nil
}

// update is the main update function.
func update(screen *ebiten.Image) error {
	displayFrame++
	if displayFrame%config.FramesPerUpdate == 0 {
		if err := modelUpdate(); err != nil {
			return err
		}
	}
	return displayedMerged.draw(screen) // One draw call.
}

// Navigate attempts to construct a path within the terrain.
func Navigate(from, to vec.I2) []vec.I2 {
	path, err := vec.FindPath(obstacles, paths, from, to, camPos, camPos.Add(camSize))
	if err != nil {
		// Go near to the cursor position.
		e, q := obstacles.NearestPoint(to)
		if config.Debug {
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
	if config.Debug {
		log.Printf("path: %v", path)
	}
	return path
}
