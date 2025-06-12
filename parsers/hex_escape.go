package parsers

import (
	"github.com/ingridhq/zebrash/printers"
)

func NewHexEscapeParser() *CommandParser {
	const code = "^FH"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			text := commandText(command, code)

			char := byte('_')

			if len(text) > 0 {
				char = text[0]
			}

			printer.NextHexEscapeChar = char

			return nil, nil
		},
	}
}
