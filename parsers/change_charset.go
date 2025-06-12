package parsers

import (
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/printers"
)

func NewChangeCharsetParser() *CommandParser {
	const code = "^CI"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)

			if len(parts) > 0 {
				if v, err := strconv.Atoi(strings.Trim(parts[0], " ")); err == nil {
					printer.CurrentCharset = v
				}
			}

			return nil, nil
		},
	}
}
