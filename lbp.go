package vlfeat

/*
#include <stdlib.h>
#include <lbp.h>
*/
import "C"

type VlLbpMappingType int

const (
	VlLbpUniform VlLbpMappingType = 0
)

type Lbp struct {
	p *C.VlLbp
}

// https://www.vlfeat.org/api/lbp_8c.html#a3e6b2fc3465c379f3acc45c9fe5b179c
func NewLbp(lbpType VlLbpMappingType, transposed bool) Lbp {
	cTransposed := 0
	if transposed {
		cTransposed = 1
	}
	p := C.vl_lbp_new(C.VlLbpMappingType(lbpType), C.int(cTransposed))
	return Lbp{p: p}
}

// https://www.vlfeat.org/api/lbp_8c.html#af4061c7ff063118d14893cafe5b55ed8
func (lbp *Lbp) Delete() {
	C.vl_lbp_delete(lbp.p)
}

// https://www.vlfeat.org/api/lbp_8c.html#a9fdc1d38de7ce1494cdd911bfd8957a7
func (lbp *Lbp) GetDimension() uint {
	return uint(C.vl_lbp_get_dimension(lbp.p))
}

// https://www.vlfeat.org/api/lbp_8c.html#a605c416d7ab3906609cbb62404a20253
func (lbp *Lbp) Process(image []float32, imgWidth, imgHeight, cellSize uint) []float32 {
	imgPtr := toCFloatArrayPtr(image)
	numCols := imgWidth / cellSize
	numRows := imgHeight / cellSize
	dimension := lbp.GetDimension()
	length := numCols * numRows * dimension
	cFeatures := make([]C.float, length)
	C.vl_lbp_process(lbp.p, &cFeatures[0], imgPtr, C.uint(imgWidth), C.uint(imgHeight), C.uint(cellSize))
	features := make([]float32, length)
	for i, data := range cFeatures {
		features[i] = float32(data)
	}
	return features
}
