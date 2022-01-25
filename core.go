package vlfeat

/*
#include <stdlib.h>
*/
import "C"

type VlErrorType int

const (
	VlErrorOK       VlErrorType = 0
	VlErrorOverflow VlErrorType = 1
	VlErrorAlloc    VlErrorType = 2
	VlErrorBadArg   VlErrorType = 3
	VlErrorIO       VlErrorType = 4
	VlErrorEOF      VlErrorType = 5
)

func toCFloatArrayPtr(img []float32) *C.float {
	datas := make([]C.float, len(img))
	for i, goData := range img {
		datas[i] = C.float(goData)
	}
	return &datas[0]
}
