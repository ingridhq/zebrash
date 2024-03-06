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
			text := commandText(command, code)

			if printer.NextElementFieldData != nil {
				switch fd := printer.NextElementFieldData.(type) {
				case *elements.Barcode128:
					return &elements.Barcode128WithData{
						Barcode128: *fd,
						Position:   printer.NextElementPosition,
						Data:       text,
					}, nil
				}
			}

			return &elements.TextField{
				Font:         printer.GetNextFontOrDefault(),
				Position:     printer.NextElementPosition,
				Orientation:  printer.DefaultOrientation,
				Alignment:    printer.DefaultAlignment,
				Text:         commandText(command, code),
				Block:        printer.NextElementFieldBlock,
				ReversePrint: printer.IsReversePrint(),
			}, nil
		},
	}
}
