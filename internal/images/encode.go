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

	const threshold = 128

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

// EncodeGrayscale encodes img as an 8-bit grayscale PNG without binarisation,
// preserving the sub-pixel anti-aliasing that gg produces during rendering.
func EncodeGrayscale(w io.Writer, img image.Image) error {
	rgba, ok := img.(*image.RGBA)
	if !ok {
		return fmt.Errorf("img is not an RGBA image")
	}

	result := image.NewGray(rgba.Rect)
	for i := 3; i < len(rgba.Pix); i += 4 {
		result.Pix[(i-3)/4] = rgba.Pix[i-3] // R channel; R=G=B for greyscale rendering
	}

	return png.Encode(w, result)
}
