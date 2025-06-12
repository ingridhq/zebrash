package parsers

import (
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewGraphicSymbolParser() *CommandParser {
	const code = "^GS"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			symbol := &elements.GraphicSymbol{
				Width:       printer.DefaultFont.Width,
				Height:      printer.DefaultFont.Height,
				Orientation: printer.DefaultOrientation,
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 && len(parts[0]) > 0 {
				symbol.Orientation = toFieldOrientation(parts[0][0])
			}

			if len(parts) > 1 {
				v, _ := strconv.Atoi(strings.Trim(parts[1], " "))
				symbol.Height = float64(v)
			}

			if len(parts) > 2 {
				v, _ := strconv.Atoi(strings.Trim(parts[2], " "))
				symbol.Width = float64(v)
			}

			printer.NextElementFieldElement = symbol

			return nil, nil
		},
	}
}
