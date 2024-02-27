package zebrash

import (
	"fmt"
	"strings"

	"github.com/ingridhq/zebrash/elements"
	"github.com/ingridhq/zebrash/parsers"
	"github.com/ingridhq/zebrash/printers"
)

type Parser struct {
	printer        *printers.VirtualPrinter
	commandParsers []*parsers.CommandParser
}

func NewParser() *Parser {
	return &Parser{
		printer: printers.NewVirtualPrinter(),
		commandParsers: []*parsers.CommandParser{
			parsers.NewGraphicBoxParser(),
			parsers.NewGraphicCircleParser(),
			parsers.NewChangeDefaultFontParser(),
			parsers.NewChangeFontParser(),
			parsers.NewFieldOriginParser(),
			parsers.NewFieldBlockParser(),
			parsers.NewFieldSeparatorParser(),
			parsers.NewFieldDataParser(),
			parsers.NewFieldValueParser(),
			parsers.NewFieldOrientationParser(),
		},
	}
}

func (p *Parser) Parse(zplData []byte) ([]elements.LabelInfo, error) {
	var results []elements.LabelInfo
	var resultElements []interface{}

	const startCode = "^XA"
	const endCode = "^XZ"

	commands, err := splitZplCommands(zplData)
	if err != nil {
		return nil, fmt.Errorf("failed to split zpl commands: %w", err)
	}

	for _, command := range commands {
		if strings.ToUpper(command) == startCode {
			p.printer.NextDownloadFormatName = ""
			continue
		}

		if strings.ToUpper(command) == endCode {
			results = append(results, elements.LabelInfo{
				DownloadFormatName: p.printer.NextDownloadFormatName,
				Elements:           resultElements,
			})

			resultElements = nil
			continue
		}

		for _, cp := range p.commandParsers {
			if !cp.CanParse(command) {
				continue
			}

			el, err := cp.Parse(command, p.printer)
			if err != nil {
				return nil, fmt.Errorf("failed to parse zpl command %v: %w", command, err)
			}

			if el != nil {
				resultElements = append(resultElements, el)
			}
		}
	}

	return results, nil
}

func splitZplCommands(zplData []byte) ([]string, error) {
	data := strings.ReplaceAll(string(zplData), "\n", "")
	data = strings.ReplaceAll(data, "\r", "")
	data = strings.ReplaceAll(data, "\t", "")

	caret := byte('^')
	tilde := byte('~')

	var buff strings.Builder
	var results []string

	const changeTildeCode = "CT"
	const changeCaretCode = "CC"

	for i := 0; i < len(data); i++ {
		c := data[i]

		if c == caret || c == tilde {
			command := buff.String()
			buff.Reset()

			if len(command) > 2 {
				command = normalizeCommand(command, tilde, caret)

				switch {
				case strings.Index(command, changeTildeCode) == 1:
					tilde = command[3]
				case strings.Index(command, changeCaretCode) == 1:
					caret = command[3]
				default:
					results = append(results, command)
				}
			}
		}

		if err := buff.WriteByte(c); err != nil {
			return nil, err
		}
	}

	command := buff.String()

	if len(command) > 0 {
		command = normalizeCommand(command, tilde, caret)
		results = append(results, command)
	}

	return results, nil
}

func normalizeCommand(command string, tilde, caret byte) string {
	if caret != '^' && command[0] == caret {
		command = "^" + command[1:]
	}

	if tilde != '~' && command[0] == tilde {
		command = "~" + command[1:]
	}

	return strings.TrimLeft(command, " ")
}
