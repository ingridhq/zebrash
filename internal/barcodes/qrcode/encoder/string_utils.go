package encoder

import (
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
)

const (
	StringUtils_ASSUME_SHIFT_JIS = false
	// Retained for ABI compatibility with earlier versions
	StringUtils_SHIFT_JIS = "SJIS"
	StringUtils_GB2312    = "GB2312"
)

var (
	StringUtils_PLATFORM_DEFAULT_ENCODING = unicode.UTF8
	StringUtils_SHIFT_JIS_CHARSET         = japanese.ShiftJIS         // "SJIS"
	StringUtils_GB2312_CHARSET            = simplifiedchinese.GB18030 // "GB2312"
	StringUtils_EUC_JP                    = japanese.EUCJP            // "EUC_JP"
)
