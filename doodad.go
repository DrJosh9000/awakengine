package awakengine

import "github.com/DrJosh9000/vec"

// BaseDoodad models a static object rendered at the player layer, but computed as
// obstacles (like terrain).
type BaseDoodad struct {
	UL, DR vec.I2 // base obstacle box (pos relative)
	*Anim
	Frame int
}

// Doodad is an instance of a BaseDoodad in a specific location.
type Doodad struct {
	P vec.I2
	*BaseDoodad
}

// Anim implements Sprite.
func (b *BaseDoodad) Anim() *Anim { return b.Anim }

// Frame implements Sprite.
func (b *BaseDoodad) Frame() int { return b.Frame }

// Pos implements Sprite.
func (d *Doodad) Pos() vec.I2 { return d.Pos }

// ByYPos orders Sprites by Y position (least to greatest).
type DoodadsByYPos []*Doodad

// Len implements sort.Interface.
func (b DoodadsByYPos) Len() int { return len(b) }

// Less implements sort.Interface.
func (b DoodadsByYPos) Less(i, j int) bool { return b[i].Pos.Y < b[j].Pos.Y }

// Swap implements sort.Interface.
func (b DoodadsByYPos) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
