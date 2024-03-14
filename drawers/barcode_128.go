package drawers

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	barcodeLib "github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/elements"
)

const barcodeLineFontSizeScaleFactor = 20

var startCodeRegex = regexp.MustCompile(`^(>[9:;])`)

var invalidStartInvocationRegex = regexp.MustCompile(`^.+>[9:;]`)

func NewBarcode128Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, options DrawerOptions) error {
			barcode, ok := element.(*elements.Barcode128WithData)
			if !ok {
				return nil
			}

			content := invalidStartInvocationRegex.ReplaceAllString(barcode.Data, "")
			switch barcode.Mode {
			case elements.BarcodeModeNormal:
				content = modifyBarcodeContentNormalMode(content)
			case elements.BarcodeModeEan:
				content = modifyBarcodeContentEanMode(content)
			case elements.BarcodeModeUcc:
				content = modifyBarcodeContentUccMode(content)
			case elements.BarcodeModeAutomatic:
				// DO NOTHING with the content
				// subset prefixes like `>;` will not be skipped and encoded as label content
			}
			encodedBarcode, err := code128.Encode(content)
			if err != nil {
				return fmt.Errorf("failed to encode barcode: %s", err.Error())
			}
			scaledBarcode, err := barcodeLib.Scale(encodedBarcode, encodedBarcode.Bounds().Size().X*barcode.Width, encodedBarcode.Bounds().Max.Y*barcode.Height)
			if err != nil {
				return fmt.Errorf("failed to scale barcode: %s", err.Error())
			}

			width := float64(scaledBarcode.Bounds().Dx())
			height := float64(scaledBarcode.Bounds().Dy())
			applyBarcodeRotationToCtx(gCtx, barcode, width, height)
			defer gCtx.Identity()

			drawImage(gCtx, scaledBarcode, barcode.Position.X, barcode.Position.Y)
			if barcode.Line {
				applyLineTextToCtx(gCtx, content, barcode, width, height)
			}

			return nil
		},
	}
}

func applyBarcodeRotationToCtx(gCtx *gg.Context, barcodeElement *elements.Barcode128WithData, width, height float64) {
	if rotate := barcodeElement.Orientation.GetDegrees(); rotate != 0 {
		gCtx.RotateAbout(gg.Radians(rotate), float64(barcodeElement.Position.X), float64(barcodeElement.Position.Y))

		switch barcodeElement.Orientation {
		case elements.FieldOrientation90:
			gCtx.Translate(0, -height)
		case elements.FieldOrientation180:
			gCtx.Translate(-width, -height)
		case elements.FieldOrientation270:
			gCtx.Translate(-width, 0)
		}
	}
}

func applyLineTextToCtx(gCtx *gg.Context, content string, barcodeElement *elements.Barcode128WithData, width, height float64) {
	gCtx.SetColor(color.Black)
	fontSize := width / barcodeLineFontSizeScaleFactor

	face := truetype.NewFace(font0, &truetype.Options{Size: fontSize})
	gCtx.SetFontFace(face)

	x := float64(barcodeElement.Position.X) + float64(width)/2
	y := float64(barcodeElement.Position.Y) - fontSize
	if !barcodeElement.LineAbove {
		y = float64(barcodeElement.Position.Y) + height + fontSize
	}
	gCtx.DrawStringAnchored(content, x, y, 0.5, 0.5)
}

func modifyBarcodeContentNormalMode(content string) string {
	// replace beginning if it's a match
	content = startCodeRegex.ReplaceAllString(content, "")
	// support hand-rolled GS1
	return strings.ReplaceAll(content, `>8`, string(code128.FNC1))
}

func modifyBarcodeContentEanMode(content string) string {
	content = strings.ReplaceAll(content, `>8`, string(code128.FNC1))
	if !strings.HasPrefix(content, string(code128.FNC1)) {
		content = string(code128.FNC1) + content
	}
	return content
}

func modifyBarcodeContentUccMode(content string) string {
	content = addZerosPrefix(content)
	content = content[:19]
	checksumDigit := calculateUccBarcodeChecksumDigit(content)
	return string(code128.FNC1) + content + strconv.Itoa(checksumDigit)
}

func addZerosPrefix(in string) string {
	prefixLen := 19 - len(in)
	var b strings.Builder
	for i := 0; i < prefixLen; i++ {
		b.WriteRune('0')
	}
	b.WriteString(in)
	return b.String()
}

func calculateUccBarcodeChecksumDigit(content string) int {
	checksum := 0
	for i := 0; i < 19; i++ {
		checksum += int(content[i]-48) * (i%2*2 + 7)
	}
	return checksum % 10
}
