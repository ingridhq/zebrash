package images

import (
	"image"

	"golang.org/x/image/draw"
)

func NewScaled(img image.Image, scaleFactorX, scaleFactorY int) image.Image {
	bounds := img.Bounds()

	if scaleFactorX == 1 && scaleFactorY == 1 {
		return img
	}

	dst := image.NewRGBA(image.Rect(0, 0, bounds.Max.X*scaleFactorX, bounds.Max.Y*scaleFactorY))

	draw.NearestNeighbor.Scale(dst, dst.Rect, img, bounds, draw.Over, nil)

	return dst
}
