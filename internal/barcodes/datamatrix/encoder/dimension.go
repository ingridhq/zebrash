package encoder

type Dimension struct {
	Width  int
	Height int
}

func NewDimension(width, height int) *Dimension {
	return &Dimension{
		Width:  width,
		Height: height,
	}
}

func (d *Dimension) GetWidth() int {
	return d.Width
}

func (d *Dimension) GetHeight() int {
	return d.Height
}
