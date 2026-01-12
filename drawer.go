package zebrash

import (
	"fmt"
	"io"
	"math"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/drawers"
	"github.com/ingridhq/zebrash/elements"
	drawers_internal "github.com/ingridhq/zebrash/internal/drawers"
	"github.com/ingridhq/zebrash/internal/images"
)

type reversePrintable interface {
	IsReversePrint() bool
}

type Drawer struct {
	elementDrawers []*drawers_internal.ElementDrawer
}

func NewDrawer() *Drawer {
	return &Drawer{
		elementDrawers: []*drawers_internal.ElementDrawer{
			drawers_internal.NewGraphicBoxDrawer(),
			drawers_internal.NewGraphicCircleDrawer(),
			drawers_internal.NewGraphicFieldDrawer(),
			drawers_internal.NewGraphicDiagonalLineDrawer(),
			drawers_internal.NewTextFieldDrawer(),
			drawers_internal.NewMaxicodeDrawer(),
			drawers_internal.NewBarcode128Drawer(),
			drawers_internal.NewBarcodeEan13Drawer(),
			drawers_internal.NewBarcode2of5Drawer(),
			drawers_internal.NewBarcode39Drawer(),
			drawers_internal.NewBarcodePdf417Drawer(),
			drawers_internal.NewBarcodeAztecDrawer(),
			drawers_internal.NewBarcodeDatamatrixDrawer(),
			drawers_internal.NewBarcodeQrDrawer(),
		},
	}
}

func (d *Drawer) DrawLabelAsPng(label elements.LabelInfo, output io.Writer, options drawers.DrawerOptions) error {
	options = options.WithDefaults()
	state := &drawers_internal.DrawerState{}

	widthMm := options.LabelWidthMm
	heightMm := options.LabelHeightMm
	dpmm := options.Dpmm

	labelWidth := int(math.Ceil(widthMm * float64(dpmm)))
	imageWidth := labelWidth
	if label.PrintWidth > 0 {
		imageWidth = min(labelWidth, label.PrintWidth)
	}

	imageHeight := int(math.Ceil(heightMm * float64(dpmm)))

	gCtx := gg.NewContext(imageWidth, imageHeight)
	gCtx.SetColor(images.ColorWhite)
	gCtx.Clear()

	var gReversePrintBuff *gg.Context

	for _, element := range label.Elements {
		reversePrint := false

		if el, ok := element.(reversePrintable); ok {
			reversePrint = el.IsReversePrint()
		}

		gCtx2 := gCtx
		if reversePrint {
			if gReversePrintBuff == nil {
				gReversePrintBuff = gg.NewContext(imageWidth, imageHeight)
			} else if err := images.Zerofill(gReversePrintBuff.Image()); err != nil {
				return fmt.Errorf("failed to clear reverse print buffer: %w", err)
			}

			gCtx2 = gReversePrintBuff
		}

		for _, drawer := range d.elementDrawers {
			err := drawer.Draw(gCtx2, element, options, state)
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

	// If print width was less than label width
	// Draw everything onto the new, wider image and center the content
	if imageWidth != labelWidth {
		imgCtx := gCtx
		gCtx = gg.NewContext(labelWidth, imageHeight)
		gCtx.SetColor(images.ColorWhite)
		gCtx.Clear()

		img := imgCtx.Image()
		gCtx.DrawImage(img, (labelWidth-imageWidth)/2, 0)
	}

	return images.EncodeMonochrome(output, gCtx.Image())
}
