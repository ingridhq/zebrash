package zebrash

import (
	"bytes"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	file, err := os.ReadFile("./testdata/fedex.zpl")
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

	err = drawer.DrawLabelAsPng(res[0], &buff, DrawerOptions{})
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("./testdata/out.png", buff.Bytes(), 0744)
	if err != nil {
		t.Fatal(err)
	}
}
