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

func DecodeEmbeddedImage(value string) ([]byte, error) {
	const z64Prefix = ":Z64:"
	if strings.HasPrefix(value, z64Prefix) {
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

	return hx.DecodeString(value)
}
