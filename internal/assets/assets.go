package assets

import (
	_ "embed"
)

// Slightly modified HelveticaBoldCondensed with some added glyphs for Polish and Turkish
//
//go:embed fonts/HelveticaBoldCondensedCustom.ttf
var FontHelveticaBold []byte

//go:embed fonts/DejaVuSansMono.ttf
var FontDejavuSansMono []byte

//go:embed fonts/DejaVuSansMonoBold.ttf
var FontDejavuSansMonoBold []byte

//go:embed fonts/ZplGSCustom.ttf
var FontZplGS []byte
