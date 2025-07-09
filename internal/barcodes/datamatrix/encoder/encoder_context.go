package encoder

import (
	"fmt"

	"golang.org/x/text/encoding/charmap"
)

type EncoderContext struct {
	msg         []byte
	opts        Options
	codewords   []byte
	pos         int
	newEncoding int
	symbolInfo  *SymbolInfo
	skipAtEnd   int
}

func NewEncoderContext(msg string) (*EncoderContext, error) {
	//From this point on Strings are not Unicode anymore!
	msgBinary, e := charmap.ISO8859_1.NewEncoder().Bytes([]byte(msg))
	if e != nil {
		return nil, fmt.Errorf("message contains characters outside ISO-8859-1 encoding. %v", e)
	}
	sb := make([]byte, 0, len(msgBinary))
	for i, c := 0, len(msgBinary); i < c; i++ {
		ch := msgBinary[i] & 0xff
		sb = append(sb, ch)
	}
	return &EncoderContext{
		msg: sb, //Not Unicode here!
		opts: Options{
			Shape: SymbolShapeHint_FORCE_NONE,
		},
		codewords:   make([]byte, 0, len(sb)),
		newEncoding: -1,
	}, nil
}

func (ec *EncoderContext) SetSymbolShape(shape SymbolShapeHint) {
	ec.opts.Shape = shape
}

func (ec *EncoderContext) SetSizeConstraints(minSize, maxSize *Dimension) {
	ec.opts.MinSize = minSize
	ec.opts.MaxSize = maxSize
}

func (ec *EncoderContext) GetMessage() []byte {
	return ec.msg
}

func (ec *EncoderContext) SetSkipAtEnd(count int) {
	ec.skipAtEnd = count
}

func (ec *EncoderContext) GetCurrentChar() byte {
	return ec.msg[ec.pos]
}

func (ec *EncoderContext) GetCurrent() byte {
	return ec.msg[ec.pos]
}

func (ec *EncoderContext) GetCodewords() []byte {
	return ec.codewords
}

func (ec *EncoderContext) WriteCodewords(codewords []byte) {
	ec.codewords = append(ec.codewords, codewords...)
}

func (ec *EncoderContext) WriteCodeword(codeword byte) {
	ec.codewords = append(ec.codewords, codeword)
}

func (ec *EncoderContext) GetCodewordCount() int {
	return len(ec.codewords)
}

func (ec *EncoderContext) GetNewEncoding() int {
	return ec.newEncoding
}

func (ec *EncoderContext) SignalEncoderChange(encoding int) {
	ec.newEncoding = encoding
}

func (ec *EncoderContext) ResetEncoderSignal() {
	ec.newEncoding = -1
}

func (ec *EncoderContext) HasMoreCharacters() bool {
	return ec.pos < ec.getTotalMessageCharCount()
}

func (ec *EncoderContext) getTotalMessageCharCount() int {
	return len(ec.msg) - ec.skipAtEnd
}

func (ec *EncoderContext) GetRemainingCharacters() int {
	return ec.getTotalMessageCharCount() - ec.pos
}

func (ec *EncoderContext) GetSymbolInfo() *SymbolInfo {
	return ec.symbolInfo
}

func (ec *EncoderContext) UpdateSymbolInfo() error {
	return ec.UpdateSymbolInfoByLength(ec.GetCodewordCount())
}

func (ec *EncoderContext) UpdateSymbolInfoByLength(len int) error {
	var e error
	if ec.symbolInfo == nil || len > ec.symbolInfo.GetDataCapacity() {
		ec.symbolInfo, e = SymbolInfo_Lookup(len, ec.opts, true)
	}
	return e
}

func (ec *EncoderContext) ResetSymbolInfo() {
	ec.symbolInfo = nil
}
