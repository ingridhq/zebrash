package drawers

import (
	"fmt"
	"image"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/barcodes/code39"
	"github.com/ingridhq/zebrash/elements"
)

func NewBarcode39Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.Barcode39WithData)
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

			img, err = code39.Encode(content, barcode.Width, barcode.Height, barcode.WidthRatio)
			if err != nil {
				return fmt.Errorf("failed to encode code39 barcode: %w", err)
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
