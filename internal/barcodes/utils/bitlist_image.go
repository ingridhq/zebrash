package utils

import (
	"image"
	"image/color"
	"math"

	"github.com/ingridhq/zebrash/internal/images"
)

type WideNarrowList struct {
	Data           [][2]bool
	narrowBarWidth int
	wideBarWidth   int
}

func (resBits *BitList) ToWideNarrowList(wideBarWidth, narrowBarWidth int) *WideNarrowList {
	var res [][2]bool

	prevB := resBits.GetBit(0)
	c := 0

	for i := range resBits.Len() {
		b := resBits.GetBit(i)
		if prevB == b {
			c++
			continue
		}

		res = append(res, [2]bool{c > 1, prevB})

		prevB = b
		c = 1
	}

	res = append(res, [2]bool{c > 1, prevB})

	return &WideNarrowList{
		Data:           res,
		wideBarWidth:   wideBarWidth,
		narrowBarWidth: narrowBarWidth,
	}
}

func (list *WideNarrowList) GetTotalWidth() int {
	width := 0

	for i := range list.Data {
		width += list.GetBarWidth(i)
	}

	return max(1, width)
}

func (list *WideNarrowList) GetBarWidth(idx int) int {
	if list.Data[idx][0] {
		return list.wideBarWidth
	}

	return list.narrowBarWidth
}

func (resBits *BitList) ToImage(width, height int, widthRatio float64) image.Image {
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
