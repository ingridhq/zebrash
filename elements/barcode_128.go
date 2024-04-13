package elements

type BarcodeMode int

const (
	BarcodeModeNo        BarcodeMode = 0
	BarcodeModeUcc       BarcodeMode = 1
	BarcodeModeAutomatic BarcodeMode = 2
	BarcodeModeEan       BarcodeMode = 3
)

type Barcode128 struct {
	// The bar code orientation to use.
	// Valid values are N (no rotation), R (rotate 90° clockwise), I (rotate 180° clockwise), and B (rotate 270° clockwise).
	// The default value is the orientation configured via the ^FW command, which itself defaults to N (no rotation).
	Orientation FieldOrientation
	// The bar code height, in dots.
	// Any number between 1 and 32,000 may be used.
	// The default value is the bar code height configured via the ^BY command, which itself defaults to 10.
	Height int
	// Whether or not to include human-readable text with the bar code.
	// Valid values are Y and N.
	// The default value is Y (include human-readable text).
	Line bool
	// Whether or not to place the human-readable text above the bar code.
	// Valid values are Y and N.
	// The default value is N (if printed, text is placed below the bar code),
	// except for mode U where the default is Y (if printed, text is placed above the bar code).
	LineAbove bool
	// Whether or not to calculate a GS1 (UCC) Mod 10 check digit.
	// Valid values are Y and N.
	// The default value is N (GS1 check digit is not calculated).
	// TODO: Figure out if it should be implemented, as it's part of the interface but reference libraries disregard this value.
	CheckDigit bool
	// The mode to use to encode the bar code data.
	// Valid values are N (no mode, subsets are specified explicitly as part of the field data),
	// U (UCC case mode, field data must contain 19 digits),
	// A (automatic mode, the ZPL engine automatically determines the subsets that are used to encode the data),
	// and D (UCC/EAN mode, field data must contain GS1 numbers).
	// The default value is N (no mode, subsets are specified explicitly as part of the field data).
	Mode BarcodeMode
}

type Barcode128WithData struct {
	Barcode128
	Width    int
	Position LabelPosition
	Data     string
}
