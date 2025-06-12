package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewBarcodeAztecParser() *CommandParser {
	const code = "^BO"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			barcode := &elements.BarcodeAztec{
				Orientation: printer.DefaultOrientation,
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 && len(parts[0]) > 0 {
				barcode.Orientation = toFieldOrientation(parts[0][0])
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					barcode.Magnification = v
				}
			}

			// TODO: handle eci (parts[2])

			if len(parts) > 3 {
				if v, err := strconv.Atoi(parts[3]); err == nil {
					barcode.Size = v
				}
			}

			printer.NextElementFieldElement = barcode

			return nil, nil
		},
	}
}
