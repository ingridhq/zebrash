package printers

import (
	"github.com/ingridhq/zebrash/elements"
)

type VirtualPrinter struct {
	StoredGraphics           map[string]elements.StoredGraphics
	LabelHomePosition        elements.LabelPosition
	NextElementPosition      elements.LabelPosition
	DefaultFont              elements.FontInfo
	DefaultOrientation       elements.FieldOrientation
	DefaultAlignment         elements.TextAlignment
	NextElementAlignment     *elements.TextAlignment
	NextElementFieldElement  any
	NextElementFieldData     string
	NextFont                 *elements.FontInfo
	NextDownloadFormatName   string
	NextHexEscapeChar        byte
	NextElementFieldReverse  bool
	LabelReverse             bool
	DefaultBarcodeDimensions elements.BarcodeDimensions
	CurrentCharset           int
	PrintWidth               int
}

func NewVirtualPrinter() *VirtualPrinter {
	return &VirtualPrinter{
		StoredGraphics: map[string]elements.StoredGraphics{},
		DefaultFont: elements.FontInfo{
			Name: "A",
		}.WithAdjustedSizes(),
		DefaultBarcodeDimensions: elements.BarcodeDimensions{
			ModuleWidth: 2,
			Height:      10,
			WidthRatio:  3,
		},
	}
}

func (p *VirtualPrinter) SetDefaultOrientation(orientation elements.FieldOrientation) {
	p.DefaultOrientation = orientation
	p.DefaultFont.Orientation = orientation
	if p.NextFont != nil {
		p.NextFont.Orientation = orientation
	}
}

func (p *VirtualPrinter) GetNextFontOrDefault() elements.FontInfo {
	if p.NextFont != nil {
		return *p.NextFont
	}

	return p.DefaultFont
}

func (p *VirtualPrinter) GetNextElementAlignmentOrDefault() elements.TextAlignment {
	if p.NextElementAlignment != nil {
		return *p.NextElementAlignment
	}

	return p.DefaultAlignment
}

func (p *VirtualPrinter) GetReversePrint() elements.ReversePrint {
	return elements.ReversePrint{
		Value: p.NextElementFieldReverse || p.LabelReverse,
	}
}

func (p *VirtualPrinter) ResetState() {
	p.NextElementPosition = elements.LabelPosition{}
	p.NextElementFieldElement = nil
	p.NextElementFieldData = ""
	p.NextElementAlignment = nil
	p.NextFont = nil
	p.NextElementFieldReverse = false
	p.NextHexEscapeChar = 0
}
