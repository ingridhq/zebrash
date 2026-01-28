package images

import (
	"image/draw"
)

// Rotate180 rotates an image 180 degrees in place, mutating the passed-in image.
func Rotate180(img draw.Image) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 1 && h <= 1 {
		return
	}

	minX, minY := b.Min.X, b.Min.Y
	total := w * h

	// Swap pixel i with pixel total-1-i in row-major order.
	for i := 0; i < total/2; i++ {
		x1, y1 := i%w, i/w
		j := total - 1 - i
		x2, y2 := j%w, j/w

		ax1, ay1 := minX+x1, minY+y1
		ax2, ay2 := minX+x2, minY+y2

		c1 := img.At(ax1, ay1)
		c2 := img.At(ax2, ay2)
		img.Set(ax1, ay1, c2)
		img.Set(ax2, ay2, c1)
	}
}
