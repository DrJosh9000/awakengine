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
	"image/png"

	"github.com/hajimehoshi/ebiten"
)

var (
	allImages map[string]*ebiten.Image
	allData   = make(map[string][]byte)
)

// RegisterImage tells the engine that a key maps to an image.
// Registered images will be loaded into texture memory in Load.
func RegisterImage(key string, png []byte) {
	allData[key] = png
}

// Image returns the Image associated with a key. Will return nil
// if the key isn't registered or the image isn't loaded yet.
func Image(key string) *ebiten.Image {
	return allImages[key]
}

func loadAllImages() error {
	allImages = make(map[string]*ebiten.Image)
	for k, d := range allData {
		i, err := loadPNG(d, ebiten.FilterNearest)
		if err != nil {
			return err
		}
		allImages[k] = i
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
