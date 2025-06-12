package parsers

import (
	"fmt"

	"github.com/ingridhq/zebrash/hex"
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldDataParser() *CommandParser {
	const code = "^FD"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			var err error

			text := commandText(command, code)

			if printer.NextHexEscapeChar != 0 {
				text, err = hex.DecodeEscapedString(text, printer.NextHexEscapeChar)
				if err != nil {
					return nil, fmt.Errorf("failed to decode escaped hex string: %w", err)
				}
			}

			printer.NextElementFieldData = text

			return nil, nil
		},
	}
}
