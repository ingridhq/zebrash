// Package twooffive can create interleaved and standard "2 of 5" barcodes.
package twooffive

import (
	"fmt"
	"image"

	"github.com/ingridhq/zebrash/barcodes/utils"
)

const patternWidth = 5

type pattern [patternWidth]bool

var (
	encodingTable = map[rune]pattern{
		'0': {false, false, true, true, false},
		'1': {true, false, false, false, true},
		'2': {false, true, false, false, true},
		'3': {true, true, false, false, false},
		'4': {false, false, true, false, true},
		'5': {true, false, true, false, false},
		'6': {false, true, true, false, false},
		'7': {false, false, false, true, true},
		'8': {true, false, false, true, false},
		'9': {false, true, false, true, false},
	}

	startPattern = []bool{true, false, true, false}
	endPattern   = []bool{true, true, true, false, true}

	widths = map[bool]int{
		true:  3,
		false: 1,
	}
)

func EncodeInterleaved(content string, width, height int, widthRatio float64, checkDigit bool) (image.Image, string, error) {
	var err error

	if checkDigit {
		content, err = addCheckDigit(content)
		if err != nil {
			return nil, "", fmt.Errorf("failed to add check digit: %w", err)
		}
	}

	// Can't encode values of odd length so they must be prepended with 0
	if len(content)%2 == 1 {
		content = "0" + content
	}

	resBits := new(utils.BitList)
	resBits.AddBit(startPattern...)

	var lastRune *rune
	for _, r := range content {
		var a, b pattern
		if lastRune == nil {
			lastRune = new(rune)
			*lastRune = r
			continue
		} else {
			var o1, o2 bool
			a, o1 = encodingTable[*lastRune]
			b, o2 = encodingTable[r]
			if !o1 || !o2 {
				return nil, "", fmt.Errorf("can not encode %q", content)
			}
			lastRune = nil
		}

		for i := range patternWidth {
			for range widths[a[i]] {
				resBits.AddBit(true)
			}
			for range widths[b[i]] {
				resBits.AddBit(false)
			}
		}
	}

	resBits.AddBit(endPattern...)

	return resBits.ToImage(width, height, widthRatio), content, nil
}

func addCheckDigit(content string) (string, error) {
	even := true
	sum := 0
	for _, r := range content {
		if _, ok := encodingTable[r]; !ok {
			return "", fmt.Errorf("can not encode %q", content)
		}

		value := utils.RuneToInt(r)
		if even {
			sum += value * 3
		} else {
			sum += value
		}

		even = !even
	}

	sum = sum % 10
	if sum > 0 {
		sum = 10 - sum
	}

	return content + string(utils.IntToRune(sum)), nil
}
