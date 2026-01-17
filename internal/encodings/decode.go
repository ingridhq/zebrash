package encodings

import (
	"strings"

	"golang.org/x/text/encoding/charmap"
)

// Encodings 0-13 are all in fact CP850 encoding
// 13 is normal CP850
// 0-12 have some characters replaced with other characters
var characterSets013 = [14][11]string{
	{"#", "0", "@", "[", "¢", "]", "^", "`", "{", "|", "}"},
	{"#", "0", "@", "⅓", "¢", "⅔", "^", "`", "¼", "½", "¾"},
	{"£", "0", "@", "[", "¢", "]", "^", "`", "{", "|", "}"},
	{"ƒ", "0", "§", "[", "IJ", "]", "^", "`", "{", "ij", "}"},
	{"#", "0", "@", "Æ", "Ø", "Å", "^", "`", "æ", "ø", "å"},
	{"Ü", "0", "É", "Ä", "Ö", "Å", "Ü", "é", "ä", "ö", "å"},
	{"#", "0", "§", "Ä", "Ö", "Ü", "^", "`", "ä", "ö", "ü"},
	{"£", "0", "à", "[", "ç", "]", "^", "`", "é", "|", "ù"},
	{"#", "0", "à", "â", "ç", "ê", "î", "ô", "é", "ù", "è"},
	{"£", "0", "§", "[", "ç", "é", "^", "ù", "à", "ò", "è"},
	{"#", "0", "§", "¡", "Ñ", "¿", "^", "`", "{", "ñ", "ç"},
	{"£", "0", "É", "Ä", "Ö", "Ü", "^", "ä", "ë", "ï", "ö"},
	{"#", "0", "@", "[", "¥", "]", "^", "`", "{", "|", "}"},
	{"#", "0", "@", "[", "\\", "]", "^", "`", "{", "|", "}"},
}

func ToUnicodeText(text string, charset int) (string, error) {
	switch {
	case charset >= 0 && charset <= 13:
		text, err := charmap.CodePage850.NewDecoder().String(text)
		if err != nil {
			return "", err
		}

		if charset < 13 {
			search := characterSets013[13]
			replace := characterSets013[charset]

			for i, v := range search {
				text = strings.ReplaceAll(text, v, replace[i])
			}
		}

		return text, nil
	case charset == 27:
		return charmap.Windows1252.NewDecoder().String(text)
	default:
		return text, nil
	}
}
