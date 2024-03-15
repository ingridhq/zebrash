package images

import (
	"image"
	"image/color"
)

func NewTransparent(img image.Image) *TransparentImgWrap {
	return &TransparentImgWrap{img: img}
}

type TransparentImgWrap struct {
	img image.Image
}

func (w *TransparentImgWrap) ColorModel() color.Model {
	return color.RGBAModel
}

func (w *TransparentImgWrap) Bounds() image.Rectangle {
	return w.img.Bounds()
}

func (w *TransparentImgWrap) At(x, y int) color.Color {
	c := w.img.At(x, y)
	if r, _, _, _ := c.RGBA(); r == 0 {
		return color.RGBA{A: 255}
	}

	return color.RGBA{A: 0}
}
