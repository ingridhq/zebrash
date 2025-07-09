package elements

type FieldOrientation int

const (
	// no rotation
	FieldOrientationNormal FieldOrientation = 0
	// rotate 90° clockwise
	FieldOrientation90 FieldOrientation = 1
	// rotate 180° clockwise
	FieldOrientation180 FieldOrientation = 2
	// rotate 270° clockwise
	FieldOrientation270 FieldOrientation = 3
)

func (v FieldOrientation) GetDegrees() float64 {
	switch v {
	case FieldOrientation90:
		return 90
	case FieldOrientation180:
		return 180
	case FieldOrientation270:
		return 270
	default:
		return 0
	}
}
