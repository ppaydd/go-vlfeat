package vlfeat

/*
#include <stdlib.h>
#include <sift.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

// sift keypoint struct
type SiftKeypoint struct {
	O     int     `json:"o"`
	Ix    int     `json:"ix"`
	Iy    int     `json:"iy"`
	Is    int     `json:"is"`
	X     float32 `json:"x"`
	Y     float32 `json:"y"`
	S     float32 `json:"s"`
	Sigma float32 `json:"sigma"`
}

// function for conversion between go and C
func getSiftKeyPoints(ret *C.VlSiftKeypoint, length int) []SiftKeypoint {
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(ret)),
		Len:  length,
		Cap:  length,
	}
	s := *(*[]C.VlSiftKeypoint)(unsafe.Pointer(&hdr))

	keys := make([]SiftKeypoint, length)
	for i, r := range s {
		keys[i] = SiftKeypoint{int(r.o), int(r.ix), int(r.iy), int(r.is), float32(r.x), float32(r.y), float32(r.s), float32(r.sigma)}
	}
	return keys
}

func toCSiftKeypoint(keypoint SiftKeypoint) C.VlSiftKeypoint {
	cKeypoint := C.VlSiftKeypoint{
		o:     C.int(keypoint.O),
		ix:    C.int(keypoint.Ix),
		iy:    C.int(keypoint.Iy),
		is:    C.int(keypoint.Is),
		x:     C.float(keypoint.X),
		y:     C.float(keypoint.Y),
		s:     C.float(keypoint.S),
		sigma: C.float(keypoint.Sigma),
	}
	return cKeypoint
}

// sift algorithm

type Sift struct {
	p *C.VlSiftFilt
}

// https://www.vlfeat.org/api/sift_8c.html#adff66a155e30ed412bc8bbb97dfa2fae
func NewSift(width, height, noctaves, nlevels, o_min int) Sift {
	p := C.vl_sift_new(C.int(width), C.int(height), C.int(noctaves), C.int(nlevels), C.int(o_min))
	return Sift{p: p}
}

// https://www.vlfeat.org/api/sift_8c.html#ab242293326626641411e7d7f43a109b2
func (sift *Sift) Delete() {
	C.vl_sift_delete(sift.p)
}

// https://www.vlfeat.org/api/sift_8c.html#a97cca9a09efaadc9dd0671912b9d5e05
func (sift *Sift) ProcessFirstOctave(img []float32) VlErrorType {
	imgPtr := toCFloatArrayPtr(img)
	return VlErrorType(C.vl_sift_process_first_octave(sift.p, imgPtr))
}

// https://www.vlfeat.org/api/sift_8c.html#a610cab1a3bf7d38e389afda9037f14da
func (sift *Sift) ProcessNextOctave() VlErrorType {
	return VlErrorType(C.vl_sift_process_next_octave(sift.p))
}

// https://www.vlfeat.org/api/sift_8c.html#a65c55820964f4f6609ca9ef1d547b2c4
func (sift *Sift) Detect() {
	C.vl_sift_detect(sift.p)
}

// https://www.vlfeat.org/api/sift_8c.html#a919c860a1c8db300a6e3b960976fad70
func (sift *Sift) CalcKeypointOrientations(keypoint SiftKeypoint) (int, []float64) {
	ckeypoint := toCSiftKeypoint(keypoint)
	cAngles := make([]C.double, 4)
	angleCount := int(C.vl_sift_calc_keypoint_orientations(sift.p, &cAngles[0], &ckeypoint))
	angles := make([]float64, 4)
	for i, angle := range cAngles {
		angles[i] = float64(angle)
	}
	return angleCount, angles
}

// https://www.vlfeat.org/api/sift_8c.html#a85f3878a53ef7151b569c1b3ea4d13b6
// descLength is descr(result) array length
// The function fills the buffer descr which must be large enough to hold the descriptor.
func (sift *Sift) CalcKeypointDescriptor(descLength int, keypoint SiftKeypoint, angle float64) []float32 {
	ckeypoint := toCSiftKeypoint(keypoint)
	cDesc := make([]C.float, descLength)
	C.vl_sift_calc_keypoint_descriptor(sift.p, &cDesc[0], &ckeypoint, C.double(angle))
	desc := make([]float32, descLength)
	for i, des := range cDesc {
		desc[i] = float32(des)
	}
	return desc
}

// https://www.vlfeat.org/api/sift_8c.html#a335f3295ba77b3bb937e5272fe1a02fc
// descLength is descr(result) array length
// // The function fills the buffer descr which must be large enough to hold the descriptor.
func (sift *Sift) CalcRawDescriptor(img []float32, descLength, width, height int, x, y, s, angle float64) []float32 {
	cDesc := make([]C.float, descLength)
	imgPtr := toCFloatArrayPtr(img)
	C.vl_sift_calc_raw_descriptor(sift.p, imgPtr, &cDesc[0], C.int(width), C.int(height), C.double(x), C.double(y), C.double(s), C.double(angle))
	desc := make([]float32, descLength)
	for i, des := range cDesc {
		desc[i] = float32(des)
	}
	return desc
}

// https://www.vlfeat.org/api/sift_8c.html#a6f3fc8e38b6c0c520cb90b1a63ddc031
func (sift *Sift) KeypointInit(x, y, sigma float64) SiftKeypoint {
	var ckeypoint C.VlSiftKeypoint
	C.vl_sift_keypoint_init(sift.p, &ckeypoint, C.double(x), C.double(y), C.double(sigma))
	keypoint := SiftKeypoint{
		O:     int(ckeypoint.o),
		Ix:    int(ckeypoint.ix),
		Iy:    int(ckeypoint.iy),
		Is:    int(ckeypoint.is),
		X:     float32(ckeypoint.x),
		Y:     float32(ckeypoint.y),
		S:     float32(ckeypoint.s),
		Sigma: float32(ckeypoint.sigma),
	}
	return keypoint
}

// Retrieve data and parameters

// https://www.vlfeat.org/api/sift_8h.html#a70186e579c8eff1bcabf408f46169cad
func (sift *Sift) GetOctaveIndex() int {
	return int(C.vl_sift_get_octave_index(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#a5e0cd96b3985635b82adabc3ce8b2242
func (sift *Sift) GetNoctaves() int {
	return int(C.vl_sift_get_noctaves(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#aa3db07e91c86f992c31b8e2335a760a9
func (sift *Sift) GetOctaveFirst() int {
	return int(C.vl_sift_get_octave_first(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#a89bd76ab5c1e584ff8e46dfdc93ea748
func (sift *Sift) GetOctaveWidth() int {
	return int(C.vl_sift_get_octave_width(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#a9769e8f6d84ec75804e873229526eb10
func (sift *Sift) GetOctaveHeight() int {
	return int(C.vl_sift_get_octave_height(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#a751c116352e72eed8a111e7c1e06a18e
func (sift *Sift) GetNlevels() int {
	return int(C.vl_sift_get_nlevels(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#aa45b8e7413384c7d6525f439e68856fe
func (sift *Sift) GetNkeypoints() int {
	return int(C.vl_sift_get_nkeypoints(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#a08959e6a90c98bf397e3430e79a6ea9c
func (sift *Sift) GetPeakThresh() float64 {
	return float64(C.vl_sift_get_peak_thresh(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#adb5b0159af92e1ce1462ddaaaa55a747
func (sift *Sift) GetEdgeThresh() float64 {
	return float64(C.vl_sift_get_edge_thresh(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#ad4b57c390ca004dc56b0b0b1abf0c7a9
func (sift *Sift) GetNormThresh() float64 {
	return float64(C.vl_sift_get_norm_thresh(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#ae0272723812d5072619475d4787be78e
func (sift *Sift) GetMagnif() float64 {
	return float64(C.vl_sift_get_magnif(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#a76053a5e655b9577995fea0fbe429078
func (sift *Sift) GetWindowSize() float64 {
	return float64(C.vl_sift_get_window_size(sift.p))
}

// https://www.vlfeat.org/api/sift_8h.html#a400759060e87dc7a6264555b90b0a221
func (sift *Sift) GetOctave(s int) []float32 {
	width := sift.GetOctaveWidth()
	height := sift.GetOctaveHeight()
	cOctave := C.vl_sift_get_octave(sift.p, C.int(s))
	length := width * height

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cOctave)),
		Len:  length,
		Cap:  length,
	}
	cOctaveSlice := *(*[]C.float)(unsafe.Pointer(&hdr))
	octave := make([]float32, length)
	for i, r := range cOctaveSlice {
		octave[i] = float32(r)
	}
	return octave
}

// https://www.vlfeat.org/api/sift_8c.html#a65c55820964f4f6609ca9ef1d547b2c4
func (sift *Sift) GetKeypoints() []SiftKeypoint {
	ckeypoints := C.vl_sift_get_keypoints(sift.p)
	length := sift.GetNkeypoints()
	return getSiftKeyPoints(ckeypoints, length)
}

// https://www.vlfeat.org/api/sift_8h.html#af69118a1c5d4d17bccac87d11fe8ce8f
func (sift *Sift) SetPeakThresh(t float64) {
	C.vl_sift_set_peak_thresh(sift.p, C.double(t))
}

// https://www.vlfeat.org/api/sift_8h.html#ab7173b402b85de43ebf36fcabde77508
func (sift *Sift) SetEdgeThresh(t float64) {
	C.vl_sift_set_edge_thresh(sift.p, C.double(t))
}

// https://www.vlfeat.org/api/sift_8h.html#a86703f33aad31638909acd9697f93115
func (sift *Sift) SetNormThresh(t float64) {
	C.vl_sift_set_norm_thresh(sift.p, C.double(t))
}

// https://www.vlfeat.org/api/sift_8h.html#a595579dd7952807c074c5311a6500121
func (sift *Sift) SetMagnif(m float64) {
	C.vl_sift_set_magnif(sift.p, C.double(m))
}

// https://www.vlfeat.org/api/sift_8h.html#af5996cc6171c6e3c8810fb400abbad21
func (sift *Sift) SetWindowSize(m float64) {
	C.vl_sift_set_window_size(sift.p, C.double(m))
}
