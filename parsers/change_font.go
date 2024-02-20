package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewChangeFontParser() *CommandParser {
	const code = "^A"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			font := elements.FontInfo{
				Name:        printer.DefaultFont.Name,
				Width:       printer.DefaultFont.Width,
				Height:      printer.DefaultFont.Height,
				Orientation: elements.FieldOrientationNormal,
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 && len(parts[0]) > 1 {
				font.Name = string(parts[0][0])
				font.Orientation = toFieldOrientation(parts[0][1])
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

			printer.NextFont = &font

			return nil, nil
		},
	}
}
