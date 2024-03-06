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
		name    string
		srcPath string
		dstPath string
	}{
		{
			name:    "UPS label",
			srcPath: "ups.zpl",
			dstPath: "ups.png",
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

			expectedPng, err := os.ReadFile("./testdata/" + tC.dstPath)
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

			err = drawer.DrawLabelAsPng(res[0], &buff, drawers.DrawerOptions{})
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(buff.Bytes(), expectedPng); diff != "" {
				t.Errorf("mismatched png output (-got,+want):\n %s", diff)
			}
		})
	}
}

func TestDrawLabelAsPng_Barcode128(t *testing.T) {
	filenames := []string{
		"barcode128_line_above",
		"barcode128_line",
		"barcode128_mode_a",
		"barcode128_mode_d",
		"barcode128_mode_n",
		"barcode128_mode_u",
		"barcode128_rotated",
	}

	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			file, err := os.ReadFile("./testdata/" + filename + ".zpl")
			if err != nil {
				t.Fatal(err)
			}

			expectedPng, err := os.ReadFile("./testdata/" + filename + ".png")
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

			err = drawer.DrawLabelAsPng(res[0], &buff, drawers.DrawerOptions{})
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(buff.Bytes(), expectedPng); diff != "" {
				t.Fatalf("unexpected output label, diff(+got,-want): %s", diff)
			}
		})
	}
}
