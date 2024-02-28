package hex

import (
	hx "encoding/hex"
	"fmt"
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
