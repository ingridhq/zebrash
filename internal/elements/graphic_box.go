package elements

type GraphicBox struct {
	ReversePrint

	Position        LabelPosition
	Width           int
	Height          int
	BorderThickness int
	CornerRounding  int
	LineColor       LineColor
}
