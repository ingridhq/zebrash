package parsers

import (
	"github.com/ingridhq/zebrash/printers"
)

func NewHexEscapeParser() *CommandParser {
	const code = "^FH"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			text := commandText(command, code)

			if len(text) > 0 {
				printer.NextHexEscapeChar = text[0]
			}

			return nil, nil
		},
	}
}
