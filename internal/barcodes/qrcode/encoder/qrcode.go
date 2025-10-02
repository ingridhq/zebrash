package encoder

import (
	"strconv"
)

const QRCode_NUM_MASK_PATERNS = 8

type QRCode struct {
	mode        *Mode
	ecLevel     ErrorCorrectionLevel
	version     *Version
	maskPattern int
	matrix      *ByteMatrix
}

func NewQRCode() *QRCode {
	return &QRCode{
		maskPattern: -1,
	}
}

func (code *QRCode) GetMode() *Mode {
	return code.mode
}

func (code *QRCode) GetECLevel() ErrorCorrectionLevel {
	return code.ecLevel
}

func (code *QRCode) GetVersion() *Version {
	return code.version
}

func (code *QRCode) GetMaskPattern() int {
	return code.maskPattern
}

func (code *QRCode) GetMatrix() *ByteMatrix {
	return code.matrix
}

func (code *QRCode) String() string {
	result := make([]byte, 0, 200)
	result = append(result, "<<\n"...)
	result = append(result, " mode: "...)
	result = append(result, code.mode.String()...)
	result = append(result, "\n ecLevel: "...)
	result = append(result, code.ecLevel.String()...)
	result = append(result, "\n version: "...)
	result = append(result, code.version.String()...)
	result = append(result, "\n maskPattern: "...)
	result = append(result, strconv.Itoa(code.maskPattern)...)
	if code.matrix == nil {
		result = append(result, "\n matrix: nil\n"...)
	} else {
		result = append(result, "\n matrix:\n"...)
		result = append(result, code.matrix.String()...)
	}
	result = append(result, ">>\n"...)
	return string(result)
}

func (code *QRCode) SetMode(value *Mode) {
	code.mode = value
}

func (code *QRCode) SetECLevel(value ErrorCorrectionLevel) {
	code.ecLevel = value
}

func (code *QRCode) SetVersion(value *Version) {
	code.version = value
}

func (code *QRCode) SetMaskPattern(value int) {
	code.maskPattern = value
}

func (code *QRCode) SetMatrix(value *ByteMatrix) {
	code.matrix = value
}

func QRCode_IsValidMaskPattern(maskPattern int) bool {
	return maskPattern >= 0 && maskPattern < QRCode_NUM_MASK_PATERNS
}
