package elements

import (
	internalElements "github.com/ingridhq/zebrash/internal/elements"
)

type LabelInfo struct {
	PrintWidth  int
	Orientation internalElements.FieldOrientation
	Elements    []any
}
