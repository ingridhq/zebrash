package parsers

import (
	"math"
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/internal/elements"
	"github.com/ingridhq/zebrash/internal/printers"
)

func NewBarcode39Parser() *CommandParser {
	const code = "^B3"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			barcode := &elements.Barcode39{
				Orientation: printer.DefaultOrientation,
				Height:      printer.DefaultBarcodeDimensions.Height,
				Line:        true,
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 && len(parts[0]) > 0 {
				barcode.Orientation = toFieldOrientation(parts[0][0])
			}

			if len(parts) > 1 && len(parts[1]) > 0 {
				barcode.CheckDigit = toBoolField(parts[1][0])
			}

			if len(parts) > 2 {
				if v, err := strconv.ParseFloat(strings.Trim(parts[2], " "), 32); err == nil {
					barcode.Height = int(math.Ceil(v))
				}
			}

			if len(parts) > 3 && len(parts[3]) > 0 {
				barcode.Line = toBoolField(parts[3][0])
			}

			if len(parts) > 4 && len(parts[4]) > 0 {
				barcode.LineAbove = toBoolField(parts[4][0])
			}

			printer.NextElementFieldElement = barcode

			return nil, nil
		},
	}
}
