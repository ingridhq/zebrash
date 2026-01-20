package parsers

import (
	"strconv"

	"github.com/DawidBury/zebrash/internal/printers"
)

func NewFieldNumberParser() *CommandParser {
	const code = "^FN"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			number := commandText(command, code)

			if v, err := strconv.Atoi(number); err == nil && v >= 0 {
				printer.NextElementFieldNumber = v
			}

			return nil, nil
		},
	}
}
