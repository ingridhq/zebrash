package elements

import (
	"fmt"
	"strings"

	"github.com/ingridhq/maxicode"
)

type Maxicode struct {
	// The mode to use to encode the bar code data.
	// Valid values are 2 (numeric postal code),
	// 3 (alphanumeric postal code),
	// 4 (standard), 5 (full EEC), and 6 (reader programming).
	// The default value is 2 (numeric postal code).
	Mode int
}

type MaxicodeWithData struct {
	Code     Maxicode
	Position LabelPosition
	Data     string
}

func (barcode *MaxicodeWithData) GetInputData() (string, error) {
	const RS = maxicode.RS
	const GS = maxicode.GS
	const codeHeader = "[)>" + RS + "01" + GS
	const headerLen = 9

	data := barcode.Data
	headerPos := strings.Index(data, codeHeader)

	if headerPos < 0 || len(data) < headerPos+headerLen {
		return "", fmt.Errorf("invalid length of maxicode data")
	}

	mainData := data[:headerPos]
	addData := data[headerPos:]
	headerData := addData[0:headerLen]

	// ZPL commands have maxicode data as 2 separate sections
	// main data - class of service (3 chars long), ship to country (chars long), postal code
	// additional data - everything else
	// sections are separate by the special maxicode header
	// we need to combine them into one data before we can produce maxicode out of it

	if len(mainData) < 7 {
		return "", fmt.Errorf("invalid length of maxicode main data")
	}

	classOfService := mainData[0:3]
	shipToCountry := mainData[3:6]
	postalCode := mainData[6:]

	inputData := strings.ReplaceAll(addData, headerData, headerData+postalCode+GS+shipToCountry+GS+classOfService+GS)

	return inputData, nil
}
