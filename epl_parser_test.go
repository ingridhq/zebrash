package zebrash

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ingridhq/zebrash/drawers"
	"github.com/ingridhq/zebrash/internal/elements"
)

func TestEPLParseSingleLabel(t *testing.T) {
	epl := "N\nA10,20,0,1,1,1,N,\"Hello\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	if got := len(labels); got != 1 {
		t.Fatalf("expected 1 label, got %d", got)
	}

	if got := len(labels[0].Elements); got != 1 {
		t.Fatalf("expected 1 element, got %d", got)
	}
}

func TestEPLParseText(t *testing.T) {
	epl := "N\nA50,100,0,2,1,1,N,\"Hello World\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	tf, ok := labels[0].Elements[0].(*elements.TextField)
	if !ok {
		t.Fatalf("expected *TextField, got %T", labels[0].Elements[0])
	}

	if tf.Text != "Hello World" {
		t.Errorf("text = %q, want %q", tf.Text, "Hello World")
	}
	if tf.Position.X != 50 || tf.Position.Y != 100 {
		t.Errorf("position = (%d,%d), want (50,100)", tf.Position.X, tf.Position.Y)
	}
	if tf.Font.Orientation != elements.FieldOrientationNormal {
		t.Errorf("orientation = %v, want Normal", tf.Font.Orientation)
	}
	// Font 2 base: 10w x 16h, mult 1x1 → Width == Height for scaleX=1.0
	if tf.Font.Width != 16 || tf.Font.Height != 16 {
		t.Errorf("font size = (%.0f,%.0f), want (16,16)", tf.Font.Width, tf.Font.Height)
	}
}

func TestEPLParseTextRotated(t *testing.T) {
	epl := "N\nA50,100,1,1,1,1,N,\"Rotated\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	tf := labels[0].Elements[0].(*elements.TextField)
	if tf.Font.Orientation != elements.FieldOrientation90 {
		t.Errorf("orientation = %v, want 90°", tf.Font.Orientation)
	}
}

func TestEPLParseTextReverse(t *testing.T) {
	epl := "N\nA50,100,0,1,1,1,R,\"Reverse\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	tf := labels[0].Elements[0].(*elements.TextField)
	if !tf.IsReversePrint() {
		t.Error("expected reverse print to be true")
	}
}

func TestEPLParseTextMultiplier(t *testing.T) {
	epl := "N\nA10,20,0,3,2,3,N,\"Big\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	tf := labels[0].Elements[0].(*elements.TextField)
	// Font 3 base: 12w x 20h, mult 2x3 → hMult≠vMult so width scaled by ratio
	if tf.Font.Width != 24 || tf.Font.Height != 60 {
		t.Errorf("font size = (%.0f,%.0f), want (24,60)", tf.Font.Width, tf.Font.Height)
	}
}

func TestEPLParseEmptyTextSkipped(t *testing.T) {
	epl := "N\nA10,20,0,1,1,1,N,\"\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	// Empty text should produce no elements, hence no label
	if len(labels) != 0 {
		t.Fatalf("expected 0 labels for empty text, got %d", len(labels))
	}
}

func TestEPLParseBarcode(t *testing.T) {
	epl := "N\nB50,100,0,1,3,6,200,B,\"12345\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	if got := len(labels[0].Elements); got != 1 {
		t.Fatalf("expected 1 element, got %d", got)
	}

	bc, ok := labels[0].Elements[0].(*elements.Barcode128WithData)
	if !ok {
		t.Fatalf("expected *Barcode128WithData, got %T", labels[0].Elements[0])
	}

	if bc.Data != "12345" {
		t.Errorf("data = %q, want %q", bc.Data, "12345")
	}
	if bc.Position.X != 50 || bc.Position.Y != 100 {
		t.Errorf("position = (%d,%d), want (50,100)", bc.Position.X, bc.Position.Y)
	}
	if bc.Height != 200 {
		t.Errorf("height = %d, want 200", bc.Height)
	}
	if !bc.Line {
		t.Error("expected human-readable line to be true")
	}
}

func TestEPLParseBarcodeCode39(t *testing.T) {
	epl := "N\nB50,100,0,0,2,5,100,N,\"ABC123\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	bc, ok := labels[0].Elements[0].(*elements.Barcode39WithData)
	if !ok {
		t.Fatalf("expected *Barcode39WithData, got %T", labels[0].Elements[0])
	}

	if bc.Data != "ABC123" {
		t.Errorf("data = %q, want %q", bc.Data, "ABC123")
	}
	if bc.Width != 2 {
		t.Errorf("width = %d, want 2", bc.Width)
	}
	if bc.WidthRatio != 2.5 {
		t.Errorf("widthRatio = %f, want 2.5", bc.WidthRatio)
	}
}

func TestEPLParseBarcodeInterleaved2of5(t *testing.T) {
	epl := "N\nB50,100,0,G,2,5,100,B,\"1234567890\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	_, ok := labels[0].Elements[0].(*elements.Barcode2of5WithData)
	if !ok {
		t.Fatalf("expected *Barcode2of5WithData, got %T", labels[0].Elements[0])
	}
}

func TestEPLParseLine(t *testing.T) {
	epl := "N\nLO10,20,300,5\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	gb, ok := labels[0].Elements[0].(*elements.GraphicBox)
	if !ok {
		t.Fatalf("expected *GraphicBox, got %T", labels[0].Elements[0])
	}

	if gb.Position.X != 10 || gb.Position.Y != 20 {
		t.Errorf("position = (%d,%d), want (10,20)", gb.Position.X, gb.Position.Y)
	}
	if gb.Width != 300 || gb.Height != 5 {
		t.Errorf("size = (%d,%d), want (300,5)", gb.Width, gb.Height)
	}
	if gb.BorderThickness != 5 {
		t.Errorf("borderThickness = %d, want 5", gb.BorderThickness)
	}
}

func TestEPLParseReferencePoint(t *testing.T) {
	epl := "N\nR40,10\nA50,100,0,1,1,1,N,\"Offset\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	tf := labels[0].Elements[0].(*elements.TextField)
	// Position should include reference offset: (50+40, 100+10) = (90, 110)
	if tf.Position.X != 90 || tf.Position.Y != 110 {
		t.Errorf("position = (%d,%d), want (90,110)", tf.Position.X, tf.Position.Y)
	}
}

func TestEPLParseMultipleLabels(t *testing.T) {
	epl := "N\nA10,20,0,1,1,1,N,\"Label1\"\nP1\nN\nA30,40,0,1,1,1,N,\"Label2\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	if got := len(labels); got != 2 {
		t.Fatalf("expected 2 labels, got %d", got)
	}

	tf1 := labels[0].Elements[0].(*elements.TextField)
	if tf1.Text != "Label1" {
		t.Errorf("label 1 text = %q, want %q", tf1.Text, "Label1")
	}

	tf2 := labels[1].Elements[0].(*elements.TextField)
	if tf2.Text != "Label2" {
		t.Errorf("label 2 text = %q, want %q", tf2.Text, "Label2")
	}
}

func TestEPLParseWithoutPCommand(t *testing.T) {
	epl := "N\nA10,20,0,1,1,1,N,\"NoPrint\"\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	if got := len(labels); got != 1 {
		t.Fatalf("expected 1 label (auto-emitted), got %d", got)
	}
}

func TestEPLParseIgnoredCommands(t *testing.T) {
	epl := "N\nQ822,24\nS4\nD15\nZB\nA10,20,0,1,1,1,N,\"Test\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	if got := len(labels); got != 1 {
		t.Fatalf("expected 1 label, got %d", got)
	}
	if got := len(labels[0].Elements); got != 1 {
		t.Fatalf("expected 1 element, got %d", got)
	}
}

func TestEPLParseMixedElements(t *testing.T) {
	epl := `N
A10,20,0,1,1,1,N,"Header"
B50,100,0,1,3,6,200,N,"12345"
LO0,300,400,2
A10,320,0,2,1,1,N,"Footer"
P1
`
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	if got := len(labels[0].Elements); got != 4 {
		t.Fatalf("expected 4 elements, got %d", got)
	}

	if _, ok := labels[0].Elements[0].(*elements.TextField); !ok {
		t.Errorf("element 0: expected *TextField, got %T", labels[0].Elements[0])
	}
	if _, ok := labels[0].Elements[1].(*elements.Barcode128WithData); !ok {
		t.Errorf("element 1: expected *Barcode128WithData, got %T", labels[0].Elements[1])
	}
	if _, ok := labels[0].Elements[2].(*elements.GraphicBox); !ok {
		t.Errorf("element 2: expected *GraphicBox, got %T", labels[0].Elements[2])
	}
	if _, ok := labels[0].Elements[3].(*elements.TextField); !ok {
		t.Errorf("element 3: expected *TextField, got %T", labels[0].Elements[3])
	}
}

func TestEPLParseDPDUK(t *testing.T) {
	file := mustReadFile("./testdata/dpduk.epl", t)

	parser := NewEPLParser()
	labels, err := parser.Parse(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(labels) == 0 {
		t.Fatal("no labels parsed from dpduk.epl")
	}

	label := labels[0]
	if len(label.Elements) == 0 {
		t.Fatal("no elements in the parsed label")
	}

	// Count element types
	var texts, barcodes, lines int
	for _, el := range label.Elements {
		switch el.(type) {
		case *elements.TextField:
			texts++
		case *elements.Barcode128WithData:
			barcodes++
		case *elements.GraphicBox:
			lines++
		}
	}

	if texts == 0 {
		t.Error("expected at least one text element")
	}
	if barcodes == 0 {
		t.Error("expected at least one barcode element")
	}
	if lines == 0 {
		t.Error("expected at least one line element")
	}

	t.Logf("DPD UK EPL: %d texts, %d barcodes, %d lines", texts, barcodes, lines)
}

func TestEPLDrawDPDUK(t *testing.T) {
	file := mustReadFile("./testdata/dpduk.epl", t)

	parser := NewEPLParser()
	labels, err := parser.Parse(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(labels) == 0 {
		t.Fatal("no labels parsed")
	}

	drawer := NewDrawer()
	buff := new(bytes.Buffer)

	err = drawer.DrawLabelAsPng(labels[0], buff, drawers.DrawerOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if buff.Len() == 0 {
		t.Fatal("empty PNG output")
	}

	t.Logf("DPD UK EPL rendered to %d bytes PNG", buff.Len())
}

func TestEPLParseNResetsReferencePoint(t *testing.T) {
	epl := "N\nR40,10\nA10,20,0,1,1,1,N,\"First\"\nP1\nN\nA10,20,0,1,1,1,N,\"Second\"\nP1\n"
	parser := NewEPLParser()

	labels, err := parser.Parse([]byte(epl))
	if err != nil {
		t.Fatal(err)
	}

	if got := len(labels); got != 2 {
		t.Fatalf("expected 2 labels, got %d", got)
	}

	tf1 := labels[0].Elements[0].(*elements.TextField)
	// First label has R40,10 offset
	if tf1.Position.X != 50 || tf1.Position.Y != 30 {
		t.Errorf("label 1 position = (%d,%d), want (50,30)", tf1.Position.X, tf1.Position.Y)
	}

	tf2 := labels[1].Elements[0].(*elements.TextField)
	// Second label: N resets reference point to (0,0)
	if tf2.Position.X != 10 || tf2.Position.Y != 20 {
		t.Errorf("label 2 position = (%d,%d), want (10,20)", tf2.Position.X, tf2.Position.Y)
	}
}

func TestEPLParseFontSizes(t *testing.T) {
	tests := []struct {
		fontNum        int
		expectedWidth  float64
		expectedHeight float64
	}{
		{1, 12, 12},  // 8x12, equal mult → Width == Height
		{2, 16, 16},  // 10x16
		{3, 20, 20},  // 12x20
		{4, 24, 24},  // 14x24
		{5, 48, 48},  // 32x48
		{9, 12, 12},  // Unknown font defaults to font 1 (8x12)
	}

	for _, tc := range tests {
		epl := fmt.Sprintf("N\nA10,20,0,%d,1,1,N,\"test\"\nP1\n", tc.fontNum)
		parser := NewEPLParser()

		labels, err := parser.Parse([]byte(epl))
		if err != nil {
			t.Fatal(err)
		}

		tf := labels[0].Elements[0].(*elements.TextField)
		if tf.Font.Width != tc.expectedWidth || tf.Font.Height != tc.expectedHeight {
			t.Errorf("font %d: size = (%.0f,%.0f), want (%.0f,%.0f)",
				tc.fontNum, tf.Font.Width, tf.Font.Height, tc.expectedWidth, tc.expectedHeight)
		}
	}
}

func TestEPLParseRotations(t *testing.T) {
	tests := []struct {
		rotation int
		expected elements.FieldOrientation
	}{
		{0, elements.FieldOrientationNormal},
		{1, elements.FieldOrientation90},
		{2, elements.FieldOrientation180},
		{3, elements.FieldOrientation270},
		{7, elements.FieldOrientationNormal}, // Invalid defaults to normal
	}

	for _, tc := range tests {
		epl := fmt.Sprintf("N\nA10,20,%d,1,1,1,N,\"test\"\nP1\n", tc.rotation)
		parser := NewEPLParser()

		labels, err := parser.Parse([]byte(epl))
		if err != nil {
			t.Fatal(err)
		}

		tf := labels[0].Elements[0].(*elements.TextField)
		if tf.Font.Orientation != tc.expected {
			t.Errorf("rotation %d: got %v, want %v", tc.rotation, tf.Font.Orientation, tc.expected)
		}
	}
}
