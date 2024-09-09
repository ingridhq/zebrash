package zebrash

import (
	"fmt"
	"io"
	"math"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/drawers"
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/images"
)

type reversePrintable interface {
	IsReversePrint() bool
}

type Drawer struct {
	elementDrawers []*drawers.ElementDrawer
}

func NewDrawer() *Drawer {
	return &Drawer{
		elementDrawers: []*drawers.ElementDrawer{
			drawers.NewGraphicBoxDrawer(),
			drawers.NewGraphicCircleDrawer(),
			drawers.NewGraphicFieldDrawer(),
			drawers.NewGraphicDiagonalLineDrawer(),
			drawers.NewTextFieldDrawer(),
			drawers.NewMaxicodeDrawer(),
			drawers.NewBarcode128Drawer(),
			drawers.NewBarcodePdf417Drawer(),
			drawers.NewBarcodeAztecDrawer(),
			drawers.NewBarcodeDatamatrixDrawer(),
			drawers.NewBarcodeQrDrawer(),
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
	gCtx.SetColor(images.ColorWhite)
	gCtx.Clear()

	for _, element := range label.Elements {
		reversePrint := false

		if el, ok := element.(reversePrintable); ok {
			reversePrint = el.IsReversePrint()
		}

		gCtx2 := gCtx
		if reversePrint {
			gCtx2 = gg.NewContext(int(imageWidth), int(imageHeight))
		}

		for _, drawer := range d.elementDrawers {
			err := drawer.Draw(gCtx2, element, options)
			if err != nil {
				return fmt.Errorf("failed to draw zpl element: %w", err)
			}
		}

		if reversePrint {
			if err := images.ReversePrint(gCtx2.Image(), gCtx.Image()); err != nil {
				return err
			}
		}
	}

	return images.EncodeMonochrome(output, gCtx.Image())
}
