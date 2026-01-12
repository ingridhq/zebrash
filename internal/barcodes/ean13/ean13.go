package ean13

import (
	"image"
	"image/color"

	"github.com/ingridhq/zebrash/internal/images"
)

type ean13 struct {
	code     []bool
	width    int
	height   int
	barWidth int
}

func newEan13(code []bool, height, barWidth int) *ean13 {
	barWidth = max(1, barWidth)
	height = max(1, height)

	return &ean13{
		code:     code,
		width:    len(code) * barWidth,
		height:   height,
		barWidth: barWidth,
	}
}

func (c *ean13) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *ean13) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.width, c.height)
}

func (c *ean13) At(x, y int) color.Color {
	x /= c.barWidth

	if x >= 0 && x < len(c.code) && c.code[x] {
		return images.ColorBlack
	}

	return images.ColorTransparent
}
