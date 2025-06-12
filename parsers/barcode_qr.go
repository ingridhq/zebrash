package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

// ^BQ orientation, model, magnification, errorCorrection, mask
func NewBarcodeQrParser() *CommandParser {
	const code = "^BQ"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			barcode := &elements.BarcodeQr{
				Magnification: 1,
			}

			parts := splitCommand(command, code, 0)

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					barcode.Magnification = min(max(v, 1), 10)
				}
			}

			printer.NextElementFieldElement = barcode

			return nil, nil
		},
	}
}
