package vlfeat

/*
#include <stdlib.h>
#include <ikmeans.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type VlIKMAlgorithms int

const (
	VL_IKM_LLOYD VlIKMAlgorithms = 0
	VL_IKM_ELKAN VlIKMAlgorithms = 1
)

type IKM struct {
	p *C.VlIKMFilt
}

/* Create and destroy */

// https://www.vlfeat.org/api/ikmeans_8h.html#af5a42441a336dd73c39c5a1ec8028444
func NewIKM(method VlIKMAlgorithms) IKM {
	p := C.vl_ikm_new(C.int(method))
	return IKM{p: p}
}

// https://www.vlfeat.org/api/ikmeans_8h.html#a33a659152e03286d390aa17b1a212a7b
func (ikm *IKM) Delete() {
	C.vl_ikm_delete(ikm.p)
}

/* Process data */

func (ikm *IKM) Init(centers []int, M, K uint) {
	cCenters := make([]C.int, len(centers))
	for i, center := range centers {
		cCenters[i] = C.int(center)
	}
	C.vl_ikm_init(ikm.p, &cCenters[0], C.uint(M), C.uint(K))
}

func (ikm *IKM) InitRand(M, K uint) {
	C.vl_ikm_init_rand(ikm.p, C.uint(M), C.uint(K))
}

func (ikm *IKM) InitRandData(data []uint8, M, N, K uint) {
	dataPtr := toCUcharArrayPtr(data)
	C.vl_ikm_init_rand_data(ikm.p, dataPtr, C.uint(M), C.uint(N), C.uint(K))
}

// https://www.vlfeat.org/api/ikmeans_8h.html#a58ccd60ab79dbfe8b9823a1caede7879
func (ikm *IKM) Train(data []uint8, N uint) {
	dataPtr := toCUcharArrayPtr(data)
	C.vl_ikm_train(ikm.p, dataPtr, C.uint(N))
}

// https://www.vlfeat.org/api/ikmeans_8h.html#ac7e2bc8d34d514a4fe43b57d12dc244a
func (ikm *IKM) Push(data []uint8, N uint) []uint {
	dataPtr := toCUcharArrayPtr(data)
	asgn := make([]uint, N)
	cAsgn := make([]C.uint, N)
	C.vl_ikm_push(ikm.p, &cAsgn[0], dataPtr, C.uint(N))
	for i, angnData := range cAsgn {
		asgn[i] = uint(angnData)
	}
	return asgn
}

// https://www.vlfeat.org/api/ikmeans_8h.html#a03f0ed5b6f3680b1f872c01e47b36cb8
func PushOne(centers []int, data []uint8, M, K uint) uint {
	cCenters := make([]C.int, len(centers))
	for i, center := range centers {
		cCenters[i] = C.int(center)
	}

	dataPtr := toCUcharArrayPtr(data)
	return uint(C.vl_ikm_push_one(&cCenters[0], dataPtr, C.uint(M), C.uint(K)))
}

/* Retrieve data and parameters */

// https://www.vlfeat.org/api/ikmeans_8h.html#a10787adeb3ac4c3d6cce7a2e3598525f
func (ikm *IKM) GetNdims() uint {
	return uint(C.vl_ikm_get_ndims(ikm.p))
}

// https://www.vlfeat.org/api/ikmeans_8h.html#aba2ce262018aae7049c53a33af0db8e1
func (ikm *IKM) GetK() uint {
	return uint(C.vl_ikm_get_K(ikm.p))
}

// https://www.vlfeat.org/api/ikmeans_8h.html#afe77d60b7c64f773b9e5edbcebca6e21
func (ikm *IKM) GetVerbosity() int {
	return int(C.vl_ikm_get_verbosity(ikm.p))
}

// https://www.vlfeat.org/api/ikmeans_8h.html#ad1dcf031fd04bdbaf9d389dcd960c953
func (ikm *IKM) GetMaxNiters() uint {
	return uint(C.vl_ikm_get_max_niters(ikm.p))
}

// https://www.vlfeat.org/api/ikmeans_8h.html#a56e13744afc5e1917f4c036fe1c36678
func (ikm *IKM) GetCenters() []int {
	M := ikm.GetNdims()
	K := ikm.GetK()
	length := M * K
	cCenterPtr := C.vl_ikm_get_centers(ikm.p)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cCenterPtr)),
		Len:  int(length),
		Cap:  int(length),
	}
	cCenterSlice := *(*[]C.int)(unsafe.Pointer(&hdr))
	centers := make([]int, length)
	for i, center := range cCenterSlice {
		centers[i] = int(center)
	}
	return centers
}

// https://www.vlfeat.org/api/ikmeans_8h.html#adbf9d13d091a5bec378023de11bfa690
func (ikm *IKM) SetVerbosity(verb int) {
	C.vl_ikm_set_verbosity(ikm.p, C.int(verb))
}

// https://www.vlfeat.org/api/ikmeans_8h.html#a25053477cd794f1a0fc3095a08640c62
func (ikm *IKM) SetMaxNiters(maxNiters uint) {
	C.vl_ikm_set_max_niters(ikm.p, C.uint(maxNiters))
}
