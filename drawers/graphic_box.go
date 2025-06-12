package drawers

import (
	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
)

func NewGraphicBoxDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
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

			setLineColor(gCtx, box.LineColor)
			gCtx.SetLineCapSquare()
			gCtx.SetLineWidth(border)

			if box.CornerRounding > 0 {
				drawRoundedRectangle(gCtx, float64(box.Position.X), float64(box.Position.Y), width, height, float64(box.CornerRounding), border)
			} else {
				drawRectangle(gCtx, x, y, w, h)
			}

			return nil
		},
	}
}

func drawRectangle(gCtx *gg.Context, x, y, w, h float64) {
	gCtx.DrawLine(x, y, x+w, y)
	gCtx.DrawLine(x+w, y, x+w, y+h)
	gCtx.DrawLine(x+w, y+h, x, y+h)
	gCtx.DrawLine(x, y+h, x, y)
	gCtx.Stroke()
}

func drawRoundedRectangle(gCtx *gg.Context, x, y, w, h, rounding, border float64) {
	defer gCtx.ResetClip()

	r2 := rounding * (min(w-2*border, h-2*border)) / 16
	if r2 > 0 {
		gCtx.DrawRoundedRectangle(x+border, y+border, w-2*border, h-2*border, r2)
		gCtx.Clip()
		gCtx.InvertMask()
	}

	r1 := rounding * (min(w, h)) / 16
	gCtx.DrawRoundedRectangle(x, y, w, h, r1)
	gCtx.Fill()
}
