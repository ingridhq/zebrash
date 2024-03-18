package images

import (
	"image"
	"image/color"
)

type ShiftedImgWrap struct {
	img    image.Image
	shiftX int
	shiftY int
}

func NewShifted(img image.Image, shiftX, shiftY int) *ShiftedImgWrap {
	return &ShiftedImgWrap{
		img:    img,
		shiftX: shiftX,
		shiftY: shiftY,
	}
}

func (w *ShiftedImgWrap) ColorModel() color.Model {
	return w.img.ColorModel()
}

func (w *ShiftedImgWrap) Bounds() image.Rectangle {
	return w.img.Bounds()
}

func (w *ShiftedImgWrap) At(x, y int) color.Color {
	return w.img.At(x-w.shiftX, y-w.shiftY)
}
