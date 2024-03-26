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

func (font FontInfo) WithAdjustedSizes() FontInfo {
	orgSize, ok := bitmapFontSizes[font.Name]
	if !ok {
		// Not a bitmap font, for example 0 which is scalable font
		// so it does not have incremental size increases
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
