package pdf417

const (
	minCols = 2
	maxCols = 30
	maxRows = 30
	minRows = 2
)

func calculateNumberOfRows(m, k, c int) int {
	r := ((m + 1 + k) / c) + 1
	if c*r >= (m + 1 + k + c) {
		r--
	}
	return r
}

func calcDimensions(targetColumns, dataWords, eccWords int) (cols, rows int) {
	cols = 0
	rows = 0

	for c := max(minCols, targetColumns); c <= maxCols; c++ {
		r := calculateNumberOfRows(dataWords, eccWords, c)

		if r < minRows {
			break
		}

		if r > maxRows {
			continue
		}

		cols = c
		rows = r

		break
	}

	if rows == 0 {
		r := calculateNumberOfRows(dataWords, eccWords, minCols)
		if r < minRows {
			rows = minRows
			cols = minCols
		}
	}

	return
}
