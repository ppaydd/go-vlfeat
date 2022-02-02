package vlfeat

/*
#include <stdlib.h>
#include <hikmeans.h>
*/
import "C"

type HIKM struct {
	p *C.VlHIKMTree
}

/*  Create and destroy */

// https://www.vlfeat.org/api/hikmeans_8h.html#ae48c89b710a84568dcac00ddd763e604
func NewHIKM(method VlIKMAlgorithms) HIKM {
	p := C.vl_hikm_new(C.int(method))
	return HIKM{p: p}
}

// https://www.vlfeat.org/api/hikmeans_8h.html#a7830ea7acf7332e0eda8369b58853808
func (hikm *HIKM) Delete() {
	C.vl_hikm_delete(hikm.p)
}

/* Retrieve data and parameters */

// https://www.vlfeat.org/api/hikmeans_8h.html#a8bbbad989ed222178d7d80a7ef8a6a8c
func (hikm *HIKM) GetNdims() uint {
	return uint(C.vl_hikm_get_ndims(hikm.p))
}

// https://www.vlfeat.org/api/hikmeans_8h.html#ad01623c87b275ba47ab9dcf6679d3bae
func (hikm *HIKM) GetK() uint {
	return uint(C.vl_hikm_get_K(hikm.p))
}

// https://www.vlfeat.org/api/hikmeans_8h.html#a51b8a3da781958f50d4ef41dcbda2d32
func (hikm *HIKM) GetDepth() uint {
	return uint(C.vl_hikm_get_depth(hikm.p))
}

// https://www.vlfeat.org/api/hikmeans_8h.html#a41f0093905ef603a8b93e921522638e4
func (hikm *HIKM) GetVerbosity() int {
	return int(C.vl_hikm_get_verbosity(hikm.p))
}

// https://www.vlfeat.org/api/hikmeans_8h.html#af305ded1d09a18af4ea3d9562a2bdb1c
func (hikm *HIKM) GetMaxNiters() uint {
	return uint(C.vl_hikm_get_max_niters(hikm.p))
}

/* Set parameters */

// https://www.vlfeat.org/api/hikmeans_8h.html#a23789e4faee416be030f2a519de096e3
func (hikm *HIKM) SetVerbosity(verb int) {
	C.vl_hikm_set_verbosity(hikm.p, C.int(verb))
}

// https://www.vlfeat.org/api/hikmeans_8h.html#a431d5ac9cb5d5ae058dcc9785ad08fe5
func (hikm *HIKM) SetMaxNiters(maxNiters int) {
	C.vl_hikm_set_max_niters(hikm.p, C.int(maxNiters))
}

/* Process data */

func (hikm *HIKM) Init(M, K, depth uint) {
	C.vl_hikm_init(hikm.p, C.uint(M), C.uint(K), C.uint(depth))
}

// https://www.vlfeat.org/api/hikmeans_8h.html#ace8de873d52287e32918999864b93fbd
func (hikm *HIKM) Train(data []uint8, N uint) {
	dataPtr := toCUcharArrayPtr(data)
	C.vl_hikm_train(hikm.p, dataPtr, C.uint(N))
}

// https://www.vlfeat.org/api/hikmeans_8h.html#aacb98ccf6f8e9cc45dcfda9ca4f648c9
func (hikm *HIKM) Push(data []uint8, N uint) []uint {
	dataPtr := toCUcharArrayPtr(data)
	asgn := make([]uint, N)
	cAsgn := make([]C.uint, N)
	C.vl_hikm_push(hikm.p, &cAsgn[0], dataPtr, C.uint(N))
	for i, angnData := range cAsgn {
		asgn[i] = uint(angnData)
	}
	return asgn
}
