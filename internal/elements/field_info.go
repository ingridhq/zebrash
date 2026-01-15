package elements

type FieldInfo struct {
	ReversePrint

	Element        any
	Font           FontInfo
	Position       LabelPosition
	Alignment      FieldAlignment
	Width          int
	WidthRatio     float64
	Height         int
	CurrentCharset int
}
