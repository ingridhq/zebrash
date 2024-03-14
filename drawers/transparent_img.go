package drawers

import (
	"image"
	"image/color"
)

type transparentImgWrap struct {
	img image.Image
}

func (b *transparentImgWrap) ColorModel() color.Model {
	return color.RGBAModel
}

func (b *transparentImgWrap) Bounds() image.Rectangle {
	return b.img.Bounds()
}

func (b *transparentImgWrap) At(x, y int) color.Color {
	c := b.img.At(x, y)
	if c == color.Black {
		return color.RGBA{A: 255}
	}

	return color.RGBA{A: 0}
}
