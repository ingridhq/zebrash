package parsers

import (
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewChangeDefaultFontParser() *CommandParser {
	const code = "^CF"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			font := elements.FontInfo{
				Name:        printer.DefaultFont.Name,
				Orientation: printer.DefaultOrientation,
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 {
				font.Name = strings.ToUpper(parts[0])
			}

			if len(parts) > 1 {
				v, _ := strconv.Atoi(parts[1])
				font.Height = float64(v)
			}

			if len(parts) > 2 {
				v, _ := strconv.Atoi(parts[2])
				font.Width = float64(v)
			}

			printer.DefaultFont = font.WithAdjustedSizes()

			return nil, nil
		},
	}
}
