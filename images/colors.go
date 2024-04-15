package images

import "image/color"

var (
	ColorBlack       = color.RGBA{A: 255}
	ColorWhite       = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	ColorTransparent = color.RGBA{A: 0}
)
