package parsers

import (
	"github.com/ingridhq/zebrash/internal/printers"
)

func NewChangeFontAliasParser() *CommandParser {
	const code = "^CW"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)

			var alias string
			if len(parts) > 0 {
				alias = toValidFontName(parts[0])
			}

			var path string
			if len(parts) > 1 {
				path = parts[1]
			}

			if alias != "" && path != "" {
				printer.StoredFontAliases[alias] = path
			}

			return nil, nil
		},
	}
}
