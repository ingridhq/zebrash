package pdf417

import (
	"fmt"
	"image"

	"github.com/ingridhq/zebrash/barcodes/utils"
	"github.com/ingridhq/zebrash/images"
)

const (
	padding_codeword = 900
)

// Encodes the given data as PDF417 barcode.
// securityLevel should be between 0 and 8. The higher the number, the more
// additional error-correction codes are added.
func Encode(data string, securityLevel byte, targetRowHeight, targetColumns int) (image.Image, error) {
	if securityLevel >= 9 {
		return nil, fmt.Errorf("invalid security level %d", securityLevel)
	}

	sl := securitylevel(securityLevel)

	dataWords, err := highlevelEncode(data)
	if err != nil {
		return nil, err
	}

	columns, rows := calcDimensions(targetColumns, len(dataWords), sl.ErrorCorrectionWordCount())
	if columns < minCols || columns > maxCols || rows < minRows || rows > maxRows {
		return nil, fmt.Errorf("unable to fit data in barcode")
	}

	codeWords, err := encodeData(dataWords, columns, sl)
	if err != nil {
		return nil, err
	}

	grid := [][]int{}
	for i := 0; i < len(codeWords); i += columns {
		grid = append(grid, codeWords[i:min(i+columns, len(codeWords))])
	}

	codes := [][]int{}

	for rowNum, row := range grid {
		table := rowNum % 3
		rowCodes := make([]int, 0, columns+4)

		rowCodes = append(rowCodes, start_word)
		rowCodes = append(rowCodes, getCodeword(table, getLeftCodeWord(rowNum, rows, columns, securityLevel)))

		for _, word := range row {
			rowCodes = append(rowCodes, getCodeword(table, word))
		}

		rowCodes = append(rowCodes, getCodeword(table, getRightCodeWord(rowNum, rows, columns, securityLevel)))
		rowCodes = append(rowCodes, stop_word)

		codes = append(codes, rowCodes)
	}

	codeData := renderBarcode(codes)
	width := (columns+4)*17 + 1
	height := codeData.Len() / width
	scaleY := 4

	if targetRowHeight > 0 {
		scaleY = max(1, rows*targetRowHeight/height)
	}

	barcode := &pdfBarcode{
		width:  (columns+4)*17 + 1,
		height: height,
		code:   codeData,
	}

	return images.NewScaled(barcode, 2, scaleY), nil
}

func encodeData(dataWords []int, columns int, sl securitylevel) ([]int, error) {
	dataCount := len(dataWords)

	ecCount := sl.ErrorCorrectionWordCount()

	padWords := getPadding(dataCount, ecCount, columns)
	dataWords = append(dataWords, padWords...)

	length := len(dataWords) + 1
	dataWords = append([]int{length}, dataWords...)

	ecWords := sl.Compute(dataWords)

	return append(dataWords, ecWords...), nil
}

func getLeftCodeWord(rowNum int, rows int, columns int, securityLevel byte) int {
	tableId := rowNum % 3

	var x int

	switch tableId {
	case 0:
		x = (rows - 1) / 3
	case 1:
		x = int(securityLevel) * 3
		x += (rows - 1) % 3
	case 2:
		x = columns - 1
	}

	return 30*(rowNum/3) + x
}

func getRightCodeWord(rowNum int, rows int, columns int, securityLevel byte) int {
	tableId := rowNum % 3

	var x int

	switch tableId {
	case 0:
		x = columns - 1
	case 1:
		x = (rows - 1) / 3
	case 2:
		x = int(securityLevel) * 3
		x += (rows - 1) % 3
	}

	return 30*(rowNum/3) + x
}

func getPadding(dataCount int, ecCount int, columns int) []int {
	totalCount := dataCount + ecCount + 1
	mod := totalCount % columns

	padding := []int{}

	if mod > 0 {
		padCount := columns - mod
		padding = make([]int, padCount)
		for i := 0; i < padCount; i++ {
			padding[i] = padding_codeword
		}
	}

	return padding
}

func renderBarcode(codes [][]int) *utils.BitList {
	bl := new(utils.BitList)
	for _, row := range codes {
		lastIdx := len(row) - 1
		for i, col := range row {
			if i == lastIdx {
				bl.AddBits(col, 18)
			} else {
				bl.AddBits(col, 17)
			}
		}
	}
	return bl
}
