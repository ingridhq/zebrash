package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/hex"
	"github.com/ingridhq/zebrash/printers"
)

func NewDownloadGraphicsParser() *CommandParser {
	const code = "~DG"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			parts := splitCommand(command, code, 0)

			graphics := elements.StoredGraphics{
				RowBytes: 1,
			}

			path := elements.StoredGraphicsDefaultPath

			if len(parts) > 0 && parts[0] != "" {
				path = parts[0]
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					graphics.TotalBytes = v
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					graphics.RowBytes = v
				}
			}

			if len(parts) > 3 {
				data, err := hex.DecodeEmbeddedImage(strings.Join(parts[3:], ""))
				if err != nil {
					return nil, fmt.Errorf("failed to decode embedded graphics: %w", err)
				}

				graphics.Data = data
			}

			printer.StoredGraphics[path] = graphics

			return nil, nil
		},
	}
}
