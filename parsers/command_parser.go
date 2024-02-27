package parsers

import (
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/printers"
)

type CommandParser struct {
	CommandCode string
	Parse       func(command string, printer *printers.VirtualPrinter) (interface{}, error)
}

func (p *CommandParser) CanParse(command string) bool {
	return strings.HasPrefix(command, p.CommandCode)
}

func splitCommand(command, prefix string, pos int) []string {
	data := command[len(prefix)+pos:]

	return strings.Split(data, ",")
}

func commandText(command, prefix string) string {
	return command[len(prefix):]
}

func toFieldOrientation(orientation byte) elements.FieldOrientation {
	switch orientation {
	case 'N':
		return elements.FieldOrientationNormal
	case 'R':
		return elements.FieldOrientation90
	case 'I':
		return elements.FieldOrientation180
	case 'B':
		return elements.FieldOrientation270
	default:
		return elements.FieldOrientationNormal
	}
}

func toTextAlignment(alignment byte) elements.TextAlignment {
	switch alignment {
	case 'L':
		return elements.TextAlignmentLeft
	case 'R':
		return elements.TextAlignmentRight
	case 'J':
		return elements.TextAlignmentJustified
	default:
		return elements.TextAlignmentLeft
	}
}
