package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/drawers"
	"github.com/ingridhq/zebrash/internal/barcodes/aztec" // Add this line
	"github.com/ingridhq/zebrash/internal/elements"
)

func NewBarcodeAztecDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ drawers.DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeAztecWithData)
			if !ok {
				return nil
			}

			var layers int = aztec.DEFAULT_LAYERS
			var minECCPercent int = aztec.DEFAULT_EC_PERCENT

			const sizeModeFullRangeOffset = 200

			if barcode.Size > 0 {
				switch {
				// Case for static full-range symbol sizes (e.g., 201-232)
				case barcode.Size >= sizeModeFullRangeOffset && barcode.Size <= sizeModeFullRangeOffset+32:
					layers = barcode.Size - sizeModeFullRangeOffset
				// Case for dynamic sizing based on error correction percentage (e.g., 1-99)
				case barcode.Size >= 1 && barcode.Size <= 99:
					minECCPercent = barcode.Size
					// layers remains aztec.DEFAULT_LAYERS (0) to let the library auto-determine
				// Case for compact modes (e.g., 101-104)
				case barcode.Size >= 101 && barcode.Size <= 104:
					layers = -(barcode.Size - 100) // Negative number indicates compact
				default:
					return fmt.Errorf("aztec barcode size/mode %d is not supported", barcode.Size)
				}
			}

			img, err := aztec.Encode([]byte(barcode.Data), minECCPercent, layers, barcode.Magnification)
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
