package parsers

import (
	"math"
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewBarcodeDatamatrixParser() *CommandParser {
	const code = "^BX"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			barcode := &elements.BarcodeDatamatrix{
				Orientation: printer.DefaultOrientation,
				Height:      printer.DefaultBarcodeDimensions.Height,
				Format:      6,
				Escape:      '~',
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

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					barcode.Quality = v
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

			if len(parts) > 5 {
				if v, err := strconv.Atoi(parts[5]); err == nil && v > 0 {
					barcode.Format = v
				}
			}

			if len(parts) > 6 && len(parts[6]) > 0 {
				barcode.Escape = parts[6][0]
			}

			if len(parts) > 7 {
				if v, err := strconv.Atoi(parts[7]); err == nil && v > 0 && v < 3 {
					barcode.Ratio = elements.DatamatrixRatio(v)
				}
			}

			printer.NextElementFieldElement = barcode

			return nil, nil
		},
	}
}
