package ean13

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSanitizeContent(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		result  string
	}{
		{
			name:    "more than 12 digits with invalid characters",
			content: "123gg123gg123gg123gg123gg123gg",
			result:  "1012300123005",
		},
		{
			name:    "less than 12 digits",
			content: "123",
			result:  "0000000001236",
		},
		{
			name:    "exactly 13 digits with valid check digit",
			content: "5901234123457",
			result:  "5901234123457",
		},
		{
			name:    "exactly 13 digits with invalid check digit",
			content: "5901234123456",
			result:  "5901234123457",
		},
		{
			name:    "exactly 12 digits",
			content: "123456789012",
			result:  "1234567890128",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			res := sanitizeContent(tC.content)

			if diff := cmp.Diff(res, tC.result); diff != "" {
				t.Errorf("mismatched result (-got,+want):\n %s", diff)
			}
		})
	}
}
