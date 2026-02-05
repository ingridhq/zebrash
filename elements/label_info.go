package elements

type LabelInfo struct {
	// Width of the label
	PrintWidth int
	// Inverted mode, which mirrors label content across a horizontal axis.
	Inverted bool
	// Label elements (barcodes, shapes, texts, etc)
	Elements []any
}
