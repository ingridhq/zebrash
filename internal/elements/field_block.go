package elements

type FieldBlock struct {
	// The maximum text block width, in dots.
	// Text longer than this width is wrapped to another line.
	// Any number between 0 and 9,999 may be used. The default value is 0, which is rarely useful.
	MaxWidth int

	// The maximum number of text lines to allow.
	// If the text does not fit on the specified number of lines, any remaining text is drawn over the previous text on the last line.
	// Any number between 1 and 9,999 may be used. The default value is 1.
	MaxLines int

	// Extra spacing to add between lines, in dots.
	// Positive numbers increase the distance between lines, negative numbers decrease the distance between lines.
	// Any number between -9,999 and 9,999 may be used. The default value is 0.
	LineSpacing int

	// The text alignment to apply to the text block.
	// Valid values are L (left), R (right), C (center) and J (justified).
	// The default value is L (left).
	Alignment TextAlignment

	// The hanging indent to apply to all lines except the first line, in dots.
	// Any number between 0 and 9,999 may be used. The default value is 0.
	HangingIndent int
}
