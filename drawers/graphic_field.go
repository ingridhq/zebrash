package drawers

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
)

func NewGraphicFieldDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, _ DrawerOptions) error {
			field, ok := element.(*elements.GraphicField)
			if !ok {
				return nil
			}

			if field.DataBytes != len(field.Data) {
				return fmt.Errorf("unexpected length of graphic field data, expected %d bytes", field.DataBytes)
			}

			width := field.RowBytes * 8
			height := field.DataBytes / field.RowBytes

			img := image.NewGray(image.Rect(0, 0, width, height))

			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					// Width for our bitmap data is in bits because each pixel is represented by one bit
					// but the actual data we have is in bytes
					// Here we access the value of each bit and check if it is 1 or 0
					val := ((field.Data[y*(width/8)+x/8]) >> (7 - x%8)) & 0xF

					if val != 0 {
						img.SetGray(x, y, color.Gray{Y: 0})
					} else {
						img.SetGray(x, y, color.Gray{Y: 255})
					}
				}
			}

			gCtx.DrawImage(img, field.Position.X, field.Position.Y)

			return nil
		},
	}
}
