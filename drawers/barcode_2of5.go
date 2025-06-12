package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/barcodes/twooffive"
	"github.com/ingridhq/zebrash/elements"
)

func NewBarcode2of5Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.Barcode2of5WithData)
			if !ok {
				return nil
			}

			// data to encode into barcode
			content := barcode.Data

			img, text, err := twooffive.EncodeInterleaved(content, barcode.Width, barcode.Height, barcode.WidthRatio, barcode.CheckDigit)
			if err != nil {
				return fmt.Errorf("failed to encode 2 of 5 barcode: %w", err)
			}

			width := float64(img.Bounds().Dx())
			height := float64(img.Bounds().Dy())
			pos := adjustImageTypeSetPosition(img, barcode.Position, barcode.Orientation)

			rotateImage(gCtx, img, pos, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, pos.X, pos.Y)
			if barcode.Line {
				applyLineTextToCtx(gCtx, text, pos, barcode.LineAbove, width, height)
			}

			return nil
		},
	}
}
