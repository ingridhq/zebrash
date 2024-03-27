package drawers

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/images"
)

func NewGraphicFieldDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, _ DrawerOptions) error {
			field, ok := element.(*elements.GraphicField)
			if !ok {
				return nil
			}

			width := field.RowBytes * 8
			height := len(field.Data) / field.RowBytes

			img := image.NewRGBA(image.Rect(0, 0, width, height))

			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					idx := y*(width/8) + x/8
					if idx >= len(field.Data) {
						continue
					}

					// Width for our bitmap data is in bits because each pixel is represented by one bit
					// but the actual data we have is in bytes
					// Here we access the value of each bit and check if it is 1 or 0
					val := ((field.Data[idx]) >> (7 - x%8)) & 1
					if val != 0 {
						img.SetRGBA(x, y, color.RGBA{A: 255})
					}
				}
			}

			imgScaled := images.NewScaled(img, float64(field.MagnificationX), float64(field.MagnificationY))
			gCtx.DrawImage(imgScaled, field.Position.X, field.Position.Y)

			return nil
		},
	}
}
