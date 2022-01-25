package vlfeat

/*
#include <stdlib.h>
#include <dsift.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

// Dsift keypoint struct
type DsiftKeypoint struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	S    float64 `json:"s"`
	Norm float64 `json:"norm"`
}

type DsiftDescriptorGeometry struct {
	NumBinT  int `json:"numBinT"`
	NumBinX  int `json:"numBinX"`
	NumBinY  int `json:"numBinY"`
	BinSizeX int `json:"binSizeX"`
	BinSizeY int `json:"binSizeY"`
}

// function for dsift keypoint struct conversion between go and C
func getDsiftKeyPoints(ret *C.VlDsiftKeypoint, length int) []DsiftKeypoint {
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(ret)),
		Len:  length,
		Cap:  length,
	}
	s := *(*[]C.VlDsiftKeypoint)(unsafe.Pointer(&hdr))

	keys := make([]DsiftKeypoint, length)
	for i, r := range s {
		keys[i] = DsiftKeypoint{float64(r.x), float64(r.y), float64(r.s), float64(r.norm)}
	}
	return keys
}

// dsift algorithm

type Dsift struct {
	p *C.VlDsiftFilter
}

// https://www.vlfeat.org/api/dsift_8c.html#aa9ba7ffaa72c137c457642ce833dab05
func NewDsift(imWidth, imHeight int) Dsift {
	p := C.vl_dsift_new(C.int(imWidth), C.int(imHeight))
	return Dsift{p: p}
}

// https://www.vlfeat.org/api/dsift_8c.html#aa025e58a852d8df078c6b74b8136c704
func NewDsiftBaic(imWidth, imHeight, step, binSize int) Dsift {
	p := C.vl_dsift_new_basic(C.int(imWidth), C.int(imHeight), C.int(step), C.int(binSize))
	return Dsift{p: p}
}

// https://www.vlfeat.org/api/dsift_8c.html#aa123f1d9e79ab01882646f713dfb4f0c
func (dsift *Dsift) Delete() {
	C.vl_dsift_delete(dsift.p)
}

// https://www.vlfeat.org/api/dsift_8c.html#a09d5525ad7e16e2b9f3f1b9d273c85f6
func (dsift *Dsift) Process(img []float32) {
	imgPtr := toCFloatArrayPtr(img)
	C.vl_dsift_process(dsift.p, imgPtr)
}

// set parameters
// https://www.vlfeat.org/api/dsift_8h.html#a42ae6bf77a9b737fd1e45ad5c43263dd
func (dsift *Dsift) SetSteps(stepX, stepY int) {
	C.vl_dsift_set_steps(dsift.p, C.int(stepX), C.int(stepY))
}

// https://www.vlfeat.org/api/dsift_8h.html#a7d34c8e257c873f2ed580b046296d1ac
func (dsift *Dsift) SetBounds(minX, minY, maxX, maxY int) {
	C.vl_dsift_set_bounds(dsift.p, C.int(minX), C.int(minY), C.int(maxX), C.int(maxY))
}

// https://www.vlfeat.org/api/dsift_8h.html#a930f20c25eab08d9490830b0a358ff2b
func (dsift *Dsift) SetGeometry(geom DsiftDescriptorGeometry) {
	cGeom := C.VlDsiftDescriptorGeometry{
		numBinT:  C.int(geom.NumBinT),
		numBinX:  C.int(geom.NumBinX),
		numBinY:  C.int(geom.NumBinY),
		binSizeX: C.int(geom.BinSizeX),
		binSizeY: C.int(geom.BinSizeY),
	}
	C.vl_dsift_set_geometry(dsift.p, &cGeom)
}

// https://www.vlfeat.org/api/dsift_8h.html#a8dfe2d20dbe9885d0c139c5b81b5f4b0
func (dsift *Dsift) SetFlatWindow(useFlatWindow bool) {
	useFlatWindowNum := 0
	if useFlatWindow {
		useFlatWindowNum = 1
	}
	C.vl_dsift_set_flat_window(dsift.p, C.int(useFlatWindowNum))
}

// https://www.vlfeat.org/api/dsift_8h.html#ae60fa31ff4df09e8025525714dac9563
func (dsift *Dsift) SetWindowSize(windowSize float64) {
	C.vl_dsift_set_window_size(dsift.p, C.double(windowSize))
}

// Retrieving data and parameters

// https://www.vlfeat.org/api/dsift_8h.html#aeade8b18c21954f00d76f5dcb12b9bfe
func (dsift *Dsift) GetDescriptorSize() int {
	return int(C.vl_dsift_get_descriptor_size(dsift.p))
}

// https://www.vlfeat.org/api/dsift_8h.html#a06f036e38fb68d2237dd60efb6f21236
func (dsift *Dsift) GetDescriptors() []float32 {
	length := dsift.GetDescriptorSize()
	cDesc := C.vl_dsift_get_descriptors(dsift.p)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cDesc)),
		Len:  length,
		Cap:  length,
	}
	cDescSlice := *(*[]C.float)(unsafe.Pointer(&hdr))
	desc := make([]float32, length)
	for i, des := range cDescSlice {
		desc[i] = float32(des)
	}
	return desc
}

// https://www.vlfeat.org/api/dsift_8h.html#a3b5fabb1496fc91a70669d4201f47a5b
func (dsift *Dsift) GetKeypointNum() int {
	return int(C.vl_dsift_get_keypoint_num(dsift.p))
}

// https://www.vlfeat.org/api/dsift_8h.html#a5b586c60d079a65ded73c4fd2387d6bf
func (dsift *Dsift) GetKeypoints() []DsiftKeypoint {
	cKeypoints := C.vl_dsift_get_keypoints(dsift.p)
	return getDsiftKeyPoints(cKeypoints, dsift.GetKeypointNum())
}

// https://www.vlfeat.org/api/dsift_8h.html#a92580ab01e5fb7967b1fa69abde717d0
func (dsift *Dsift) GetBounds() (int, int, int, int) {
	var cMinX, cMinY, cMaxX, cMaxY C.int
	C.vl_dsift_get_bounds(dsift.p, &cMinX, &cMinY, &cMaxX, &cMaxY)
	return int(cMinX), int(cMinY), int(cMaxX), int(cMaxY)
}

// https://www.vlfeat.org/api/dsift_8h.html#ac2ba5f78ac2675e547a4ef36c1bf4654
func (dsift *Dsift) GetSteps() (int, int) {
	var cStepX, cStepY C.int
	C.vl_dsift_get_steps(dsift.p, &cStepX, &cStepY)
	return int(cStepX), int(cStepY)
}

// https://www.vlfeat.org/api/dsift_8h.html#a05dc468116dc64d75b45853ead4c5031
func (dsift *Dsift) GetGeometry() DsiftDescriptorGeometry {
	cGeom := C.vl_dsift_get_geometry(dsift.p)
	return DsiftDescriptorGeometry{
		NumBinT:  int(cGeom.numBinT),
		NumBinX:  int(cGeom.numBinX),
		NumBinY:  int(cGeom.numBinY),
		BinSizeX: int(cGeom.binSizeX),
		BinSizeY: int(cGeom.binSizeY),
	}
}

// https://www.vlfeat.org/api/dsift_8h.html#a36450893fad7e5d9bb4ea003b05ea6d2
func (dsift *Dsift) GetFlatWindow() bool {
	flatWindow := C.vl_dsift_get_flat_window(dsift.p)
	return int(flatWindow) != 0
}

// https://www.vlfeat.org/api/dsift_8h.html#afa8d4a02e6f7a0d8b89b585c689e164e
func (dsift *Dsift) GetWindowSize() float64 {
	return float64(C.vl_dsift_get_window_size(dsift.p))
}
