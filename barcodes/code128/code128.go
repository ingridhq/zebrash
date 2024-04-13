package code128

import (
	"image"
	"image/color"

	"github.com/ingridhq/zebrash/images"
)

type code128 struct {
	code     []bool
	width    int
	height   int
	barWidth int
}

func newCode128(code []bool, height, barWidth int) *code128 {
	barWidth = max(1, barWidth)
	height = max(1, height)

	return &code128{
		code:     code,
		width:    len(code) * barWidth,
		height:   height,
		barWidth: barWidth,
	}
}

func (c *code128) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *code128) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.width, c.height)
}

func (c *code128) At(x, y int) color.Color {
	x /= c.barWidth

	if x >= 0 && x < len(c.code) && c.code[x] {
		return images.ColorBlack
	}

	return images.ColorTransparent
}
