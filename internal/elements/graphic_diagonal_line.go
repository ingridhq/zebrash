package elements

type GraphicDiagonalLine struct {
	ReversePrint

	Position        LabelPosition
	Width           int
	Height          int
	BorderThickness int
	LineColor       LineColor
	TopToBottom     bool
}
