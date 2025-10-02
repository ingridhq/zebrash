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

func (m *ByteMatrix) GetHeight() int {
	return m.height
}

func (m *ByteMatrix) GetWidth() int {
	return m.width
}

func (m *ByteMatrix) Get(x, y int) int8 {
	return m.bytes[y][x]
}

func (m *ByteMatrix) GetArray() [][]int8 {
	return m.bytes
}

func (m *ByteMatrix) Set(x, y int, value int8) {
	m.bytes[y][x] = value
}

func (m *ByteMatrix) SetBool(x, y int, value bool) {
	if value {
		m.bytes[y][x] = 1
	} else {
		m.bytes[y][x] = 0
	}
}

func (m *ByteMatrix) Clear(value int8) {
	for y := range m.bytes {
		for x := range m.bytes[y] {
			m.bytes[y][x] = value
		}
	}
}

func (m *ByteMatrix) String() string {
	result := make([]byte, 0, 2*(m.width+1)*m.height)
	for _, row := range m.bytes {
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
