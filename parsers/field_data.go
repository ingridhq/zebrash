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
						Code:     *fd,
						Position: printer.NextElementPosition,
						Data:     text,
					}, nil
				case *elements.BarcodePdf417:
					return &elements.BarcodePdf417WithData{
						BarcodePdf417: *fd,
						Position:      printer.NextElementPosition,
						Data:          text,
					}, nil
				case *elements.Barcode128:
					return &elements.Barcode128WithData{
						Barcode128: *fd,
						Width:      printer.DefaultBarcodeDimensions.ModuleWidth,
						Position:   printer.NextElementPosition,
						Data:       text,
					}, nil
				case *elements.BarcodeAztec:
					return &elements.BarcodeAztecWithData{
						BarcodeAztec: *fd,
						Position:     printer.NextElementPosition,
						Data:         text,
					}, nil
				}
			}

			reversePrint := printer.IsReversePrint()
			font := printer.GetNextFontOrDefault()
			pos := printer.NextElementPosition

			return &elements.TextField{
				Font:         font,
				Position:     pos,
				Orientation:  printer.DefaultOrientation,
				Alignment:    printer.GetNextElementAlignmentOrDefault(),
				Text:         text,
				Block:        printer.NextElementFieldBlock,
				ReversePrint: reversePrint,
			}, nil
		},
	}
}
