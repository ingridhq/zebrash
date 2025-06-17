package code39

import (
	"fmt"
	"image"
	"strings"

	"github.com/ingridhq/zebrash/barcodes/utils"
)

const code39AlphabetString = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ-. $/+%"

// code39CharacterEncodings These represent the encodings of characters, as patterns of wide and narrow bars.
// The 9 least-significant bits of each int correspond to the pattern of wide and narrow,
// with 1s representing "wide" and 0s representing narrow.
var code39CharacterEncodings = []int{
	0x034, 0x121, 0x061, 0x160, 0x031, 0x130, 0x070, 0x025, 0x124, 0x064, // 0-9
	0x109, 0x049, 0x148, 0x019, 0x118, 0x058, 0x00D, 0x10C, 0x04C, 0x01C, // A-J
	0x103, 0x043, 0x142, 0x013, 0x112, 0x052, 0x007, 0x106, 0x046, 0x016, // K-T
	0x181, 0x0C1, 0x1C0, 0x091, 0x190, 0x0D0, 0x085, 0x184, 0x0C4, 0x0A8, // U-$
	0x0A2, 0x08A, 0x02A, // /-%
}

const code39AsteriskEncoding = 0x094

func Encode(contents string, width, height int, widthRatio float64) (image.Image, error) {
	length := len(contents)
	if length > 80 {
		return nil, fmt.Errorf("requested contents should be less than 80 digits long, but got %v", length)
	}

	for i := range length {
		indexInString := strings.Index(code39AlphabetString, string(contents[i]))
		if indexInString < 0 {
			var e error
			contents, e = code39TryToConvertToExtendedMode(contents)
			if e != nil {
				return nil, e
			}
			length = len(contents)
			if length > 80 {
				return nil, fmt.Errorf("requested contents should be less than 80 digits long, but got %v (extended full ASCII mode)", length)
			}

			break
		}
	}

	widths := make([]int, 9)
	codeWidth := 24 + 1 + (13 * length)
	result := utils.NewBitList(codeWidth)

	code39ToIntArray(code39AsteriskEncoding, widths)
	appendPattern(result, widths, true)
	narrowWhite := []int{1}
	appendPattern(result, narrowWhite, false)

	// append next character to byte matrix
	for i := range length {
		indexInString := strings.Index(code39AlphabetString, string(contents[i]))
		code39ToIntArray(code39CharacterEncodings[indexInString], widths)
		appendPattern(result, widths, true)
		appendPattern(result, narrowWhite, false)
	}

	code39ToIntArray(code39AsteriskEncoding, widths)
	appendPattern(result, widths, true)

	return result.ToImage(width, height, widthRatio), nil
}

func appendPattern(target *utils.BitList, pattern []int, startColor bool) {
	color := startColor
	for _, len := range pattern {
		for range len {
			target.AddBit(color)
		}

		color = !color
	}
}

func code39ToIntArray(a int, toReturn []int) {
	for i := range 9 {
		temp := a & (1 << uint(8-i))
		if temp == 0 {
			toReturn[i] = 1
		} else {
			toReturn[i] = 2
		}
	}
}

func code39TryToConvertToExtendedMode(contents string) (string, error) {
	length := len(contents)
	extendedContent := make([]byte, 0)
	for i := range length {
		character := contents[i]
		switch character {
		case '\u0000':
			extendedContent = append(extendedContent, []byte("%U")...)
		case ' ', '-', '.':
			extendedContent = append(extendedContent, character)
		case '@':
			extendedContent = append(extendedContent, []byte("%V")...)
		case '`':
			extendedContent = append(extendedContent, []byte("%W")...)
		default:
			if character <= 26 {
				extendedContent = append(extendedContent, '$')
				extendedContent = append(extendedContent, 'A'+(character-1))
			} else if character < ' ' {
				extendedContent = append(extendedContent, '%')
				extendedContent = append(extendedContent, 'A'+(character-27))
			} else if character <= ',' || character == '/' || character == ':' {
				extendedContent = append(extendedContent, '/')
				extendedContent = append(extendedContent, 'A'+(character-33))
			} else if character <= '9' {
				extendedContent = append(extendedContent, '0'+(character-48))
			} else if character <= '?' {
				extendedContent = append(extendedContent, '%')
				extendedContent = append(extendedContent, 'F'+(character-59))
			} else if character <= 'Z' {
				extendedContent = append(extendedContent, 'A'+(character-65))
			} else if character <= '_' {
				extendedContent = append(extendedContent, '%')
				extendedContent = append(extendedContent, 'K'+(character-91))
			} else if character <= 'z' {
				extendedContent = append(extendedContent, '+')
				extendedContent = append(extendedContent, 'A'+(character-97))
			} else if character <= 127 {
				extendedContent = append(extendedContent, '%')
				extendedContent = append(extendedContent, 'P'+(character-123))
			} else {
				return string(extendedContent), fmt.Errorf("requested content contains a non-encodable character: '%v'", contents[i])
			}
		}
	}

	return string(extendedContent), nil
}
