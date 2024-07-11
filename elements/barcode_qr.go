package elements

import (
	"fmt"
)

type BarcodeQr struct {
	// The bar code magnification to use.
	// Any number between 1 and 10 may be used.
	// The default value depends on the print density being used.
	Magnification int
}

type BarcodeQrWithData struct {
	BarcodeQr
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

func (barcode *BarcodeQrWithData) GetInputData() (string, QrErrorCorrectionLevel, error) {
	if len(barcode.Data) < 4 {
		return "", 0, fmt.Errorf("invalid qr barcode data")
	}

	return barcode.Data[3:], QrErrorCorrectionLevel(barcode.Data[0]), nil
}
