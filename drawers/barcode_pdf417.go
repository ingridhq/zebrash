package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/barcodes/pdf417"
	"github.com/ingridhq/zebrash/elements"
)

func NewBarcodePdf417Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodePdf417WithData)
			if !ok {
				return nil
			}

			img, err := pdf417.Encode(barcode.Data, byte(barcode.Security), barcode.RowHeight, barcode.Columns)
			if err != nil {
				return fmt.Errorf("failed to encode pdf417 barcode: %w", err)
			}

			pos := adjustImageTypeSetPosition(img, barcode.Position, barcode.Orientation)

			rotateImage(gCtx, img, pos, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, pos.X, pos.Y)

			return nil
		},
	}
}
