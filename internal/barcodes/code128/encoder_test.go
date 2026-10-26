package code128

import "testing"

func TestEncodeInvalidCharacter(t *testing.T) {
	content := "JOÃO"

	if _, err := EncodeAuto(content, 100, 2); err == nil {
		t.Fatalf("EncodeAuto(%q) returned no error", content)
	}

	if _, _, err := EncodeNoMode(content, 100, 2); err == nil {
		t.Fatalf("EncodeNoMode(%q) returned no error", content)
	}
}
