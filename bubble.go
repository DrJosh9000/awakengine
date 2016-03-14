package awakengine

import (
	"fmt"

	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
)

const bubblePartSize = 5

// Bubble is an ImageParts that renders a bubble at any size larger than 15x15.
type Bubble struct {
	ul, dr vec.I2
	flat   *ebiten.Image
}

// AllBubbleParts is an ImageParts that renders the bubble piecewise.
type AllBubbleParts vec.I2

// NewBubble prepares a bubble of the correct size.
func NewBubble(pos, size vec.I2) (*Bubble, error) {
	f, err := ebiten.NewImage(size.X, size.Y, ebiten.FilterNearest)
	if err != nil {
		return nil, fmt.Errorf("creating image: %v", err)
	}
	if err := f.DrawImage(allImages["bubble"], &ebiten.DrawImageOptions{ImageParts: AllBubbleParts(size)}); err != nil {
		return nil, fmt.Errorf("drawing bubble: %v", err)
	}
	return &Bubble{
		ul:   pos,
		dr:   pos.Add(size),
		flat: f,
	}, nil
}

// Len implements ImageParts.
func (b *Bubble) Len() int { return 1 }

// Src implements ImageParts.
func (b *Bubble) Src(int) (x0, y0, x1, y1 int) {
	x1, y1 = b.dr.Sub(b.ul).C()
	return
}

// Dst implements ImageParts.
func (b *Bubble) Dst(int) (x0, y0, x1, y1 int) {
	return b.ul.X, b.ul.Y, b.dr.X, b.dr.Y
}

// Draw draws the bubble to the screen.
func (b *Bubble) Draw(screen *ebiten.Image) error {
	return screen.DrawImage(b.flat, &ebiten.DrawImageOptions{ImageParts: b})
}

// Len implements ImageParts.
func (b AllBubbleParts) Len() int { return 9 }

// Src implements ImageParts.
func (b AllBubbleParts) Src(i int) (x0, y0, x1, y1 int) {
	x0, y0 = vec.Div(i, 3).Mul(bubblePartSize).C()
	return x0, y0, x0 + bubblePartSize, y0 + bubblePartSize
}

// Dst implements ImageParts.
func (b AllBubbleParts) Dst(i int) (x0, y0, x1, y1 int) {
	j, k := vec.Div(i, 3).C()
	switch j {
	case 0:
		x1 = bubblePartSize
	case 1:
		x0 = bubblePartSize
		x1 = b.X - bubblePartSize
	case 2:
		x0 = b.X - bubblePartSize
		x1 = b.X
	}
	switch k {
	case 0:
		y1 = bubblePartSize
	case 1:
		y0 = bubblePartSize
		y1 = b.Y - bubblePartSize
	case 2:
		y0 = b.Y - bubblePartSize
		y1 = b.Y
	}
	return
}
