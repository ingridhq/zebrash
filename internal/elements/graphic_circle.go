package elements

type GraphicCircle struct {
	ReversePrint

	Position LabelPosition

	// The diameter of the circle, in dots.
	// Any number between 3 and 4,095 may be used.
	// The default value is 3.
	CircleDiameter int

	// The line thickness to use to draw the circle, in dots.
	// Any number between 1 and 4,095 may be used.
	// The default value is 1.
	BorderThickness int

	// The line color to use to draw the circle.
	// Valid values are B (black) and W (white).
	// The default value is B (black).
	LineColor LineColor
}
