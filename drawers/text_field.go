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
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, state *DrawerState) error {
			text, ok := element.(*elements.TextField)
			if !ok {
				return nil
			}

			text = adjustTextField(text)

			fontSize := text.Font.GetSize()
			scaleX := text.Font.GetScaleX()
			face := truetype.NewFace(getTffFont(text.Font), &truetype.Options{Size: fontSize})
			gCtx.SetFontFace(face)

			setLineColor(gCtx, elements.LineColorBlack)

			w, h := gCtx.MeasureString(text.Text)

			x, y := getTextTopLeftPos(text, w, h, state)
			state.UpdateAutomaticTextPosition(text, w, scaleX)

			ax, ay := getTextAxAy(text)

			if scaleX != 1.0 {
				gCtx.ScaleAbout(scaleX, 1, x, y)
			}

			if rotate := text.Font.Orientation.GetDegrees(); rotate != 0 {
				gCtx.RotateAbout(gg.Radians(rotate), x, y)
			}

			defer gCtx.Identity()

			if text.Block != nil {
				maxWidth := float64(text.Block.MaxWidth) / scaleX
				align := getTextBlockAlign(text.Block)
				gCtx.DrawStringWrapped(text.Text, x, y-h, ax, ay, maxWidth, 1+float64(text.Block.LineSpacing), align)
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

func getTextTopLeftPos(text *elements.TextField, w, h float64, state *DrawerState) (float64, float64) {
	x, y := state.GetTextPosition(text)

	if !text.Position.CalculateFromBottom {
		switch text.Font.Orientation {
		case elements.FieldOrientation90:
			return x + h/4, y
		case elements.FieldOrientation180:
			return x + w, y + h/4
		case elements.FieldOrientation270:
			return x + 3*h/4, y + w
		default:
			return x, y + 3*h/4
		}
	}

	lines := 1.0
	spacing := 0.0

	if text.Block != nil {
		lines = float64(max(text.Block.MaxLines, 1))
		spacing = float64(text.Block.LineSpacing)
	}

	offset := (lines - 1) * (h + spacing)

	switch text.Font.Orientation {
	case elements.FieldOrientation90:
		return x + offset, y
	case elements.FieldOrientation180:
		return x, y + offset
	case elements.FieldOrientation270:
		return x - offset, y
	default:
		return x, y - offset
	}
}

func getTextAxAy(text *elements.TextField) (float64, float64) {
	ax := 0.0
	ay := 0.0

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

func mustLoadFont(fontData []byte) *truetype.Font {
	font, err := truetype.Parse(fontData)
	if err != nil {
		panic(err.Error())
	}

	return font
}
