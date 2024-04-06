package parsers

import (
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewChangeFontParser() *CommandParser {
	const code = "^A"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			parts := splitCommand(command, code, 0)
			if len(parts) == 0 || len(parts[0]) == 0 {
				// Use default font
				printer.NextFont = nil
				return nil, nil
			}

			font := elements.FontInfo{
				Name: strings.ToUpper(string(parts[0][0])),
			}

			if len(parts[0]) > 1 {
				font.Orientation = toFieldOrientation(parts[0][1])
			}

			if len(parts) > 1 {
				v, _ := strconv.Atoi(parts[1])
				font.Height = float64(v)
			}

			if len(parts) > 2 {
				v, _ := strconv.Atoi(parts[2])
				font.Width = float64(v)
			}

			font = font.WithAdjustedSizes()
			printer.NextFont = &font

			return nil, nil
		},
	}
}
