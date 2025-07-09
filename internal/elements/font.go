package elements

import (
	"math"
)

type FontInfo struct {
	Name        string
	Width       float64
	Height      float64
	Orientation FieldOrientation
}

func (font FontInfo) GetSize() float64 {
	switch font.Orientation {
	case FieldOrientation90, FieldOrientation270:
		return font.getWidthToHeightRatio() * font.Width
	default:
		return font.Height
	}
}

func (font FontInfo) GetScaleX() float64 {
	scaleX := 1.0
	if font.Width != 0 {
		scaleX = font.getWidthToHeightRatio() * font.Width / font.GetSize()
	}

	return scaleX
}

var bitmapFontSizes = map[string][2]float64{
	"A":  {9, 5},
	"B":  {11, 7},
	"C":  {18, 10},
	"D":  {18, 10},
	"E":  {28, 15},
	"F":  {26, 13},
	"G":  {60, 40},
	"H":  {21, 13},
	"GS": {24, 24},
}

// Bitmap fonts (everything other than font 0)
// cannot be freely scaled
// their size should always divide by their base size without remainder
// so we need to adjust them
// NOTE: in order to emulate Zebra fonts 0, A-H we use only 2 different TTF fonts
// so don't confuse Zebra and our fonts, they are not the same thing
func (font FontInfo) WithAdjustedSizes() FontInfo {
	orgSize, ok := bitmapFontSizes[font.Name]
	// Scalable font
	// Just set width and height to the same value if one of them is zero
	if !ok {
		if font.Width == 0 {
			font.Width = font.Height
		}

		if font.Height == 0 {
			font.Height = font.Width
		}

		return font
	}

	if font.Width == 0 && font.Height == 0 {
		font.Width = orgSize[1]
		font.Height = orgSize[0]
		return font
	}

	if font.Width == 0 {
		font.Width = orgSize[1] * max(math.Round(font.Height/orgSize[0]), 1)
	} else {
		font.Width = orgSize[1] * max(math.Round(font.Width/orgSize[1]), 1)
	}

	if font.Height == 0 {
		font.Height = orgSize[0] * max(math.Round(font.Width/orgSize[1]), 1)
	} else {
		font.Height = orgSize[0] * max(math.Round(font.Height/orgSize[0]), 1)
	}

	return font
}

func (font FontInfo) getWidthToHeightRatio() float64 {
	if font.Name == "0" || font.Name == "GS" {
		return 1.0
	}

	return 2.0
}
