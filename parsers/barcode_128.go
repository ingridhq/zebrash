package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewBarcode128Parser() *CommandParser {
	const code = "^BC"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			barcode := &elements.Barcode128{}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 && len(parts[0]) > 0 {
				barcode.Orientation = toFieldOrientation(parts[0][0])
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					barcode.Height = v
				}
			}

			printer.NextElementFieldData = barcode

			return nil, nil
		},
	}
}
