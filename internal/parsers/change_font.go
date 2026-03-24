package parsers

import (
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/zebrash/internal/elements"
	"github.com/ingridhq/zebrash/internal/printers"
)

func NewChangeFontParser() *CommandParser {
	const code = "^A"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)
			if len(parts) == 0 || len(parts[0]) == 0 {
				// Use default font
				printer.NextFont = nil
				return nil, nil
			}

			fontName := toValidFontName(parts[0])

			var customFont *truetype.Font

			switch {
			// Font is referenced by alias
			case parts[0][0] != '@':
				customFont = printer.StoredFonts[printer.StoredFontAliases[fontName]]
			// Font is referenced directly by filename
			case parts[0][0] == '@' && len(parts) > 3:
				fontPath := strings.Trim(parts[3], " ")
				customFont = printer.StoredFonts[fontPath]
			}

			font := &elements.FontInfo{
				Name:        fontName,
				Orientation: printer.DefaultFont.Orientation,
				CustomFont:  customFont,
			}

			if !font.Exists() {
				df := printer.DefaultFont
				font = &df
			}

			if len(parts[0]) > 1 {
				font.Orientation = toFieldOrientation(parts[0][1])
			}

			if len(parts) > 1 {
				v, _ := strconv.Atoi(strings.Trim(parts[1], " "))
				font.Height = float64(v)
			}

			if len(parts) > 2 {
				v, _ := strconv.Atoi(strings.Trim(parts[2], " "))
				font.Width = float64(v)
			}

			printer.NextFont = font

			return nil, nil
		},
	}
}
