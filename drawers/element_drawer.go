package drawers

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/elements"
)

type ElementDrawer struct {
	Draw func(gCtx *gg.Context, element interface{}, options DrawerOptions) error
}

func drawRectangle(gCtx *gg.Context, x, y, w, h float64) {
	gCtx.DrawLine(x, y, x+w, y)
	gCtx.DrawLine(x+w, y, x+w, y+h)
	gCtx.DrawLine(x+w, y+h, x, y+h)
	gCtx.DrawLine(x, y+h, x, y)
}

func setLineColor(gCtx *gg.Context, color elements.LineColor) {
	switch color {
	case elements.LineColorBlack:
		gCtx.SetRGB(0, 0, 0)
	case elements.LineColorWhile:
		gCtx.SetRGB(1, 1, 1)
	}
}

func mustLoadFont(fontData []byte) *truetype.Font {
	font, err := truetype.Parse(fontData)
	if err != nil {
		panic(err.Error())
	}

	return font
}
