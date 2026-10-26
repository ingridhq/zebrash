package hex

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"testing"
)

func z64(t *testing.T, data []byte) string {
	t.Helper()

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	return z64Prefix + base64.StdEncoding.EncodeToString(buf.Bytes()) + ":0000"
}

func TestDecodeZ64SizeLimit(t *testing.T) {
	testCases := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name: "data at the size limit",
			size: maxEmbeddedImageSizeMb,
		},
		{
			name:    "data over the size limit",
			size:    maxEmbeddedImageSizeMb + 1,
			wantErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			res, err := DecodeGraphicFieldData(z64(t, make([]byte, tC.size)), 1)

			if tC.wantErr {
				if err == nil {
					t.Fatalf("got %d bytes and no error", len(res))
				}

				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if len(res) != tC.size {
				t.Fatalf("got %d bytes, want %d", len(res), tC.size)
			}
		})
	}
}
