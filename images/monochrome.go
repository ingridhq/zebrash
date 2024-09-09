package images

import (
	"fmt"
	"image"
	"image/png"
	"io"
)

func EncodeMonochrome(w io.Writer, img image.Image) error {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		return fmt.Errorf("img is not an RGBA image")
	}

	const threshold = 125

	result := image.NewGray(rgba.Rect)

	for i := 3; i < len(rgba.Pix); i += 4 {
		val := rgba.Pix[i-3]

		if val > threshold {
			val = 255
		} else {
			val = 0
		}

		result.Pix[(i-3)/4] = val
	}

	return png.Encode(w, result)
}
