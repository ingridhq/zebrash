package drawers

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/assets"
	"github.com/ingridhq/zebrash/elements"
)

var (
	font0 = mustLoadFont(assets.FontHelveticaBold)
	font1 = mustLoadFont(assets.FontDejavuSansMono)
)

func NewTextFieldDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, _ DrawerOptions) error {
			text, ok := element.(*elements.TextField)
			if !ok {
				return nil
			}

			fontWidth := float64(text.Font.Width)
			fontSize := float64(text.Font.GetSize())

			font := font0
			if text.Font.Name != "0" {
				font = font1
			}

			face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
			gCtx.SetFontFace(face)
			setLineColor(gCtx, elements.LineColorBlack)

			x := float64(text.Position.X)
			y := float64(text.Position.Y) + fontSize/4.0
			ax := 0.0
			ay := 0.5

			if fontWidth != fontSize {
				gCtx.ScaleAbout(fontWidth/fontSize, 1, x, y)
			}

			if rotate := text.Orientation.GetDegrees(); rotate != 0 {
				switch text.Orientation {
				case elements.FieldOrientation90:
					ay = 0
				case elements.FieldOrientation270:
					ax = 1.0
				case elements.FieldOrientation180:
					ax = 1.0
					ay = 0
				}

				gCtx.RotateAbout(gg.Radians(rotate), x, y)
			}

			defer gCtx.Identity()

			if text.Block != nil {
				align := gg.AlignLeft

				switch text.Block.Alignment {
				case elements.TextAlignmentRight:
					align = gg.AlignRight
				case elements.TextAlignmentJustified:
					align = gg.AlignCenter
				}

				gCtx.DrawStringWrapped(text.Text, x, y, ax, ay, float64(text.Block.MaxWidth), float64(text.Block.LineSpacing), align)
			} else {
				gCtx.DrawStringAnchored(text.Text, x, y, ax, ay)
			}

			return nil
		},
	}
}
