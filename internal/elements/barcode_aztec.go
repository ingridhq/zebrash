package elements

type BarcodeAztec struct {
	// The bar code orientation to use.
	// Valid values are N (no rotation),
	// R (rotate 90° clockwise),
	// I (rotate 180° clockwise), and B (rotate 270° clockwise).
	// The default value is the orientation configured via the ^FW command, which itself defaults to N (no rotation).
	Orientation FieldOrientation

	// The bar code magnification to use.
	// Any number between 1 and 10 may be used.
	// The default value depends on the print density being used.
	Magnification int

	// The Aztec bar code size to use.
	// Valid values are 101-104 (compact Aztec code sizes),
	// 201-232 (full-range Aztec code sizes),
	// 300 (Aztec runes),
	// and 1-99 (dynamic sizing for a specific minimum error correction percentage).
	// By default, the bar code is sized dynamically to fit the encoded data.
	Size int
}

type BarcodeAztecWithData struct {
	BarcodeAztec
	Position LabelPosition
	Data     string
}
