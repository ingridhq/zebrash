package drawers

import (
	"fmt"
	"strings"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/barcodes/datamatrix"
	"github.com/ingridhq/zebrash/barcodes/datamatrix/encoder"
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/images"
)

func NewBarcodeDatamatrixDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeDatamatrixWithData)
			if !ok {
				return nil
			}

			columns := max(barcode.Columns, 1)
			rows := max(barcode.Rows, 1)

			enc := datamatrix.NewDataMatrixWriter()

			opts := encoder.Options{
				MinSize: encoder.NewDimension(columns, rows),
			}

			switch barcode.Ratio {
			case elements.DatamatrixRatioSquare:
				opts.Shape = encoder.SymbolShapeHint_FORCE_SQUARE
			case elements.DatamatrixRatioRectangular:
				opts.Shape = encoder.SymbolShapeHint_FORCE_RECTANGLE
			}

			data := barcode.Data

			const (
				fnc1 = "_1"
				GS   = byte(29)
			)

			// First occurrence of FNC1 triggers GS1 mode
			if strings.HasPrefix(data, fnc1) {
				opts.Gs1 = true
				data = strings.TrimPrefix(data, fnc1)
			}

			// All subsequent occurrences of FNC1 are encoded as GS character
			data = strings.ReplaceAll(data, fnc1, string(GS))

			img, err := enc.Encode(data, columns, rows, opts)
			if err != nil {
				return fmt.Errorf("failed to encode datamatrix barcode: %w", err)
			}

			scale := max(barcode.Height, 1)
			scaledImg := images.NewScaled(img, scale, scale)
			pos := adjustImageTypeSetPosition(scaledImg, barcode.Position, barcode.Orientation)

			rotateImage(gCtx, scaledImg, pos, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(scaledImg, pos.X, pos.Y)

			return nil
		},
	}
}
