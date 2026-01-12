package ean13

import (
	"fmt"
	"image"
	"strconv"
)

const (
	// Total width of EAN-13 barcode
	CODE_WIDTH = 3 + // start guard
		(7 * 6) + // left bars
		5 + // middle guard
		(7 * 6) + // right bars
		3 // end guard
)

var (
	// Start and end guard pattern
	START_END_PATTERN = []int{1, 1, 1}

	// Middle guard pattern
	MIDDLE_PATTERN = []int{1, 1, 1, 1, 1}

	// L patterns (used on right side, and left side with odd parity)
	// These represent: 0-9
	L_PATTERNS = [][]int{
		{3, 2, 1, 1}, // 0
		{2, 2, 2, 1}, // 1
		{2, 1, 2, 2}, // 2
		{1, 4, 1, 1}, // 3
		{1, 1, 3, 2}, // 4
		{1, 2, 3, 1}, // 5
		{1, 1, 1, 4}, // 6
		{1, 3, 1, 2}, // 7
		{1, 2, 1, 3}, // 8
		{3, 1, 1, 2}, // 9
	}

	// L and G patterns combined (index 0-9 are L patterns, 10-19 are G patterns)
	// G patterns are the inverse of L patterns (bars become spaces, spaces become bars)
	L_AND_G_PATTERNS = [][]int{
		{3, 2, 1, 1}, // 0 L
		{2, 2, 2, 1}, // 1 L
		{2, 1, 2, 2}, // 2 L
		{1, 4, 1, 1}, // 3 L
		{1, 1, 3, 2}, // 4 L
		{1, 2, 3, 1}, // 5 L
		{1, 1, 1, 4}, // 6 L
		{1, 3, 1, 2}, // 7 L
		{1, 2, 1, 3}, // 8 L
		{3, 1, 1, 2}, // 9 L
		{1, 1, 2, 3}, // 0 G
		{1, 2, 2, 2}, // 1 G
		{2, 2, 1, 2}, // 2 G
		{1, 1, 4, 1}, // 3 G
		{2, 3, 1, 1}, // 4 G
		{1, 3, 2, 1}, // 5 G
		{4, 1, 1, 1}, // 6 G
		{2, 1, 3, 1}, // 7 G
		{3, 1, 2, 1}, // 8 G
		{2, 1, 1, 3}, // 9 G
	}

	// First digit encoding determines the parity pattern for left side
	// Each entry is a 6-bit value where 0=L pattern, 1=G pattern
	FIRST_DIGIT_ENCODINGS = []int{
		0x00, // 0: LLLLLL
		0x0B, // 1: LLGLGG
		0x0D, // 2: LLGGLG
		0x0E, // 3: LLGGGL
		0x13, // 4: LGLLGG
		0x19, // 5: LGGLLG
		0x1C, // 6: LGGGLL
		0x15, // 7: LGLGLG
		0x16, // 8: LGLGGL
		0x1A, // 9: LGGLGL
	}
)

// Encode creates an EAN-13 barcode image
// contents should be 12 or 13 digits (check digit is auto-calculated if 12 digits provided)
// barWidth specifies the width of each module in pixels
func Encode(contents string, height, barWidth int) (image.Image, error) {
	length := len(contents)

	// Validate and prepare contents
	switch length {
	case 12:
		// No check digit present, calculate it and add it
		check, err := calculateChecksum(contents)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate checksum: %w", err)
		}
		contents += strconv.Itoa(check)
	case 13:
		// Verify checksum
		ok, err := verifyChecksum(contents)
		if err != nil {
			return nil, fmt.Errorf("failed to verify checksum: %w", err)
		}
		if !ok {
			return nil, fmt.Errorf("contents do not pass checksum")
		}
	default:
		return nil, fmt.Errorf("contents should be 12 or 13 digits long, but got %d", length)
	}

	// Verify all characters are numeric
	for _, c := range contents {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("contents must be numeric, found: %c", c)
		}
	}

	// Build the barcode pattern
	result := make([]bool, CODE_WIDTH)
	pos := 0

	// Start guard
	pos += appendPattern(result, pos, START_END_PATTERN, true)

	// First digit determines the parity pattern
	firstDigit := int(contents[0] - '0')
	parities := FIRST_DIGIT_ENCODINGS[firstDigit]

	// Left side (digits 1-6)
	for i := 1; i <= 6; i++ {
		digit := int(contents[i] - '0')
		// Check if we should use G pattern (bit is set) or L pattern (bit is clear)
		if ((parities >> uint(6-i)) & 1) == 1 {
			digit += 10 // Use G pattern
		}
		pos += appendPattern(result, pos, L_AND_G_PATTERNS[digit], false)
	}

	// Middle guard
	pos += appendPattern(result, pos, MIDDLE_PATTERN, false)

	// Right side (digits 7-12)
	for i := 7; i <= 12; i++ {
		digit := int(contents[i] - '0')
		pos += appendPattern(result, pos, L_PATTERNS[digit], true)
	}

	// End guard
	appendPattern(result, pos, START_END_PATTERN, true)

	return newEan13(result, height, barWidth), nil
}

// appendPattern appends a pattern to the result array
// startColor indicates if we should start with a black bar (true) or white space (false)
func appendPattern(target []bool, pos int, pattern []int, startColor bool) int {
	color := startColor
	numAdded := 0

	for _, width := range pattern {
		for i := 0; i < width; i++ {
			target[pos] = color
			pos++
			numAdded++
		}
		color = !color
	}

	return numAdded
}

// calculateChecksum calculates the EAN-13 checksum digit
func calculateChecksum(contents string) (int, error) {
	if len(contents) != 12 {
		return 0, fmt.Errorf("contents must be 12 digits for checksum calculation")
	}

	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(contents[i] - '0')
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}

	checksum := (10 - (sum % 10)) % 10
	return checksum, nil
}

// verifyChecksum verifies the EAN-13 checksum
func verifyChecksum(contents string) (bool, error) {
	if len(contents) != 13 {
		return false, fmt.Errorf("contents must be 13 digits for checksum verification")
	}

	expectedCheck, err := calculateChecksum(contents[:12])
	if err != nil {
		return false, err
	}

	actualCheck := int(contents[12] - '0')
	return expectedCheck == actualCheck, nil
}
