package hex

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	hx "encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func DecodeEscapedString(value string, escapeChar byte) (string, error) {
	replaceChar := string(escapeChar)

	hexEscapeRegex, err := regexp.Compile(regexp.QuoteMeta(replaceChar) + `([0-9A-Fa-f]{2})`)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %w", err)
	}

	res := hexEscapeRegex.ReplaceAllStringFunc(value, func(s string) string {
		v, err := hx.DecodeString(strings.TrimPrefix(s, replaceChar))
		if err == nil {
			return string(v)
		}

		return s
	})

	return res, nil
}

var compressCounts = map[byte]int{
	'G': 1,
	'H': 2,
	'I': 3,
	'J': 4,
	'K': 5,
	'L': 6,
	'M': 7,
	'N': 8,
	'O': 9,
	'P': 10,
	'Q': 11,
	'R': 12,
	'S': 13,
	'T': 14,
	'U': 15,
	'V': 16,
	'W': 17,
	'X': 18,
	'Y': 19,
	'g': 20,
	'h': 40,
	'i': 60,
	'j': 80,
	'k': 100,
	'l': 120,
	'm': 140,
	'n': 160,
	'o': 180,
	'p': 200,
	'q': 220,
	'r': 240,
	's': 260,
	't': 280,
	'u': 300,
	'v': 320,
	'w': 340,
	'x': 360,
	'y': 380,
	'z': 400,
}

func DecodeGraphicFieldData(data string, rowBytes int) ([]byte, error) {
	if z64Encoded(data) {
		return decodeZ64(data)
	}

	var (
		result strings.Builder
		line   strings.Builder
	)

	prevLine := ""
	compressCount := 0
	rowHex := rowBytes * 2

	for i := 0; i < len(data); i++ {
		char := data[i]

		if line.Len() >= rowHex {
			prevLine = line.String()
			line.Reset()
			result.WriteString(prevLine)
		}

		if c, ok := compressCounts[char]; ok {
			compressCount += c
			continue
		}

		switch char {
		case ',':
			if rowHex > line.Len() {
				line.WriteString(strings.Repeat("0", rowHex-line.Len()))
			}

			continue
		case '!':
			if rowHex > line.Len() {
				line.WriteString(strings.Repeat("1", rowHex-line.Len()))
			}

			continue
		case ':':
			line.WriteString(prevLine)
			continue
		}

		line.WriteString(strings.Repeat(string(char), max(compressCount, 1)))
		compressCount = 0
	}

	if line.Len() > 0 {
		result.WriteString(line.String())
	}

	return hx.DecodeString(result.String())
}

const z64Prefix = ":Z64:"

func z64Encoded(value string) bool {
	return strings.HasPrefix(value, z64Prefix)
}

func decodeZ64(value string) ([]byte, error) {
	value = strings.TrimPrefix(value, z64Prefix)

	idx := strings.LastIndex(value, ":")
	if idx >= 0 {
		value = value[:idx]
	}

	dec, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(dec)
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}

	defer z.Close()

	return io.ReadAll(z)
}
