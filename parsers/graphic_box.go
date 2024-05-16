package parsers

import (
	"math"
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
				ReversePrint:    printer.GetReversePrint(),
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 {
				if v, err := strconv.ParseFloat(parts[0], 32); err == nil && v > 0 {
					result.Width = int(math.Ceil(v))
				}
			}

			if len(parts) > 1 {
				if v, err := strconv.ParseFloat(parts[1], 32); err == nil && v > 0 {
					result.Height = int(math.Ceil(v))
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil && v > 0 {
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
