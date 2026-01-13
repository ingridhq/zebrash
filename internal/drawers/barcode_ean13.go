package drawers

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

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

			// data to encode into barcode
			content := barcode.Data

			// Ensure human-readable text has checksum (convert to 13 digits if needed)
			text := barcode.Data
			if len(text) == 12 {
				// Calculate and append checksum
				checksum, err := ean13.CalculateChecksum(text)
				if err == nil {
					text = text + strconv.Itoa(checksum)
				}
			}

			var (
				img image.Image
				err error
			)

			img, err = ean13.Encode(content, barcode.Height, barcode.Width)
			if err != nil {
				return fmt.Errorf("failed to encode EAN-13 barcode: %w", err)
			}

			width := float64(img.Bounds().Dx())
			height := float64(img.Bounds().Dy())
			pos := adjustImageTypeSetPosition(img, barcode.Position, barcode.Orientation)

			rotateImage(gCtx, img, pos, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, pos.X, pos.Y)
			if barcode.Line {
				applyEan13TextToCtx(gCtx, text, pos, barcode.LineAbove, width, height, barcode.Width)
			}

			return nil
		},
	}
}

// applyEan13TextToCtx renders the human-readable text for EAN-13 barcodes
// The text should be positioned in the guard extension area and use a larger font
// EAN-13 standard layout:
// - First digit (number system) separated from the rest
// - Digits 2-7 centered under the left half (between start and middle guard)
// - Digits 8-13 centered under the right half (between middle and end guard)
func applyEan13TextToCtx(gCtx *gg.Context, content string, pos elements.LabelPosition, lineAbove bool, width, height float64, barWidth int) {
	gCtx.SetColor(color.Black)

	// EAN-13 uses text size that is not relative to barcode width
	// Font size is twice the guard extension height
	guardExtension := float64(barWidth * 5)
	fontSize := guardExtension * 2

	face := truetype.NewFace(font0, &truetype.Options{Size: fontSize})
	gCtx.SetFontFace(face)

	// Calculate text Y position
	var y float64
	if lineAbove {
		y = float64(pos.Y) - fontSize
	} else {
		// Place text in the guard extension area (below main bars, between guards)
		// Leave a gap width of one module between regular bars and text
		barcodeHeightWithoutGuard := height - guardExtension
		y = float64(pos.Y+barWidth) + barcodeHeightWithoutGuard + guardExtension/2
	}

	// Format the text with spaces to create gaps:
	// Original: "1234567890128"
	// Formatted: "1 234567 890128"
	// This creates natural gaps at the guard bars
	var formattedText string
	var x float64
	if len(content) == 13 && !lineAbove {
		formattedText = content[0:1] + "     " + content[1:7] + "       " + content[7:13]
		x = float64(pos.X) + width/2 - 2*guardExtension
	} else {
		formattedText = content
		x = float64(pos.X) + width/2
	}

	gCtx.DrawStringAnchored(formattedText, x, y, 0.5, 0.5)
}
