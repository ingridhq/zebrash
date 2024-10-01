package parsers

import (
	"fmt"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
	"golang.org/x/text/encoding/charmap"
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
				case *elements.GraphicSymbol:
					return toGraphicSymbolTextField(text, printer, fe)
				case *elements.FieldBlock:
					return toTextField(text, printer, fe)
				}
			}

			return toTextField(text, printer, nil)
		},
	}
}

var graphicSymbols = map[byte]string{
	'A': "®",
	'B': "©",
	'C': "™",
}

func toGraphicSymbolTextField(text string, printer *printers.VirtualPrinter, fe *elements.GraphicSymbol) (*elements.TextField, error) {
	if text == "" {
		return nil, nil
	}

	return &elements.TextField{
		Font: elements.FontInfo{
			Name:        "GS",
			Width:       fe.Width,
			Height:      fe.Height,
			Orientation: fe.Orientation,
		}.WithAdjustedSizes(),
		Position:     printer.NextElementPosition,
		Alignment:    printer.GetNextElementAlignmentOrDefault(),
		Text:         graphicSymbols[text[0]],
		ReversePrint: printer.GetReversePrint(),
	}, nil
}

func toTextField(text string, printer *printers.VirtualPrinter, fe *elements.FieldBlock) (*elements.TextField, error) {
	unicodeText, err := toUnicodeText(text, printer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to unicode text: %w", err)
	}

	return &elements.TextField{
		Font:         printer.GetNextFontOrDefault(),
		Position:     printer.NextElementPosition,
		Alignment:    printer.GetNextElementAlignmentOrDefault(),
		Text:         unicodeText,
		Block:        fe,
		ReversePrint: printer.GetReversePrint(),
	}, nil
}

func toUnicodeText(text string, printer *printers.VirtualPrinter) (string, error) {
	switch printer.CurrentCharset {
	case 0:
		return charmap.CodePage850.NewDecoder().String(text)
	case 27:
		return charmap.Windows1252.NewDecoder().String(text)
	default:
		return text, nil
	}
}
