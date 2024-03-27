package parsers

import (
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/printers"
)

func NewChangeFontParser() *CommandParser {
	const code = "^A"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			font := printer.DefaultFont

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 && len(parts[0]) > 1 {
				font.Name = strings.ToUpper(string(parts[0][0]))
				font.Orientation = toFieldOrientation(parts[0][1])
			}

			if len(parts) > 1 {
				v, _ := strconv.Atoi(parts[1])
				font.Height = float64(v)
			}

			if len(parts) > 2 {
				v, _ := strconv.Atoi(parts[2])
				font.Width = float64(v)

				if font.Height == 0 {
					font.Height = font.Width
				}
			}

			font = font.WithAdjustedSizes()
			printer.NextFont = &font

			return nil, nil
		},
	}
}
