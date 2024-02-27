package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewChangeDefaultFontParser() *CommandParser {
	const code = "^CF"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			font := elements.FontInfo{
				Name:        "0",
				Width:       0,
				Height:      9,
				Orientation: elements.FieldOrientationNormal,
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 {
				font.Name = parts[0]
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					font.Height = v
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					font.Width = v
				}
			}

			printer.DefaultFont = font

			return nil, nil
		},
	}
}
