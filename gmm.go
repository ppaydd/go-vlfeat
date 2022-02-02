package vlfeat

/*
#include <stdlib.h>
#include <gmm.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type VlGMMInitialization int

const (
	VlGMMKMeans VlGMMInitialization = 0
	VlGMMRand   VlGMMInitialization = 1
	VlGMMCustom VlGMMInitialization = 2
)

type GMM struct {
	p *C.VlGMM
}

// https://www.vlfeat.org/api/gmm_8c.html#afe0bdce1cf97a7b64011ae58bc8b9697
func NewGMM(dataType VlType, dimension, numComponents uint) GMM {
	p := C.vl_gmm_new(C.uint(dataType), C.uint(dimension), C.uint(numComponents))
	return GMM{p: p}
}

// https://www.vlfeat.org/api/gmm_8c.html#a2c6889a0569271e72096d18735f28211
func (gmm *GMM) Copy() GMM {
	p := C.vl_gmm_new_copy(gmm.p)
	return GMM{p: p}
}

// https://www.vlfeat.org/api/gmm_8c.html#adbde533172047f780e2713d18387a9d0
func (gmm *GMM) Delete() {
	C.vl_gmm_delete(gmm.p)
}

// https://www.vlfeat.org/api/gmm_8c.html#a1baeb175e0cdb68c548addc837d7baae
func (gmm *GMM) Reset() {
	C.vl_gmm_reset(gmm.p)
}

// https://www.vlfeat.org/api/gmm_8c.html#a856cdfe758b7c48c3309859853f38e36
func (gmm *GMM) Cluster(data interface{}) (float64, error) {
	vltype := gmm.GetDataType()
	dataPtr, length, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return 0, err
	}
	return float64(C.vl_gmm_cluster(gmm.p, dataPtr, C.uint(length))), nil
}

// https://www.vlfeat.org/api/gmm_8c.html#aab9d461e2fca63960f2751ae86946804
func (gmm *GMM) InitWithRandData(data interface{}) error {
	vltype := gmm.GetDataType()
	dataPtr, length, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return err
	}
	C.vl_gmm_init_with_rand_data(gmm.p, dataPtr, C.uint(length))
	return nil
}

// https://www.vlfeat.org/api/gmm_8c.html#a21934aa27cd02d67734d311c4207829e
func (gmm *GMM) InitWithKmeans(data interface{}, kmeansInit Kmeans) error {
	vltype := gmm.GetDataType()
	dataPtr, length, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return err
	}
	C.vl_gmm_init_with_kmeans(gmm.p, dataPtr, C.uint(length), kmeansInit.p)
	return nil
}

// https://www.vlfeat.org/api/gmm_8c.html#a4f8f3fc91d284a2866fc5fd5d5b48dfd
func (gmm *GMM) Em(data interface{}) (float64, error) {
	vltype := gmm.GetDataType()
	dataPtr, length, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return 0, err
	}
	em := C.vl_gmm_em(gmm.p, dataPtr, C.uint(length))
	return float64(em), nil
}

// https://www.vlfeat.org/api/gmm_8c.html#a001f4a994fbb997e7b7b38d56e69e495
func (gmm *GMM) SetMeans(means interface{}) error {
	vltype := gmm.GetDataType()
	dataPtr, _, err := ToCVlTypeArrayPtr(means, vltype)
	if err != nil {
		return err
	}
	C.vl_gmm_set_means(gmm.p, dataPtr)
	return nil
}

// https://www.vlfeat.org/api/gmm_8c.html#a9467941c38e0eb70f1e7de2cf8242716
func (gmm *GMM) SetCovariances(covariances interface{}) error {
	vltype := gmm.GetDataType()
	dataPtr, _, err := ToCVlTypeArrayPtr(covariances, vltype)
	if err != nil {
		return err
	}
	C.vl_gmm_set_covariances(gmm.p, dataPtr)
	return nil
}

// https://www.vlfeat.org/api/gmm_8c.html#a4c0e4759f7b082400cfe4da0991139c8
func (gmm *GMM) SetPriors(priors interface{}) error {
	vltype := gmm.GetDataType()
	dataPtr, _, err := ToCVlTypeArrayPtr(priors, vltype)
	if err != nil {
		return err
	}
	C.vl_gmm_set_priors(gmm.p, dataPtr)
	return nil
}

// https://www.vlfeat.org/api/gmm_8c.html#a9801c187ad6f0561275a715d3e589a59
func (gmm *GMM) SetNumRepetitions(numRepetitions uint) {
	C.vl_gmm_set_num_repetitions(gmm.p, C.uint(numRepetitions))
}

// https://www.vlfeat.org/api/gmm_8c.html#a606ce33d200101fa6a60e10a3e89cec1
func (gmm *GMM) SetMaxNumIterations(maxNumIterations uint) {
	C.vl_gmm_set_max_num_iterations(gmm.p, C.uint(maxNumIterations))
}

// https://www.vlfeat.org/api/gmm_8c.html#a96708e3c10107c49f136a0fe1082684a
func (gmm *GMM) SetVerbosity(verbosity int) {
	C.vl_gmm_set_verbosity(gmm.p, C.int(verbosity))
}

// https://www.vlfeat.org/api/gmm_8c.html#a3f34a10ef70b880a81eabce9fc29cc2e
func (gmm *GMM) SetInitialization(init VlGMMInitialization) {
	C.vl_gmm_set_verbosity(gmm.p, C.int(init))
}

// https://www.vlfeat.org/api/gmm_8c.html#ae740ca4d9c354ac9d83c89127bce744c
func (gmm *GMM) SetKmeansInitObject(kmeans Kmeans) {
	C.vl_gmm_set_kmeans_init_object(gmm.p, kmeans.p)
}

// https://www.vlfeat.org/api/gmm_8c.html#ac5b9aa600e348e99907cdb9f160c8d1d
func (gmm *GMM) SetCovarianceLowerBounds(bounds []float64) {
	cBounds := make([]C.double, len(bounds))
	for i, bound := range bounds {
		cBounds[i] = C.double(bound)
	}
	C.vl_gmm_set_covariance_lower_bounds(gmm.p, &cBounds[0])
}

// https://www.vlfeat.org/api/gmm_8c.html#a607ba2f82f3bfe5949f71acca0d332c8
func (gmm *GMM) SetCovarianceLowerBound(bound float64) {
	C.vl_gmm_set_covariance_lower_bound(gmm.p, C.double(bound))
}

// https://www.vlfeat.org/api/gmm_8c.html#ae8598b36a92ca066e26e03ea46dee466
func (gmm *GMM) GetMeans() unsafe.Pointer {
	return C.vl_gmm_get_means(gmm.p)
}

// https://www.vlfeat.org/api/gmm_8c.html#a319648cad7d83d2a31c7beb79b1ba125
func (gmm *GMM) GetCovariances() unsafe.Pointer {
	return C.vl_gmm_get_covariances(gmm.p)
}

// https://www.vlfeat.org/api/gmm_8c.html#a556605dddee3aa3f2c35898436c5e811
func (gmm *GMM) GetPriors() unsafe.Pointer {
	return C.vl_gmm_get_priors(gmm.p)
}

// https://www.vlfeat.org/api/gmm_8c.html#a54b1da9397fdb22c1036690e80816e76
func (gmm *GMM) GetPosteriors() unsafe.Pointer {
	return C.vl_gmm_get_posteriors(gmm.p)
}

// https://www.vlfeat.org/api/gmm_8c.html#ac8980d73a2bac3f4aba64108ea9b7986
func (gmm *GMM) GetDataType() VlType {
	return VlType(C.vl_gmm_get_data_type(gmm.p))
}

// https://www.vlfeat.org/api/gmm_8c.html#a5dd301238834c161a8ca142f32b2a83c
func (gmm *GMM) GetDimension() uint {
	return uint(C.vl_gmm_get_dimension(gmm.p))
}

// https://www.vlfeat.org/api/gmm_8c.html#aba79f6a785a9c33b84f5dc38c9b670aa
func (gmm *GMM) GetNumRepetitions() uint {
	return uint(C.vl_gmm_get_num_repetitions(gmm.p))
}

// https://www.vlfeat.org/api/gmm_8c.html#a83a3fb14322069092b4618bcabaf5fda
func (gmm *GMM) GetNumData() uint {
	return uint(C.vl_gmm_get_num_data(gmm.p))
}

// https://www.vlfeat.org/api/gmm_8c.html#a6d4f30630fd47af3dec762a58df7b653
func (gmm *GMM) GetNumClusters() uint {
	return uint(C.vl_gmm_get_num_clusters(gmm.p))
}

// https://www.vlfeat.org/api/gmm_8c.html#a618c40e1338f2433505f83d0745fce7f
func (gmm *GMM) GetLoglikelihood() float64 {
	return float64(C.vl_gmm_get_loglikelihood(gmm.p))
}

// https://www.vlfeat.org/api/gmm_8c.html#a167e3a1208bc756e1681dfc44dab5928
func (gmm *GMM) GetVerbosity() int {
	return int(C.vl_gmm_get_verbosity(gmm.p))
}

// https://www.vlfeat.org/api/gmm_8c.html#aaa311e17e0190fa3a247349d4e127e3d
func (gmm *GMM) GetMaxNumIterations() int {
	return int(C.vl_gmm_get_max_num_iterations(gmm.p))
}

// https://www.vlfeat.org/api/gmm_8c.html#a51e3a8e372e5d019c215d126823986b0
func (gmm *GMM) GetInitialization() VlGMMInitialization {
	return VlGMMInitialization(C.vl_gmm_get_initialization(gmm.p))
}

// https://www.vlfeat.org/api/gmm_8c.html#a6e149ef463f029ba96788f0aa7442025
func (gmm *GMM) GetCovarianceLowerBounds() []float64 {
	cBounds := C.vl_gmm_get_covariance_lower_bounds(gmm.p)
	length := int(gmm.GetDimension())
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cBounds)),
		Len:  int(length),
		Cap:  int(length),
	}
	cBoundSlice := *(*[]C.double)(unsafe.Pointer(&hdr))

	bounds := make([]float64, length)
	for i, bound := range cBoundSlice {
		bounds[i] = float64(bound)
	}
	return bounds
}

// https://www.vlfeat.org/api/gmm_8c.html#a6c35a5179c1a1061b0ed243d56e12dc7
func (gmm *GMM) GetKmeansInitObject() Kmeans {
	return Kmeans{p: C.vl_gmm_get_kmeans_init_object(gmm.p)}
}
