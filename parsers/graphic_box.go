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
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			result := &elements.GraphicBox{
				Position:        printer.NextElementPosition,
				Width:           1,
				Height:          1,
				BorderThickness: 1,
				CornerRounding:  0,
				LineColor:       elements.LineColorBlack,
				ReversePrint:    printer.IsReversePrint(),
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
				if v, err := strconv.Atoi(parts[4]); err == nil {
					result.CornerRounding = v
				}
			}

			return result, nil
		},
	}
}
