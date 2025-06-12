package drawers

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/images"
)

func NewGraphicFieldDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
			field, ok := element.(*elements.GraphicField)
			if !ok {
				return nil
			}

			dataLen := len(field.Data)
			if field.TotalBytes > 0 {
				dataLen = min(field.TotalBytes, dataLen)
			}

			width := field.RowBytes * 8
			height := dataLen / field.RowBytes

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
						img.SetRGBA(x, y, images.ColorBlack)
					}
				}
			}

			imgScaled := images.NewScaled(img, field.MagnificationX, field.MagnificationY)
			gCtx.DrawImage(imgScaled, field.Position.X, field.Position.Y)

			return nil
		},
	}
}
