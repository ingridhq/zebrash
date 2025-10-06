package parsers

import (
	"github.com/ingridhq/zebrash/internal/elements"
	"github.com/ingridhq/zebrash/internal/printers"
)

func NewFieldTypesetParser() *CommandParser {
	const code = "^FT"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)

			pos := elements.LabelPosition{
				CalculateFromBottom: true,
				AutomaticPosition:   true,
			}

			if len(parts) > 0 {
				if v, err := toPositiveIntField(parts[0]); err == nil {
					pos.X = v
					pos.AutomaticPosition = false
				}
			}

			if len(parts) > 1 {
				if v, err := toPositiveIntField(parts[1]); err == nil {
					pos.Y = v
					pos.AutomaticPosition = false
				}
			}

			if len(parts) > 2 {
				if val, ok := toFieldAlignment(parts[2]); ok {
					printer.NextElementAlignment = &val
				}
			}

			printer.NextElementPosition = pos.Add(printer.LabelHomePosition)

			return nil, nil
		},
	}
}
