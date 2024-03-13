package drawers

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/elements"
)

type ElementDrawer struct {
	Draw func(gCtx *gg.Context, element interface{}, options DrawerOptions) error
}

func drawImage(gCtx *gg.Context, img image.Image, xo, yo int) {
	imgWrap := &transparentImgWrap{img: img}
	gCtx.DrawImage(imgWrap, xo, yo)
}

func rotateImage(gCtx *gg.Context, img image.Image, pos elements.LabelPosition, ori elements.FieldOrientation) {
	width := float64(img.Bounds().Dx())
	height := float64(img.Bounds().Dy())
	rotate := ori.GetDegrees()

	if rotate == 0 {
		return
	}

	gCtx.RotateAbout(gg.Radians(rotate), float64(pos.X), float64(pos.Y))

	switch ori {
	case elements.FieldOrientation90:
		gCtx.Translate(0, -height)
	case elements.FieldOrientation180:
		gCtx.Translate(-width, -height)
	case elements.FieldOrientation270:
		gCtx.Translate(-width, 0)
	}
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
	case elements.LineColorWhite:
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
