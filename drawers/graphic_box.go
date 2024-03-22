package drawers

import (
	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
)

func NewGraphicBoxDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, _ DrawerOptions) error {
			box, ok := element.(*elements.GraphicBox)
			if !ok {
				return nil
			}

			width := float64(box.Width)
			height := float64(box.Height)
			border := float64(box.BorderThickness)

			if border > width {
				width = border
			}

			if border > height {
				height = border
			}

			offset := border / 2.0
			x := float64(box.Position.X) + offset
			y := float64(box.Position.Y) + offset
			w := width - border
			h := height - border

			if w == 0 && h == 0 {
				w, h = 1, 1
			}

			drawRectangle(gCtx, x, y, w, h)
			setLineColor(gCtx, box.LineColor)

			gCtx.SetLineCapSquare()
			gCtx.SetLineWidth(border)
			gCtx.Stroke()

			return nil
		},
	}
}
