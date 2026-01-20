package drawers

import (
	"github.com/fogleman/gg"
	"github.com/DawidBury/zebrash/drawers"
	"github.com/DawidBury/zebrash/internal/elements"
)

func NewGraphicCircleDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ drawers.DrawerOptions, _ *DrawerState) error {
			circle, ok := element.(*elements.GraphicCircle)
			if !ok {
				return nil
			}

			setLineColor(gCtx, circle.LineColor)
			gCtx.SetLineWidth(float64(circle.BorderThickness))

			radius := float64(circle.CircleDiameter) / 2.0
			gCtx.DrawCircle(float64(circle.Position.X)+radius, float64(circle.Position.Y)+radius, radius)

			gCtx.Stroke()

			return nil
		},
	}
}
