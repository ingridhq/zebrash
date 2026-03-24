package parsers

import (
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/internal/elements"
	"github.com/ingridhq/zebrash/internal/printers"
)

func NewChangeDefaultFontParser() *CommandParser {
	const code = "^CF"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			fontName := printer.DefaultFont.Name

			font := elements.FontInfo{
				Name:        fontName,
				Orientation: printer.DefaultOrientation,
				CustomFont:  printer.StoredFonts[printer.StoredFontAliases[fontName]],
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 {
				font.Name = toValidFontName(parts[0])
			}

			if len(parts) > 1 {
				v, _ := strconv.Atoi(strings.Trim(parts[1], " "))
				font.Height = float64(v)
			}

			if len(parts) > 2 {
				v, _ := strconv.Atoi(strings.Trim(parts[2], " "))
				font.Width = float64(v)
			}

			printer.DefaultFont = font

			return nil, nil
		},
	}
}
