package awakengine

import (
	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
)

// We go at 30 fps... to be more precise, 1 frame per 2 frames at 60fps.
// Personally I quite liked 20fps, but it was too choppy on iOS.
const animationPeriod = 3

// AnimPlayback describes the playback modes for animations.
type AnimPlayback int

// Playback modes.
const (
	AnimOneShot = AnimPlayback(iota)
	AnimLoop
)

// Anim describes an animated sprite that might play.
type Anim struct {
	Key       string
	Offset    vec.I2
	Frames    int
	FrameSize vec.I2
	Mode      AnimPlayback
}

// Image returns the image associated with this Anim.
func (a *Anim) Image() *ebiten.Image {
	return allImages[a.Key]
}
