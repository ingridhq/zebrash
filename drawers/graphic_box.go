package drawers

import (
	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
)

func NewGraphicBoxDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}) error {
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

			offsetX := border / 2.0
			offsetY := border / 2.0

			x := float64(box.Position.X) + offsetX
			y := float64(box.Position.Y) + offsetY

			drawRectangle(gCtx, x, y, width-border, height-border)
			setLineColor(gCtx, box.LineColor)

			gCtx.SetLineCapSquare()
			gCtx.SetLineWidth(border)
			gCtx.Stroke()

			return nil
		},
	}
}
