package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/printers"
)

func NewLabelHomeParser() *CommandParser {
	const code = "^LH"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			pos := printer.LabelHomePosition

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

			printer.LabelHomePosition = pos

			return nil, nil
		},
	}
}
