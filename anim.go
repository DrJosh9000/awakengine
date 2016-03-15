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
	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
)

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
