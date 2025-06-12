package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/images"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
)

func NewBarcodeQrDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeQrWithData)
			if !ok {
				return nil
			}

			enc := qrcode.NewQRCodeWriter()

			inputData, ec, err := barcode.GetInputData()
			if err != nil {
				return err
			}

			img, err := enc.Encode(inputData, gozxing.BarcodeFormat_QR_CODE, 1, 1, map[gozxing.EncodeHintType]any{
				gozxing.EncodeHintType_ERROR_CORRECTION: mapQrErrorCorrectionLevel(ec),
				gozxing.EncodeHintType_MARGIN:           0,
			})
			if err != nil {
				return fmt.Errorf("failed to encode qr barcode: %w", err)
			}

			scaledImg := images.NewScaled(img, barcode.Magnification, barcode.Magnification)
			pos := adjustImageTypeSetPosition(scaledImg, barcode.Position, elements.FieldOrientationNormal)

			gCtx.DrawImage(scaledImg, pos.X, pos.Y)

			return nil
		},
	}
}

func mapQrErrorCorrectionLevel(ec elements.QrErrorCorrectionLevel) decoder.ErrorCorrectionLevel {
	switch ec {
	case elements.QrErrorCorrectionL:
		return decoder.ErrorCorrectionLevel_L
	case elements.QrErrorCorrectionQ:
		return decoder.ErrorCorrectionLevel_Q
	case elements.QrErrorCorrectionH:
		return decoder.ErrorCorrectionLevel_H
	default:
		return decoder.ErrorCorrectionLevel_M
	}
}
