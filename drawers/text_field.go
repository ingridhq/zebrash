package drawers

import (
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/assets"
	"github.com/ingridhq/zebrash/elements"
)

var (
	font0 = mustLoadFont(assets.FontHelveticaBold)
	font1 = mustLoadFont(assets.FontDejavuSansMono)
	fontB = mustLoadFont(assets.FontDejavuSansMonoBold)
)

func NewTextFieldDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, _ DrawerOptions) error {
			text, ok := element.(*elements.TextField)
			if !ok {
				return nil
			}

			text = adjustTextField(text)

			fontSize := text.Font.GetSize()
			face := truetype.NewFace(getTffFont(text.Font), &truetype.Options{Size: fontSize})
			gCtx.SetFontFace(face)

			setLineColor(gCtx, elements.LineColorBlack)

			x, y := getTextTopLeftPos(gCtx, text)
			ax, ay := getTextAxAy(text)
			scaleX := 1.0

			if text.Font.Width != 0 {
				scaleX = text.Font.Width / fontSize
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

func adjustTextField(text *elements.TextField) *elements.TextField {
	fontName := text.Font.Name
	res := *text

	if fontName != "0" {
		// For some reason width of DejavuSansMono needs to be scaled by 2
		res.Font.Width *= 2
	}

	switch fontName {
	case "B":
		// Bold font, text in all uppercase
		res.Text = strings.ToUpper(res.Text)
	}

	return &res
}

func getTffFont(font elements.FontInfo) *truetype.Font {
	switch font.Name {
	case "0":
		return font0
	case "B":
		return fontB
	default:
		return font1
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
