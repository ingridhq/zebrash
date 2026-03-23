package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/internal/hex"
	"github.com/ingridhq/zebrash/internal/printers"
)

func NewDownloadUnboundedTtfParser() *CommandParser {
	const code = "~DU"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			data := command[len(code):]
			parts := strings.SplitN(data, ",", 3)

			path := printers.StoredFontDefaultPath

			if len(parts) > 0 && parts[0] != "" {
				path = parts[0]
			}

			if err := printers.ValidateDevice(path); err != nil {
				return nil, err
			}

			var totalBytes int

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					totalBytes = v
				}
			}

			var fontData []byte

			if len(parts) > 2 {
				data, err := hex.DecodeFontData(parts[2], totalBytes)
				if err != nil {
					return nil, fmt.Errorf("failed to decode embedded true type font data: %w", err)
				}

				fontData = data
			}

			if len(fontData) == 0 {
				return nil, nil
			}

			font, err := truetype.Parse(fontData)
			if err != nil {
				return nil, fmt.Errorf("failed to parse embedded true type font data: %w", err)
			}

			path = printers.EnsureExtensions(path, "TTF", "FNT")
			printer.StoredFonts[path] = font

			return nil, nil
		},
	}
}
