package parsers

import (
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/printers"
)

func NewBarcodeFieldDefaults() *CommandParser {
	const code = "^BY"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)
			if len(parts) > 0 {
				if v, err := strconv.Atoi(parts[0]); err == nil {
					printer.DefaultBarcodeDimensions.ModuleWidth = v
				}
			}

			if len(parts) > 1 {
				if v, err := strconv.ParseFloat(strings.Trim(parts[1], " "), 32); err == nil {
					printer.DefaultBarcodeDimensions.WidthRatio = max(2, min(v, 3))
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					printer.DefaultBarcodeDimensions.Height = v
				}
			}

			return nil, nil
		},
	}
}
