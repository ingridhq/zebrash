package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewGraphicDiagonalLineParser() *CommandParser {
	const code = "^GD"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			result := &elements.GraphicDiagonalLine{
				Position:        printer.NextElementPosition,
				Width:           1,
				Height:          1,
				BorderThickness: 1,
				LineColor:       elements.LineColorBlack,
				ReversePrint:    printer.GetReversePrint(),
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 {
				if v, err := strconv.Atoi(parts[0]); err == nil {
					result.Width = v
				}
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					result.Height = v
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					result.BorderThickness = v
				}
			}

			if len(parts) > 3 && parts[3] == "W" {
				result.LineColor = elements.LineColorWhite
			}

			if len(parts) > 4 {
				result.TopToBottom = (parts[4] == "L")
			}

			return result, nil
		},
	}
}
