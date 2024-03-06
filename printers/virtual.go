package printers

import "github.com/ingridhq/zebrash/elements"

type VirtualPrinter struct {
	StoredGraphics            map[string]elements.StoredGraphics
	LabelHomePosition         elements.LabelPosition
	NextElementPosition       elements.LabelPosition
	DefaultFont               elements.FontInfo
	DefaultOrientation        elements.FieldOrientation
	DefaultAlignment          elements.TextAlignment
	NextElementFieldBlock     *elements.FieldBlock
	NextElementFieldData      interface{}
	NextFont                  *elements.FontInfo
	NextDownloadFormatName    string
	NextHexEscapeChar         byte
	NextElementFieldReverse   bool
	LabelReverse              bool
	DefaultBarcodeModuleWidth int
	DefaultBarcodeHeight      int
}

func NewVirtualPrinter() *VirtualPrinter {
	return &VirtualPrinter{
		StoredGraphics: map[string]elements.StoredGraphics{},
		DefaultFont: elements.FontInfo{
			Name:   "0",
			Width:  0,
			Height: 9,
		},
		DefaultBarcodeModuleWidth: 2,
		DefaultBarcodeHeight:      10,
	}
}

func (p *VirtualPrinter) GetNextFontOrDefault() elements.FontInfo {
	if p.NextFont != nil {
		return *p.NextFont
	}

	return p.DefaultFont
}

func (p *VirtualPrinter) IsReversePrint() bool {
	return p.NextElementFieldReverse || p.LabelReverse
}
