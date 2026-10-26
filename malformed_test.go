package zebrash

import (
	"bytes"
	"testing"

	"github.com/ingridhq/zebrash/drawers"
)

func TestMalformedInput(t *testing.T) {
	testCases := []struct {
		name        string
		zpl         string
		wantDrawErr bool
	}{
		{
			name: "graphic field without parameters",
			zpl:  "^XA^FO10,10^GB100,100,10^FS^GF^XZ",
		},
		{
			name: "graphic field with 0 bytes per row",
			zpl:  "^XA^FO10,10^GB100,100,10^FS^FO10,10^GFA,0,0,0,^FS^XZ",
		},
		{
			name:        "code 128 with a character it cannot encode",
			zpl:         "^XA^FO10,10^BCN,100^FDJOÃO^FS^XZ",
			wantDrawErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			labels, err := NewParser().Parse([]byte(tC.zpl))
			if err != nil {
				t.Fatal(err)
			}

			if len(labels) != 1 {
				t.Fatalf("got %d labels, want 1", len(labels))
			}

			var buf bytes.Buffer
			err = NewDrawer().DrawLabelAsPng(labels[0], &buf, drawers.DrawerOptions{})

			if tC.wantDrawErr != (err != nil) {
				t.Fatalf("DrawLabelAsPng error = %v, want error: %v", err, tC.wantDrawErr)
			}
		})
	}
}
