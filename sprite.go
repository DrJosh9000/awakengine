package awakengine

import (
	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
)

// Transient is a sprite that starts at a given birth.
type Transient struct {
	P     vec.I2
	Birth int
	A     *Anim
}

// Anim implements Sprite.
func (t *Transient) Anim() *Anim { return t.A }

// Frame implements Sprite.
func (t *Transient) Frame() int { return (gameFrame - t.Birth) / animationPeriod }

// Pos implements Sprite.
func (t *Transient) Pos() vec.I2 { return t.P }

// Static just draws a frame.
type Static struct {
	P vec.I2
	F int
	A *Anim
}

// Anim implements Sprite.
func (s *Static) Anim() *Anim { return s.A }

// Frame implements Sprite.
func (s *Static) Frame() int { return s.Frame }

// Pos implements Sprite.
func (s *Static) Pos() vec.I2 { return s.P }

// Sprite is all the information required to draw an animated thingy at a point on screen.
type Sprite interface {
	Anim() *Anim
	Frame() int
	Pos() vec.I2 // logical / world position.
}

// ByYPos orders Sprites by Y position (least to greatest).
type ByYPos []Sprite

// Len implements sort.Interface.
func (b ByYPos) Len() int { return len(b) }

// Less implements sort.Interface.
func (b ByYPos) Less(i, j int) bool { return b[i].Pos().Y < b[j].Pos().Y }

// Swap implements sort.Interface.
func (b ByYPos) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// SpriteParts implements ebiten.ImageParts for sprite drawing.
type SpriteParts struct {
	Sprite
	InWorld bool
}

// Draw draws the sprite to the screen.
func (s SpriteParts) Draw(screen *ebiten.Image) error {
	return screen.DrawImage(s.Anim().Image(), &ebiten.DrawImageOptions{ImageParts: s})
}

// Len implements ebiten.ImageParts.
func (s SpriteParts) Len() int { return 1 }

// Dst implements ebiten.ImageParts.
func (s SpriteParts) Dst(i int) (x0, y0, x1, y1 int) {
	a := s.Anim()
	b := s.Pos().Sub(a.offset)
	if s.InWorld {
		b = b.Sub(camPos)
	}
	c := b.Add(a.frameSize)
	return b.X, b.Y, c.X, c.Y
}

// Src implements ebiten.ImageParts.
func (s SpriteParts) Src(i int) (x0, y0, x1, y1 int) {
	a, f := s.Anim(), s.Frame()
	switch a.mode {
	case AnimOneShot:
		if f >= a.frames {
			return
		}
	case AnimLoop:
		f %= a.frames
	}
	x0 = f * a.frameSize.X
	return x0, 0, x0 + a.frameSize.X, a.frameSize.Y
}
