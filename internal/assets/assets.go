package assets

import (
	_ "embed"
)

//go:embed fonts/HelveticaBoldCondensed.ttf
var FontHelveticaBold []byte

//go:embed fonts/DejaVuSansMono.ttf
var FontDejavuSansMono []byte

//go:embed fonts/DejaVuSansMonoBold.ttf
var FontDejavuSansMonoBold []byte
