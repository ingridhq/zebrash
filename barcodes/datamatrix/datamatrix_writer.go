package datamatrix

import (
	"fmt"

	"github.com/ingridhq/zebrash/barcodes/datamatrix/encoder"
)

// DataMatrixWriter This object renders a Data Matrix code as a BitMatrix 2D array of greyscale values.
type DataMatrixWriter struct{}

func NewDataMatrixWriter() *DataMatrixWriter {
	return &DataMatrixWriter{}
}

func (writer *DataMatrixWriter) Encode(contents string, width, height int, opts encoder.Options) (*BitMatrix, error) {
	if contents == "" {
		return nil, fmt.Errorf("found empty contents")
	}

	//1. step: Data encodation
	encoded, e := encoder.EncodeHighLevel(contents, opts)
	if e != nil {
		return nil, e
	}

	symbolInfo, _ := encoder.SymbolInfo_Lookup(len(encoded), opts, true)

	//2. step: ECC generation
	codewords, _ := encoder.ErrorCorrection_EncodeECC200(encoded, symbolInfo)

	//3. step: Module placement in Matrix
	placement := encoder.NewDefaultPlacement(codewords,
		symbolInfo.GetSymbolDataWidth(), symbolInfo.GetSymbolDataHeight())
	placement.Place()

	//4. step: low-level encoding
	return encodeLowLevel(placement, symbolInfo, width, height), nil
}

// encodeLowLevel Encode the given symbol info to a bit matrix.
func encodeLowLevel(placement *encoder.DefaultPlacement, symbolInfo *encoder.SymbolInfo, width, height int) *BitMatrix {

	symbolWidth := symbolInfo.GetSymbolDataWidth()
	symbolHeight := symbolInfo.GetSymbolDataHeight()

	matrix := encoder.NewByteMatrix(symbolInfo.GetSymbolWidth(), symbolInfo.GetSymbolHeight())

	matrixY := 0

	for y := 0; y < symbolHeight; y++ {
		// Fill the top edge with alternate 0 / 1
		var matrixX int
		if (y % symbolInfo.GetMatrixHeight()) == 0 {
			matrixX = 0
			for x := 0; x < symbolInfo.GetSymbolWidth(); x++ {
				matrix.SetBool(matrixX, matrixY, (x%2) == 0)
				matrixX++
			}
			matrixY++
		}
		matrixX = 0
		for x := 0; x < symbolWidth; x++ {
			// Fill the right edge with full 1
			if (x % symbolInfo.GetMatrixWidth()) == 0 {
				matrix.SetBool(matrixX, matrixY, true)
				matrixX++
			}
			matrix.SetBool(matrixX, matrixY, placement.GetBit(x, y))
			matrixX++
			// Fill the right edge with alternate 0 / 1
			if (x % symbolInfo.GetMatrixWidth()) == symbolInfo.GetMatrixWidth()-1 {
				matrix.SetBool(matrixX, matrixY, (y%2) == 0)
				matrixX++
			}
		}
		matrixY++
		// Fill the bottom edge with full 1
		if (y % symbolInfo.GetMatrixHeight()) == symbolInfo.GetMatrixHeight()-1 {
			matrixX = 0
			for x := 0; x < symbolInfo.GetSymbolWidth(); x++ {
				matrix.SetBool(matrixX, matrixY, true)
				matrixX++
			}
			matrixY++
		}
	}

	return convertByteMatrixToBitMatrix(matrix, width, height)
}

// convertByteMatrixToBitMatrix Convert the ByteMatrix to BitMatrix.
func convertByteMatrixToBitMatrix(matrix *encoder.ByteMatrix, reqWidth, reqHeight int) *BitMatrix {
	matrixWidth := matrix.GetWidth()
	matrixHeight := matrix.GetHeight()
	outputWidth := reqWidth
	if outputWidth < matrixWidth {
		outputWidth = matrixWidth
	}
	outputHeight := reqHeight
	if outputHeight < matrixHeight {
		outputHeight = matrixHeight
	}

	multiple := outputWidth / matrixWidth
	if mh := outputHeight / matrixHeight; mh < multiple {
		multiple = mh
	}

	leftPadding := (outputWidth - (matrixWidth * multiple)) / 2
	topPadding := (outputHeight - (matrixHeight * multiple)) / 2

	var output *BitMatrix

	// remove padding if requested width and height are too small
	if reqHeight < matrixHeight || reqWidth < matrixWidth {
		leftPadding = 0
		topPadding = 0
		output, _ = NewBitMatrix(matrixWidth, matrixHeight)
	} else {
		output, _ = NewBitMatrix(reqWidth, reqHeight)
	}

	output.Clear()
	for inputY, outputY := 0, topPadding; inputY < matrixHeight; inputY, outputY = inputY+1, outputY+multiple {
		// Write the contents of this row of the bytematrix
		for inputX, outputX := 0, leftPadding; inputX < matrixWidth; inputX, outputX = inputX+1, outputX+multiple {
			if matrix.Get(inputX, inputY) == 1 {
				output.SetRegion(outputX, outputY, multiple, multiple)
			}
		}
	}

	return output
}
