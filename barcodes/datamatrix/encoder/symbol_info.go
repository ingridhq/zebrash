package encoder

import (
	"fmt"
)

// Symbol info table for DataMatrix.

var symbols = []*SymbolInfo{
	NewSymbolInfo(false, 3, 5, 8, 8, 1),
	NewSymbolInfo(false, 5, 7, 10, 10, 1),
	/*rect*/ NewSymbolInfo(true, 5, 7, 16, 6, 1),
	NewSymbolInfo(false, 8, 10, 12, 12, 1),
	/*rect*/ NewSymbolInfo(true, 10, 11, 14, 6, 2),
	NewSymbolInfo(false, 12, 12, 14, 14, 1),
	/*rect*/ NewSymbolInfo(true, 16, 14, 24, 10, 1),

	NewSymbolInfo(false, 18, 14, 16, 16, 1),
	NewSymbolInfo(false, 22, 18, 18, 18, 1),
	/*rect*/ NewSymbolInfo(true, 22, 18, 16, 10, 2),
	NewSymbolInfo(false, 30, 20, 20, 20, 1),
	/*rect*/ NewSymbolInfo(true, 32, 24, 16, 14, 2),
	NewSymbolInfo(false, 36, 24, 22, 22, 1),
	NewSymbolInfo(false, 44, 28, 24, 24, 1),
	/*rect*/ NewSymbolInfo(true, 49, 28, 22, 14, 2),

	NewSymbolInfo(false, 62, 36, 14, 14, 4),
	NewSymbolInfo(false, 86, 42, 16, 16, 4),
	NewSymbolInfo(false, 114, 48, 18, 18, 4),
	NewSymbolInfo(false, 144, 56, 20, 20, 4),
	NewSymbolInfo(false, 174, 68, 22, 22, 4),

	NewSymbolInfoRS(false, 204, 84, 24, 24, 4, 102, 42),
	NewSymbolInfoRS(false, 280, 112, 14, 14, 16, 140, 56),
	NewSymbolInfoRS(false, 368, 144, 16, 16, 16, 92, 36),
	NewSymbolInfoRS(false, 456, 192, 18, 18, 16, 114, 48),
	NewSymbolInfoRS(false, 576, 224, 20, 20, 16, 144, 56),
	NewSymbolInfoRS(false, 696, 272, 22, 22, 16, 174, 68),
	NewSymbolInfoRS(false, 816, 336, 24, 24, 16, 136, 56),
	NewSymbolInfoRS(false, 1050, 408, 18, 18, 36, 175, 68),
	NewSymbolInfoRS(false, 1304, 496, 20, 20, 36, 163, 62),
	NewDataMatrixSymbolInfo144(),
}

type SymbolInfo struct {
	rectangular    bool
	dataCapacity   int
	errorCodewords int
	matrixWidth    int
	matrixHeight   int
	dataRegions    int
	rsBlockData    int
	rsBlockError   int

	funcGetInterleavedBlockCount         func(*SymbolInfo) int
	funcGetDataLengthForInterleavedBlock func(*SymbolInfo, int) int
}

func NewSymbolInfo(rectangular bool, dataCapacity, errorCodewords,
	matrixWidth, matrixHeight, dataRegions int) *SymbolInfo {
	return NewSymbolInfoRS(rectangular, dataCapacity, errorCodewords,
		matrixWidth, matrixHeight, dataRegions, dataCapacity, errorCodewords)
}

func NewSymbolInfoRS(rectangular bool, dataCapacity, errorCodewords,
	matrixWidth, matrixHeight, dataRegions, rsBlockData, rsBlockError int) *SymbolInfo {
	return &SymbolInfo{
		rectangular:    rectangular,
		dataCapacity:   dataCapacity,
		errorCodewords: errorCodewords,
		matrixWidth:    matrixWidth,
		matrixHeight:   matrixHeight,
		dataRegions:    dataRegions,
		rsBlockData:    rsBlockData,
		rsBlockError:   rsBlockError,

		funcGetInterleavedBlockCount:         defaultGetInterleavedBlockCount,
		funcGetDataLengthForInterleavedBlock: defaultGetDataLengthForInterleavedBlock,
	}
}

func SymbolInfo_Lookup(dataCodewords int, opts Options, fail bool) (*SymbolInfo, error) {

	for _, symbol := range symbols {
		if opts.Shape == SymbolShapeHint_FORCE_SQUARE && symbol.rectangular {
			continue
		}
		if opts.Shape == SymbolShapeHint_FORCE_RECTANGLE && !symbol.rectangular {
			continue
		}
		if opts.MinSize != nil && (symbol.GetSymbolWidth() < opts.MinSize.GetWidth() || symbol.GetSymbolHeight() < opts.MinSize.GetHeight()) {
			continue
		}
		if opts.MaxSize != nil && (symbol.GetSymbolWidth() > opts.MaxSize.GetWidth() || symbol.GetSymbolHeight() > opts.MaxSize.GetHeight()) {
			continue
		}
		if dataCodewords <= symbol.dataCapacity {
			return symbol, nil
		}
	}
	if fail {
		return nil, fmt.Errorf("can't find a symbol arrangement that matches the message. Data codewords: %d", dataCodewords)
	}
	return nil, nil
}

func (si *SymbolInfo) getHorizontalDataRegions() int {
	switch si.dataRegions {
	case 1:
		return 1
	case 2, 4:
		return 2
	case 16:
		return 4
	case 36:
		return 6
	default:
		return 0 // Cannot handle this number of data regions
	}
}

func (si *SymbolInfo) getVerticalDataRegions() int {
	switch si.dataRegions {
	case 1, 2:
		return 1
	case 4:
		return 2
	case 16:
		return 4
	case 36:
		return 6
	default:
		return 0 // Cannot handle this number of data regions
	}
}

func (si *SymbolInfo) GetSymbolDataWidth() int {
	return si.getHorizontalDataRegions() * si.matrixWidth
}

func (si *SymbolInfo) GetSymbolDataHeight() int {
	return si.getVerticalDataRegions() * si.matrixHeight
}

func (si *SymbolInfo) GetSymbolWidth() int {
	return si.GetSymbolDataWidth() + (si.getHorizontalDataRegions() * 2)
}

func (si *SymbolInfo) GetSymbolHeight() int {
	return si.GetSymbolDataHeight() + (si.getVerticalDataRegions() * 2)
}

func (si *SymbolInfo) GetCodewordCount() int {
	return si.dataCapacity + si.errorCodewords
}

func (si *SymbolInfo) GetInterleavedBlockCount() int {
	return si.funcGetInterleavedBlockCount(si)
}

func defaultGetInterleavedBlockCount(si *SymbolInfo) int {
	return si.dataCapacity / si.rsBlockData
}

func (si *SymbolInfo) GetDataCapacity() int {
	return si.dataCapacity
}

func (si *SymbolInfo) GetErrorCodewords() int {
	return si.errorCodewords
}

func (si *SymbolInfo) GetMatrixWidth() int {
	return si.matrixWidth
}

func (si *SymbolInfo) GetMatrixHeight() int {
	return si.matrixHeight
}

func (si *SymbolInfo) GetDataLengthForInterleavedBlock(index int) int {
	return si.funcGetDataLengthForInterleavedBlock(si, index)
}

func defaultGetDataLengthForInterleavedBlock(si *SymbolInfo, index int) int {
	return si.rsBlockData
}

func (si *SymbolInfo) GetErrorLengthForInterleavedBlock(index int) int {
	return si.rsBlockError
}

func (si *SymbolInfo) String() string {
	shape := "Square"
	if si.rectangular {
		shape = "Rectangular"
	}
	return fmt.Sprintf(
		"%s Symbpl: data region %dx%d, symbol size %dx%d, symbol data size %dx%d, codewords %d+%d",
		shape, si.matrixWidth, si.matrixHeight,
		si.GetSymbolWidth(), si.GetSymbolHeight(),
		si.GetSymbolDataWidth(), si.GetSymbolDataHeight(),
		si.dataCapacity, si.errorCodewords)
}
