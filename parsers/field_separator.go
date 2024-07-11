package parsers

import (
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

func NewFieldSeparatorParser() *CommandParser {
	const code = "^FS"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (interface{}, error) {
			defer printer.ResetState()

			text := printer.NextElementFieldData
			if text == "" {
				return nil, nil
			}

			if printer.NextElementFieldElement != nil {
				switch fe := printer.NextElementFieldElement.(type) {
				case *elements.Maxicode:
					return &elements.MaxicodeWithData{
						Code:     *fe,
						Position: printer.NextElementPosition,
						Data:     text,
					}, nil
				case *elements.BarcodePdf417:
					return &elements.BarcodePdf417WithData{
						BarcodePdf417: *fe,
						Position:      printer.NextElementPosition,
						Data:          text,
					}, nil
				case *elements.Barcode128:
					return &elements.Barcode128WithData{
						Barcode128: *fe,
						Width:      printer.DefaultBarcodeDimensions.ModuleWidth,
						Position:   printer.NextElementPosition,
						Data:       text,
					}, nil
				case *elements.BarcodeAztec:
					return &elements.BarcodeAztecWithData{
						BarcodeAztec: *fe,
						Position:     printer.NextElementPosition,
						Data:         text,
					}, nil
				case *elements.BarcodeDatamatrix:
					return &elements.BarcodeDatamatrixWithData{
						BarcodeDatamatrix: *fe,
						Position:          printer.NextElementPosition,
						Data:              text,
					}, nil
				case *elements.BarcodeQr:
					return &elements.BarcodeQrWithData{
						BarcodeQr: *fe,
						Position:  printer.NextElementPosition,
						Data:      text,
					}, nil
				case *elements.FieldBlock:
					return &elements.TextField{
						Font:         printer.GetNextFontOrDefault(),
						Position:     printer.NextElementPosition,
						Alignment:    printer.GetNextElementAlignmentOrDefault(),
						Text:         text,
						Block:        fe,
						ReversePrint: printer.GetReversePrint(),
					}, nil
				}
			}

			return &elements.TextField{
				Font:         printer.GetNextFontOrDefault(),
				Position:     printer.NextElementPosition,
				Alignment:    printer.GetNextElementAlignmentOrDefault(),
				Text:         text,
				ReversePrint: printer.GetReversePrint(),
			}, nil
		},
	}
}
