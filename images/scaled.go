package images

import (
	"image"
	"image/color"
	"math"
)

type ScaledImgWrap struct {
	img          image.Image
	scaledRect   image.Rectangle
	scaleFactorX float64
	scaleFactorY float64
}

func NewScaled(img image.Image, scaleFactorX, scaleFactorY float64) image.Image {
	bounds := img.Bounds()

	if scaleFactorX == 1 && scaleFactorY == 1 {
		return img
	}

	width := int(math.Round(float64(bounds.Dx()) * scaleFactorX))
	height := int(math.Round(float64(bounds.Dy()) * scaleFactorY))

	return &ScaledImgWrap{
		img:          img,
		scaledRect:   image.Rect(0, 0, width, height),
		scaleFactorX: scaleFactorX,
		scaleFactorY: scaleFactorY,
	}
}

func (w *ScaledImgWrap) ColorModel() color.Model {
	return w.img.ColorModel()
}

func (w *ScaledImgWrap) Bounds() image.Rectangle {
	return w.scaledRect
}

func (w *ScaledImgWrap) At(x, y int) color.Color {
	x = w.downscaleX(x)
	y = w.downscaleY(y)

	return w.img.At(x, y)
}

func (w *ScaledImgWrap) downscaleX(v int) int {
	return int(math.Round(float64(v) / w.scaleFactorX))
}

func (w *ScaledImgWrap) downscaleY(v int) int {
	return int(math.Round(float64(v) / w.scaleFactorY))
}
