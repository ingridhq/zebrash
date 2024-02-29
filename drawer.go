package zebrash

import (
	"fmt"
	"io"
	"math"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/drawers"
	"github.com/ingridhq/zebrash/elements"
)

type Drawer struct {
	elementDrawers []*drawers.ElementDrawer
}

func NewDrawer() *Drawer {
	return &Drawer{
		elementDrawers: []*drawers.ElementDrawer{
			drawers.NewGraphicBoxDrawer(),
			drawers.NewGraphicCircleDrawer(),
			drawers.NewGraphicFieldDrawer(),
			drawers.NewTextFieldDrawer(),
			drawers.NewMaxicodeDrawer(),
		},
	}
}

func (d *Drawer) DrawLabelAsPng(label elements.LabelInfo, output io.Writer, options drawers.DrawerOptions) error {
	options = options.WithDefaults()

	widthMm := options.LabelWidthMm
	heightMm := options.LabelHeightMm
	dpmm := options.Dpmm

	imageWidth := math.Ceil(widthMm * float64(dpmm))
	imageHeight := math.Ceil(heightMm * float64(dpmm))

	gCtx := gg.NewContext(int(imageWidth), int(imageHeight))
	gCtx.SetRGB(1, 1, 1)
	gCtx.Clear()

	for _, element := range label.Elements {
		for _, drawer := range d.elementDrawers {
			err := drawer.Draw(gCtx, element, options)
			if err != nil {
				return fmt.Errorf("failed to draw zpl element: %w", err)
			}
		}
	}

	return gCtx.EncodePNG(output)
}
