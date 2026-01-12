package drawers

import (
	"fmt"
	"image"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/drawers"
	"github.com/ingridhq/zebrash/internal/barcodes/ean13"
	"github.com/ingridhq/zebrash/internal/elements"
)

func NewBarcodeEan13Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ drawers.DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeEan13WithData)
			if !ok {
				return nil
			}

			// data to encode into barcode
			content := barcode.Data
			// human-readable text
			text := barcode.Data

			var (
				img image.Image
				err error
			)

			img, err = ean13.Encode(content, barcode.Height, barcode.Width)
			if err != nil {
				return fmt.Errorf("failed to encode EAN-13 barcode: %w", err)
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
