package parsers

import (
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewImageLoadParser() *CommandParser {
	const code = "^IL"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)

			result := &elements.GraphicField{
				MagnificationX: 1,
				MagnificationY: 1,
			}

			path := elements.StoredGraphicsDefaultPath

			if len(parts) > 0 && parts[0] != "" {
				path = parts[0]
			}

			v, ok := printer.StoredGraphics[path]
			if !ok {
				return nil, nil
			}

			result.Data = v.Data
			result.DataBytes = v.TotalBytes
			result.TotalBytes = v.TotalBytes
			result.RowBytes = v.RowBytes

			return result, nil
		},
	}
}
