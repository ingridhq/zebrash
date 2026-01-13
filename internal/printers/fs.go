package printers

import (
	"fmt"
	"slices"
	"strings"
)

const (
	StoredFormatDefaultPath   = "R:UNKNOWN.ZPL"
	StoredGraphicsDefaultPath = "R:UNKNOWN.GRF"
)

var validDevices = []string{
	"R", "E", "B", "A", "Z",
}

func ValidateDevice(path string) error {
	parts := strings.SplitN(path, ":", 2)
	if len(parts) < 2 {
		return fmt.Errorf("path does not contain device name")
	}

	if !slices.Contains(validDevices, parts[0]) {
		return fmt.Errorf("invalid device name %s, must be one of %s", parts[0], strings.Join(validDevices, ", "))
	}

	return nil
}

func EnsureExtension(path, ext string) string {
	pathParts := strings.SplitN(path, ".", 2)
	return fmt.Sprintf("%s.%s", pathParts[0], ext)
}
