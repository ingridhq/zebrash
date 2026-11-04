package zebrash

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ingridhq/zebrash/drawers"
)

func renderZplForFbTest(t *testing.T, zpl string) []byte {
	t.Helper()

	labels, err := NewParser().Parse([]byte(zpl))
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 1 {
		t.Fatalf("expected 1 label, got %d", len(labels))
	}

	buff := new(bytes.Buffer)
	err = NewDrawer().DrawLabelAsPng(labels[0], buff, drawers.DrawerOptions{
		LabelWidthMm:  50,
		LabelHeightMm: 25,
		Dpmm:          8,
	})
	if err != nil {
		t.Fatal(err)
	}

	return buff.Bytes()
}

// Every word is wider than half of the block, so each word occupies its own line.
func TestFieldBlockMaxLinesClipsExcessText(t *testing.T) {
	words := []string{"WWWW", "MMMM", "HHHH", "NNNN", "KKKK"}

	join := func(n int) string {
		res := words[0]
		for _, w := range words[1:n] {
			res += " " + w
		}
		return res
	}

	testCases := []struct {
		name      string
		origin    string // ^FO or ^FT
		font      string
		alignment string
		maxLines  int
		words     int // words in the text under test
	}{
		{"max 1 of 3 left", "^FO10,20", "^A0,30,30", "L", 1, 3},
		{"max 2 of 4 left", "^FO10,20", "^A0,30,30", "L", 2, 4},
		{"max 2 of 5 center", "^FO10,20", "^A0,30,30", "C", 2, 5},
		{"max 2 of 4 right", "^FO10,20", "^A0,30,30", "R", 2, 4},
		{"max 2 of 4 justified", "^FO10,20", "^A0,30,30", "J", 2, 4},
		{"max 1 of 3 bottom anchored", "^FT10,150", "^A0,30,30", "L", 1, 3},
		{"max 2 of 4 bottom anchored", "^FT10,150", "^A0,30,30", "L", 2, 4},
		{"max 2 of 4 rotated 90", "^FO100,20", "^A0R,30,30", "L", 2, 4},
		{"max 2 of 2 fits", "^FO10,20", "^A0,30,30", "L", 2, 2},
		{"max 3 of 2 fits with spare lines", "^FO10,20", "^A0,30,30", "L", 3, 2},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			build := func(maxLines int, text string) string {
				return fmt.Sprintf("^XA%s%s^FB150,%d,0,%s,0^FD%s^FS^XZ", tC.origin, tC.font, maxLines, tC.alignment, text)
			}

			got := renderZplForFbTest(t, build(tC.maxLines, join(tC.words)))

			// Same block, but the text contains only the lines that may be printed.
			want := renderZplForFbTest(t, build(tC.maxLines, join(min(tC.words, tC.maxLines))))

			if !bytes.Equal(got, want) {
				t.Fatalf("text beyond maxLines=%d was drawn: output differs from the render of the first %d lines only", tC.maxLines, min(tC.words, tC.maxLines))
			}
		})
	}
}
