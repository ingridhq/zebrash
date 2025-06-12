package images

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

func NewScaled(img image.Image, scaleFactorX, scaleFactorY int) image.Image {
	if scaleFactorX == 1 && scaleFactorY == 1 {
		return img
	}

	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Max.X*scaleFactorX, bounds.Max.Y*scaleFactorY))

	draw.NearestNeighbor.Scale(dst, dst.Rect, img, bounds, draw.Over, nil)

	return dst
}

func NewScaledFloat(img image.Image, scaleFactorX, scaleFactorY float64) image.Image {
	if scaleFactorX == 1 && scaleFactorY == 1 {
		return img
	}

	bounds := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, int(math.Ceil(float64(bounds.Max.X)*scaleFactorX)), int(math.Ceil(float64(bounds.Max.Y)*scaleFactorY))))

	draw.NearestNeighbor.Scale(dst, dst.Rect, img, bounds, draw.Over, nil)

	return dst
}
