package parsers

import (
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldValueParser() *CommandParser {
	const code = "^FV"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			return &elements.TextField{
				Font:         printer.GetNextFontOrDefault(),
				Pos:          printer.NextElementPosition,
				Orientation:  printer.DefaultOrientation,
				Alignment:    printer.DefaultAlignment,
				Text:         commandText(command, code),
				Block:        printer.NextElementFieldBlock,
				ReversePrint: printer.IsReversePrint(),
			}, nil
		},
	}
}
