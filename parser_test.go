package zebrash

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"

	"github.com/ingridhq/zebrash/drawers"
)

var (
	writeImageDiff bool
	updateDstFiles bool
)

func init() {
	flag.BoolVar(&writeImageDiff, "write-image-diff", true, "save PNG output diff for each test case")
	flag.BoolVar(&updateDstFiles, "update-dst-files", false, "replace destination PNGs with the ones produced during test run")
}

func TestDrawLabelAsPng(t *testing.T) {
	testCases := []struct {
		name           string
		srcPath        string
		dstPath        string
		labelIdx       int
		widthMm        float64
		heightMm       float64
		enableInverted bool
	}{
		{
			name:    "Amazon label",
			srcPath: "amazon.zpl",
			dstPath: "amazon.png",
		},
		{
			name:    "UPS label",
			srcPath: "ups.zpl",
			dstPath: "ups.png",
		},
		{
			name:           "UPS label with orientation inversion enabled",
			srcPath:        "ups.zpl",
			dstPath:        "ups_inverted.png",
			enableInverted: true,
		},
		{
			name:    "UPS Surepost label",
			srcPath: "ups_surepost.zpl",
			dstPath: "ups_surepost.png",
		},
		{
			name:    "USPS label",
			srcPath: "usps.zpl",
			dstPath: "usps.png",
		},
		{
			name:    "Fedex label",
			srcPath: "fedex.zpl",
			dstPath: "fedex.png",
		},
		{
			name:    "Label converted from PDF via zplgrf",
			srcPath: "bstc.zpl",
			dstPath: "bstc.png",
		},
		{
			name:    "Postnord DPD label",
			srcPath: "pnldpd.zpl",
			dstPath: "pnldpd.png",
		},
		{
			name:     "Postnord DPD label Page 2",
			srcPath:  "pnldpd.zpl",
			dstPath:  "pnldpd_page_2.png",
			labelIdx: 1,
		},
		{
			name:    "DPD Poland",
			srcPath: "dpdpl.zpl",
			dstPath: "dpdpl.png",
		},
		{
			name:    "DHL Paket",
			srcPath: "dhlpaket.zpl",
			dstPath: "dhlpaket.png",
		},
		{
			name:    "DHL Parcel UK",
			srcPath: "dhlparceluk.zpl",
			dstPath: "dhlparceluk.png",
		},
		{
			name:    "DHL ECommerce TR",
			srcPath: "dhlecommercetr.zpl",
			dstPath: "dhlecommercetr.png",
		},
		{
			name:    "ICA Paket",
			srcPath: "icapaket.zpl",
			dstPath: "icapaket.png",
		},
		{
			name:    "DB Schenker",
			srcPath: "dbs.zpl",
			dstPath: "dbs.png",
		},
		{
			name:    "GLS DK Return",
			srcPath: "glsdk_return.zpl",
			dstPath: "glsdk_return.png",
		},
		{
			name:    "GLS MyGLS CZ API",
			srcPath: "glscz.zpl",
			dstPath: "glscz.png",
		},
		{
			name:    "JCPenney",
			srcPath: "jcpenney.zpl",
			dstPath: "jcpenney.png",
		},
		{
			name:    "Kmart",
			srcPath: "kmart.zpl",
			dstPath: "kmart.png",
		},
		{
			name:    "Labelary",
			srcPath: "labelary.zpl",
			dstPath: "labelary.png",
		},
		{
			name:    "Polish Post Pocztex",
			srcPath: "pocztex.zpl",
			dstPath: "pocztex.png",
		},
		{
			name:    "Porterbuddy",
			srcPath: "porterbuddy.zpl",
			dstPath: "porterbuddy.png",
		},
		{
			name:    "Bring Posten",
			srcPath: "posten.zpl",
			dstPath: "posten.png",
		},
		{
			name:    "Return With QR Code",
			srcPath: "return_qrcode.zpl",
			dstPath: "return_qrcode.png",
		},
		{
			name:    "Reverse print",
			srcPath: "reverse.zpl",
			dstPath: "reverse.png",
		},
		{
			name:    "Reverse print with QR code",
			srcPath: "reverse_qr.zpl",
			dstPath: "reverse_qr.png",
		},
		{
			name:    "Swiss Post",
			srcPath: "swisspost.zpl",
			dstPath: "swisspost.png",
		},
		{
			name:    "Text Field Typeset Normal rotation",
			srcPath: "text_ft_n.zpl",
			dstPath: "text_ft_n.png",
		},
		{
			name:    "Text Field Typeset 90 degrees rotation",
			srcPath: "text_ft_r.zpl",
			dstPath: "text_ft_r.png",
		},
		{
			name:    "Text Field Typeset 180 degrees rotation",
			srcPath: "text_ft_i.zpl",
			dstPath: "text_ft_i.png",
		},
		{
			name:    "Text Field Typeset 270 degrees rotation",
			srcPath: "text_ft_b.zpl",
			dstPath: "text_ft_b.png",
		},
		{
			name:    "Text Field Typeset automatic position",
			srcPath: "text_ft_auto_pos.zpl",
			dstPath: "text_ft_auto_pos.png",
		},
		{
			name:    "Text Field Origin Normal rotation",
			srcPath: "text_fo_n.zpl",
			dstPath: "text_fo_n.png",
		},
		{
			name:    "Text Field Origin 90 degrees rotation",
			srcPath: "text_fo_r.zpl",
			dstPath: "text_fo_r.png",
		},
		{
			name:    "Text Field Origin 180 degrees rotation",
			srcPath: "text_fo_i.zpl",
			dstPath: "text_fo_i.png",
		},
		{
			name:    "Text Field Origin 270 degrees rotation",
			srcPath: "text_fo_b.zpl",
			dstPath: "text_fo_b.png",
		},
		{
			name:    "Text Field Multi-line",
			srcPath: "text_multiline.zpl",
			dstPath: "text_multiline.png",
		},
		{
			name:     "Text with non-existing font that falls back to default",
			srcPath:  "text_fallback_default.zpl",
			dstPath:  "text_fallback_default.png",
			widthMm:  160,
			heightMm: 230,
		},
		{
			name:    "Barcode128 with 'line' and 'line above' set",
			srcPath: "barcode128_line_above.zpl",
			dstPath: "barcode128_line_above.png",
		},
		{
			name:    "Barcode128 with 'line' set",
			srcPath: "barcode128_line.zpl",
			dstPath: "barcode128_line.png",
		},
		{
			name:    "Barcode128 mode a",
			srcPath: "barcode128_mode_a.zpl",
			dstPath: "barcode128_mode_a.png",
		},
		{
			name:    "Barcode128 mode d",
			srcPath: "barcode128_mode_d.zpl",
			dstPath: "barcode128_mode_d.png",
		},
		{
			name:    "Barcode128 mode u",
			srcPath: "barcode128_mode_u.zpl",
			dstPath: "barcode128_mode_u.png",
		},
		{
			name:    "Barcode128 mode n",
			srcPath: "barcode128_mode_n.zpl",
			dstPath: "barcode128_mode_n.png",
		},
		{
			name:    "Barcode128 mode n switching from subset C to B to A",
			srcPath: "barcode128_mode_n_cba_sets.zpl",
			dstPath: "barcode128_mode_n_cba_sets.png",
		},
		{
			name:    "Barcode128 default width set",
			srcPath: "barcode128_default_width.zpl",
			dstPath: "barcode128_default_width.png",
		},
		{
			name:    "Barcode128 180 degrees rotated",
			srcPath: "barcode128_rotated.zpl",
			dstPath: "barcode128_rotated.png",
		},
		{
			name:    "Graphic box normal",
			srcPath: "gb_normal.zpl",
			dstPath: "gb_normal.png",
		},
		{
			name:    "Graphic box with 0 height",
			srcPath: "gb_0_height.zpl",
			dstPath: "gb_0_height.png",
		},
		{
			name:    "Graphic box with 0 width",
			srcPath: "gb_0_width.zpl",
			dstPath: "gb_0_width.png",
		},
		{
			name:    "Graphic box with corner-rounding",
			srcPath: "gb_rounded.zpl",
			dstPath: "gb_rounded.png",
		},
		{
			name:    "Text encodings 0-13",
			srcPath: "encodings_013.zpl",
			dstPath: "encodings_013.png",
		},
		{
			name:    "Aztec barcode error correction",
			srcPath: "aztec_ec.zpl",
			dstPath: "aztec_ec.png",
		},
		{
			name:    "QR code with offset",
			srcPath: "qr_code_offset.zpl",
			dstPath: "qr_code_offset.png",
		},
		{
			name:    "QR code with field typeset and manual binary mode",
			srcPath: "qr_code_ft_manual.zpl",
			dstPath: "qr_code_ft_manual.png",
		},
		{
			name:    "EAN-13 barcode test",
			srcPath: "ean13.zpl",
			dstPath: "ean13.png",
		},
		{
			name:    "Templating",
			srcPath: "templating.zpl",
			dstPath: "templating.png",
		},
		{
			name:    "Graphic symbol",
			srcPath: "gs.zpl",
			dstPath: "gs.png",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			fullSrcPath := "./testdata/" + tC.srcPath
			fullDstPath := "./testdata/" + tC.dstPath
			fullDiffPath := "./testdata/diff/" + tC.dstPath
			file := mustReadFile(fullSrcPath, t)

			parser := NewParser()

			res, err := parser.Parse(file)
			if err != nil {
				t.Fatal(err)
			}

			if len(res) == 0 {
				t.Fatal("no labels in the response")
			}

			drawer := NewDrawer()

			buff := new(bytes.Buffer)
			err = drawer.DrawLabelAsPng(res[tC.labelIdx], buff, drawers.DrawerOptions{
				LabelWidthMm:         tC.widthMm,
				LabelHeightMm:        tC.heightMm,
				EnableInvertedLabels: tC.enableInverted,
			})
			if err != nil {
				t.Fatal(err)
			}

			if updateDstFiles {
				mustWriteFile(fullDstPath, buff.Bytes(), t)
			}

			gotImg := mustDecodePng(buff.Bytes(), t)
			wantImg := mustDecodePng(mustReadFile(fullDstPath, t), t)

			compareImages(gotImg, wantImg, fullDiffPath, t)
		})
	}
}

func compareImages(got, want image.Image, fullDiffPath string, t *testing.T) {
	if !got.Bounds().Eq(want.Bounds()) {
		t.Fatalf("Image bounds differ: got=%v want=%v", got.Bounds(), want.Bounds())
	}

	gotGray, ok := got.(*image.Gray)
	if !ok {
		t.Fatalf("Got is not grayscale image")
	}

	wantGray, ok := want.(*image.Gray)
	if !ok {
		t.Fatalf("Want is not grayscale image")
	}

	bounds := got.Bounds()
	diffImg := image.NewRGBA(bounds)

	const maxReportedSampleMismatches = 5
	var sampleMismatches []string

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gotV := gotGray.GrayAt(x, y)
			wantV := wantGray.GrayAt(x, y)

			if gotV == wantV {
				diffImg.Set(x, y, color.RGBA{R: gotV.Y, G: gotV.Y, B: gotV.Y, A: 255})
				continue
			}

			if len(sampleMismatches) < maxReportedSampleMismatches {
				sampleMismatches = append(sampleMismatches, fmt.Sprintf("(x: %d, y: %d): got=GRAY(%d) want=GRAY(%d)", x, y, gotV.Y, wantV.Y))
			}

			if wantV.Y > gotV.Y {
				diffImg.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			} else {
				diffImg.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
			}
		}
	}

	if len(sampleMismatches) == 0 {
		return
	}

	t.Errorf("First %d sample mismatches: \n%s", maxReportedSampleMismatches, strings.Join(sampleMismatches, "\n"))

	if !writeImageDiff {
		return
	}

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, diffImg); err != nil {
		t.Fatalf("Failed to encode diff image: %v", err)
	}

	mustWriteFile(fullDiffPath, buff.Bytes(), t)
	t.Logf("Diff image is saved to %s", fullDiffPath)
}

func mustDecodePng(data []byte, t *testing.T) image.Image {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode png: %v", err)
	}

	return img
}

func mustReadFile(path string, t *testing.T) []byte {
	res, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	return res
}

func mustWriteFile(path string, data []byte, t *testing.T) {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write file %s: %v", path, err)
	}
}
