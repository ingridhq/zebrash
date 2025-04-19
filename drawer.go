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
			drawers.NewBarcode39Drawer(),
			drawers.NewBarcodePdf417Drawer(),
			drawers.NewBarcodeAztecDrawer(),
			drawers.NewBarcodeDatamatrixDrawer(),
			drawers.NewBarcodeQrDrawer(),
		},
	}
}

func (d *Drawer) DrawLabelAsPng(label elements.LabelInfo, output io.Writer, options drawers.DrawerOptions) error {
	options = options.WithDefaults()
	state := &drawers.DrawerState{}

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

	for _, element := range label.Elements {
		reversePrint := false

		if el, ok := element.(reversePrintable); ok {
			reversePrint = el.IsReversePrint()
		}

		gCtx2 := gCtx
		if reversePrint {
			gCtx2 = gg.NewContext(imageWidth, imageHeight)
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
