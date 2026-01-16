package images

import (
	"image"
	"image/color"

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

type scaled1DHeightWrap struct {
	src    image.Image
	width  int
	height int
}

func NewScaled1DHeight(img image.Image, scaleFactorY int) image.Image {
	return &scaled1DHeightWrap{
		src:    img,
		width:  img.Bounds().Dx(),
		height: scaleFactorY,
	}
}

func (w *scaled1DHeightWrap) ColorModel() color.Model {
	return w.src.ColorModel()
}

func (w *scaled1DHeightWrap) Bounds() image.Rectangle {
	return image.Rect(0, 0, w.width, w.height)
}

func (w *scaled1DHeightWrap) At(x, y int) color.Color {
	srcBounds := w.src.Bounds()
	return w.src.At(srcBounds.Min.X+x, srcBounds.Min.Y)
}
