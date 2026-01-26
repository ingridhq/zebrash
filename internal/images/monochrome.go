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

	result := image.NewGray(rgba.Rect)

	for i := 3; i < len(rgba.Pix); i += 4 {
		val := rgba.Pix[i-3]

		result.Pix[(i-3)/4] = val
	}

	return png.Encode(w, result)
}
