package drawers

import (
	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
)

func NewGraphicDiagonalLineDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, _ DrawerOptions) error {
			line, ok := element.(*elements.GraphicDiagonalLine)
			if !ok {
				return nil
			}

			x := float64(line.Position.X)
			y := float64(line.Position.Y)
			width := float64(line.Width)
			height := float64(line.Height)
			border := float64(line.BorderThickness)

			if line.TopToBottom {
				gCtx.DrawLine(x, y, x+width, y+height)
			} else {
				gCtx.DrawLine(x+width, y, x, y+height)
			}

			setLineColor(gCtx, line.LineColor)

			gCtx.SetLineCapSquare()
			gCtx.SetLineWidth(border)
			gCtx.Stroke()

			return nil
		},
	}
}
