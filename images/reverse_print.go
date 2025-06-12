package images

import (
	"fmt"
	"image"
)

func ReversePrint(mask, background image.Image) error {
	rgba1, ok := mask.(*image.RGBA)
	if !ok {
		return fmt.Errorf("mask is not an RGBA image")
	}

	rgba2, ok := background.(*image.RGBA)
	if !ok {
		return fmt.Errorf("img is not an RGBA image")
	}

	const alphaThreshold = 30

	for i := 3; i < len(rgba1.Pix); i += 4 {
		a1 := rgba1.Pix[i]

		if a1 < alphaThreshold {
			continue
		}

		rgba2.Pix[i-3] = 255 - rgba2.Pix[i-3]
		rgba2.Pix[i-2] = 255 - rgba2.Pix[i-2]
		rgba2.Pix[i-1] = 255 - rgba2.Pix[i-1]
	}

	return nil
}
