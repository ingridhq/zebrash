package parsers

import (
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldSeparatorParser() *CommandParser {
	const code = "^FS"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			printer.NextElementPosition = elements.LabelPosition{}
			printer.NextElementFieldBlock = nil
			printer.NextElementFieldData = nil
			printer.NextFont = nil
			printer.NextElementFieldReverse = false
			printer.DefaultAlignment = elements.TextAlignmentLeft
			printer.DefaultOrientation = elements.FieldOrientationNormal

			return nil, nil
		},
	}
}
