package elements

type LabelPosition struct {
	X                   int
	Y                   int
	CalculateFromBottom bool
	AutomaticPosition   bool
}

func (p LabelPosition) Add(pos LabelPosition) LabelPosition {
	return LabelPosition{
		X:                   p.X + pos.X,
		Y:                   p.Y + pos.Y,
		CalculateFromBottom: p.CalculateFromBottom,
		AutomaticPosition:   p.AutomaticPosition,
	}
}
