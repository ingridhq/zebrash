package utils

import (
	"image"
	"image/color"
	"math"

	"github.com/ingridhq/zebrash/internal/images"
)

func (resBits *BitList) ToImageWithBarsRatio(width, height int, widthRatio float64) image.Image {
	widthRatio = max(min(3, widthRatio), 2)
	wideBarWidth := int(math.Round(widthRatio * float64(width)))

	barsList := resBits.ToWideNarrowList(wideBarWidth, width)
	img := image.NewRGBA(image.Rect(0, 0, barsList.GetTotalWidth(), 1))

	px := 0
	for i, v := range barsList.Data {
		for range barsList.GetBarWidth(i) {
			img.Set(px, 0, getColor(v[1]))
			px++
		}
	}

	return images.NewScaled1DHeight(img, height)
}

func getColor(b bool) color.RGBA {
	if b {
		return images.ColorBlack
	}

	return images.ColorTransparent
}
