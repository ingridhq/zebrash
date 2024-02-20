package zebrash

import (
	"fmt"
	"io"
	"math"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/drawers"
	"github.com/ingridhq/zebrash/elements"
)

type DrawerOptions struct {
	LabelWidthMm  float64
	LabelHeightMm float64
	Dpmm          int
}

type Drawer struct {
	elementDrawers []*drawers.ElementDrawer
}

func NewDrawer() *Drawer {
	return &Drawer{
		elementDrawers: []*drawers.ElementDrawer{
			drawers.NewGraphicBoxDrawer(),
			drawers.NewTextFieldDrawer(),
		},
	}
}

func (d *Drawer) DrawLabelAsPng(label elements.LabelInfo, output io.Writer, options DrawerOptions) error {
	widthMm := options.LabelWidthMm
	if widthMm == 0 {
		widthMm = 101.6
	}

	heightMm := options.LabelHeightMm
	if heightMm == 0 {
		heightMm = 152.4
	}

	dpmm := options.Dpmm
	if dpmm == 0 {
		dpmm = 8
	}

	imageWidth := math.Ceil(widthMm * float64(dpmm))
	imageHeight := math.Ceil(heightMm * float64(dpmm))

	gCtx := gg.NewContext(int(imageWidth), int(imageHeight))
	gCtx.SetRGB(1, 1, 1)
	gCtx.Clear()

	for _, element := range label.Elements {
		for _, drawer := range d.elementDrawers {
			err := drawer.Draw(gCtx, element)
			if err != nil {
				return fmt.Errorf("failed to draw zpl element: %w", err)
			}
		}
	}

	return gCtx.EncodePNG(output)
}
