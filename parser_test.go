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
