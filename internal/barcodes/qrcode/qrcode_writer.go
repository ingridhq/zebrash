package qrcode

import (
	"fmt"

	"github.com/DawidBury/zebrash/internal/barcodes/qrcode/encoder"
	"github.com/DawidBury/zebrash/internal/barcodes/utils"
)

func Encode(contents string, width, height int, errorCorrectionLevel encoder.ErrorCorrectionLevel, opts encoder.Options) (*utils.BitMatrix, error) {
	code, e := encoder.Encoder_encode(contents, errorCorrectionLevel, opts)
	if e != nil {
		return nil, e
	}

	return renderResult(code, width, height, opts)
}

// renderResult Note that the input matrix uses 0 == white, 1 == black, while the output matrix uses
// transparent and black rgb colors
func renderResult(code *encoder.QRCode, width, height int, opts encoder.Options) (*utils.BitMatrix, error) {
	input := code.GetMatrix()
	if input == nil {
		return nil, fmt.Errorf("IllegalStateException")
	}

	quietZone := opts.QuietZone
	inputWidth := input.GetWidth()
	inputHeight := input.GetHeight()
	qrWidth := inputWidth + (quietZone * 2)
	qrHeight := inputHeight + (quietZone * 2)
	outputWidth := qrWidth
	if outputWidth < width {
		outputWidth = width
	}
	outputHeight := qrHeight
	if outputHeight < height {
		outputHeight = height
	}

	multiple := outputWidth / qrWidth
	if h := outputHeight / qrHeight; multiple > h {
		multiple = h
	}
	// Padding includes both the quiet zone and the extra white pixels to accommodate the requested
	// dimensions. For example, if input is 25x25 the QR will be 33x33 including the quiet zone.
	// If the requested size is 200x160, the multiple will be 4, for a QR of 132x132. These will
	// handle all the padding from 100x100 (the actual QR) up to 200x160.
	leftPadding := (outputWidth - (inputWidth * multiple)) / 2
	topPadding := (outputHeight - (inputHeight * multiple)) / 2

	output, e := utils.NewBitMatrix(outputWidth, outputHeight)
	if e != nil {
		return nil, e
	}

	for inputY, outputY := 0, topPadding; inputY < inputHeight; inputY, outputY = inputY+1, outputY+multiple {
		// Write the contents of this row of the barcode
		for inputX, outputX := 0, leftPadding; inputX < inputWidth; inputX, outputX = inputX+1, outputX+multiple {
			if input.Get(inputX, inputY) == 1 {
				output.SetRegion(outputX, outputY, multiple, multiple)
			}
		}
	}

	return output, nil
}
