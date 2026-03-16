package zebrash

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	ie "github.com/ingridhq/zebrash/internal/elements"
)

// EPLParser parses EPL (Eltron Programming Language) data into label elements
// that can be rendered using the existing Drawer.
type EPLParser struct{}

func NewEPLParser() *EPLParser {
	return &EPLParser{}
}

// Parse parses EPL data and returns a slice of LabelInfo, one per label.
// EPL is line-based: N starts a label, P ends it. Element commands (A, B, LO)
// create text, barcodes, and lines respectively.
func (p *EPLParser) Parse(eplData []byte) ([]elements.LabelInfo, error) {
	lines := strings.Split(string(eplData), "\n")

	var results []elements.LabelInfo
	var currentElements []any
	var refX, refY int

	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		switch {
		case line == "N":
			currentElements = nil
			refX, refY = 0, 0

		case isEPLReferencePoint(line):
			parts := strings.SplitN(line[1:], ",", 2)
			if len(parts) >= 1 {
				refX, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
			}
			if len(parts) >= 2 {
				refY, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}

		case strings.HasPrefix(line, "A"):
			el, err := parseEPLText(line, refX, refY)
			if err != nil {
				return nil, fmt.Errorf("failed to parse EPL A command %q: %w", line, err)
			}
			if el != nil {
				currentElements = append(currentElements, el)
			}

		case strings.HasPrefix(line, "B"):
			el, err := parseEPLBarcode(line, refX, refY)
			if err != nil {
				return nil, fmt.Errorf("failed to parse EPL B command %q: %w", line, err)
			}
			if el != nil {
				currentElements = append(currentElements, el)
			}

		case strings.HasPrefix(line, "LO"):
			el, err := parseEPLLine(line, refX, refY)
			if err != nil {
				return nil, fmt.Errorf("failed to parse EPL LO command %q: %w", line, err)
			}
			if el != nil {
				currentElements = append(currentElements, el)
			}

		case isEPLPrintCommand(line):
			if len(currentElements) > 0 {
				results = append(results, elements.LabelInfo{
					Elements: currentElements,
				})
			}
			currentElements = nil

			// Ignored commands: Q (form length), S (speed), D (density), ZB/ZT (direction)
		}
	}

	// Handle labels without a trailing P command
	if len(currentElements) > 0 {
		results = append(results, elements.LabelInfo{
			Elements: currentElements,
		})
	}

	return results, nil
}

// isEPLReferencePoint checks if line is an EPL R (reference point) command.
func isEPLReferencePoint(line string) bool {
	return len(line) > 1 && line[0] == 'R' && (line[1] >= '0' && line[1] <= '9')
}

// isEPLPrintCommand checks if line is an EPL P (print) command.
func isEPLPrintCommand(line string) bool {
	if len(line) == 0 || line[0] != 'P' {
		return false
	}
	if len(line) == 1 {
		return true
	}
	// P followed by digits (copy count)
	for _, c := range line[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// EPL font base dimensions (width, height) at 203 DPI.
var eplFontSizes = map[int][2]int{
	1: {8, 12},
	2: {10, 16},
	3: {12, 20},
	4: {14, 24},
	5: {32, 48},
}

func eplFontSize(fontNum int) (width, height int) {
	if size, ok := eplFontSizes[fontNum]; ok {
		return size[0], size[1]
	}
	return 8, 12 // default to font 1
}

func eplRotation(rotation int) ie.FieldOrientation {
	switch rotation {
	case 1:
		return ie.FieldOrientation90
	case 2:
		return ie.FieldOrientation180
	case 3:
		return ie.FieldOrientation270
	default:
		return ie.FieldOrientationNormal
	}
}

// parseEPLText parses an EPL A (ASCII text) command.
// Format: Ap1,p2,p3,p4,p5,p6,p7,"DATA"
func parseEPLText(line string, refX, refY int) (*ie.TextField, error) {
	dataStart := strings.Index(line, "\"")
	if dataStart == -1 {
		return nil, nil
	}
	dataEnd := strings.LastIndex(line, "\"")
	if dataEnd <= dataStart {
		return nil, nil
	}

	text := line[dataStart+1 : dataEnd]
	if text == "" {
		return nil, nil
	}

	paramStr := strings.TrimRight(line[1:dataStart], ",")
	parts := strings.Split(paramStr, ",")

	if len(parts) < 7 {
		return nil, fmt.Errorf("EPL A command requires at least 7 parameters, got %d", len(parts))
	}

	x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	rotation, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
	fontNum, _ := strconv.Atoi(strings.TrimSpace(parts[3]))
	hMult, _ := strconv.Atoi(strings.TrimSpace(parts[4]))
	vMult, _ := strconv.Atoi(strings.TrimSpace(parts[5]))
	reverse := strings.TrimSpace(parts[6])

	if hMult < 1 {
		hMult = 1
	}
	if vMult < 1 {
		vMult = 1
	}

	baseW, baseH := eplFontSize(fontNum)

	// EPL bitmap cell dimensions (e.g. 8x12) include inter-character spacing
	// and don't represent the actual glyph aspect ratio of the TTF substitute
	// font (Helvetica Bold).  Using raw baseW/baseH as Width/Height would
	// produce scaleX = 8/12 ≈ 0.67, squishing text horizontally and making
	// it look thin and blurry.
	//
	// Instead we set Width == Height so scaleX == 1.0 at equal multipliers,
	// then apply the multiplier ratio so that independent h/v scaling still
	// works (e.g. hMult=2, vMult=1 doubles the width).
	fontHeight := float64(baseH * vMult)
	fontWidth := fontHeight
	if hMult != vMult {
		fontWidth = fontHeight * float64(hMult*baseW) / float64(vMult*baseH)
	}

	return &ie.TextField{
		ReversePrint: ie.ReversePrint{Value: reverse == "R"},
		Font: ie.FontInfo{
			Name:        "0",
			Width:       fontWidth,
			Height:      fontHeight,
			Orientation: eplRotation(rotation),
		},
		Position: ie.LabelPosition{
			X: x + refX,
			Y: y + refY,
		},
		Text: text,
	}, nil
}

// parseEPLBarcode parses an EPL B (barcode) command.
// Format: Bp1,p2,p3,p4,p5,p6,p7,p8,"DATA"
func parseEPLBarcode(line string, refX, refY int) (any, error) {
	dataStart := strings.Index(line, "\"")
	if dataStart == -1 {
		return nil, nil
	}
	dataEnd := strings.LastIndex(line, "\"")
	if dataEnd <= dataStart {
		return nil, nil
	}

	data := line[dataStart+1 : dataEnd]
	if data == "" {
		return nil, nil
	}

	paramStr := strings.TrimRight(line[1:dataStart], ",")
	parts := strings.Split(paramStr, ",")

	if len(parts) < 8 {
		return nil, fmt.Errorf("EPL B command requires at least 8 parameters, got %d", len(parts))
	}

	x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	rotation, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
	bcType := strings.TrimSpace(parts[3])
	narrowBar, _ := strconv.Atoi(strings.TrimSpace(parts[4]))
	wideBar, _ := strconv.Atoi(strings.TrimSpace(parts[5]))
	height, _ := strconv.Atoi(strings.TrimSpace(parts[6]))
	humanReadable := strings.TrimSpace(parts[7])

	if narrowBar < 1 {
		narrowBar = 1
	}
	if height < 1 {
		height = 10
	}

	pos := ie.LabelPosition{
		X: x + refX,
		Y: y + refY,
	}
	orient := eplRotation(rotation)
	showLine := humanReadable == "B"

	widthRatio := float64(wideBar) / float64(narrowBar)
	if widthRatio < 2 {
		widthRatio = 2
	}

	switch bcType {
	case "0":
		return &ie.Barcode39WithData{
			Barcode39: ie.Barcode39{
				Orientation: orient,
				Height:      height,
				Line:        showLine,
			},
			Width:      narrowBar,
			WidthRatio: widthRatio,
			Position:   pos,
			Data:       data,
		}, nil

	case "1", "1A", "1B", "1C":
		// Code 128 variants
		return &ie.Barcode128WithData{
			Barcode128: ie.Barcode128{
				Orientation: orient,
				Height:      height,
				Line:        showLine,
				Mode:        ie.BarcodeModeAutomatic,
			},
			Width:    narrowBar,
			Position: pos,
			Data:     data,
		}, nil

	case "3", "4", "5", "6", "7", "8":
		// Code 128 sub-types
		return &ie.Barcode128WithData{
			Barcode128: ie.Barcode128{
				Orientation: orient,
				Height:      height,
				Line:        showLine,
				Mode:        ie.BarcodeModeAutomatic,
			},
			Width:    narrowBar,
			Position: pos,
			Data:     data,
		}, nil

	case "9", "A":
		// UPC-A / UPC-E — fall through to Code 128 as best-effort
		return &ie.Barcode128WithData{
			Barcode128: ie.Barcode128{
				Orientation: orient,
				Height:      height,
				Line:        showLine,
				Mode:        ie.BarcodeModeAutomatic,
			},
			Width:    narrowBar,
			Position: pos,
			Data:     data,
		}, nil

	case "B":
		return &ie.BarcodeEan13WithData{
			BarcodeEan13: ie.BarcodeEan13{
				Orientation: orient,
				Height:      height,
				Line:        showLine,
			},
			Width:    narrowBar,
			Position: pos,
			Data:     data,
		}, nil

	case "E", "F":
		// Codabar
		return &ie.Barcode128WithData{
			Barcode128: ie.Barcode128{
				Orientation: orient,
				Height:      height,
				Line:        showLine,
				Mode:        ie.BarcodeModeAutomatic,
			},
			Width:    narrowBar,
			Position: pos,
			Data:     data,
		}, nil

	case "G", "H":
		// Interleaved 2 of 5
		return &ie.Barcode2of5WithData{
			Barcode2of5: ie.Barcode2of5{
				Orientation: orient,
				Height:      height,
				Line:        showLine,
			},
			Width:      narrowBar,
			WidthRatio: widthRatio,
			Position:   pos,
			Data:       data,
		}, nil

	default:
		// Default to Code 128 Auto for unknown types
		return &ie.Barcode128WithData{
			Barcode128: ie.Barcode128{
				Orientation: orient,
				Height:      height,
				Line:        showLine,
				Mode:        ie.BarcodeModeAutomatic,
			},
			Width:    narrowBar,
			Position: pos,
			Data:     data,
		}, nil
	}
}

// parseEPLLine parses an EPL LO (line draw) command.
// Format: LOp1,p2,p3,p4
func parseEPLLine(line string, refX, refY int) (*ie.GraphicBox, error) {
	paramStr := line[2:] // Skip "LO"
	parts := strings.Split(paramStr, ",")

	if len(parts) < 4 {
		return nil, fmt.Errorf("EPL LO command requires 4 parameters, got %d", len(parts))
	}

	x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	width, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
	height, _ := strconv.Atoi(strings.TrimSpace(parts[3]))

	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	return &ie.GraphicBox{
		Position: ie.LabelPosition{
			X: x + refX,
			Y: y + refY,
		},
		Width:           width,
		Height:          height,
		BorderThickness: min(width, height),
		LineColor:       ie.LineColorBlack,
	}, nil
}
