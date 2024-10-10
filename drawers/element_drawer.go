package drawers

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/images"
)

type ElementDrawer struct {
	Draw func(gCtx *gg.Context, element interface{}, options DrawerOptions, state *DrawerState) error
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

func setLineColor(gCtx *gg.Context, color elements.LineColor) {
	switch color {
	case elements.LineColorBlack:
		gCtx.SetColor(images.ColorBlack)
	case elements.LineColorWhite:
		gCtx.SetColor(images.ColorWhite)
	}
}
