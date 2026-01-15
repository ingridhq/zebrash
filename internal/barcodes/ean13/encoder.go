package ean13

import (
	"fmt"
	"image"
	"strings"
	"unicode"

	"github.com/ingridhq/zebrash/internal/barcodes/utils"
)

type encodedNumber struct {
	LeftOdd  []bool
	LeftEven []bool
	Right    []bool
	CheckSum []bool
}

var encoderTable = map[rune]encodedNumber{
	'0': {
		[]bool{false, false, false, true, true, false, true},
		[]bool{false, true, false, false, true, true, true},
		[]bool{true, true, true, false, false, true, false},
		[]bool{false, false, false, false, false, false},
	},
	'1': {
		[]bool{false, false, true, true, false, false, true},
		[]bool{false, true, true, false, false, true, true},
		[]bool{true, true, false, false, true, true, false},
		[]bool{false, false, true, false, true, true},
	},
	'2': {
		[]bool{false, false, true, false, false, true, true},
		[]bool{false, false, true, true, false, true, true},
		[]bool{true, true, false, true, true, false, false},
		[]bool{false, false, true, true, false, true},
	},
	'3': {
		[]bool{false, true, true, true, true, false, true},
		[]bool{false, true, false, false, false, false, true},
		[]bool{true, false, false, false, false, true, false},
		[]bool{false, false, true, true, true, false},
	},
	'4': {
		[]bool{false, true, false, false, false, true, true},
		[]bool{false, false, true, true, true, false, true},
		[]bool{true, false, true, true, true, false, false},
		[]bool{false, true, false, false, true, true},
	},
	'5': {
		[]bool{false, true, true, false, false, false, true},
		[]bool{false, true, true, true, false, false, true},
		[]bool{true, false, false, true, true, true, false},
		[]bool{false, true, true, false, false, true},
	},
	'6': {
		[]bool{false, true, false, true, true, true, true},
		[]bool{false, false, false, false, true, false, true},
		[]bool{true, false, true, false, false, false, false},
		[]bool{false, true, true, true, false, false},
	},
	'7': {
		[]bool{false, true, true, true, false, true, true},
		[]bool{false, false, true, false, false, false, true},
		[]bool{true, false, false, false, true, false, false},
		[]bool{false, true, false, true, false, true},
	},
	'8': {
		[]bool{false, true, true, false, true, true, true},
		[]bool{false, false, false, true, false, false, true},
		[]bool{true, false, false, true, false, false, false},
		[]bool{false, true, false, true, true, false},
	},
	'9': {
		[]bool{false, false, false, true, false, true, true},
		[]bool{false, false, true, false, true, true, true},
		[]bool{true, true, true, false, true, false, false},
		[]bool{false, true, true, false, true, false},
	},
}

func calcCheckNum(code string) rune {
	x3 := len(code) == 7
	sum := 0
	for _, r := range code {
		curNum := utils.RuneToInt(r)
		if curNum < 0 || curNum > 9 {
			return 'B'
		}
		if x3 {
			curNum = curNum * 3
		}
		x3 = !x3
		sum += curNum
	}

	return utils.IntToRune((10 - (sum % 10)) % 10)
}

func sanitizeContent(content string) string {
	content = strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return '0'
	}, content)

	content = fmt.Sprintf("%012s", content)

	if len(content) > 12 {
		content = content[0:1] + content[len(content)-11:]
	}

	return content + string(calcCheckNum(content))
}

func Encode(content string, height, barWidth int) (image.Image, string, error) {
	content = sanitizeContent(content)
	code := encodeEAN13(content)
	return newEan13(code, height, barWidth), content, nil
}

func encodeEAN13(code string) *utils.BitList {
	result := new(utils.BitList)
	result.AddBit(true, false, true)

	var firstNum []bool
	for cpos, r := range code {
		num, ok := encoderTable[r]
		if !ok {
			return nil
		}
		if cpos == 0 {
			firstNum = num.CheckSum
			continue
		}

		var data []bool
		if cpos < 7 { // Left
			if firstNum[cpos-1] {
				data = num.LeftEven
			} else {
				data = num.LeftOdd
			}
		} else {
			data = num.Right
		}

		if cpos == 7 {
			result.AddBit(false, true, false, true, false)
		}
		result.AddBit(data...)
	}
	result.AddBit(true, false, true)
	return result
}
