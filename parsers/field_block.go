package parsers

import (
	"strconv"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldBlockParser() *CommandParser {
	const code = "^FB"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			block := elements.FieldBlock{
				MaxWidth:      0,
				MaxLines:      1,
				LineSpacing:   0,
				Alignment:     printer.GetNextElementAlignmentOrDefault(),
				HangingIndent: 0,
			}

			parts := splitCommand(command, code, 0)

			if len(parts) > 0 {
				if v, err := strconv.Atoi(parts[0]); err == nil {
					block.MaxWidth = v
				}
			}

			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					block.MaxLines = v
				}
			}

			if len(parts) > 2 {
				if v, err := strconv.Atoi(parts[2]); err == nil {
					block.LineSpacing = v
				}
			}

			if len(parts) > 3 && len(parts[3]) > 0 {
				block.Alignment = toTextAlignment(parts[3][0])
			}

			if len(parts) > 4 {
				if v, err := strconv.Atoi(parts[4]); err == nil {
					block.HangingIndent = v
				}
			}

			printer.NextElementFieldElement = &block

			return nil, nil
		},
	}
}
