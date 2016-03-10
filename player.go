package awakengine

import (
	"log"

	"github.com/DrJosh9000/vec"
)

const playerSpeed = 1

var (
	playerAnims = map[PlayerState]*Anim{
		{playerWalking, vec.Left}: &Anim{
			imgKey:    "walk_l",
			offset:    vec.I2{8, 31},
			frames:    7,
			frameSize: vec.I2{16, 32},
			mode:      AnimLoop,
		},
		{playerWalking, vec.Right}: &Anim{
			imgKey:    "walk_r",
			offset:    vec.I2{7, 31},
			frames:    7,
			frameSize: vec.I2{16, 32},
			mode:      AnimLoop,
		},
		{playerWalking, vec.Up}: &Anim{
			imgKey:    "walk_u",
			offset:    vec.I2{7, 27},
			frames:    14,
			frameSize: vec.I2{16, 33},
			mode:      AnimLoop,
		},
		{playerWalking, vec.Down}: &Anim{
			imgKey:    "walk_d",
			offset:    vec.I2{7, 27},
			frames:    14,
			frameSize: vec.I2{16, 33},
			mode:      AnimLoop,
		},
		{playerIdle, vec.Left}: &Anim{
			imgKey:    "idle_l",
			offset:    vec.I2{7, 31},
			frames:    1,
			frameSize: vec.I2{16, 32},
			mode:      AnimLoop,
		},
		{playerIdle, vec.Up}: &Anim{
			imgKey:    "idle_l",
			offset:    vec.I2{7, 31},
			frames:    1,
			frameSize: vec.I2{16, 32},
			mode:      AnimLoop,
		},
		{playerIdle, vec.Right}: &Anim{
			imgKey:    "idle_r",
			offset:    vec.I2{9, 31},
			frames:    1,
			frameSize: vec.I2{16, 32},
			mode:      AnimLoop,
		},
		{playerIdle, vec.Down}: &Anim{
			imgKey:    "idle_r",
			offset:    vec.I2{9, 31},
			frames:    1,
			frameSize: vec.I2{16, 32},
			mode:      AnimLoop,
		},
	}

	player = &Player{
		pos: vec.F2{32 * 9, 32 * 10},
	}
	playerFatUL = vec.I2{-4, -1}
	playerFatDR = vec.I2{3, 5}
)

// PlayerState describes the current simple state of the player: what are they doing,
// and what direction are they facing.
type PlayerState struct {
	a PlayerActivity
	d vec.Direction
}

// PlayerActivity describes what the player is doing.
type PlayerActivity int

const (
	playerIdle = PlayerActivity(iota)
	playerWalking
)

// Player encapsulates all the state of the player character.
type Player struct {
	pos   vec.F2
	frame int
	path  []vec.I2
	state PlayerState
}

// Anim implements Sprite.
func (p *Player) Anim() *Anim { return playerAnims[p.state] }

// Frame implements Sprite.
func (p *Player) Frame() int { return p.frame }

// Pos implements Sprite.
func (p *Player) Pos() vec.I2 { return p.pos.I2() }

// Update responds to the passage of time, or user input by adjusting the player state.
func (p *Player) Update(frame int, event Event) {
	if event.Type == EventMouseUp {
		c := event.Pos
		goalAckMarker.birth = frame
		goalAckMarker.pos = c
		path, err := vec.FindPath(obstacles, paths, p.Pos(), c, camPos, camPos.Add(camSize))
		if err != nil {
			// Go near to the cursor position.
			e, q := obstacles.NearestPoint(c)
			if Debug {
				log.Printf("nearest edge: %v point: %v", e, q)
			}
			q = q.Add(e.V.Sub(e.U).Normal().Sgn()) // Adjust it slightly...
			path2, err2 := vec.FindPath(obstacles, paths, p.Pos(), q, camPos, camPos.Add(camSize))
			if err2 != nil {
				// Ok... Go as far as we can go.
				p2, y := obstacles.NearestBlock(p.Pos(), c)
				if y {
					c = p2.Sub(p2.Sub(p.Pos()).Sgn())
				}
				path2 = []vec.I2{c}
			}
			path = path2
		}
		p.path = path
		if Debug {
			log.Printf("path: %v", path)
		}
	}

	if len(p.path) > 0 {
		d := p.path[0].F2().Sub(p.pos)
		np := d.Norm()

		// If at the current point, either aim at the next point, or stop.
		if np < 1 {
			if len(p.path) > 1 {
				p.pos = p.path[0].F2() // Important to stop the player clipping through the wall...
				p.path = p.path[1:]
				d = p.path[0].F2().Sub(p.pos)
				np = d.Norm()
			} else {
				p.pos = p.path[0].F2()
				p.path = nil
				p.frame = 0
				p.state.a = playerIdle
				return
			}
		}
		p.state.a = playerWalking
		p.state.d = d.Dir()
		if frame%animationPeriod == 0 {
			// Player walks straight there.
			p.pos = p.pos.Add(d.Mul(playerSpeed / np))
			p.frame++
		}
	} else {
		p.frame = 0
		p.state.a = playerIdle // preserve direction
	}
}
