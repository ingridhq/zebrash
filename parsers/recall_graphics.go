package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewRecallGraphicsParser() *CommandParser {
	const code = "^XG"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)

			result := &elements.GraphicField{
				Position:       printer.NextElementPosition,
				MagnificationX: 1,
				MagnificationY: 1,
				ReversePrint:   printer.GetReversePrint(),
			}

			path := elements.StoredGraphicsDefaultPath

			if len(parts) > 0 && parts[0] != "" {
				path = parts[0]
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil && v >= 0 && v <= 10 {
					result.MagnificationX = v
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil && v >= 0 && v <= 10 {
					result.MagnificationY = v
				}
			}

			if v, ok := printer.StoredGraphics[path]; ok {
				result.Data = v.Data
				result.DataBytes = v.TotalBytes
				result.TotalBytes = v.TotalBytes
				result.RowBytes = v.RowBytes
			}

			return result, nil
		},
	}
}
