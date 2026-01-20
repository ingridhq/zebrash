package parsers

import (
	"github.com/DawidBury/zebrash/internal/elements"
	"github.com/DawidBury/zebrash/internal/printers"
)

func NewFieldOriginParser() *CommandParser {
	const code = "^FO"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			pos := elements.LabelPosition{
				CalculateFromBottom: false,
			}

			parts := splitCommand(command, code, 0)

			if len(parts) > 0 {
				if v, err := toPositiveIntField(parts[0]); err == nil {
					pos.X = v
				}
			}

			if len(parts) > 1 {
				if v, err := toPositiveIntField(parts[1]); err == nil {
					pos.Y = v
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
