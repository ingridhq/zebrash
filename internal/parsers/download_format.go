package parsers

import (
	"cmp"

	"github.com/DawidBury/zebrash/internal/printers"
)

func NewDownloadFormatParser() *CommandParser {
	const code = "^DF"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			path := commandText(command, code)
			path = cmp.Or(path, printers.StoredFormatDefaultPath)

			if err := printers.ValidateDevice(path); err != nil {
				return nil, err
			}

			printer.NextDownloadFormatName = printers.EnsureExtension(path, "ZPL")

			return nil, nil
		},
	}
}
