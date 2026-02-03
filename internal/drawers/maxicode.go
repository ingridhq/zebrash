package drawers

import (
	"fmt"
	"math"

	"github.com/ingridhq/gg"
	"github.com/ingridhq/maxicode"
	"github.com/ingridhq/zebrash/drawers"
	"github.com/ingridhq/zebrash/internal/elements"
)

func NewMaxicodeDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, options drawers.DrawerOptions, _ *DrawerState) error {
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

			dpmm := float64(options.Dpmm)
			hexRectW := int(math.Round(0.76 * dpmm))
			hexRectH := int(math.Round(0.88 * dpmm))

			img := grid.Draw(dpmm).Image()
			pos := adjustImageTypeSetPosition(img, barcode.Position, elements.FieldOrientationNormal)

			gCtx.DrawImage(img, pos.X-hexRectW, pos.Y-hexRectH)

			return nil
		},
	}
}
