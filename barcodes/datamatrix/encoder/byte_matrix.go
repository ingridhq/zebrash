package encoder

type ByteMatrix struct {
	bytes  [][]int8
	width  int
	height int
}

func NewByteMatrix(width, height int) *ByteMatrix {
	bytes := make([][]int8, height)
	for i := 0; i < height; i++ {
		bytes[i] = make([]int8, width)
	}
	return &ByteMatrix{bytes, width, height}
}

func (matrix *ByteMatrix) GetHeight() int {
	return matrix.height
}

func (matrix *ByteMatrix) GetWidth() int {
	return matrix.width
}

func (matrix *ByteMatrix) Get(x, y int) int8 {
	return matrix.bytes[y][x]
}

func (matrix *ByteMatrix) GetArray() [][]int8 {
	return matrix.bytes
}

func (matrix *ByteMatrix) Set(x, y int, value int8) {
	matrix.bytes[y][x] = value
}

func (matrix *ByteMatrix) SetBool(x, y int, value bool) {
	if value {
		matrix.bytes[y][x] = 1
	} else {
		matrix.bytes[y][x] = 0
	}
}

func (matrix *ByteMatrix) Clear(value int8) {
	for y := range matrix.bytes {
		for x := range matrix.bytes[y] {
			matrix.bytes[y][x] = value
		}
	}
}

func (matrix *ByteMatrix) String() string {
	result := make([]byte, 0, 2*(matrix.width+1)*matrix.height)
	for _, row := range matrix.bytes {
		for _, b := range row {
			switch b {
			case 0:
				result = append(result, " 0"...)
			case 1:
				result = append(result, " 1"...)
			default:
				result = append(result, "  "...)
			}
		}
		result = append(result, '\n')
	}
	return string(result)
}
