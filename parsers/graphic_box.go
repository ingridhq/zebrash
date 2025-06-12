package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewGraphicBoxParser() *CommandParser {
	const code = "^GB"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			result := &elements.GraphicBox{
				Position:        printer.NextElementPosition,
				Width:           1,
				Height:          1,
				BorderThickness: 1,
				CornerRounding:  0,
				LineColor:       elements.LineColorBlack,
				ReversePrint:    printer.GetReversePrint(),
			}

			parts := splitCommand(command, code, 0)

			if len(parts) > 2 {
				if v, err := toPositiveIntField(parts[2]); err == nil && v > 0 {
					result.BorderThickness = v
				}
			}

			if len(parts) > 0 {
				if v, err := toPositiveIntField(parts[0]); err == nil && v > 0 {
					result.Width = max(v, result.BorderThickness)
				}
			}

			if len(parts) > 1 {
				if v, err := toPositiveIntField(parts[1]); err == nil && v > 0 {
					result.Height = max(v, result.BorderThickness)
				}
			}

			if len(parts) > 3 && parts[3] == "W" {
				result.LineColor = elements.LineColorWhite
			}

			if len(parts) > 4 {
				if v, err := strconv.Atoi(parts[4]); err == nil && v > 0 && v < 9 {
					result.CornerRounding = v
				}
			}

			return result, nil
		},
	}
}
