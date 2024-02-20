package elements

type FontInfo struct {
	Name        string
	Width       int
	Height      int
	Orientation FieldOrientation
}

func (f FontInfo) GetSize() int {
	if f.Height > 0 {
		return f.Height
	}

	return f.Width
}
