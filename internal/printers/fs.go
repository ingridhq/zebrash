package printers

import (
	"fmt"
	"slices"
	"strings"
)

const (
	StoredFormatDefaultPath   = "R:UNKNOWN.ZPL"
	StoredGraphicsDefaultPath = "R:UNKNOWN.GRF"
	StoredFontDefaultPath     = "R:UNKNOWN.FNT"
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

func EnsureExtensions(path string, exts ...string) string {
	if len(exts) == 0 {
		return path
	}

	pathParts := strings.SplitN(path, ".", 2)
	if len(pathParts) > 1 && slices.Contains(exts, pathParts[1]) {
		return path
	}

	return fmt.Sprintf("%s.%s", pathParts[0], exts[0])
}
