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
	"github.com/ingridhq/zebrash/barcodes/code128"
	"github.com/ingridhq/zebrash/elements"
)

const (
	barcodeLineFontSizeScaleFactor = 20
)

var (
	barcode128FNC1 = string(code128.ESCAPE_FNC_1)

	parenthesisAndSpacesRegex = regexp.MustCompile(`[\(\)\s]`)
)

func NewBarcode128Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element interface{}, options DrawerOptions) error {
			barcode, ok := element.(*elements.Barcode128WithData)
			if !ok {
				return nil
			}

			// data to encode into barcode
			content := barcode.Data
			// human-readable text
			text := barcode.Data

			switch barcode.Mode {
			case elements.BarcodeModeEan:
				content, text = modifyBarcodeContentEanMode(content)
			case elements.BarcodeModeUcc:
				content = modifyBarcodeContentUccMode(content)
			case elements.BarcodeModeAutomatic:
				// DO NOTHING with the content
				// invocation codes like `>;` will not be skipped and instead will be encoded as label content
			}

			var (
				img image.Image
				err error
			)

			switch barcode.Mode {
			case elements.BarcodeModeNo:
				img, text, err = code128.EncodeNoMode(content, barcode.Height, barcode.Width)
			default:
				img, err = code128.EncodeAuto(content, barcode.Height, barcode.Width)
			}

			if err != nil {
				return fmt.Errorf("failed to encode barcode: %w", err)
			}

			width := float64(img.Bounds().Dx())
			height := float64(img.Bounds().Dy())

			rotateImage(gCtx, img, barcode.Position, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, barcode.Position.X, barcode.Position.Y)
			if barcode.Line {
				applyLineTextToCtx(gCtx, text, barcode, width, height)
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

// Allows dealing with UCC/EAN with and without chained application identifiers.
// The code starts in the appropriate subset followed by FNC1 to indicate a UCC/EAN 128 barcode.
// The printer automatically strips out parenthesis and spaces for encoding but prints them in the human-readable section.
// The printer automatically determines if a check digit is required, calculates it, and prints it.
func modifyBarcodeContentEanMode(content string) (string, string) {
	// Don't show special functions in human-readable text
	text := strings.ReplaceAll(content, ">8", "")

	content = parenthesisAndSpacesRegex.ReplaceAllString(content, "")
	content = strings.ReplaceAll(content, ">8", barcode128FNC1)
	if !strings.HasPrefix(content, barcode128FNC1) {
		content = barcode128FNC1 + content
	}

	return content, text
}

// Content must contain 19 numeric digits.
// Subset C using FNC1 values is automatically selected.
// Excess digits (above 19) in ^FD or ^SN will be eliminated.
// Below 19 digits in ^FD or ^SN adds zeros to right to bring count to 19.
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
