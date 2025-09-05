package assets

import (
	_ "embed"
)

//go:embed fonts/ArialCondensedBoldRegular.ttf
var FontArialCondensedBold []byte

//go:embed fonts/DejaVuSansMono.ttf
var FontDejavuSansMono []byte

//go:embed fonts/DejaVuSansMonoBold.ttf
var FontDejavuSansMonoBold []byte
