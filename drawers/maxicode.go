package drawers

import (
	"bytes"
	"fmt"
	"image"

	"github.com/fogleman/gg"
	"github.com/ingridhq/maxicode"
	"github.com/ingridhq/zebrash/elements"
)

func NewMaxicodeDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, options DrawerOptions) error {
			barcode, ok := element.(*elements.MaxicodeWithData)
			if !ok {
				return nil
			}

			inputData, err := barcode.GetInputData()
			if err != nil {
				return err
			}

			grid, err := maxicode.Encode(barcode.Code.Mode, 0, inputData)
			if err != nil {
				return fmt.Errorf("failed to encode maxicode grid: %w", err)
			}

			var buff bytes.Buffer

			err = grid.Draw(float64(options.Dpmm)).EncodePNG(&buff)
			if err != nil {
				return fmt.Errorf("failed to encode maxicode png: %w", err)
			}

			img, _, err := image.Decode(&buff)
			if err != nil {
				return fmt.Errorf("failed to convert maxicode png to image: %w", err)
			}

			gCtx.DrawImage(img, barcode.Pos.X, barcode.Pos.Y)

			return nil
		},
	}
}
