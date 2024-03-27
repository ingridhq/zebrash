package drawers

type DrawerOptions struct {
	LabelWidthMm  float64
	LabelHeightMm float64
	Dpmm          int
}

func (d DrawerOptions) WithDefaults() DrawerOptions {
	res := d

	// by default produce 4x8 inches 203 dpi label

	if res.LabelWidthMm == 0 {
		res.LabelWidthMm = 101.6
	}

	if res.LabelHeightMm == 0 {
		res.LabelHeightMm = 203.2
	}

	if res.Dpmm == 0 {
		res.Dpmm = 8
	}

	return res
}
