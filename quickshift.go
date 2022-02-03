package vlfeat

/*
#include <stdlib.h>
#include <quickshift.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type QuickShift struct {
	p *C.VlQS
}

// https://www.vlfeat.org/api/quickshift_8c.html#a9fc34955bf121df6d1d5bd991b8d2b13
func NewQuickShift(img []float64, height, width, channles int) QuickShift {
	cImg := make([]C.double, len(img))
	for i, pix := range img {
		cImg[i] = C.double(pix)
	}
	p := C.vl_quickshift_new(&cImg[0], C.int(height), C.int(width), C.int(channles))
	return QuickShift{p: p}
}

// https://www.vlfeat.org/api/quickshift_8c.html#a9c2a39344fb684d899f22faf358425e8
func (qs *QuickShift) Delete() {
	C.vl_quickshift_delete(qs.p)
}

// https://www.vlfeat.org/api/quickshift_8c.html#aebfc7337283b7a8f8be1a1c0fd4f5f92
func (qs *QuickShift) Process() {
	C.vl_quickshift_process(qs.p)
}

/* Set parameters */

func (qs *QuickShift) SetMaxDist(tau float64) {
	C.vl_quickshift_set_max_dist(qs.p, C.double(tau))
}

func (qs *QuickShift) SetKernelSize(sigma float64) {
	C.vl_quickshift_set_kernel_size(qs.p, C.double(sigma))
}

func (qs *QuickShift) SetMedoid(medoid bool) {
	cMedoid := 0
	if medoid {
		cMedoid = 1
	}
	C.vl_quickshift_set_medoid(qs.p, C.int(cMedoid))
}

/* Retrieve data and parameters */

// https://www.vlfeat.org/api/quickshift_8h.html#a0a6c066b205bf0382d19b7c9fe8628d5
func (qs *QuickShift) GetMaxDist() float64 {
	return float64(C.vl_quickshift_get_max_dist(qs.p))
}

// https://www.vlfeat.org/api/quickshift_8h.html#a929d3f6d37b578d40edddb0e6aa9d4e3
func (qs *QuickShift) GetKernelSize() float64 {
	return float64(C.vl_quickshift_get_kernel_size(qs.p))
}

// https://www.vlfeat.org/api/quickshift_8h.html#a199679a6a594f25a690c97d5c67515a0
func (qs *QuickShift) GetMedoid() bool {
	return C.vl_quickshift_get_medoid(qs.p) != 0
}

// https://www.vlfeat.org/api/quickshift_8h.html#ab116d2dbad717ce889e56ea8335a1463
func (qs *QuickShift) GetParents() [][]int {
	cParents := C.vl_quickshift_get_parents(qs.p)
	width := qs.GetWidth()
	height := qs.GetHeight()
	length := width * height
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cParents)),
		Len:  length,
		Cap:  length,
	}
	cParentsSlice := *(*[]C.int)(unsafe.Pointer(&hdr))
	parents := make([][]int, height)
	for i := 0; i < height; i++ {
		parents[i] = make([]int, width)
		for j := 0; j < width; j++ {
			parents[i][j] = int(cParentsSlice[i*width+j])
		}
	}
	return parents
}

// https://www.vlfeat.org/api/quickshift_8h.html#a675666418f0bc2665be61f0b243a8c1d
func (qs *QuickShift) GetDists() [][]float64 {
	cDists := C.vl_quickshift_get_dists(qs.p)
	width := qs.GetWidth()
	height := qs.GetHeight()
	length := width * height
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cDists)),
		Len:  length,
		Cap:  length,
	}
	cDistsSlice := *(*[]C.double)(unsafe.Pointer(&hdr))
	parents := make([][]float64, height)
	for i := 0; i < height; i++ {
		parents[i] = make([]float64, width)
		for j := 0; j < width; j++ {
			parents[i][j] = float64(cDistsSlice[i*width+j])
		}
	}
	return parents
}

// https://www.vlfeat.org/api/quickshift_8h.html#acceb1733e008e0542164429e92253837
func (qs *QuickShift) GetDensity() [][]float64 {
	cDensity := C.vl_quickshift_get_dists(qs.p)
	width := qs.GetWidth()
	height := qs.GetHeight()
	length := width * height
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cDensity)),
		Len:  length,
		Cap:  length,
	}
	cDensitySlice := *(*[]C.double)(unsafe.Pointer(&hdr))
	parents := make([][]float64, height)
	for i := 0; i < height; i++ {
		parents[i] = make([]float64, width)
		for j := 0; j < width; j++ {
			parents[i][j] = float64(cDensitySlice[i*width+j])
		}
	}
	return parents
}

func (qs *QuickShift) GetWidth() int {
	return int(qs.p.width)
}

func (qs *QuickShift) GetHeight() int {
	return int(qs.p.height)
}
