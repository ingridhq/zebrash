package parsers

import (
	"fmt"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/hex"
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldDataParser() *CommandParser {
	const code = "^FD"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			var err error

			text := commandText(command, code)

			if printer.NextHexEscapeChar != 0 {
				text, err = hex.DecodeEscapedString(text, printer.NextHexEscapeChar)
				if err != nil {
					return nil, fmt.Errorf("failed to decode escaped hex string: %w", err)
				}
			}

			if printer.NextElementFieldData != nil {
				switch fd := printer.NextElementFieldData.(type) {
				case *elements.Maxicode:
					return &elements.MaxicodeWithData{
						Code: *fd,
						Pos:  printer.NextElementPosition,
						Data: text,
					}, nil
				}
			}

			reversePrint := printer.IsReversePrint()
			font := printer.GetNextFontOrDefault()
			pos := printer.NextElementPosition

			return &elements.TextField{
				Font:         font,
				Pos:          pos,
				Orientation:  printer.DefaultOrientation,
				Alignment:    printer.DefaultAlignment,
				Text:         text,
				Block:        printer.NextElementFieldBlock,
				ReversePrint: reversePrint,
			}, nil
		},
	}
}
