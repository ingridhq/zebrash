package printers

import (
	"github.com/ingridhq/zebrash/internal/elements"
)

type VirtualPrinter struct {
	StoredGraphics          map[string]elements.StoredGraphics
	StoredFormats           map[string]*elements.StoredFormat
	LabelHomePosition       elements.LabelPosition
	NextElementPosition     elements.LabelPosition
	DefaultFont             elements.FontInfo
	DefaultOrientation      elements.FieldOrientation
	DefaultAlignment        elements.FieldAlignment
	NextElementAlignment    *elements.FieldAlignment
	NextElementFieldElement any
	NextElementFieldData    string
	NextElementFieldNumber  int
	NextFont                *elements.FontInfo
	// When not empty ZPL elements are parsed into the stored template instead of the output
	NextDownloadFormatName   string
	NextHexEscapeChar        byte
	NextElementFieldReverse  bool
	LabelReverse             bool
	DefaultBarcodeDimensions elements.BarcodeDimensions
	CurrentCharset           int
	PrintWidth               int
	LabelInverted            bool
}

func NewVirtualPrinter() *VirtualPrinter {
	return &VirtualPrinter{
		StoredGraphics: make(map[string]elements.StoredGraphics),
		StoredFormats:  make(map[string]*elements.StoredFormat),
		DefaultFont: elements.FontInfo{
			Name: "A",
		},
		DefaultAlignment: elements.FieldAlignmentLeft,
		DefaultBarcodeDimensions: elements.BarcodeDimensions{
			ModuleWidth: 2,
			Height:      10,
			WidthRatio:  3,
		},
		NextElementFieldNumber: -1,
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

func (p *VirtualPrinter) GetNextElementAlignmentOrDefault() elements.FieldAlignment {
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

func (p *VirtualPrinter) GetFieldInfo() elements.FieldInfo {
	return elements.FieldInfo{
		ReversePrint:   p.GetReversePrint(),
		Element:        p.NextElementFieldElement,
		Font:           p.GetNextFontOrDefault(),
		Position:       p.NextElementPosition,
		Alignment:      p.GetNextElementAlignmentOrDefault(),
		Width:          p.DefaultBarcodeDimensions.ModuleWidth,
		WidthRatio:     p.DefaultBarcodeDimensions.WidthRatio,
		Height:         p.DefaultBarcodeDimensions.Height,
		CurrentCharset: p.CurrentCharset,
	}
}

func (p *VirtualPrinter) ResetFieldState() {
	p.NextElementPosition = elements.LabelPosition{}
	p.NextElementFieldElement = nil
	p.NextElementFieldData = ""
	p.NextElementFieldNumber = -1
	p.NextElementAlignment = nil
	p.NextFont = nil
	p.NextElementFieldReverse = false
	p.NextHexEscapeChar = 0
}

func (p *VirtualPrinter) ResetLabelState() {
	p.NextDownloadFormatName = ""
	p.LabelInverted = false
}
