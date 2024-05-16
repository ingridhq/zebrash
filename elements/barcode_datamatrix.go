package elements

type DatamatrixRatio int

const (
	DatamatrixRatioSquare      DatamatrixRatio = 1
	DatamatrixRatioRectangular DatamatrixRatio = 2
)

type BarcodeDatamatrix struct {
	// The bar code orientation to use.
	// Valid values are N (no rotation), R (rotate 90° clockwise), I (rotate 180° clockwise), and B (rotate 270° clockwise).
	// The default value is the orientation configured via the ^FW command, which itself defaults to N (no rotation).
	Orientation FieldOrientation
	// The bar code element height, in dots. The individual elements are square, so the element height and width will be the same.
	// Any number between 1 and the label width may be used.
	// The default value is the element height necessary for the total bar code height to match the bar code height configured via the ^BY command.
	Height int

	// The level of error correction to apply.
	// Valid values are 0 (ECC 0), 50 (ECC 50), 80 (ECC 80), 100 (ECC 100), 140 (ECC 140) and 200 (ECC 200).
	// The default value is 0 (scan errors are detected but not corrected). Always use quality level 200 (ECC 200).
	Quality int

	// The number of columns to encode. For ECC 200 bar codes, even numbers between 1 and 144 may be used.
	// This parameter can be used to control the bar code width. The default value depends on the amount of data encoded.
	Columns int

	// The number of rows to encode. For ECC 200 bar codes, even numbers between 1 and 144 may be used.
	// This parameter can be used to control the bar code height. The default value depends on the amount of data encoded.
	Rows int

	// The type of data that needs to be encoded.
	// Valid values are 1, 2, 3, 4, 5 and 6. The default value is 6.
	// This parameter is ignored for ECC 200 bar codes (the recommended quality level).
	Format int

	// The escape character used to escape control sequences in the field data. The default value is "~" (tilde).
	Escape byte

	// The desired aspect ratio, if any. Valid values are 1 (square) and 2 (rectangular).
	Ratio DatamatrixRatio
}

type BarcodeDatamatrixWithData struct {
	BarcodeDatamatrix
	Position LabelPosition
	Data     string
}
