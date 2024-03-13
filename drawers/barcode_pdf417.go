package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/barcodes/pdf417"
	"github.com/ingridhq/zebrash/elements"
)

func NewBarcodePdf417Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, options DrawerOptions) error {
			barcode, ok := element.(*elements.BarcodePdf417WithData)
			if !ok {
				return nil
			}

			img, err := pdf417.Encode(barcode.Data, byte(barcode.Security), barcode.RowHeight, barcode.Columns)
			if err != nil {
				return fmt.Errorf("failed to encode pdf417 barcode: %w", err)
			}

			rotateImage(gCtx, img, barcode.Position, barcode.Orientation)

			defer gCtx.Identity()

			drawImage(gCtx, img, barcode.Position.X, barcode.Position.Y)

			return nil
		},
	}
}
