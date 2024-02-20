package elements

type TextField struct {
	Font         FontInfo
	Pos          LabelPosition
	Orientation  FieldOrientation
	Alignment    TextAlignment
	Text         string
	Block        *FieldBlock
	ReversePrint bool
}
