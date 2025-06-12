package parsers

import (
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldOrientationParser() *CommandParser {
	const code = "^FW"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)

			if len(parts) > 0 && len(parts[0]) > 0 {
				printer.SetDefaultOrientation(toFieldOrientation(parts[0][0]))
			}

			if len(parts) > 1 && len(parts[1]) > 0 {
				printer.DefaultAlignment = toTextAlignment(parts[1][0])
			}

			return nil, nil
		},
	}
}
