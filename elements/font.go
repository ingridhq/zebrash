package elements

import "math"

type FontInfo struct {
	Name        string
	Width       float64
	Height      float64
	Orientation FieldOrientation
}

func (font FontInfo) GetSize() float64 {
	switch font.Orientation {
	case FieldOrientation90, FieldOrientation270:
		return font.Width
	default:
		return font.Height
	}
}

var bitmapFontSizes = map[string][2]float64{
	"A": {9, 5},
	"B": {11, 7},
	"C": {18, 10},
	"D": {18, 10},
	"E": {28, 15},
	"F": {26, 13},
	"G": {60, 40},
	"H": {21, 13},
}

// Bitmap fonts (everything other than font 0)
// cannot be freely scaled
// their size should always divide by their base size without remainder
// so we need to adjust them
// NOTE: in order to emulate Zebra fonts 0, A-H we use only 2 different TTF fonts
// so don't confuse Zebra and our fonts, they are not the same thing
func (font FontInfo) WithAdjustedSizes() FontInfo {
	orgSize, ok := bitmapFontSizes[font.Name]
	if !ok {
		return font
	}

	if font.Width == 0 && font.Height == 0 {
		font.Width = orgSize[1]
		font.Height = orgSize[0]
		return font
	}

	if font.Width == 0 {
		font.Width = orgSize[1] * math.Round(font.Height/orgSize[0])
	} else {
		font.Width = orgSize[1] * math.Round(font.Width/orgSize[1])
	}

	if font.Height == 0 {
		font.Height = orgSize[0] * math.Round(font.Width/orgSize[1])
	} else {
		font.Height = orgSize[0] * math.Round(font.Height/orgSize[0])
	}

	return font
}
