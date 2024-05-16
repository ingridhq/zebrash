package parsers

import (
	"math"
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldOriginParser() *CommandParser {
	const code = "^FO"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			pos := elements.LabelPosition{
				CalculateFromBottom: false,
			}

			parts := splitCommand(command, code, 0)

			if len(parts) > 0 {
				if v, err := strconv.ParseFloat(strings.Trim(parts[0], " "), 32); err == nil {
					pos.X = int(math.Abs(math.Ceil(v)))
				}
			}

			if len(parts) > 1 {
				if v, err := strconv.ParseFloat(strings.Trim(parts[1], " "), 32); err == nil {
					pos.Y = int(math.Abs(math.Ceil(v)))
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					switch v {
					case 0:
						val := elements.TextAlignmentLeft
						printer.NextElementAlignment = &val
					case 1:
						val := elements.TextAlignmentRight
						printer.NextElementAlignment = &val
					case 2:
						val := elements.TextAlignmentJustified
						printer.NextElementAlignment = &val
					}
				}
			}

			printer.NextElementPosition = pos.Add(printer.LabelHomePosition)

			return nil, nil
		},
	}
}
