package elements

type BarcodeDimensions struct {
	ModuleWidth int
	Height      int
	// widthRatio: The default ratio between wide bars and narrow bars. Any decimal number between 2 and 3 may be used.
	// The number must be a multiple of 0.1 (i.e. 2.0, 2.1, 2.2, 2.3, ... , 2.9, 3.0).
	// Larger numbers generally result in fewer bar code scan failures.
	// The default value is the previously configured value, or 3 if no value has been set.
	WidthRatio float64
}
