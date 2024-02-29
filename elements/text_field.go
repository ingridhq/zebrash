package elements

type TextField struct {
	Font         FontInfo
	Position     LabelPosition
	Orientation  FieldOrientation
	Alignment    TextAlignment
	Text         string
	Block        *FieldBlock
	ReversePrint bool
}
