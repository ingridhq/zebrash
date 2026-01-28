package parsers

import (
	"github.com/ingridhq/zebrash/internal/printers"
)

func NewPrintOrientationParser() *CommandParser {
	return &CommandParser{
		CommandCode: "^PO",
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, "^PO", 0)
			if len(parts) > 0 && len(parts[0]) > 0 {
				printer.LabelOrientation = toFieldOrientation(parts[0][0])
			}
			return nil, nil
		},
	}
}
