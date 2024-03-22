package parsers

import (
	"github.com/ingridhq/zebrash/printers"
)

func NewLabelReversePrintParser() *CommandParser {
	const code = "^LR"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			text := commandText(command, code)

			printer.LabelReverse = (text == "Y")

			return nil, nil
		},
	}
}
