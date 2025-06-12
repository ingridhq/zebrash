package drawers

import (
	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
)

func NewGraphicCircleDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
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
