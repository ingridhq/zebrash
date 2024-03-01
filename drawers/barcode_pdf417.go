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

			img, err := pdf417.Encode(barcode.Data, byte(barcode.Code.Security), barcode.Code.RowHeight, barcode.Code.Columns)
			if err != nil {
				return fmt.Errorf("failed to encode pdf417 barcode: %w", err)
			}

			width := float64(img.Bounds().Dx())
			height := float64(img.Bounds().Dy())

			if rotate := barcode.Code.Orientation.GetDegrees(); rotate != 0 {
				gCtx.RotateAbout(gg.Radians(rotate), float64(barcode.Position.X), float64(barcode.Position.Y))

				switch barcode.Code.Orientation {
				case elements.FieldOrientation90:
					gCtx.Translate(0, -height)
				case elements.FieldOrientation180:
					gCtx.Translate(-width, -height)
				case elements.FieldOrientation270:
					gCtx.Translate(-width, 0)
				}
			}

			defer gCtx.Identity()

			gCtx.DrawImage(img, barcode.Position.X, barcode.Position.Y)

			return nil
		},
	}
}
