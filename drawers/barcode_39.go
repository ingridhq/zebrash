package drawers

import (
	"fmt"
	"image"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/images"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

func NewBarcode39Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, _ DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.Barcode39WithData)
			if !ok {
				return nil
			}

			// data to encode into barcode
			content := barcode.Data
			// human-readable text
			text := barcode.Data

			var (
				img image.Image
				err error
			)

			enc := oned.NewCode39Writer()
			img, err = enc.Encode(content, gozxing.BarcodeFormat_CODE_39, 1, 1, nil)
			if err != nil {
				return fmt.Errorf("failed to encode code39 barcode: %w", err)
			}

			img = images.NewScaled(img, barcode.Width, barcode.Height)
			width := float64(img.Bounds().Dx())
			height := float64(img.Bounds().Dy())
			pos := adjustImageTypeSetPosition(img, barcode.Position, barcode.Orientation)

			rotateImage(gCtx, img, pos, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, pos.X, pos.Y)
			if barcode.Line {
				applyLineTextToCtx(gCtx, text, pos, barcode.LineAbove, width, height)
			}

			return nil
		},
	}
}
