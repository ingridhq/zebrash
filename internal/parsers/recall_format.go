package parsers

import (
	"cmp"

	"github.com/ingridhq/zebrash/internal/printers"
)

func NewRecallFormatParser() *CommandParser {
	const code = "^XF"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			path := commandText(command, code)
			path = cmp.Or(path, printers.StoredFormatDefaultPath)

			if err := printers.ValidateDevice(path); err != nil {
				return nil, err
			}

			if v, ok := printer.StoredFormats[printers.EnsureExtension(path, "ZPL")]; ok {
				return v.ToRecalledFormat(), nil
			}

			return nil, nil
		},
	}
}
