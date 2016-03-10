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

func RegisterImage(key string, png []bytes) {
	allData[key] = png
}

func Image(k string) *ebiten.Image {
	return allImages[k]
}

func LoadAllImages() error {
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
