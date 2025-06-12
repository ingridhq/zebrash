package parsers

import (
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldValueParser() *CommandParser {
	const code = "^FV"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			printer.NextElementFieldData = commandText(command, code)

			return nil, nil
		},
	}
}
