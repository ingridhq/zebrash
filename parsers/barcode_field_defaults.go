package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/printers"
)

func NewBarcodeFieldDefaults() *CommandParser {
	const code = "^BY"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			parts := splitCommand(command, code, 0)
			if len(parts) > 0 {
				if v, err := strconv.Atoi(parts[0]); err == nil {
					printer.BarcodeInfo.DefaultModuleWidth = v
				}
			}
			// TODO: We should be parsing WideBarToNarrowBarWidthRatio from parts[1], but:
			// 1. The libs we use for encoding barcodes does not support this
			// 2. We don't need it at the moment
			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					printer.BarcodeInfo.DefaultHeight = v
				}
			}

			return nil, nil
		},
	}
}
