package elements

type BarcodePdf417 struct {
	// The bar code orientation to use.
	// Valid values are N (no rotation),
	// R (rotate 90° clockwise),
	// I (rotate 180° clockwise), and B (rotate 270° clockwise).
	// The default value is the orientation configured via the ^FW command, which itself defaults to N (no rotation).
	Orientation FieldOrientation

	// The bar code row height, in dots.
	// Any number between 1 and the label height may be used.
	// The default value is the row height necessary for the total bar code height
	// to match the bar code height configured via the ^BY command.
	RowHeight int

	// The level of error correction to apply.
	// Any number between 0 and 8 may be used.
	// The higher the number, the larger the generated bar code and the more resilient it is to scan errors.
	// The default value is 0 (scan errors are detected but not corrected).
	Security int

	// The number of data columns to encode.
	// Any number between 1 and 30 may be used.
	// This parameter can be used to control the bar code width.
	// The default value depends on the amount of data encoded.
	Columns int

	// The number of rows to encode.
	// Any number between 3 and 90 may be used.
	// This parameter can be used to control the bar code height.
	// The default value depends on the amount of data encoded.
	Rows int

	// Whether or not to generate a truncated PDF417 bar code, also known as compact PDF417.
	// Truncated PDF417 bar codes are narrower because they do not include right row indicators,
	// but should only be used when label damage is unlikely. Valid values are Y and N.
	// The default value is N (do not truncate).
	Truncate bool
}

type BarcodePdf417WithData struct {
	BarcodePdf417
	Position LabelPosition
	Data     string
}
