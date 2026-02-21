package drawers

import "testing"

func TestBarcodeCaptionFontSize(t *testing.T) {
	tests := []struct {
		barWidth float64
		want    float64
	}{
		{1, 10},
		{2, 20},
		{3, 30},
		{4, 40},
		{5, 45},
		{6, 50},
		{7, 60},
		{8, 70},
		{9, 80},
		{10, 90},
		{11, 100},
		{0.5, 10},   // clamped to 1
		{-1, 10},    // clamped to 1
	}
	for _, tt := range tests {
		got := barcodeCaptionFontSize(tt.barWidth)
		if got != tt.want {
			t.Errorf("barcodeCaptionFontSize(%v) = %v, want %v", tt.barWidth, got, tt.want)
		}
	}
}
