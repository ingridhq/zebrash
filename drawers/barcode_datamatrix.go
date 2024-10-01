package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"github.com/makiuchi-d/gozxing/datamatrix/encoder"
)

func NewBarcodeDatamatrixDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, _ DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeDatamatrixWithData)
			if !ok {
				return nil
			}

			enc := datamatrix.NewDataMatrixWriter()

			dims, err := gozxing.NewDimension(barcode.Columns, barcode.Rows)
			if err != nil {
				return fmt.Errorf("failed to create datamatrix dimensions: %w", err)
			}

			hints := map[gozxing.EncodeHintType]interface{}{
				gozxing.EncodeHintType_MIN_SIZE: dims,
			}

			switch barcode.Ratio {
			case elements.DatamatrixRatioSquare:
				hints[gozxing.EncodeHintType_DATA_MATRIX_SHAPE] = encoder.SymbolShapeHint_FORCE_SQUARE
			case elements.DatamatrixRatioRectangular:
				hints[gozxing.EncodeHintType_DATA_MATRIX_SHAPE] = encoder.SymbolShapeHint_FORCE_RECTANGLE
			}

			ratio := max(barcode.Height, 1)

			img, err := enc.Encode(barcode.Data, gozxing.BarcodeFormat_DATA_MATRIX, ratio*barcode.Columns, ratio*barcode.Rows, hints)
			if err != nil {
				return fmt.Errorf("failed to encode datamatrix barcode: %w", err)
			}

			rotateImage(gCtx, img, barcode.Position, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, barcode.Position.X, barcode.Position.Y)

			return nil
		},
	}
}
