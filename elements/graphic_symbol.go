package elements

type GraphicSymbol struct {
	Width       float64
	Height      float64
	Orientation FieldOrientation
}

type GraphicSymbolWithData struct {
	GraphicSymbol
	Position LabelPosition
	Data     string
}
