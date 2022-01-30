package vlfeat

/*
#include <stdlib.h>
#include <vlad.h>
*/
import "C"
import "unsafe"

type VladFlag int

const (
	VladFlagNormalizeComponents VladFlag = 1
	VladFlagSquareRoot          VladFlag = 2
	VladFlagUnnormalized        VladFlag = 3
	VladFlagNormalizeMass       VladFlag = 4
)

// https://www.vlfeat.org/api/vlad_8c.html#a6ee2926e14d4a76e9d99ba128bfe5a80
func VladEncode(dataType VlType, means interface{}, dimension, numClusters uint, data interface{}, numData uint, assignments interface{}, flag VladFlag) ([]float64, error) {
	encLength := int(dimension * numClusters)
	enc := make([]float64, encLength)
	cEnc := make([]C.double, encLength)
	meansPtr, _, err := ToCVlTypeArrayPtr(means, dataType)
	if err != nil {
		return enc, err
	}
	assignmentsPtr, _, err := ToCVlTypeArrayPtr(assignments, dataType)
	if err != nil {
		return enc, err
	}
	dataPtr, _, err := ToCVlTypeArrayPtr(data, dataType)
	if err != nil {
		return enc, err
	}
	encPtr := unsafe.Pointer(&cEnc)
	C.vl_vlad_encode(encPtr, C.vl_type(VlTypeFloat), meansPtr, C.uint(dimension), C.uint(numClusters), dataPtr, C.uint(numData), assignmentsPtr, C.int(flag))
	for i, des := range cEnc {
		enc[i] = float64(des)
	}
	return enc, nil
}
