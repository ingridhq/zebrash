package code128

import (
	"reflect"
	"testing"
)

// Expected codewords are independent of the input scanner. Comparing complete
// bars (including checksum/stop) also catches truncation hidden by HRI alone.
func TestEncodeNoModePayloadBoundaries(t *testing.T) {
	tests := []struct {
		name, input, text string
		words             []byte
	}{
		{"empty", "", "", []byte{104}},
		{"one character", "A", "A", []byte{104, 33}},
		{"two characters", "AB", "AB", []byte{104, 33, 34}},
		{"explicit B one character", ">:A", "A", []byte{104, 33}},
		{"escaped greater than", ">0", ">", []byte{104, 30}},
		{"literal zero after escape", ">00", ">0", []byte{104, 30, 16}},
		{"character after escape", "A>0B", "A>B", []byte{104, 33, 30, 34}},
		{"zero after middle escape", "A>00B", "A>0B", []byte{104, 33, 30, 16, 34}},
		{"repeated escapes", "A>0B>0C", "A>B>C", []byte{104, 33, 30, 34, 30, 35}},
		{"trailing escaped greater than", "A>0", "A>", []byte{104, 33, 30}},
		{"circumflex escape tail", "><A", "^A", []byte{104, 62, 33}},
		{"tilde escape tail", ">=A", "~A", []byte{104, 94, 33}},
		{"function one tail", "A>8B", "AB", []byte{104, 33, 102, 34}},
		{"function three tail", "A>2B", "AB", []byte{104, 33, 96, 34}},
		{"function two tail", "A>3B", "AB", []byte{104, 33, 97, 34}},
		{"subset A pairs", ">90001", " !", []byte{103, 0, 1}},
		{"subset A boundary", ">96364", "_\x00", []byte{103, 63, 64}},
		{"subset A incomplete pair", ">9001", " ", []byte{103, 0}},
		{"subset C pairs", ">;1234", "1234", []byte{105, 12, 34}},
		{"subset C incomplete pair", ">;123", "12", []byte{105, 12}},
		{"subset C ignored first noninteger", ">;D2", "02", []byte{105, 2}},
		{"subset C invalid second digit", ">;2D34", "34", []byte{105, 34}},
		{"subset C to B tail", ">;12>6A", "12A", []byte{105, 12, 100, 33}},
		{"subset A to B tail", ">912>6A", ",A", []byte{103, 12, 100, 33}},
		{"subset B to A pair", "A>700", "A ", []byte{104, 33, 101, 0}},
		{"incomplete invocation", ">", "", []byte{104}},
		{"incomplete invocation after data", "A>", "A", []byte{104, 33}},
		{"incomplete invocation after numeric pair", ">;12>", "12", []byte{105, 12}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			img, text, err := EncodeNoMode(test.input, 10, 1)
			if err != nil {
				t.Fatal(err)
			}
			if text != test.text {
				t.Errorf("HRI = %q, want %q", text, test.text)
			}
			var want []bool
			checksum := int(test.words[0])
			for i, word := range test.words {
				want = appendPattern(want, CODE_PATTERNS[word])
				checksum += i * int(word)
			}
			want = appendPattern(want, CODE_PATTERNS[checksum%103])
			want = appendPattern(want, CODE_PATTERNS[STOP])
			if !reflect.DeepEqual(img.(*code128).code, want) {
				t.Errorf("bars differ from complete codewords %v + checksum/stop", test.words)
			}
		})
	}
}
