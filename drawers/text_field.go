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
			fontSize := getTextFontSize(text)

			font := font0
			if text.Font.Name != "0" {
				font = font1
			}

			face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
			gCtx.SetFontFace(face)

			setLineColor(gCtx, elements.LineColorBlack)

			x, y := getTextTopLeftPos(gCtx, text)
			ax, ay := getTextAxAy(text)
			scaleX := 1.0

			if fontWidth != fontSize {
				scaleX = fontWidth / fontSize
				gCtx.ScaleAbout(scaleX, 1, x, y)
			}

			if rotate := text.Font.Orientation.GetDegrees(); rotate != 0 {
				gCtx.RotateAbout(gg.Radians(rotate), x, y)
			}

			defer gCtx.Identity()

			if text.Block != nil {
				maxWidth := float64(text.Block.MaxWidth) / scaleX
				align := getTextBlockAlign(text.Block)
				gCtx.DrawStringWrapped(text.Text, x, y, ax, ay, maxWidth, float64(text.Block.LineSpacing), align)
			} else {
				gCtx.DrawStringAnchored(text.Text, x, y, ax, ay)
			}

			return nil
		},
	}
}

func getTextTopLeftPos(gCtx *gg.Context, text *elements.TextField) (float64, float64) {
	x := float64(text.Position.X)
	y := float64(text.Position.Y)

	w, h := getTextDimensions(gCtx, text)

	if !text.Position.CalculateFromBottom {
		switch text.Font.Orientation {
		case elements.FieldOrientation90:
			return x + 3*h/4, y
		case elements.FieldOrientation180:
			return x + w, y + 3*h/4
		case elements.FieldOrientation270:
			return x + h/4, y + w
		default:
			return x, y + h/4
		}
	}

	switch text.Font.Orientation {
	case elements.FieldOrientation90:
		return x + h/2, y
	case elements.FieldOrientation180:
		return x, y + h/2
	case elements.FieldOrientation270:
		return x - h/2, y
	default:
		return x, y - h/2
	}
}

func getTextDimensions(gCtx *gg.Context, text *elements.TextField) (float64, float64) {
	if text.Block != nil {
		return gCtx.MeasureMultilineString(text.Text, float64(text.Block.LineSpacing))
	}

	return gCtx.MeasureString(text.Text)
}

func getTextFontSize(text *elements.TextField) float64 {
	w := float64(text.Font.Width)
	h := float64(text.Font.Height)

	switch text.Font.Orientation {
	case elements.FieldOrientation90, elements.FieldOrientation270:
		return w
	default:
		return h
	}
}

func getTextAxAy(text *elements.TextField) (float64, float64) {
	ax := 0.0
	ay := 0.5

	switch text.Alignment {
	case elements.TextAlignmentLeft:
		ax = 0
	case elements.TextAlignmentRight:
		ax = 1
	case elements.TextAlignmentJustified, elements.TextAlignmentCenter:
		ax = 0.5
	}

	return ax, ay
}

func getTextBlockAlign(block *elements.FieldBlock) gg.Align {
	switch block.Alignment {
	case elements.TextAlignmentRight:
		return gg.AlignRight
	case elements.TextAlignmentJustified, elements.TextAlignmentCenter:
		return gg.AlignCenter
	default:
		return gg.AlignLeft
	}
}
