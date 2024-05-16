package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/images"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"github.com/makiuchi-d/gozxing/datamatrix/encoder"
)

func NewBarcodeDatamatrixDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, options DrawerOptions) error {
			barcode, ok := element.(*elements.BarcodeDatamatrixWithData)
			if !ok {
				return nil
			}

			enc := datamatrix.NewDataMatrixWriter()

			hints := make(map[gozxing.EncodeHintType]interface{})

			switch barcode.Ratio {
			case elements.DatamatrixRatioSquare:
				hints[gozxing.EncodeHintType_DATA_MATRIX_SHAPE] = encoder.SymbolShapeHint_FORCE_SQUARE
			case elements.DatamatrixRatioRectangular:
				hints[gozxing.EncodeHintType_DATA_MATRIX_SHAPE] = encoder.SymbolShapeHint_FORCE_RECTANGLE
			}

			img, err := enc.Encode(barcode.Data, gozxing.BarcodeFormat_DATA_MATRIX, barcode.Columns, barcode.Rows, hints)
			if err != nil {
				return fmt.Errorf("failed to encode datamatrix barcode: %w", err)
			}

			ratio := float64(max(barcode.Height, 1))

			scaledImg := images.NewScaled(img, ratio, ratio)

			rotateImage(gCtx, scaledImg, barcode.Position, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(scaledImg, barcode.Position.X, barcode.Position.Y)

			return nil
		},
	}
}
