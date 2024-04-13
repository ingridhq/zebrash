package aztec

import (
	"image"
	"image/color"

	"github.com/ingridhq/zebrash/barcodes/utils"
	"github.com/ingridhq/zebrash/images"
)

type aztecCode struct {
	*utils.BitList
	size    int
	content []byte
}

func newAztecCode(size int) *aztecCode {
	return &aztecCode{utils.NewBitList(size * size), size, nil}
}

func (c *aztecCode) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *aztecCode) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.size-1, c.size-1)
}

func (c *aztecCode) At(x, y int) color.Color {
	if c.GetBit(x*c.size + y) {
		return images.ColorBlack
	}
	return images.ColorTransparent
}

func (c *aztecCode) set(x, y int) {
	c.SetBit(x*c.size+y, true)
}
