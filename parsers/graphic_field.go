package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/hex"
	"github.com/ingridhq/zebrash/printers"
)

func NewGraphicFieldParser() *CommandParser {
	const code = "^GF"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			result := &elements.GraphicField{
				Position:       printer.NextElementPosition,
				MagnificationX: 1,
				MagnificationY: 1,
				ReversePrint:   printer.GetReversePrint(),
			}

			parts := splitCommand(command, code, 0)
			if len(parts) > 0 && len(parts[0]) > 0 {
				switch parts[0][0] {
				case 'A':
					result.Format = elements.GraphicFieldFormatHex
				case 'B':
					result.Format = elements.GraphicFieldFormatRaw
				case 'C':
					result.Format = elements.GraphicFieldFormatAR
				}
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					result.DataBytes = v
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					result.TotalBytes = v
				}
			}

			if len(parts) > 3 {
				if v, err := strconv.Atoi(parts[3]); err == nil {
					result.RowBytes = v
				}
			}

			if len(parts) > 4 {
				data := strings.Trim(strings.Join(parts[4:], ","), " ")

				switch result.Format {
				case elements.GraphicFieldFormatHex:
					v, err := hex.DecodeGraphicFieldData(data, result.RowBytes)
					if err != nil {
						return nil, fmt.Errorf("failed to decode hex string: %w", err)
					}

					result.Data = v
				case elements.GraphicFieldFormatRaw:
					result.Data = []byte(data)
				}

			}

			return result, nil
		},
	}
}
