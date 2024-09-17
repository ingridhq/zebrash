package zebrash

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ingridhq/zebrash/drawers"
)

func TestDrawLabelAsPng(t *testing.T) {
	testCases := []struct {
		name     string
		srcPath  string
		dstPath  string
		labelIdx int
	}{
		{
			name:    "UPS label",
			srcPath: "ups.zpl",
			dstPath: "ups.png",
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
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			file, err := os.ReadFile("./testdata/" + tC.srcPath)
			if err != nil {
				t.Fatal(err)
			}

			parser := NewParser()

			res, err := parser.Parse(file)
			if err != nil {
				t.Fatal(err)
			}

			var buff bytes.Buffer

			drawer := NewDrawer()

			if len(res) == 0 {
				t.Fatal("no labels in the response")
			}

			err = drawer.DrawLabelAsPng(res[tC.labelIdx], &buff, drawers.DrawerOptions{})
			if err != nil {
				t.Fatal(err)
			}

			expectedPng, err := os.ReadFile("./testdata/" + tC.dstPath)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(buff.Bytes(), expectedPng); diff != "" {
				t.Errorf("mismatched png output (-got,+want):\n %s", diff)
			}
		})
	}
}
