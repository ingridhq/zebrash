package images

import (
	"fmt"
	"image"
)

func Zerofill(img image.Image) error {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		return fmt.Errorf("img is not an RGBA image")
	}

	for i := range rgba.Pix {
		rgba.Pix[i] = 0
	}

	return nil
}
