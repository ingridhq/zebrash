package ean13

import (
	"image"
	"image/color"

	"github.com/ingridhq/zebrash/internal/barcodes/utils"
	"github.com/ingridhq/zebrash/internal/images"
)

type ean13 struct {
	code           *utils.BitList
	width          int
	height         int
	barWidth       int
	guardExtension int // Extra height for guard bars
}

func newEan13(code *utils.BitList, height, barWidth int) *ean13 {
	barWidth = max(1, barWidth)
	height = max(1, height)

	return &ean13{
		code:           code,
		width:          code.Len() * barWidth,
		height:         height,
		barWidth:       barWidth,
		guardExtension: CalculateGuardExtension(barWidth),
	}
}

func (c *ean13) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *ean13) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.width, c.height+c.guardExtension)
}

func (c *ean13) At(x, y int) color.Color {
	x /= c.barWidth

	if x < 0 || x >= c.code.Len() {
		return images.ColorTransparent
	}

	if !c.code.GetBit(x) {
		return images.ColorTransparent
	}

	if isGuardBar(x) {
		// Guard bars extend the full height
		return images.ColorBlack
	}

	// Regular bars only draw in the upper portion (leaving room for guard extension at bottom)
	if y < c.height {
		return images.ColorBlack
	}
	return images.ColorTransparent
}

// isGuardBar checks if a module position is part of a guard pattern
// EAN-13 structure:
// - Start guard: modules 0-2 (3 modules)
// - Left digits: modules 3-44 (42 modules = 6 digits * 7)
// - Middle guard: modules 45-49 (5 modules)
// - Right digits: modules 50-91 (42 modules = 6 digits * 7)
// - End guard: modules 92-94 (3 modules)
func isGuardBar(x int) bool {
	// Start guard (first 3 modules)
	if x >= 0 && x <= 2 {
		return true
	}
	// Middle guard (5 modules after 3 start + 42 left = position 45-49)
	if x >= 45 && x <= 49 {
		return true
	}
	// End guard (last 3 modules = positions 92-94)
	if x >= 92 && x <= 94 {
		return true
	}
	return false
}

func CalculateGuardExtension(barWidth int) int {
	// Guard bars are typically extended by about 5 times the X-dimension (module width)
	// This makes them visually distinctive per EAN-13 standard
	return barWidth * 5
}
