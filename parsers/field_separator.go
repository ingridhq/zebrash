package parsers

import (
	"fmt"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
	"golang.org/x/text/encoding/charmap"
)

func NewFieldSeparatorParser() *CommandParser {
	const code = "^FS"

	return &CommandParser{
		CommandCode: code,
		Parse: func(command string, printer *printers.VirtualPrinter) (any, error) {
			defer printer.ResetState()

			text := printer.NextElementFieldData

			if printer.NextElementFieldElement == nil && text == "" {
				return nil, nil
			}

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
			case *elements.Barcode2of5:
				return &elements.Barcode2of5WithData{
					Barcode2of5: *fe,
					Width:       printer.DefaultBarcodeDimensions.ModuleWidth,
					WidthRatio:  printer.DefaultBarcodeDimensions.WidthRatio,
					Position:    printer.NextElementPosition,
					Data:        text,
				}, nil
			case *elements.Barcode39:
				return &elements.Barcode39WithData{
					Barcode39:  *fe,
					Width:      printer.DefaultBarcodeDimensions.ModuleWidth,
					WidthRatio: printer.DefaultBarcodeDimensions.WidthRatio,
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

	// \& = carriage return/line feed
	unicodeText = strings.ReplaceAll(unicodeText, `\&`, "\n")

	return &elements.TextField{
		Font:         printer.GetNextFontOrDefault(),
		Position:     printer.NextElementPosition,
		Alignment:    printer.GetNextElementAlignmentOrDefault(),
		Text:         unicodeText,
		Block:        fe,
		ReversePrint: printer.GetReversePrint(),
	}, nil
}

// Encodings 0-13 are all in fact CP850 encoding
// 13 is normal CP850
// 0-12 have some characters replaced with other characters
var characterSets013 = [14][11]string{
	{"#", "0", "@", "[", "¢", "]", "^", "`", "{", "|", "}"},
	{"#", "0", "@", "⅓", "¢", "⅔", "^", "`", "¼", "½", "¾"},
	{"£", "0", "@", "[", "¢", "]", "^", "`", "{", "|", "}"},
	{"ƒ", "0", "§", "[", "IJ", "]", "^", "`", "{", "ij", "}"},
	{"#", "0", "@", "Æ", "Ø", "Å", "^", "`", "æ", "ø", "å"},
	{"Ü", "0", "É", "Ä", "Ö", "Å", "Ü", "é", "ä", "ö", "å"},
	{"#", "0", "§", "Ä", "Ö", "Ü", "^", "`", "ä", "ö", "ü"},
	{"£", "0", "à", "[", "ç", "]", "^", "`", "é", "|", "ù"},
	{"#", "0", "à", "â", "ç", "ê", "î", "ô", "é", "ù", "è"},
	{"£", "0", "§", "[", "ç", "é", "^", "ù", "à", "ò", "è"},
	{"#", "0", "§", "¡", "Ñ", "¿", "^", "`", "{", "ñ", "ç"},
	{"£", "0", "É", "Ä", "Ö", "Ü", "^", "ä", "ë", "ï", "ö"},
	{"#", "0", "@", "[", "¥", "]", "^", "`", "{", "|", "}"},
	{"#", "0", "@", "[", "\\", "]", "^", "`", "{", "|", "}"},
}

func toUnicodeText(text string, printer *printers.VirtualPrinter) (string, error) {
	switch {
	case printer.CurrentCharset >= 0 && printer.CurrentCharset <= 13:
		text, err := charmap.CodePage850.NewDecoder().String(text)
		if err != nil {
			return "", err
		}

		if printer.CurrentCharset < 13 {
			search := characterSets013[13]
			replace := characterSets013[printer.CurrentCharset]

			for i, v := range search {
				text = strings.ReplaceAll(text, v, replace[i])
			}
		}

		return text, nil
	case printer.CurrentCharset == 27:
		return charmap.Windows1252.NewDecoder().String(text)
	default:
		return text, nil
	}
}
