package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/barcodes/aztec"
	"github.com/ingridhq/zebrash/elements"
)

func NewBarcodeAztecDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeAztecWithData)
			if !ok {
				return nil
			}

			layers := 0

			// Aztec barcode full range modes in ZPL are in range 200, 232
			// Library that we use works in modes 0,32
			// So we need to translate ZPL mode into library mode by removing the offset
			const sizeModeFullRangeOffset = 200

			if barcode.Size > 0 {
				switch {
				case barcode.Size >= sizeModeFullRangeOffset && barcode.Size <= sizeModeFullRangeOffset+32:
					layers = barcode.Size - sizeModeFullRangeOffset
				default:
					return fmt.Errorf("aztec barcode size %d is not supported", barcode.Size)
				}
			}

			img, err := aztec.Encode([]byte(barcode.Data), 0, layers, barcode.Magnification)
			if err != nil {
				return fmt.Errorf("failed to encode aztec barcode: %w", err)
			}

			pos := adjustImageTypeSetPosition(img, barcode.Position, barcode.Orientation)

			rotateImage(gCtx, img, pos, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, pos.X, pos.Y)

			return nil
		},
	}
}
