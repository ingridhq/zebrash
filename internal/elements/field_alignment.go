package elements

// The origin alignment to use.
// Not the same as text alignment for text block
type FieldAlignment int

const (
	FieldAlignmentLeft  FieldAlignment = 0
	FieldAlignmentRight FieldAlignment = 1
	// automatic alignment based on the direction of the field data text
	FieldAlignmentAuto FieldAlignment = 2
)
