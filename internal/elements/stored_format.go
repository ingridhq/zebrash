package elements

import (
	"fmt"
	"strings"

	"github.com/ingridhq/zebrash/internal/encodings"
)

type StoredFormat struct {
	Inverted bool
	Elements []any
}

func (sf *StoredFormat) ToRecalledFormat() *RecalledFormat {
	res := &RecalledFormat{
		Inverted:  sf.Inverted,
		fieldRefs: make(map[int][]*RecalledField),
	}

	for _, el := range sf.Elements {
		res.AddElement(el)
	}

	return res
}

type RecalledFormat struct {
	Inverted  bool
	elements  []any
	fieldRefs map[int][]*RecalledField
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
			ReversePrint: field.ReversePrint,
			Code:         *fe,
			Position:     field.Position,
			Data:         text,
		}, nil
	case *BarcodePdf417:
		return &BarcodePdf417WithData{
			ReversePrint:  field.ReversePrint,
			BarcodePdf417: *fe,
			Position:      field.Position,
			Data:          text,
		}, nil
	case *Barcode128:
		return &Barcode128WithData{
			ReversePrint: field.ReversePrint,
			Barcode128:   *fe,
			Width:        field.Width,
			Position:     field.Position,
			Data:         text,
		}, nil
	case *BarcodeEan13:
		return &BarcodeEan13WithData{
			ReversePrint: field.ReversePrint,
			BarcodeEan13: *fe,
			Width:        field.Width,
			Position:     field.Position,
			Data:         text,
		}, nil
	case *Barcode2of5:
		return &Barcode2of5WithData{
			ReversePrint: field.ReversePrint,
			Barcode2of5:  *fe,
			Width:        field.Width,
			WidthRatio:   field.WidthRatio,
			Position:     field.Position,
			Data:         text,
		}, nil
	case *Barcode39:
		return &Barcode39WithData{
			ReversePrint: field.ReversePrint,
			Barcode39:    *fe,
			Width:        field.Width,
			WidthRatio:   field.WidthRatio,
			Position:     field.Position,
			Data:         text,
		}, nil
	case *BarcodeAztec:
		return &BarcodeAztecWithData{
			ReversePrint: field.ReversePrint,
			BarcodeAztec: *fe,
			Position:     field.Position,
			Data:         text,
		}, nil
	case *BarcodeDatamatrix:
		return &BarcodeDatamatrixWithData{
			ReversePrint:      field.ReversePrint,
			BarcodeDatamatrix: *fe,
			Position:          field.Position,
			Data:              text,
		}, nil
	case *BarcodeQr:
		return &BarcodeQrWithData{
			ReversePrint: field.ReversePrint,
			BarcodeQr:    *fe,
			Height:       field.Height,
			Position:     field.Position,
			Data:         text,
		}, nil
	case *GraphicSymbol:
		return toGraphicSymbolTextField(text, field, fe)
	case *FieldBlock:
		return toTextField(text, field, fe)
	}

	return toTextField(text, field, nil)
}

func toGraphicSymbolTextField(text string, field FieldInfo, fe *GraphicSymbol) (any, error) {
	text = toGSText(text)
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
		Text:         text,
		ReversePrint: field.ReversePrint,
	}, nil
}

func toGSText(text string) string {
	var res strings.Builder

	for _, r := range text {
		// Keep leading spaces
		if r == ' ' {
			res.WriteRune(r)
			continue
		}

		if r >= 'A' && r <= 'E' {
			res.WriteRune(r)
		}

		// We stop after the first non-space character
		break
	}

	return res.String()
}

func toTextField(text string, field FieldInfo, fe *FieldBlock) (any, error) {
	// \& = carriage return/line feed
	text = strings.ReplaceAll(text, `\&`, "\n")

	unicodeText, err := encodings.ToUnicodeText(text, field.CurrentCharset)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to unicode text: %w", err)
	}

	return &TextField{
		Font:         field.Font.WithAdjustedSizes(),
		Position:     field.Position,
		Alignment:    field.Alignment,
		Text:         unicodeText,
		Block:        fe,
		ReversePrint: field.ReversePrint,
	}, nil
}
