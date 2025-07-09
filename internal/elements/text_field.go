package elements

type TextField struct {
	ReversePrint

	Font      FontInfo
	Position  LabelPosition
	Alignment TextAlignment
	Text      string
	Block     *FieldBlock
}
