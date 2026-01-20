package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DawidBury/zebrash/internal/elements"
	"github.com/DawidBury/zebrash/internal/hex"
	"github.com/DawidBury/zebrash/internal/printers"
)

func NewDownloadGraphicsParser() *CommandParser {
	const code = "~DG"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			data := command[len(code):]
			parts := strings.SplitN(data, ",", 4)

			graphics := elements.StoredGraphics{
				RowBytes: 1,
			}

			path := printers.StoredGraphicsDefaultPath

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
					graphics.RowBytes = min(v, 9999999)
				}
			}

			if len(parts) > 3 {
				data, err := hex.DecodeGraphicFieldData(parts[3], graphics.RowBytes)
				if err != nil {
					return nil, fmt.Errorf("failed to decode embedded graphics: %w", err)
				}

				graphics.Data = data
			}

			path = printers.EnsureExtension(path, "GRF")
			printer.StoredGraphics[path] = graphics

			return nil, nil
		},
	}
}
