// Package aztec can create Aztec Code barcodes
package aztec

import (
	"fmt"
	"image"
	"math"

	"github.com/DawidBury/zebrash/internal/barcodes/utils"
	"github.com/DawidBury/zebrash/internal/images"
)

const (
	DEFAULT_EC_PERCENT  = 23 // ZPL Default
	DEFAULT_LAYERS      = 0
	max_nb_bits         = 32
	max_nb_bits_compact = 4
)

var (
	word_size = []int{
		4, 6, 6, 8, 8, 8, 8, 8, 8, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
		12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	}
)

func totalBitsInLayer(layers int, compact bool) int {
	tmp := 112
	if compact {
		tmp = 88
	}
	return (tmp + 16*layers) * layers
}

func stuffBits(bits *utils.BitList, wordSize int) *utils.BitList {
	out := new(utils.BitList)
	n := bits.Len()
	mask := (1 << uint(wordSize)) - 2
	for i := 0; i < n; i += wordSize {
		word := 0
		for j := 0; j < wordSize; j++ {
			if i+j >= n || bits.GetBit(i+j) {
				word |= 1 << uint(wordSize-1-j)
			}
		}

		switch word & mask {
		case mask:
			out.AddBits(word&mask, byte(wordSize))
			i--
		case 0:
			out.AddBits(word|1, byte(wordSize))
			i--
		default:
			out.AddBits(word, byte(wordSize))
		}
	}
	return out
}

func generateModeMessage(compact bool, layers, messageSizeInWords int) *utils.BitList {
	modeMessage := new(utils.BitList)
	if compact {
		modeMessage.AddBits(layers-1, 2)
		modeMessage.AddBits(messageSizeInWords-1, 6)
		modeMessage = generateCheckWords(modeMessage, 28, 4)
	} else {
		modeMessage.AddBits(layers-1, 5)
		modeMessage.AddBits(messageSizeInWords-1, 11)
		modeMessage = generateCheckWords(modeMessage, 40, 4)
	}
	return modeMessage
}

func drawModeMessage(matrix *aztecCode, compact bool, matrixSize int, modeMessage *utils.BitList) {
	center := matrixSize / 2
	if compact {
		for i := 0; i < 7; i++ {
			offset := center - 3 + i
			if modeMessage.GetBit(i) {
				matrix.set(offset, center-5)
			}
			if modeMessage.GetBit(i + 7) {
				matrix.set(center+5, offset)
			}
			if modeMessage.GetBit(20 - i) {
				matrix.set(offset, center+5)
			}
			if modeMessage.GetBit(27 - i) {
				matrix.set(center-5, offset)
			}
		}
	} else {
		for i := 0; i < 10; i++ {
			offset := center - 5 + i + i/5
			if modeMessage.GetBit(i) {
				matrix.set(offset, center-7)
			}
			if modeMessage.GetBit(i + 10) {
				matrix.set(center+7, offset)
			}
			if modeMessage.GetBit(29 - i) {
				matrix.set(offset, center+7)
			}
			if modeMessage.GetBit(39 - i) {
				matrix.set(center-7, offset)
			}
		}
	}
}

func drawBullsEye(matrix *aztecCode, center, size int) {
	for i := 0; i < size; i += 2 {
		for j := center - i; j <= center+i; j++ {
			matrix.set(j, center-i)
			matrix.set(j, center+i)
			matrix.set(center-i, j)
			matrix.set(center+i, j)
		}
	}
	matrix.set(center-size, center-size)
	matrix.set(center-size+1, center-size)
	matrix.set(center-size, center-size+1)
	matrix.set(center+size, center-size)
	matrix.set(center+size, center-size+1)
	matrix.set(center+size, center-size-1)
}

// Encode returns an aztec barcode with the given content
func Encode(data []byte, minECCPercent, userSpecifiedLayers, magnification int) (image.Image, error) {
	bits := highlevelEncode(data)
	var layers, TotalBitsInLayer, wordSize int
	var compact bool
	var stuffedBits *utils.BitList

	if userSpecifiedLayers != DEFAULT_LAYERS {
		// This branch handles fixed layers (compact/full-range modes).
		compact = userSpecifiedLayers < 0
		if compact {
			layers = -userSpecifiedLayers
		} else {
			layers = userSpecifiedLayers
		}
		if (compact && layers > max_nb_bits_compact) || (!compact && layers > max_nb_bits) {
			return nil, fmt.Errorf("illegal value %d for layers", userSpecifiedLayers)
		}
		TotalBitsInLayer = totalBitsInLayer(layers, compact)
		wordSize = word_size[layers]
		usableBitsInLayers := TotalBitsInLayer - (TotalBitsInLayer % wordSize)
		stuffedBits = stuffBits(bits, wordSize)

		if stuffedBits.Len()*4/3 > usableBitsInLayers {
			return nil, fmt.Errorf("data too large for user specified layer")
		}
		if compact && stuffedBits.Len() > wordSize*64 {
			return nil, fmt.Errorf("data too large for user specified layer")
		}
	} else {
		// This branch handles dynamic sizing based on the ZPL ECC formula.
		wordSize = 0
		stuffedBits = nil
		for i := 0; ; i++ {
			if i > max_nb_bits {
				return nil, fmt.Errorf("data too large for an aztec code")
			}
			compact = i <= 3
			layers = i
			if compact {
				layers = i + 1
			}
			currentWordSize := word_size[layers]

			if wordSize != currentWordSize {
				wordSize = currentWordSize
				stuffedBits = stuffBits(bits, wordSize)
			}

			if compact && (stuffedBits.Len()/wordSize) > 64 {
				continue
			}

			TotalBitsInLayer = totalBitsInLayer(layers, compact)
			usableBitsInLayers := TotalBitsInLayer - (TotalBitsInLayer % wordSize)
			totalSymbolWords := usableBitsInLayers / wordSize

			if totalSymbolWords == 0 {
				continue
			}

			requiredDataWords := (stuffedBits.Len() + wordSize - 1) / wordSize

			// ZPL Logic: Find smallest symbol where data fits into the space *not* reserved for ECC.
			// ECC is calculated based on the *total capacity* of the candidate symbol, using floating point math.
			eccWordsToReserve := int(math.Ceil(float64(totalSymbolWords*minECCPercent)/100.0)) + 3
			availableDataWords := totalSymbolWords - eccWordsToReserve

			if requiredDataWords <= availableDataWords {
				// We found the smallest symbol that satisfies the condition.
				break
			}
		}
	}

	messageBits := generateCheckWords(stuffedBits, TotalBitsInLayer, wordSize)
	messageSizeInWords := stuffedBits.Len() / wordSize
	modeMessage := generateModeMessage(compact, layers, messageSizeInWords)

	var baseMatrixSize int
	if compact {
		baseMatrixSize = 11 + layers*4
	} else {
		baseMatrixSize = 14 + layers*4
	}
	alignmentMap := make([]int, baseMatrixSize)
	var matrixSize int

	if compact {
		matrixSize = baseMatrixSize
		for i := 0; i < len(alignmentMap); i++ {
			alignmentMap[i] = i
		}
	} else {
		matrixSize = baseMatrixSize + 1 + 2*((baseMatrixSize/2-1)/15)
		origCenter := baseMatrixSize / 2
		center := matrixSize / 2
		for i := 0; i < origCenter; i++ {
			newOffset := i + i/15
			alignmentMap[origCenter-i-1] = center - newOffset - 1
			alignmentMap[origCenter+i] = center + newOffset + 1
		}
	}
	code := newAztecCode(matrixSize)
	code.content = data

	// draw data bits
	for i, rowOffset := 0, 0; i < layers; i++ {
		rowSize := (layers - i) * 4
		if compact {
			rowSize += 9
		} else {
			rowSize += 12
		}

		for j := 0; j < rowSize; j++ {
			columnOffset := j * 2
			for k := 0; k < 2; k++ {
				if messageBits.GetBit(rowOffset + columnOffset + k) {
					code.set(alignmentMap[i*2+k], alignmentMap[i*2+j])
				}
				if messageBits.GetBit(rowOffset + rowSize*2 + columnOffset + k) {
					code.set(alignmentMap[i*2+j], alignmentMap[baseMatrixSize-1-i*2-k])
				}
				if messageBits.GetBit(rowOffset + rowSize*4 + columnOffset + k) {
					code.set(alignmentMap[baseMatrixSize-1-i*2-k], alignmentMap[baseMatrixSize-1-i*2-j])
				}
				if messageBits.GetBit(rowOffset + rowSize*6 + columnOffset + k) {
					code.set(alignmentMap[baseMatrixSize-1-i*2-j], alignmentMap[i*2+k])
				}
			}
		}
		rowOffset += rowSize * 8
	}

	drawModeMessage(code, compact, matrixSize, modeMessage)

	if compact {
		drawBullsEye(code, matrixSize/2, 5)
	} else {
		drawBullsEye(code, matrixSize/2, 7)
		for i, j := 0, 0; i < baseMatrixSize/2-1; i, j = i+15, j+16 {
			for k := (matrixSize / 2) & 1; k < matrixSize; k += 2 {
				code.set(matrixSize/2-j, k)
				code.set(matrixSize/2+j, k)
				code.set(k, matrixSize/2-j)
				code.set(k, matrixSize/2+j)
			}
		}
	}
	return images.NewScaled(code, magnification, magnification), nil
}
