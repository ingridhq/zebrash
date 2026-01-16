package utils

type WideNarrowList struct {
	Data           [][2]bool
	narrowBarWidth int
	wideBarWidth   int
}

// Converts BitList to useful representation for barcodes that have narrow / wide bars with variable bar ratio
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
