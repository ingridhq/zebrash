package printers

import "github.com/ingridhq/zebrash/elements"

type VirtualPrinter struct {
	LabelHomePosition       elements.LabelPosition
	NextElementPosition     elements.LabelPosition
	DefaultFont             elements.FontInfo
	DefaultOrientation      elements.FieldOrientation
	DefaultAlignment        elements.TextAlignment
	NextElementFieldBlock   *elements.FieldBlock
	NextElementFieldData    interface{}
	NextFont                *elements.FontInfo
	NextDownloadFormatName  string
	NextHexEscapeChar       byte
	NextElementFieldReverse bool
	LabelReverse            bool
}

func NewVirtualPrinter() *VirtualPrinter {
	return &VirtualPrinter{
		DefaultFont: elements.FontInfo{
			Name:   "0",
			Width:  0,
			Height: 9,
		},
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
