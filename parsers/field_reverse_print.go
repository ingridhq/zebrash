package parsers

import (
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldReversePrintParser() *CommandParser {
	const code = "^FR"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			printer.NextElementFieldReverse = true

			return nil, nil
		},
	}
}
