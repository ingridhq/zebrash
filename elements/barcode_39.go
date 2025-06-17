package elements

type Barcode39 struct {
	// The bar code orientation to use.
	// Valid values are N (no rotation), R (rotate 90° clockwise), I (rotate 180° clockwise), and B (rotate 270° clockwise).
	// The default value is the orientation configured via the ^FW command, which itself defaults to N (no rotation).
	Orientation FieldOrientation
	// The bar code height, in dots.
	// Any number between 1 and 32,000 may be used.
	// The default value is the bar code height configured via the ^BY command, which itself defaults to 10.
	Height int
	// Whether or not to include human-readable text with the bar code.
	// Valid values are Y and N. The default value is Y (include human-readable text).
	Line bool
	// Whether or not to place the human-readable text above the bar code.
	// Valid values are Y and N. The default value is N (if printed, text is placed below the bar code).
	LineAbove bool
	// Whether or not to add a check digit to the bar code. Valid values are Y and N.
	// The default value is N (no check digit).
	CheckDigit bool
}

type Barcode39WithData struct {
	Barcode39
	Width      int
	WidthRatio float64
	Position   LabelPosition
	Data       string
}
