package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/ingridhq/maxicode"
	"github.com/ingridhq/zebrash/elements"
)

func NewMaxicodeDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, options DrawerOptions, _ *DrawerState) error {
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

			img := grid.Draw(float64(options.Dpmm)).Image()

			gCtx.DrawImage(img, barcode.Position.X, barcode.Position.Y)

			return nil
		},
	}
}
