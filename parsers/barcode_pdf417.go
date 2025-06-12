package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewBarcodePdf417Parser() *CommandParser {
	const code = "^B7"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			barcode := &elements.BarcodePdf417{
				Orientation: printer.DefaultOrientation,
			}

			parts := splitCommand(command, code, 0)

			if len(parts) > 0 && len(parts[0]) > 0 {
				barcode.Orientation = toFieldOrientation(parts[0][0])
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					barcode.RowHeight = v
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					barcode.Security = v
				}
			}

			if len(parts) > 3 {
				if v, err := strconv.Atoi(parts[3]); err == nil {
					barcode.Columns = v
				}
			}

			if len(parts) > 4 {
				if v, err := strconv.Atoi(parts[4]); err == nil {
					barcode.Rows = v
				}
			}

			if len(parts) > 5 && len(parts[5]) > 0 {
				barcode.Truncate = toBoolField(parts[5][0])
			}

			printer.NextElementFieldElement = barcode

			return nil, nil
		},
	}
}
