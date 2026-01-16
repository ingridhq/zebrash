package drawers

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/drawers"
	"github.com/ingridhq/zebrash/internal/barcodes/ean13"
	"github.com/ingridhq/zebrash/internal/elements"
)

func NewBarcodeEan13Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ drawers.DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeEan13WithData)
			if !ok {
				return nil
			}

			img, text, err := ean13.Encode(barcode.Data, barcode.Height, barcode.Width)
			if err != nil {
				return fmt.Errorf("failed to encode EAN-13 barcode: %w", err)
			}

			width := float64(img.Bounds().Dx())
			height := float64(img.Bounds().Dy())
			guardExtension := ean13.CalculateGuardExtension(barcode.Width)

			pos := adjustImageTypeSetPosition(img, barcode.Position, barcode.Orientation)
			pos = adjustEan13Position(pos, barcode.Orientation, guardExtension)

			rotateImage(gCtx, img, pos, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, pos.X, pos.Y)
			if barcode.Line {
				applyEan13TextToCtx(gCtx, text, pos, barcode.LineAbove, width, height, barcode.Width, guardExtension)
			}

			return nil
		},
	}
}

func adjustEan13Position(pos elements.LabelPosition, ori elements.FieldOrientation, guardExtension int) elements.LabelPosition {
	if pos.CalculateFromBottom {
		return pos
	}

	x := pos.X
	y := pos.Y

	switch ori {
	case elements.FieldOrientation90:
		x -= guardExtension
	case elements.FieldOrientation180:
		y -= guardExtension
	}

	return elements.LabelPosition{
		X: x,
		Y: y,
	}
}

func applyEan13TextToCtx(gCtx *gg.Context, text string, pos elements.LabelPosition, lineAbove bool, width, height float64, barWidth, guardExtension int) {
	gCtx.SetColor(color.Black)

	fontSize := float64(guardExtension) * 2
	face := truetype.NewFace(font0, &truetype.Options{Size: fontSize})
	gCtx.SetFontFace(face)

	if len(text) == 13 && !lineAbove {
		// Put some hidden characts between EAN-13 sections
		formattedText := fmt.Sprintf("%s||%s||%s", text[0:1], text[1:7], text[7:13])
		// Insert spaces between every digit so justified alignment works properly
		formattedText = wrapWithSpaces(formattedText)
		x := float64(pos.X) - float64(barWidth)*10
		y := float64(pos.Y) + height + float64(guardExtension)/2 + float64(barWidth)
		w := width + float64(barWidth)*5
		drawStringJustified(gCtx, formattedText, x, y, 0, 0, w, []string{"|"})
	} else {
		x := float64(pos.X) + float64(barWidth)*8
		y := float64(pos.Y) - float64(guardExtension)/2
		w := width - float64(barWidth)*16
		formattedText := wrapWithSpaces(text)
		drawStringJustified(gCtx, formattedText, x, y, 0, 0, w, nil)
	}
}

func wrapWithSpaces(text string) string {
	var b strings.Builder

	for _, r := range text {
		b.WriteRune(' ')
		b.WriteRune(r)
		b.WriteRune(' ')
	}

	return b.String()
}
