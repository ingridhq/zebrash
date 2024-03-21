package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldTypesetParser() *CommandParser {
	const code = "^FT"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			pos := elements.LabelPosition{
				CalculateFromBottom: true,
			}

			parts := splitCommand(command, code, 0)

			if len(parts) > 0 {
				if v, err := strconv.Atoi(parts[0]); err == nil {
					pos.X = v
				}
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					pos.Y = v
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
