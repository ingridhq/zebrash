package drawers

import (
	"fmt"
	"image"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/images"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

const (
	barcodeLineFontSizeScaleFactor = 20
	// FNC1 - Special Function 1
	barcode128FNC1 = "\u00f1"
)

var (
	startCodeRegex              = regexp.MustCompile(`^(>[9:;])`)
	invalidStartInvocationRegex = regexp.MustCompile(`^.+>[9:;]`)
)

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

			var (
				img image.Image
				err error
			)

			enc := oned.NewCode128Writer()
			img, err = enc.Encode(content, gozxing.BarcodeFormat_CODE_128, 0, barcode.Height, nil)
			if err != nil {
				return fmt.Errorf("failed to encode barcode: %s", err.Error())
			}

			img = images.NewShifted(img, -4, 0)
			img = images.NewScaled(img, float64(barcode.Width), 1)

			width := float64(img.Bounds().Dx())
			height := float64(img.Bounds().Dy())

			rotateImage(gCtx, img, barcode.Position, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(images.NewTransparent(img), barcode.Position.X, barcode.Position.Y)
			if barcode.Line {
				applyLineTextToCtx(gCtx, content, barcode, width, height)
			}

			return nil
		},
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
	return strings.ReplaceAll(content, `>8`, barcode128FNC1)
}

func modifyBarcodeContentEanMode(content string) string {
	content = strings.ReplaceAll(content, `>8`, barcode128FNC1)
	if !strings.HasPrefix(content, barcode128FNC1) {
		content = barcode128FNC1 + content
	}
	return content
}

func modifyBarcodeContentUccMode(content string) string {
	content = addZerosPrefix(content)
	content = content[:19]
	checksumDigit := calculateUccBarcodeChecksumDigit(content)
	return barcode128FNC1 + content + strconv.Itoa(checksumDigit)
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
