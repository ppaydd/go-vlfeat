package vlfeat

/*
#include <stdlib.h>
#include <aib.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type AIB struct {
	p *C.VlAIB
}

/* Create and destroy */

// https://www.vlfeat.org/api/aib_8h.html#a4bb02325ebe150348976dbb46183f8f3
func NewAIB(pcx [][]float64, rows, cols int) AIB {
	length := rows * cols
	cPcx := make([]C.double, length)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			cPcx[i*cols+j] = C.double(pcx[i][j])
		}
	}
	p := C.vl_aib_new(&cPcx[0], C.uint(rows), C.uint(cols))
	return AIB{p: p}
}

// https://www.vlfeat.org/api/aib_8h.html#a145bb0e2d8f512613e8fd3fb92fe7ace
func (aib *AIB) Delete() {
	C.vl_aib_delete(aib.p)
}

/* Process data */

// https://www.vlfeat.org/api/aib_8h.html#a2adbc69469a6200896e6582474c398ee
func (aib *AIB) Process() {
	C.vl_aib_process(aib.p)
}

/* Retrieve results */

func (aib *AIB) GetNvalues() int {
	return int(aib.p.nvalues)
}

// https://www.vlfeat.org/api/aib_8h.html#ab5f34e685e826902284748f40328f27e
func (aib *AIB) GetParents() []uint {
	nvalues := aib.GetNvalues()
	length := nvalues*2 - 1
	cParentsPtr := C.vl_aib_get_parents(aib.p)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cParentsPtr)),
		Len:  int(length),
		Cap:  int(length),
	}
	cParentsSlice := *(*[]C.uint)(unsafe.Pointer(&hdr))
	parents := make([]uint, length)
	for i, parent := range cParentsSlice {
		parents[i] = uint(parent)
	}
	return parents
}

// https://www.vlfeat.org/api/aib_8h.html#a9fca2d098bd88a72132c5eca55db4f7b
func (aib *AIB) GetCost() []float64 {
	length := aib.GetNvalues()
	cParentsPtr := C.vl_aib_get_costs(aib.p)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cParentsPtr)),
		Len:  int(length),
		Cap:  int(length),
	}
	cParentsSlice := *(*[]C.double)(unsafe.Pointer(&hdr))
	parents := make([]float64, length)
	for i, parent := range cParentsSlice {
		parents[i] = float64(parent)
	}
	return parents
}

// https://www.vlfeat.org/api/aib_8h.html#a7b284789083a9644e96f66da1319ec96
func (aib *AIB) GetVerbosity() int {
	return int(C.vl_aib_get_verbosity(aib.p))
}

// https://www.vlfeat.org/api/aib_8h.html#a452b9ca04f5c7a46a6b82df3f3ed9bf5
func (aib *AIB) SetVerbosity(verbosity int) {
	C.vl_aib_set_verbosity(aib.p, C.int(verbosity))
}
