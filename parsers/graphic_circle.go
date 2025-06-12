package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewGraphicCircleParser() *CommandParser {
	const code = "^GC"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			result := &elements.GraphicCircle{
				Position:        printer.NextElementPosition,
				CircleDiameter:  3,
				BorderThickness: 1,
				LineColor:       elements.LineColorBlack,
				ReversePrint:    printer.GetReversePrint(),
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 {
				if v, err := strconv.Atoi(parts[0]); err == nil {
					result.CircleDiameter = v
				}
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					result.BorderThickness = v
				}
			}

			if len(parts) > 2 && parts[2] == "W" {
				result.LineColor = elements.LineColorWhite
			}

			return result, nil
		},
	}
}
