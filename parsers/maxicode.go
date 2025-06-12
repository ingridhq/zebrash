package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewMaxicodeParser() *CommandParser {
	const code = "^BD"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			barcode := &elements.Maxicode{}

			parts := splitCommand(command, code, 0)

			if len(parts) > 0 {
				if v, err := strconv.Atoi(parts[0]); err == nil {
					barcode.Mode = v
				}
			}

			printer.NextElementFieldElement = barcode

			return nil, nil
		},
	}
}
