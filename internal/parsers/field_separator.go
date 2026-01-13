package parsers

import (
	"github.com/ingridhq/zebrash/internal/elements"
	"github.com/ingridhq/zebrash/internal/printers"
)

func NewFieldSeparatorParser() *CommandParser {
	const code = "^FS"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			defer printer.ResetState()

			if printer.NextElementFieldNumber < 0 {
				f := &elements.RecalledField{
					StoredField: elements.StoredField{
						Number: printer.NextElementFieldNumber,
						Field:  printer.GetFieldInfo(),
					},
					Data: printer.NextElementFieldData,
				}

				// Not a part of template, we can just resolve it immediately
				return f.Resolve()
			}

			if printer.NextDownloadFormatName == "" {
				return &elements.RecalledFieldData{
					Number: printer.NextElementFieldNumber,
					Data:   printer.NextElementFieldData,
				}, nil
			}

			return &elements.StoredField{
				Number: printer.NextElementFieldNumber,
				Field:  printer.GetFieldInfo(),
			}, nil
		},
	}
}
