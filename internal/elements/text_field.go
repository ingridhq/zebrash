package elements

type TextField struct {
	ReversePrint

	Font      FontInfo
	Position  LabelPosition
	Alignment FieldAlignment
	Text      string
	Block     *FieldBlock
}
