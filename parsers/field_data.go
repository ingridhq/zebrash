package parsers

import (
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldDataParser() *CommandParser {
	const code = "^FD"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			text := commandText(command, code)
			reversePrint := printer.IsReversePrint()
			font := printer.GetNextFontOrDefault()
			pos := printer.NextElementPosition

			return &elements.TextField{
				Font:         font,
				Pos:          pos,
				Orientation:  printer.DefaultOrientation,
				Alignment:    printer.DefaultAlignment,
				Text:         text,
				Block:        printer.NextElementFieldBlock,
				ReversePrint: reversePrint,
			}, nil
		},
	}
}
