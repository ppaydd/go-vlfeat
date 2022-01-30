package vlfeat

/*
#include <stdlib.h>
#include <fisher.h>
*/
import "C"
import "unsafe"

type FisherFlag int

const (
	FinsherFlagSquareRoot FisherFlag = 1
	FinsherFlagNormalized FisherFlag = 2
	FinsherFlagImproved   FisherFlag = 3
	FinsherFlagFast       FisherFlag = 4
)

// https://www.vlfeat.org/api/fisher_8h.html#a4c13fe11e9847f9046c2636a6e77c1bd
func FisherEncode(dataType VlType, means interface{}, dimension, numClusters uint, covariances, priors, data interface{}, numData uint, flag FisherFlag) (uint, []float64, error) {
	encLength := 2 * int(dimension*numClusters)
	enc := make([]float64, encLength)
	cEnc := make([]C.double, encLength)
	meansPtr, _, err := ToCVlTypeArrayPtr(means, dataType)
	if err != nil {
		return 0, enc, err
	}
	covariancesPtr, _, err := ToCVlTypeArrayPtr(covariances, dataType)
	if err != nil {
		return 0, enc, err
	}
	priorsPtr, _, err := ToCVlTypeArrayPtr(priors, dataType)
	if err != nil {
		return 0, enc, err
	}
	dataPtr, _, err := ToCVlTypeArrayPtr(data, dataType)
	if err != nil {
		return 0, enc, err
	}
	encPtr := unsafe.Pointer(&cEnc)
	result := C.vl_fisher_encode(encPtr, C.vl_type(VlTypeFloat), meansPtr, C.uint(dimension), C.uint(numClusters), covariancesPtr, priorsPtr, dataPtr, C.uint(numData), C.int(flag))
	for i, des := range cEnc {
		enc[i] = float64(des)
	}
	return uint(result), enc, nil
}
