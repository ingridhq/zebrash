package utils

import (
	"image"
	"image/color"

	"github.com/ingridhq/zebrash/images"
)

// Allocate at least 10 pixels for each bar in the barcode
// then we downscale the result image to 1/10th
// It is done because widthRatio increases with 0.1 steps
const minBarWidth = 10

func (resBits *BitList) ToImage(width, height int, widthRatio float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, minBarWidth*resBits.Len(), 1))

	widthRatio = max(min(3, widthRatio), 2)
	wideBarWidth := int(widthRatio * minBarWidth)
	prevB := resBits.GetBit(0)
	px := 0
	c := 0

	for i := range resBits.Len() {
		b := resBits.GetBit(i)
		if prevB == b {
			c++
			continue
		}

		for range getBarWidth(c > 1, wideBarWidth) {
			img.Set(px, 0, getColor(prevB))
			px++
		}

		prevB = b
		c = 1
	}

	for range getBarWidth(c > 1, wideBarWidth) {
		img.Set(px, 0, getColor(prevB))
		px++
	}

	return images.NewScaledFloat(img, float64(width)*0.1, float64(height))
}

func getColor(b bool) color.RGBA {
	if b {
		return images.ColorBlack
	}

	return images.ColorTransparent
}

func getBarWidth(isWide bool, wideBarWidth int) int {
	if isWide {
		return wideBarWidth
	}

	return minBarWidth
}
