package vlfeat

/*
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"reflect"
	"unsafe"
)

type VlErrorType int

const (
	VlErrorOK       VlErrorType = 0
	VlErrorOverflow VlErrorType = 1
	VlErrorAlloc    VlErrorType = 2
	VlErrorBadArg   VlErrorType = 3
	VlErrorIO       VlErrorType = 4
	VlErrorEOF      VlErrorType = 5
)

type VlType uint

const (
	VlTypeFloat  VlType = 1
	VlTypeDouble VlType = 2
	VlTypeInt8   VlType = 3
	VlTypeUint8  VlType = 4
	VlTypeInt16  VlType = 5
	VlTypeUint16 VlType = 6
	VlTypeInt32  VlType = 7
	VlTypeUint32 VlType = 8
	VlTypeInt64  VlType = 9
	VlTypeUint64 VlType = 10
)

// img data switch to C
func toCFloatArrayPtr(img []float32) *C.float {
	datas := make([]C.float, len(img))
	for i, goData := range img {
		datas[i] = C.float(goData)
	}
	return &datas[0]
}

func toCUcharArrayPtr(img []uint8) *C.uchar {
	datas := make([]C.uchar, len(img))
	for i, goData := range img {
		datas[i] = C.uchar(goData)
	}
	return &datas[0]
}

func ToCVlTypeArrayPtr(data interface{}, vlType VlType) (unsafe.Pointer, int, error) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice, reflect.Array:
		slice := reflect.ValueOf(data)
		length := slice.Len()
		switch vlType {
		case VlTypeFloat:
			cData := make([]C.float, length)
			for i := 0; i < length; i++ {
				cData[i] = C.float(slice.Index(i).Float())
			}
			return unsafe.Pointer(&cData), length, nil
		case VlTypeDouble:
			cData := make([]C.double, length)
			for i := 0; i < length; i++ {
				cData[i] = C.double(slice.Index(i).Float())
			}
			return unsafe.Pointer(&cData), length, nil
		case VlTypeInt8:
			cData := make([]C.char, length)
			for i := 0; i < length; i++ {
				cData[i] = C.char(slice.Index(i).Int())
			}
			return unsafe.Pointer(&cData), length, nil
		case VlTypeUint8:
			cData := make([]C.uchar, length)
			for i := 0; i < length; i++ {
				cData[i] = C.uchar(slice.Index(i).Int())
			}
			return unsafe.Pointer(&cData), length, nil
		case VlTypeInt16:
			cData := make([]C.short, length)
			for i := 0; i < length; i++ {
				cData[i] = C.short(slice.Index(i).Int())
			}
			return unsafe.Pointer(&cData), length, nil
		case VlTypeUint16:
			cData := make([]C.ushort, length)
			for i := 0; i < length; i++ {
				cData[i] = C.ushort(slice.Index(i).Int())
			}
			return unsafe.Pointer(&cData), length, nil
		case VlTypeInt32:
			cData := make([]C.int, length)
			for i := 0; i < length; i++ {
				cData[i] = C.int(slice.Index(i).Int())
			}
			return unsafe.Pointer(&cData), length, nil
		case VlTypeUint32:
			cData := make([]C.uint, length)
			for i := 0; i < length; i++ {
				cData[i] = C.uint(slice.Index(i).Int())
			}
			return unsafe.Pointer(&cData), length, nil
		case VlTypeInt64:
			cData := make([]C.longlong, length)
			for i := 0; i < length; i++ {
				cData[i] = C.longlong(slice.Index(i).Int())
			}
			return unsafe.Pointer(&cData), length, nil
		case VlTypeUint64:
			cData := make([]C.ulonglong, length)
			for i := 0; i < length; i++ {
				cData[i] = C.ulonglong(slice.Index(i).Int())
			}
			return unsafe.Pointer(&cData), length, nil
		}
	}
	return nil, 0, errors.New("must be support slice or array")
}
