package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/printers"
)

func NewPrintWidthParser() *CommandParser {
	const code = "^PW"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)

			if len(parts) > 0 {
				if v, err := strconv.Atoi(parts[0]); err == nil {
					printer.PrintWidth = max(v, 2)
				}
			}

			return nil, nil
		},
	}
}
