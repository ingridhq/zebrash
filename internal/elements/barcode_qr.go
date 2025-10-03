package elements

import (
	"fmt"
	"strconv"
)

type BarcodeQr struct {
	// The bar code magnification to use.
	// Any number between 1 and 10 may be used.
	// The default value depends on the print density being used.
	Magnification int
}

type BarcodeQrWithData struct {
	BarcodeQr
	Height   int
	Position LabelPosition
	Data     string
}

type QrErrorCorrectionLevel byte

const (
	QrErrorCorrectionH QrErrorCorrectionLevel = 'H'
	QrErrorCorrectionQ QrErrorCorrectionLevel = 'Q'
	QrErrorCorrectionM QrErrorCorrectionLevel = 'M'
	QrErrorCorrectionL QrErrorCorrectionLevel = 'L'
)

type QrCharacterMode byte

const (
	QrCharacterModeAutomatic    QrCharacterMode = 0
	QrCharacterModeBinary       QrCharacterMode = 'B'
	QrCharacterModeNumeric      QrCharacterMode = 'N'
	QrCharacterModeAlphanumeric QrCharacterMode = 'A'
	QrCharacterModeKanji        QrCharacterMode = 'K'
)

func (barcode *BarcodeQrWithData) GetInputData() (string, QrErrorCorrectionLevel, QrCharacterMode, error) {
	if len(barcode.Data) < 4 {
		return "", 0, 0, fmt.Errorf("invalid qr barcode data")
	}

	data := barcode.Data[3:]
	mode := QrCharacterModeAutomatic
	level := QrErrorCorrectionLevel(barcode.Data[0])

	// First character of the data in manual mode defines character mode
	if barcode.Data[1] == 'M' && len(data) > 0 {
		mode = QrCharacterMode(data[0])
		data = data[1:]
	}

	if mode != QrCharacterModeBinary {
		return data, level, mode, nil
	}

	if len(data) < 5 {
		return "", 0, 0, fmt.Errorf("invalid qr barcode byte mode data")
	}

	dataLen, err := strconv.Atoi(data[0:4])
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid qr barcode byte mode data length: %w", err)
	}

	data = data[4:]
	dataLen = min(len(data), dataLen)

	return data[:dataLen], level, mode, nil
}
