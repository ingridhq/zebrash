package encoder

type DefaultPlacement struct {
	codewords []byte
	numrows   int
	numcols   int
	bits      []int8
}

func NewDefaultPlacement(codewords []byte, numcols, numrows int) *DefaultPlacement {
	p := &DefaultPlacement{
		codewords: codewords,
		numcols:   numcols,
		numrows:   numrows,
		bits:      make([]int8, numcols*numrows),
	}
	for i := range p.bits {
		p.bits[i] = -1
	}
	return p
}

func (dp *DefaultPlacement) GetBit(col, row int) bool {
	return dp.bits[row*dp.numcols+col] == 1
}

func (dp *DefaultPlacement) setBit(col, row int, bit bool) {
	b := int8(0)
	if bit {
		b = 1
	}
	dp.bits[row*dp.numcols+col] = b
}

func (dp *DefaultPlacement) hasBit(col, row int) bool {
	return dp.bits[row*dp.numcols+col] >= 0
}

func (dp *DefaultPlacement) Place() {
	pos := 0
	row := 4
	col := 0

	for {
		// repeatedly first check for one of the special corner cases, then...
		if (row == dp.numrows) && (col == 0) {
			dp.corner1(pos)
			pos++
		}
		if (row == dp.numrows-2) && (col == 0) && ((dp.numcols % 4) != 0) {
			dp.corner2(pos)
			pos++
		}
		if (row == dp.numrows-2) && (col == 0) && (dp.numcols%8 == 4) {
			dp.corner3(pos)
			pos++
		}
		if (row == dp.numrows+4) && (col == 2) && ((dp.numcols % 8) == 0) {
			dp.corner4(pos)
			pos++
		}
		// sweep upward diagonally, inserting successive characters...
		for {
			if (row < dp.numrows) && (col >= 0) && !dp.hasBit(col, row) {
				dp.utah(row, col, pos)
				pos++
			}
			row -= 2
			col += 2
			if row < 0 || (col >= dp.numcols) {
				break
			}
		}
		row++
		col += 3

		// and then sweep downward diagonally, inserting successive characters, ...
		for {
			if (row >= 0) && (col < dp.numcols) && !dp.hasBit(col, row) {
				dp.utah(row, col, pos)
				pos++
			}
			row += 2
			col -= 2
			if row >= dp.numrows || col < 0 {
				break
			}
		}
		row += 3
		col++

		// ...until the entire array is scanned
		if row >= dp.numrows && col >= dp.numcols {
			break
		}
	}

	// Lastly, if the lower righthand corner is untouched, fill in fixed pattern
	if !dp.hasBit(dp.numcols-1, dp.numrows-1) {
		dp.setBit(dp.numcols-1, dp.numrows-1, true)
		dp.setBit(dp.numcols-2, dp.numrows-2, true)
	}
}

func (dp *DefaultPlacement) module(row, col, pos, bit int) {
	if row < 0 {
		row += dp.numrows
		col += 4 - ((dp.numrows + 4) % 8)
	}
	if col < 0 {
		col += dp.numcols
		row += 4 - ((dp.numcols + 4) % 8)
	}
	// Note the conversion:
	v := dp.codewords[pos]
	v &= 1 << uint(8-bit)
	dp.setBit(col, row, v != 0)
}

// utah Places the 8 bits of a utah-shaped symbol character in ECC200.
//
// @param row the row
// @param col the column
// @param pos character position
func (dp *DefaultPlacement) utah(row, col, pos int) {
	dp.module(row-2, col-2, pos, 1)
	dp.module(row-2, col-1, pos, 2)
	dp.module(row-1, col-2, pos, 3)
	dp.module(row-1, col-1, pos, 4)
	dp.module(row-1, col, pos, 5)
	dp.module(row, col-2, pos, 6)
	dp.module(row, col-1, pos, 7)
	dp.module(row, col, pos, 8)
}

func (dp *DefaultPlacement) corner1(pos int) {
	dp.module(dp.numrows-1, 0, pos, 1)
	dp.module(dp.numrows-1, 1, pos, 2)
	dp.module(dp.numrows-1, 2, pos, 3)
	dp.module(0, dp.numcols-2, pos, 4)
	dp.module(0, dp.numcols-1, pos, 5)
	dp.module(1, dp.numcols-1, pos, 6)
	dp.module(2, dp.numcols-1, pos, 7)
	dp.module(3, dp.numcols-1, pos, 8)
}

func (dp *DefaultPlacement) corner2(pos int) {
	dp.module(dp.numrows-3, 0, pos, 1)
	dp.module(dp.numrows-2, 0, pos, 2)
	dp.module(dp.numrows-1, 0, pos, 3)
	dp.module(0, dp.numcols-4, pos, 4)
	dp.module(0, dp.numcols-3, pos, 5)
	dp.module(0, dp.numcols-2, pos, 6)
	dp.module(0, dp.numcols-1, pos, 7)
	dp.module(1, dp.numcols-1, pos, 8)
}

func (dp *DefaultPlacement) corner3(pos int) {
	dp.module(dp.numrows-3, 0, pos, 1)
	dp.module(dp.numrows-2, 0, pos, 2)
	dp.module(dp.numrows-1, 0, pos, 3)
	dp.module(0, dp.numcols-2, pos, 4)
	dp.module(0, dp.numcols-1, pos, 5)
	dp.module(1, dp.numcols-1, pos, 6)
	dp.module(2, dp.numcols-1, pos, 7)
	dp.module(3, dp.numcols-1, pos, 8)
}

func (dp *DefaultPlacement) corner4(pos int) {
	dp.module(dp.numrows-1, 0, pos, 1)
	dp.module(dp.numrows-1, dp.numcols-1, pos, 2)
	dp.module(0, dp.numcols-3, pos, 3)
	dp.module(0, dp.numcols-2, pos, 4)
	dp.module(0, dp.numcols-1, pos, 5)
	dp.module(1, dp.numcols-3, pos, 6)
	dp.module(1, dp.numcols-2, pos, 7)
	dp.module(1, dp.numcols-1, pos, 8)
}
