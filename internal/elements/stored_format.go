package elements

import (
	"fmt"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type StoredFormat struct {
	Elements []any
}

func (sf *StoredFormat) ToRecalledFormat() *RecalledFormat {
	res := NewRecalledFormat()

	for _, el := range sf.Elements {
		res.AddElement(el)
	}

	return res
}

type RecalledFormat struct {
	elements  []any
	fieldRefs map[int][]*RecalledField
}

func NewRecalledFormat() *RecalledFormat {
	return &RecalledFormat{
		fieldRefs: make(map[int][]*RecalledField),
	}
}

func (rf *RecalledFormat) AddElement(element any) bool {
	if rf == nil {
		return false
	}

	switch e := element.(type) {
	case *StoredField:
		// Convert field to recalled so it can be populated with data
		field := &RecalledField{
			StoredField: *e,
			Data:        "",
		}

		rf.elements = append(rf.elements, field)
		rf.fieldRefs[field.Number] = append(rf.fieldRefs[field.Number], field)
	case *RecalledFieldData:
		refs, ok := rf.fieldRefs[e.Number]
		// If there are currently unpopulated refs lets populate them with data
		if ok {
			for _, ref := range refs {
				ref.Data = e.Data
			}

			delete(rf.fieldRefs, e.Number)
		} else {
			// No more refs to populate, we simply add a new field
			rf.elements = append(rf.elements, &RecalledField{
				Data: e.Data,
			})
		}

	default:
		rf.elements = append(rf.elements, element)
	}

	return true
}

func (rf *RecalledFormat) ResolveElements() ([]any, error) {
	if rf == nil {
		return nil, nil
	}

	res := make([]any, 0, len(rf.elements))

	for _, element := range rf.elements {
		switch el := element.(type) {
		case *RecalledField:
			re, err := el.Resolve()
			if err != nil {
				return nil, err
			}

			res = append(res, re)
		default:
			res = append(res, el)
		}
	}

	return res, nil
}

type StoredField struct {
	Number int
	Field  FieldInfo
}

type RecalledFieldData struct {
	Number int
	Data   string
}

type RecalledField struct {
	StoredField
	Data string
}

func (f *RecalledField) Resolve() (any, error) {
	field := f.Field
	text := f.Data

	if field.Element == nil && text == "" {
		return nil, nil
	}

	switch fe := field.Element.(type) {
	case *Maxicode:
		return &MaxicodeWithData{
			Code:     *fe,
			Position: field.Position,
			Data:     text,
		}, nil
	case *BarcodePdf417:
		return &BarcodePdf417WithData{
			BarcodePdf417: *fe,
			Position:      field.Position,
			Data:          text,
		}, nil
	case *Barcode128:
		return &Barcode128WithData{
			Barcode128: *fe,
			Width:      field.Width,
			Position:   field.Position,
			Data:       text,
		}, nil
	case *Barcode2of5:
		return &Barcode2of5WithData{
			Barcode2of5: *fe,
			Width:       field.Width,
			WidthRatio:  field.WidthRatio,
			Position:    field.Position,
			Data:        text,
		}, nil
	case *Barcode39:
		return &Barcode39WithData{
			Barcode39:  *fe,
			Width:      field.Width,
			WidthRatio: field.WidthRatio,
			Position:   field.Position,
			Data:       text,
		}, nil
	case *BarcodeAztec:
		return &BarcodeAztecWithData{
			BarcodeAztec: *fe,
			Position:     field.Position,
			Data:         text,
		}, nil
	case *BarcodeDatamatrix:
		return &BarcodeDatamatrixWithData{
			BarcodeDatamatrix: *fe,
			Position:          field.Position,
			Data:              text,
		}, nil
	case *BarcodeQr:
		return &BarcodeQrWithData{
			BarcodeQr: *fe,
			Height:    field.Height,
			Position:  field.Position,
			Data:      text,
		}, nil
	case *GraphicSymbol:
		return toGraphicSymbolTextField(text, field, fe)
	case *FieldBlock:
		return toTextField(text, field, fe)
	}

	return toTextField(text, field, nil)
}

var graphicSymbols = map[byte]string{
	'A': "®",
	'B': "©",
	'C': "™",
}

func toGraphicSymbolTextField(text string, field FieldInfo, fe *GraphicSymbol) (*TextField, error) {
	if text == "" {
		return nil, nil
	}

	return &TextField{
		Font: FontInfo{
			Name:        "GS",
			Width:       fe.Width,
			Height:      fe.Height,
			Orientation: fe.Orientation,
		}.WithAdjustedSizes(),
		Position:     field.Position,
		Alignment:    field.Alignment,
		Text:         graphicSymbols[text[0]],
		ReversePrint: field.ReversePrint,
	}, nil
}

func toTextField(text string, field FieldInfo, fe *FieldBlock) (*TextField, error) {
	// \& = carriage return/line feed
	text = strings.ReplaceAll(text, `\&`, "\n")

	unicodeText, err := toUnicodeText(text, field)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to unicode text: %w", err)
	}

	return &TextField{
		Font:         field.Font,
		Position:     field.Position,
		Alignment:    field.Alignment,
		Text:         unicodeText,
		Block:        fe,
		ReversePrint: field.ReversePrint,
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

func toUnicodeText(text string, field FieldInfo) (string, error) {
	switch {
	case field.CurrentCharset >= 0 && field.CurrentCharset <= 13:
		text, err := charmap.CodePage850.NewDecoder().String(text)
		if err != nil {
			return "", err
		}

		if field.CurrentCharset < 13 {
			search := characterSets013[13]
			replace := characterSets013[field.CurrentCharset]

			for i, v := range search {
				text = strings.ReplaceAll(text, v, replace[i])
			}
		}

		return text, nil
	case field.CurrentCharset == 27:
		return charmap.Windows1252.NewDecoder().String(text)
	default:
		return text, nil
	}
}
