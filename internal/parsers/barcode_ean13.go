package parsers

import (
	"math"
	"strconv"
	"strings"

	"github.com/DawidBury/zebrash/internal/elements"
	"github.com/DawidBury/zebrash/internal/printers"
)

func NewBarcodeEan13Parser() *CommandParser {
	const code = "^BE"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			barcode := &elements.BarcodeEan13{
				Orientation: printer.DefaultOrientation,
				Height:      printer.DefaultBarcodeDimensions.Height,
				Line:        true,
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 && len(parts[0]) > 0 {
				barcode.Orientation = toFieldOrientation(parts[0][0])
			}

			if len(parts) > 1 {
				if v, err := strconv.ParseFloat(strings.Trim(parts[1], " "), 32); err == nil {
					barcode.Height = int(math.Ceil(v))
				}
			}

			if len(parts) > 2 && len(parts[2]) > 0 {
				barcode.Line = toBoolField(parts[2][0])
			}

			if len(parts) > 3 && len(parts[3]) > 0 {
				barcode.LineAbove = toBoolField(parts[3][0])
			}

			printer.NextElementFieldElement = barcode

			return nil, nil
		},
	}
}
