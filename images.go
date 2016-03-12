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
