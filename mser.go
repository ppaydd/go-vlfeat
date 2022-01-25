package vlfeat

/*
#include <stdlib.h>
#include <mser.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

// Mser stats struct
type MserStats struct {
	NumExtremal    int `json:"num_extremal"`
	NumUnstable    int `json:"num_unstable"`
	NumAbsUnstable int `json:"num_abs_unstable"`
	NumTooBig      int `json:"num_too_big"`
	NumTooSmall    int `json:"num_too_small"`
	NumDuplicates  int `json:"num_duplicates"`
}

// function for mser stats struct conversion between go and C
func getMserStats(r *C.VlMserStats) MserStats {
	stats := MserStats{
		int(r.num_extremal),
		int(r.num_unstable),
		int(r.num_abs_unstable),
		int(r.num_too_big),
		int(r.num_too_small),
		int(r.num_duplicates)}
	return stats
}

// mser algorithm
type Mser struct {
	p *C.VlMserFilt
}

// https://www.vlfeat.org/api/mser_8c.html#af6dbdcb894693e90c43d51140d17cb9c
func NewMser(dims []int) Mser {
	dimLength := len(dims)
	cDims := make([]C.int, dimLength)
	for i, dim := range dims {
		cDims[i] = C.int(dim)
	}
	p := C.vl_mser_new(C.int(dimLength), &cDims[0])
	return Mser{p: p}
}

// https://www.vlfeat.org/api/mser_8c.html#a3d94ff216cb9389b49dc5799e26ad3ba
func (mser *Mser) Delete() {
	C.vl_mser_delete(mser.p)
}

// https://www.vlfeat.org/api/mser_8c.html#ae50c576bc27a1ee7837cbbe6e5089583
func (mser *Mser) Process(img []uint8) {
	imgPtr := toCUcharArrayPtr(img)
	C.vl_mser_process(mser.p, imgPtr)
}

// https://www.vlfeat.org/api/mser_8c.html#aeeee08edd486e41126316f1d3bf90013
func (mser *Mser) EllFit() {
	C.vl_mser_ell_fit(mser.p)
}

// Retrieving data

// https://www.vlfeat.org/api/mser_8h.html#a40ccc21a2849c0aec7d5f05205843402
func (mser *Mser) GetRegionsNum() uint {
	return uint(C.vl_mser_get_regions_num(mser.p))
}

// https://www.vlfeat.org/api/mser_8h.html#a3e79dcb72a6ff1b0d7695257149017fa
func (mser *Mser) GetRegions() []uint {
	length := mser.GetRegionsNum()
	cRegions := C.vl_mser_get_regions(mser.p)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cRegions)),
		Len:  int(length),
		Cap:  int(length),
	}
	regionsSlice := *(*[]C.uint)(unsafe.Pointer(&hdr))

	regions := make([]uint, length)
	for i, region := range regionsSlice {
		regions[i] = uint(region)
	}
	return regions
}

// https://www.vlfeat.org/api/mser_8h.html#ad2e8f40c8cd872941cb1756cac334cf0
func (mser *Mser) GetEll() []float32 {
	length := int(mser.GetEllNum())
	cEll := C.vl_mser_get_ell(mser.p)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cEll)),
		Len:  length,
		Cap:  length,
	}
	cEllSlice := *(*[]C.float)(unsafe.Pointer(&hdr))

	ell := make([]float32, length)
	for i, e := range cEllSlice {
		ell[i] = float32(e)
	}
	return ell
}

// https://www.vlfeat.org/api/mser_8h.html#a1f92bdb8142fdc6f256ebb98adac1d23
func (mser *Mser) GetEllNum() uint {
	return uint(C.vl_mser_get_ell_num(mser.p))
}

// https://www.vlfeat.org/api/mser_8h.html#a275a5df25b8b628f0e5315758f4e9b44
func (mser *Mser) GetEllDof() uint {
	return uint(C.vl_mser_get_ell_dof(mser.p))
}

// https://www.vlfeat.org/api/mser_8h.html#abe6d998cfec12cb5a5eeb7d3f05e3245
func (mser *Mser) GetStats() MserStats {
	statsPtr := C.vl_mser_get_stats(mser.p)
	return getMserStats(statsPtr)
}

// Retrieving parameters

// https://www.vlfeat.org/api/mser_8h.html#a9aa9f3041186a969cb35268d4ab656bb
func (mser *Mser) GetDelta() uint8 {
	return uint8(C.vl_mser_get_delta(mser.p))
}

// https://www.vlfeat.org/api/mser_8h.html#a82b28392b17bfd5bd6765e1cab61f049
func (mser *Mser) GetMinArea() float64 {
	return float64(C.vl_mser_get_min_area(mser.p))
}

// https://www.vlfeat.org/api/mser_8h.html#a9598f59824914b8cd39ff384cdf2d708
func (mser *Mser) GetMaxArea() float64 {
	return float64(C.vl_mser_get_max_area(mser.p))
}

// https://www.vlfeat.org/api/mser_8h.html#a546850260c72fa3ef1d5b5887a1240ad
func (mser *Mser) GetMaxVariation() float64 {
	return float64(C.vl_mser_get_max_variation(mser.p))
}

// https://www.vlfeat.org/api/mser_8h.html#a420d10de58b8e7878c6e4708f223493f
func (mser *Mser) GetMinDiversity() float64 {
	return float64(C.vl_mser_get_min_diversity(mser.p))
}

// Setting parameters

// https://www.vlfeat.org/api/mser_8h.html#a9aa9f3041186a969cb35268d4ab656bb
func (mser *Mser) SetDelta(x uint8) {
	C.vl_mser_set_delta(mser.p, C.uchar(x))
}

// https://www.vlfeat.org/api/mser_8h.html#a25744f49395d441a8769b9f605480f4c
func (mser *Mser) SetMinArea(x float64) {
	C.vl_mser_set_min_area(mser.p, C.double(x))
}

// https://www.vlfeat.org/api/mser_8h.html#af517861a633fbd7cc5637ea4baef0db6
func (mser *Mser) SetMaxArea(x float64) {
	C.vl_mser_set_max_area(mser.p, C.double(x))
}

// https://www.vlfeat.org/api/mser_8h.html#a7feb4b477fe6340639d40461174cfabe
func (mser *Mser) SetMaxVariation(x float64) {
	C.vl_mser_set_max_variation(mser.p, C.double(x))
}

// https://www.vlfeat.org/api/mser_8h.html#abb7fe6fc94bb7503bdc0dd160a549de1
func (mser *Mser) SetMinDiversity(x float64) {
	C.vl_mser_set_min_diversity(mser.p, C.double(x))
}
