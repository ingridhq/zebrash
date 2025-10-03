package encoder

type BlockPair struct {
	dataBytes            []byte
	errorCorrectionBytes []byte
}

func NewBlockPair(data []byte, errorCorrection []byte) *BlockPair {
	return &BlockPair{data, errorCorrection}
}

func (bp *BlockPair) GetDataBytes() []byte {
	return bp.dataBytes
}

func (bp *BlockPair) GetErrorCorrectionBytes() []byte {
	return bp.errorCorrectionBytes
}
