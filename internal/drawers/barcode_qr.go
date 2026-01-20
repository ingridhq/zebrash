package drawers

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/DawidBury/zebrash/drawers"
	"github.com/DawidBury/zebrash/internal/barcodes/qrcode"
	"github.com/DawidBury/zebrash/internal/barcodes/qrcode/encoder"
	"github.com/DawidBury/zebrash/internal/elements"
	"github.com/DawidBury/zebrash/internal/images"
)

func NewBarcodeQrDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ drawers.DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeQrWithData)
			if !ok {
				return nil
			}

			inputData, ec, _, err := barcode.GetInputData()
			if err != nil {
				return err
			}

			img, err := qrcode.Encode(inputData, 1, 1, mapQrErrorCorrectionLevel(ec), encoder.Options{
				QuietZone: 0,
			})
			if err != nil {
				return fmt.Errorf("failed to encode qr barcode: %w", err)
			}

			scaledImg := images.NewScaled(img, barcode.Magnification, barcode.Magnification)

			pos := barcode.Position
			// Weird behavior when height set by ^BY shifts QR code vertically
			// Only works when  CalculateFromBottom is set to false (position is set via ^FO and not ^FT)
			if !pos.CalculateFromBottom {
				pos.Y += barcode.Height
			} else {
				// TODO: figure out the proper formula for ftOffset; it seems to depend on the QR code version.
				ftOffset := barcode.Magnification * 7
				pos.Y = max(pos.Y-scaledImg.Bounds().Dy(), 0) - ftOffset
			}

			gCtx.DrawImage(scaledImg, pos.X, pos.Y)

			return nil
		},
	}
}

func mapQrErrorCorrectionLevel(ec elements.QrErrorCorrectionLevel) encoder.ErrorCorrectionLevel {
	switch ec {
	case elements.QrErrorCorrectionL:
		return encoder.ErrorCorrectionLevel_L
	case elements.QrErrorCorrectionQ:
		return encoder.ErrorCorrectionLevel_Q
	case elements.QrErrorCorrectionH:
		return encoder.ErrorCorrectionLevel_H
	default:
		return encoder.ErrorCorrectionLevel_M
	}
}
