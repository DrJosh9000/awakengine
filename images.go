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
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"log"

	"github.com/DrJosh9000/vec"
	"github.com/hajimehoshi/ebiten"
)

var (
	allData = make(map[string][]byte)

	// Hey guess what? We're going to draw all the source images into one giant texture,
	// then do a single epic draw call during the game. Wheeee!
	composite       *ebiten.Image
	compositeOffset = make(map[string]vec.I2)
	compositeSize   = vec.I2{1024, 1024}

	sizes = make(map[string]vec.I2)
)

// RegisterImage tells the engine that a key maps to an image.
// Registered images will be loaded into a texture atlas with loadAllImages.
func RegisterImage(key string, png []byte) {
	allData[key] = png
}

func loadAllImages() error {
	// Prerender terrain layers to a texture.
	f, err := ebiten.NewImage(compositeSize.X, compositeSize.Y, ebiten.FilterNearest)
	if err != nil {
		return fmt.Errorf("creating composite texture: %v", err)
	}
	if err := f.Fill(color.Transparent); err != nil {
		return fmt.Errorf("filling composite texture with transparent color: %v", err)
	}
	composite = f
	p := vec.I2{0, 0}
	my := 0
	for k, d := range allData {
		i, err := loadPNG(d, ebiten.FilterNearest)
		if err != nil {
			return err
		}
		w, h := i.Size()
		sizes[k] = vec.I2{w, h}
		if w >= compositeSize.X {
			return fmt.Errorf("source image %q too wide [%d >= %d]", k, w, compositeSize.X)
		}
		if p.X+w >= compositeSize.X {
			p = vec.I2{0, my}
		}
		y := p.Y + h
		if y >= compositeSize.Y {
			return errors.New("too much source image (TODO josh: fix)")
		}
		if y > my {
			my = y
		}
		compositeOffset[k] = p
		if config.Debug {
			log.Printf("placing %q at (%d, %d)-(%d, %d)", k, p.X, p.Y, p.X+w, p.Y+h)
		}
		if err := f.DrawImage(i, &ebiten.DrawImageOptions{ImageParts: &wholeImageAt{p, vec.I2{w, h}}}); err != nil {
			return err
		}
		p.X += w
	}
	return nil
}

func loadPNG(img []byte, filter ebiten.Filter) (*ebiten.Image, error) {
	i, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(i, filter)
}

type wholeImageAt struct {
	p, sz vec.I2
}

func (a *wholeImageAt) Len() int { return 1 }
func (a *wholeImageAt) Dst(int) (x0, y0, x1, y1 int) {
	return a.p.X, a.p.Y, a.p.X + a.sz.X, a.p.Y + a.sz.Y
}
func (a *wholeImageAt) Src(int) (x0, y0, x1, y1 int) {
	return 0, 0, a.sz.X, a.sz.Y
}
