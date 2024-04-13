package code128

import (
	"fmt"
	"image"
	"strings"
)

var CODE_PATTERNS = [][]int{
	{2, 1, 2, 2, 2, 2}, // 0
	{2, 2, 2, 1, 2, 2},
	{2, 2, 2, 2, 2, 1},
	{1, 2, 1, 2, 2, 3},
	{1, 2, 1, 3, 2, 2},
	{1, 3, 1, 2, 2, 2}, // 5
	{1, 2, 2, 2, 1, 3},
	{1, 2, 2, 3, 1, 2},
	{1, 3, 2, 2, 1, 2},
	{2, 2, 1, 2, 1, 3},
	{2, 2, 1, 3, 1, 2}, // 10
	{2, 3, 1, 2, 1, 2},
	{1, 1, 2, 2, 3, 2},
	{1, 2, 2, 1, 3, 2},
	{1, 2, 2, 2, 3, 1},
	{1, 1, 3, 2, 2, 2}, // 15
	{1, 2, 3, 1, 2, 2},
	{1, 2, 3, 2, 2, 1},
	{2, 2, 3, 2, 1, 1},
	{2, 2, 1, 1, 3, 2},
	{2, 2, 1, 2, 3, 1}, // 20
	{2, 1, 3, 2, 1, 2},
	{2, 2, 3, 1, 1, 2},
	{3, 1, 2, 1, 3, 1},
	{3, 1, 1, 2, 2, 2},
	{3, 2, 1, 1, 2, 2}, // 25
	{3, 2, 1, 2, 2, 1},
	{3, 1, 2, 2, 1, 2},
	{3, 2, 2, 1, 1, 2},
	{3, 2, 2, 2, 1, 1},
	{2, 1, 2, 1, 2, 3}, // 30
	{2, 1, 2, 3, 2, 1},
	{2, 3, 2, 1, 2, 1},
	{1, 1, 1, 3, 2, 3},
	{1, 3, 1, 1, 2, 3},
	{1, 3, 1, 3, 2, 1}, // 35
	{1, 1, 2, 3, 1, 3},
	{1, 3, 2, 1, 1, 3},
	{1, 3, 2, 3, 1, 1},
	{2, 1, 1, 3, 1, 3},
	{2, 3, 1, 1, 1, 3}, // 40
	{2, 3, 1, 3, 1, 1},
	{1, 1, 2, 1, 3, 3},
	{1, 1, 2, 3, 3, 1},
	{1, 3, 2, 1, 3, 1},
	{1, 1, 3, 1, 2, 3}, // 45
	{1, 1, 3, 3, 2, 1},
	{1, 3, 3, 1, 2, 1},
	{3, 1, 3, 1, 2, 1},
	{2, 1, 1, 3, 3, 1},
	{2, 3, 1, 1, 3, 1}, // 50
	{2, 1, 3, 1, 1, 3},
	{2, 1, 3, 3, 1, 1},
	{2, 1, 3, 1, 3, 1},
	{3, 1, 1, 1, 2, 3},
	{3, 1, 1, 3, 2, 1}, // 55
	{3, 3, 1, 1, 2, 1},
	{3, 1, 2, 1, 1, 3},
	{3, 1, 2, 3, 1, 1},
	{3, 3, 2, 1, 1, 1},
	{3, 1, 4, 1, 1, 1}, // 60
	{2, 2, 1, 4, 1, 1},
	{4, 3, 1, 1, 1, 1},
	{1, 1, 1, 2, 2, 4},
	{1, 1, 1, 4, 2, 2},
	{1, 2, 1, 1, 2, 4}, // 65
	{1, 2, 1, 4, 2, 1},
	{1, 4, 1, 1, 2, 2},
	{1, 4, 1, 2, 2, 1},
	{1, 1, 2, 2, 1, 4},
	{1, 1, 2, 4, 1, 2}, // 70
	{1, 2, 2, 1, 1, 4},
	{1, 2, 2, 4, 1, 1},
	{1, 4, 2, 1, 1, 2},
	{1, 4, 2, 2, 1, 1},
	{2, 4, 1, 2, 1, 1}, // 75
	{2, 2, 1, 1, 1, 4},
	{4, 1, 3, 1, 1, 1},
	{2, 4, 1, 1, 1, 2},
	{1, 3, 4, 1, 1, 1},
	{1, 1, 1, 2, 4, 2}, // 80
	{1, 2, 1, 1, 4, 2},
	{1, 2, 1, 2, 4, 1},
	{1, 1, 4, 2, 1, 2},
	{1, 2, 4, 1, 1, 2},
	{1, 2, 4, 2, 1, 1}, // 85
	{4, 1, 1, 2, 1, 2},
	{4, 2, 1, 1, 1, 2},
	{4, 2, 1, 2, 1, 1},
	{2, 1, 2, 1, 4, 1},
	{2, 1, 4, 1, 2, 1}, // 90
	{4, 1, 2, 1, 2, 1},
	{1, 1, 1, 1, 4, 3},
	{1, 1, 1, 3, 4, 1},
	{1, 3, 1, 1, 4, 1},
	{1, 1, 4, 1, 1, 3}, // 95
	{1, 1, 4, 3, 1, 1},
	{4, 1, 1, 1, 1, 3},
	{4, 1, 1, 3, 1, 1},
	{1, 1, 3, 1, 4, 1},
	{1, 1, 4, 1, 3, 1}, // 100
	{3, 1, 1, 1, 4, 1},
	{4, 1, 1, 1, 3, 1},
	{2, 1, 1, 4, 1, 2},
	{2, 1, 1, 2, 1, 4},
	{2, 1, 1, 2, 3, 2}, // 105
	{2, 3, 3, 1, 1, 1, 2},
}

const (
	CODE_C = 99
	CODE_B = 100
	CODE_A = 101

	FNC_1   = 102
	FNC_2   = 97
	FNC_3   = 96
	FNC_4_A = 101
	FNC_4_B = 100

	START_A = 103
	START_B = 104
	START_C = 105
	STOP    = 106

	DEL_B        = 95
	CIRCUMFLEX   = 62
	GREATER_THAN = 30
	TILDE_B      = 94
)

const (
	// Dummy characters used to specify control characters in input
	ESCAPE_FNC_1 = '\u00f1'
	ESCAPE_FNC_2 = '\u00f2'
	ESCAPE_FNC_3 = '\u00f3'
	ESCAPE_FNC_4 = '\u00f4'
)

// Results of minimal lookahead for code C
type code128CType int

const (
	code128CType_UNCODABLE code128CType = iota
	code128CType_ONE_DIGIT
	code128CType_TWO_DIGITS
	code128CType_FNC_1
)

func EncodeNoMode(content string, height, barWidth int) (image.Image, string, error) {
	var humanReadable strings.Builder

	var patternsIdx []byte

	// Force set B by default if no invocation codes was found
	currSet := CODE_B

	patternsIdx = append(patternsIdx, START_B)

	for i := 1; i < len(content); i++ {
		if content[i-1] != '>' {
			// Code 128 subsets A and C are programmed as pairs of digits, 00-99, in the field data string.
			// In subset A, each pair of digits results in a single character being encoded in the barcode
			// in subset C, they are printed as entered.
			// Non-integers programmed as the ﬁrst character of a digit pair are ignored.
			// However, non-integers programmed as the second character of a digit pair invalidate the entire digit pair, and the pair is ignored
			if currSet == CODE_A || currSet == CODE_C {
				digit1 := content[i] - '0'

				// Skip entire pair
				if digit1 >= 10 {
					i++
					continue
				}

				digit0 := content[i-1] - '0'
				if digit0 >= 10 {
					digit0 = 0
				}

				// Index of character in subset A/C
				patternIdx := digit0*10 + digit1
				patternsIdx = append(patternsIdx, patternIdx)

				switch {
				case currSet == CODE_C:
					humanReadable.WriteString(fmt.Sprintf("%02d", patternIdx))
				case patternIdx < 64:
					humanReadable.WriteByte(patternIdx + ' ')
				case patternIdx < 96:
					humanReadable.WriteByte(patternIdx - 64)
				}

				i++
				continue
			}

			// Handle subset B which does not operate on pairs
			patternsIdx = append(patternsIdx, content[i-1]-' ')
			humanReadable.WriteByte(content[i-1])

			// If last iteration add remaining character
			if i == len(content)-1 {
				patternsIdx = append(patternsIdx, content[i]-' ')
				humanReadable.WriteByte(content[i])
			}

			continue
		}

		switch content[i] {
		// Escape characters
		case '<':
			if currSet != CODE_C {
				patternsIdx = append(patternsIdx, CIRCUMFLEX)
				humanReadable.WriteByte(CIRCUMFLEX + ' ')
			}
		case '0':
			if currSet != CODE_C {
				patternsIdx = append(patternsIdx, GREATER_THAN)
				humanReadable.WriteByte(GREATER_THAN + ' ')
			}
		case '=':
			if currSet != CODE_C {
				patternsIdx = append(patternsIdx, TILDE_B)

				if currSet == CODE_B {
					humanReadable.WriteByte(TILDE_B + ' ')
				}
			}
		case '1':
			if currSet != CODE_C {
				patternsIdx = append(patternsIdx, DEL_B)
			}
		// Special functions
		case '8':
			patternsIdx = append(patternsIdx, FNC_1)
		case '2':
			if currSet != CODE_C {
				patternsIdx = append(patternsIdx, FNC_3)
			}
		case '3':
			if currSet != CODE_C {
				patternsIdx = append(patternsIdx, FNC_2)
			}
		// Start characters
		case '9':
			if i == 1 {
				currSet = CODE_A
				patternsIdx[0] = START_A
			}
		case ':':
			if i == 1 {
				currSet = CODE_B
				patternsIdx[0] = START_B
			}
		case ';':
			if i == 1 {
				currSet = CODE_C
				patternsIdx[0] = START_C
			}
		// Change set invocations
		case '7':
			switch currSet {
			case CODE_A:
				patternsIdx = append(patternsIdx, FNC_4_A)
			default:
				currSet = CODE_A
				patternsIdx = append(patternsIdx, byte(currSet))
			}
		case '6':
			switch currSet {
			case CODE_B:
				patternsIdx = append(patternsIdx, FNC_4_B)
			default:
				currSet = CODE_B
				patternsIdx = append(patternsIdx, byte(currSet))
			}
		case '5':
			if currSet != CODE_C {
				patternsIdx = append(patternsIdx, byte(currSet))
			}
		}

		i++
	}

	img, err := encode(patternsIdx, height, barWidth)
	if err != nil {
		return nil, "", err
	}

	return img, humanReadable.String(), nil
}

func EncodeAuto(content string, height, barWidth int) (image.Image, error) {
	contents := []rune(content)
	length := len(contents)
	// Check length
	if length < 1 || length > 80 {
		return nil, fmt.Errorf("contents length should be between 1 and 80 characters, but got %v", length)
	}

	forcedCodeSet := -1

	// Check content
	for i := 0; i < length; i++ {
		c := contents[i]
		// check for non ascii characters that are not special GS1 characters
		switch c {
		// special function characters
		case ESCAPE_FNC_1, ESCAPE_FNC_2, ESCAPE_FNC_3, ESCAPE_FNC_4:
		// non ascii characters
		default:
			if c > 127 {
				// no full Latin-1 character set available at the moment
				// shift and manual code change are not supported
				return nil, fmt.Errorf("bad character in input: ASCII value=%v", int(c))
			}
		}

		// check characters for compatibility with forced code set
		switch forcedCodeSet {
		case CODE_A:
			// allows no ascii above 95 (no lower caps, no special symbols)
			if c > 95 && c <= 127 {
				return nil, fmt.Errorf("bad character in input for forced code set A: ASCII value=%v", int(c))
			}
		case CODE_B:
			// allows no ascii below 32 (terminal symbols)
			if c <= 32 {
				return nil, fmt.Errorf("bad character in input for forced code set B: ASCII value=%v", int(c))
			}
		case CODE_C:
			// allows only numbers and no FNC 2/3/4
			if c < 48 || (c > 57 && c <= 127) || c == ESCAPE_FNC_2 || c == ESCAPE_FNC_3 || c == ESCAPE_FNC_4 {
				return nil, fmt.Errorf("bad character in input for forced code set C: ASCII value=%v", int(c))
			}
		}
	}

	var patternsIdx []byte

	codeSet := 0  // selected code (CODE_B or CODE_C)
	position := 0 // position in contents

	for position < length {
		//Select code to use
		newCodeSet := forcedCodeSet
		if newCodeSet == -1 {
			newCodeSet = code128ChooseCode(contents, position, codeSet)
		}

		//Get the pattern index
		var patternIndex int
		if newCodeSet == codeSet {
			// Encode the current character
			// First handle escapes
			switch contents[position] {
			case ESCAPE_FNC_1:
				patternIndex = FNC_1
			case ESCAPE_FNC_2:
				patternIndex = FNC_2
			case ESCAPE_FNC_3:
				patternIndex = FNC_3
			case ESCAPE_FNC_4:
				if codeSet == CODE_A {
					patternIndex = FNC_4_A
				} else {
					patternIndex = FNC_4_B
				}
			default:
				// Then handle normal characters otherwise
				switch codeSet {
				case CODE_A:
					patternIndex = int(contents[position]) - ' '
					if patternIndex < 0 {
						// everything below a space character comes behind the underscore in the code patterns table
						patternIndex += '`'
					}
				case CODE_B:
					patternIndex = int(contents[position]) - ' '
				default:
					// CODE_CODE_C
					if position+1 == length {
						// this is the last character, but the encoding is C, which always encodes two characers
						return nil, fmt.Errorf("bad number of characters for digit only encoding")
					}
					patternIndex = (int(contents[position])-'0')*10 + (int(contents[position+1]) - '0')
					position++ // Also incremented below
				}
			}
			position++
		} else {
			// Should we change the current code?
			// Do we have a code set?
			if codeSet == 0 {
				// No, we don't have a code set
				switch newCodeSet {
				case CODE_A:
					patternIndex = START_A
				case CODE_B:
					patternIndex = START_B
				default:
					patternIndex = START_C
				}
			} else {
				// Yes, we have a code set
				patternIndex = newCodeSet
			}
			codeSet = newCodeSet
		}

		// Get the pattern
		patternsIdx = append(patternsIdx, byte(patternIndex))
	}

	return encode(patternsIdx, height, barWidth)
}

func encode(patternsIdx []byte, height, barWidth int) (image.Image, error) {
	if len(patternsIdx) == 0 {
		return nil, fmt.Errorf("no data to encode")
	}

	var result []bool

	checkSum := int(patternsIdx[0])

	for i, patternIdx := range patternsIdx {
		result = appendPattern(result, CODE_PATTERNS[patternIdx])
		checkSum += int(patternIdx) * i
	}

	// Compute and append checksum
	checkSum = checkSum % 103
	result = appendPattern(result, CODE_PATTERNS[checkSum])

	// Append stop code
	result = appendPattern(result, CODE_PATTERNS[STOP])

	return newCode128(result, height, barWidth), nil
}

func code128FindCType(value []rune, start int) code128CType {
	last := len(value)
	if start >= last {
		return code128CType_UNCODABLE
	}
	c := value[start]
	if c == ESCAPE_FNC_1 {
		return code128CType_FNC_1
	}
	if c < '0' || c > '9' {
		return code128CType_UNCODABLE
	}
	if start+1 >= last {
		return code128CType_ONE_DIGIT
	}
	c = value[start+1]
	if c < '0' || c > '9' {
		return code128CType_ONE_DIGIT
	}
	return code128CType_TWO_DIGITS
}

func code128ChooseCode(value []rune, start, oldCode int) int {
	lookahead := code128FindCType(value, start)
	if lookahead == code128CType_ONE_DIGIT {
		if oldCode == CODE_A {
			return CODE_A
		}
		return CODE_B
	}
	if lookahead == code128CType_UNCODABLE {
		if start < len(value) {
			c := value[start]
			if c < ' ' || (oldCode == CODE_A && (c < '`' ||
				(c >= ESCAPE_FNC_1 && c <= ESCAPE_FNC_4))) {
				// can continue in code A, encodes ASCII 0 to 95 or FNC1 to FNC4
				return CODE_A
			}
		}
		return CODE_B // no choice
	}
	if oldCode == CODE_A && lookahead == code128CType_FNC_1 {
		return CODE_A
	}
	if oldCode == CODE_C { // can continue in code C
		return CODE_C
	}
	if oldCode == CODE_B {
		if lookahead == code128CType_FNC_1 {
			return CODE_B // can continue in code B
		}
		// Seen two consecutive digits, see what follows
		lookahead = code128FindCType(value, start+2)
		if lookahead == code128CType_UNCODABLE || lookahead == code128CType_ONE_DIGIT {
			return CODE_B // not worth switching now
		}
		if lookahead == code128CType_FNC_1 { // two digits, then FNC_1...
			lookahead = code128FindCType(value, start+3)
			if lookahead == code128CType_TWO_DIGITS { // then two more digits, switch
				return CODE_C
			} else {
				return CODE_B // otherwise not worth switching
			}
		}
		// At this point, there are at least 4 consecutive digits.
		// Look ahead to choose whether to switch now or on the next round.
		index := start + 4
		for {
			lookahead = code128FindCType(value, index)
			if lookahead != code128CType_TWO_DIGITS {
				break
			}
			index += 2
		}
		if lookahead == code128CType_ONE_DIGIT { // odd number of digits, switch later
			return CODE_B
		}
		return CODE_C // even number of digits, switch now
	}
	// Here oldCode == 0, which means we are choosing the initial code
	if lookahead == code128CType_FNC_1 { // ignore FNC_1
		lookahead = code128FindCType(value, start+1)
	}
	if lookahead == code128CType_TWO_DIGITS { // at least two digits, start in code C
		return CODE_C
	}
	return CODE_B
}

func appendPattern(target []bool, pattern []int) []bool {
	color := true
	for _, len := range pattern {
		for j := 0; j < len; j++ {
			target = append(target, color)
		}
		color = !color
	}

	return target
}
