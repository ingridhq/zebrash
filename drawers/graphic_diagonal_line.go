package drawers

import (
	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
)

func NewGraphicDiagonalLineDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
			line, ok := element.(*elements.GraphicDiagonalLine)
			if !ok {
				return nil
			}

			x := float64(line.Position.X)
			y := float64(line.Position.Y)
			width := float64(line.Width)
			height := float64(line.Height)
			border := float64(line.BorderThickness)

			drawDiagonalLine(gCtx, x, y, width, height, border, line.TopToBottom)

			setLineColor(gCtx, line.LineColor)
			gCtx.SetLineWidth(1)
			gCtx.SetLineCapSquare()
			gCtx.Fill()

			return nil
		},
	}
}

func drawDiagonalLine(gCtx *gg.Context, x, y, w, h, b float64, bottomToTop bool) {
	gCtx.NewSubPath()
	defer gCtx.ClosePath()

	if bottomToTop {
		gCtx.MoveTo(x, y)
		gCtx.LineTo(x+b, y)
		gCtx.LineTo(x+b+w, y+h)
		gCtx.LineTo(x+w, y+h)
		gCtx.LineTo(x, y)
		return
	}

	gCtx.MoveTo(x, y+h)
	gCtx.LineTo(x+b, y+h)
	gCtx.LineTo(x+b+w, y)
	gCtx.LineTo(x+w, y)
	gCtx.LineTo(x, y+h)
}
