package pdf417

import (
	"image"
	"image/color"
	"math"
)

type pdfBarcode struct {
	width        int
	height       int
	code         *BitList
	scaleFactorX float64
	scaleFactorY float64
}

func (c *pdfBarcode) ColorModel() color.Model {
	return color.Gray16Model
}

func (c *pdfBarcode) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.scaleX(c.width), c.scaleY(c.height))
}

func (c *pdfBarcode) At(x, y int) color.Color {
	x = c.downscaleX(x)
	y = c.downscaleY(y)

	if c.code.GetBit(y*c.width + x) {
		return color.Black
	}
	return color.White
}

func (c *pdfBarcode) scaleX(v int) int {
	return int(math.Round(float64(v) * c.scaleFactorX))
}

func (c *pdfBarcode) scaleY(v int) int {
	return int(math.Round(float64(v) * c.scaleFactorY))
}

func (c *pdfBarcode) downscaleX(v int) int {
	return int(math.Round(float64(v) / c.scaleFactorX))
}

func (c *pdfBarcode) downscaleY(v int) int {
	return int(math.Round(float64(v) / c.scaleFactorY))
}
