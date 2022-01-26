package vlfeat

/*
#include <stdlib.h>
#include <covdet.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

/* covdet method */
type CovDetMethod int

const (
	COVDET_METHOD_DOG                CovDetMethod = 1
	COVDET_METHOD_HESSIAN            CovDetMethod = 2
	COVDET_METHOD_HESSIAN_LAPLACE    CovDetMethod = 3
	COVDET_METHOD_HARRIS_LAPLACE     CovDetMethod = 4
	COVDET_METHOD_MULTISCALE_HESSIAN CovDetMethod = 5
	COVDET_METHOD_MULTISCALE_HARRIS  CovDetMethod = 6
	COVDET_METHOD_NUM                CovDetMethod = 7
)

/* covdet struct and conversion between go and C*/

type CovDetFrameOrientedEllipse struct {
	X   float32 `json:"x"`
	Y   float32 `json:"y"`
	A11 float32 `json:"a11"`
	A12 float32 `json:"a12"`
	A21 float32 `json:"a21"`
	A22 float32 `json:"a22"`
}

type CovDetFeature struct {
	Frame               CovDetFrameOrientedEllipse `json:"frame"`
	PeakScore           float32                    `json:"peakScore"`
	EdgeScore           float32                    `json:"edgeScore"`
	OrientationScore    float32                    `json:"orientationScore"`
	LaplacianScaleScore float32                    `json:"laplacianScaleScore"`
}

type CovDetFeatureOrientation struct {
	Angle float64 `json:"angle"`
	Score float64 `json:"score"`
}

type CovDetFeatureLaplacianScale struct {
	Scale float64 `json:"scale"`
	Score float64 `json:"score"`
}

func getCovDetFeatures(ret *C.VlCovDetFeature, length int) []CovDetFeature {
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(ret)),
		Len:  length,
		Cap:  length,
	}
	s := *(*[]C.VlCovDetFeature)(unsafe.Pointer(&hdr))

	keys := make([]CovDetFeature, length)
	for i, r := range s {
		frame := CovDetFrameOrientedEllipse{
			float32(r.frame.x),
			float32(r.frame.y),
			float32(r.frame.a11),
			float32(r.frame.a12),
			float32(r.frame.a21),
			float32(r.frame.a22),
		}
		keys[i] = CovDetFeature{
			frame,
			float32(r.peakScore),
			float32(r.edgeScore),
			float32(r.orientationScore),
			float32(r.laplacianScaleScore),
		}
	}
	return keys
}

/* Create and destroy */

type CovDet struct {
	p *C.VlCovDet
}

// https://www.vlfeat.org/api/covdet_8c.html#adff732c569785b7dff7f15601bc77a68
func NewCovDet(method CovDetMethod) CovDet {
	p := C.vl_covdet_new(C.VlCovDetMethod(method))
	return CovDet{p: p}
}

// https://www.vlfeat.org/api/covdet_8c.html#a7abfa72a8bb2a05c12376052c29ea17f
func (covdet *CovDet) Delete() {
	C.vl_covdet_delete(covdet.p)
}

// https://www.vlfeat.org/api/covdet_8c.html#a7cfe6bd396893127c69f02f189e7f359
func (covdet *CovDet) Reset() {
	C.vl_covdet_reset(covdet.p)
}

/* Process data */

//
func (covdet *CovDet) PutImage(img []float32, imgWidth, imgHeight uint) VlErrorType {
	imgPtr := toCFloatArrayPtr(img)
	return VlErrorType(C.vl_covdet_put_image(covdet.p, imgPtr, C.uint(imgWidth), C.uint(imgHeight)))
}

// https://www.vlfeat.org/api/covdet_8c.html#abfe553fdb25132cbcf9cbb08839b7f98
func (covdet *CovDet) Detect() {
	C.vl_covdet_detect(covdet.p)
}

// https://www.vlfeat.org/api/covdet_8c.html#af266ce65ae19bff2c5d9ef1ad499b5fd
func (covdet *CovDet) AppendFeature(feature CovDetFeature) VlErrorType {
	cFeature := C.VlCovDetFeature{
		peakScore:           C.float(feature.PeakScore),
		edgeScore:           C.float(feature.EdgeScore),
		orientationScore:    C.float(feature.OrientationScore),
		laplacianScaleScore: C.float(feature.LaplacianScaleScore),
		frame: C.VlFrameOrientedEllipse{
			x:   C.float(feature.Frame.X),
			y:   C.float(feature.Frame.Y),
			a11: C.float(feature.Frame.A11),
			a12: C.float(feature.Frame.A12),
			a21: C.float(feature.Frame.A21),
			a22: C.float(feature.Frame.A22),
		},
	}
	return VlErrorType(C.vl_covdet_append_feature(covdet.p, &cFeature))
}

// https://www.vlfeat.org/api/covdet_8c.html#a0f9823a4cffd55760cf6c16352381f1a
func (covdet *CovDet) ExtractOrientations() {
	C.vl_covdet_extract_orientations(covdet.p)
}

// https://www.vlfeat.org/api/covdet_8c.html#a1858e75418d7c6364a66bbbd8899fdd6
func (covdet *CovDet) ExtractLaplacianScales() {
	C.vl_covdet_extract_laplacian_scales(covdet.p)
}

// https://www.vlfeat.org/api/covdet_8c.html#aee6fc0f35fb23b32e558e8d0131abb96
func (covdet *CovDet) ExtractAffineShape() {
	C.vl_covdet_extract_affine_shape(covdet.p)
}

// https://www.vlfeat.org/api/covdet_8c.html#adec4bf7848db7577c51a1e875a63beaa
func (covdet *CovDet) ExtractOrientationsForFrame(numScales uint, frame CovDetFrameOrientedEllipse) CovDetFeatureOrientation {
	cFrame := C.VlFrameOrientedEllipse{
		x:   C.float(frame.X),
		y:   C.float(frame.Y),
		a11: C.float(frame.A11),
		a12: C.float(frame.A12),
		a21: C.float(frame.A21),
		a22: C.float(frame.A22),
	}
	cNumScales := C.uint(numScales)
	cOrientation := C.vl_covdet_extract_orientations_for_frame(covdet.p, &cNumScales, cFrame)
	orientation := CovDetFeatureOrientation{
		Angle: float64(cOrientation.angle),
		Score: float64(cOrientation.score),
	}
	return orientation
}

// https://www.vlfeat.org/api/covdet_8c.html#a6c8b9af2827291c930c6b1d6cce1794f
func (covdet *CovDet) ExtractLaplacianScalesForFrame(numScales uint, frame CovDetFrameOrientedEllipse) CovDetFeatureLaplacianScale {
	cFrame := C.VlFrameOrientedEllipse{
		x:   C.float(frame.X),
		y:   C.float(frame.Y),
		a11: C.float(frame.A11),
		a12: C.float(frame.A12),
		a21: C.float(frame.A21),
		a22: C.float(frame.A22),
	}
	cNumScales := C.uint(numScales)
	cLaplacianScale := C.vl_covdet_extract_laplacian_scales_for_frame(covdet.p, &cNumScales, cFrame)
	laplacianScale := CovDetFeatureLaplacianScale{
		Scale: float64(cLaplacianScale.scale),
		Score: float64(cLaplacianScale.score),
	}
	return laplacianScale
}

// https://www.vlfeat.org/api/covdet_8c.html#ad3c1402a759e6056b6bd58a27ba4799c
func (covdet *CovDet) ExtractAffineShapeForFrame(frame CovDetFrameOrientedEllipse) (VlErrorType, CovDetFrameOrientedEllipse) {
	cFrame := C.VlFrameOrientedEllipse{
		x:   C.float(frame.X),
		y:   C.float(frame.Y),
		a11: C.float(frame.A11),
		a12: C.float(frame.A12),
		a21: C.float(frame.A21),
		a22: C.float(frame.A22),
	}
	// i think adapted is return value for this function
	var cAdapted C.VlFrameOrientedEllipse
	vlerr := C.vl_covdet_extract_affine_shape_for_frame(covdet.p, &cAdapted, cFrame)
	adapted := CovDetFrameOrientedEllipse{
		X:   float32(cAdapted.x),
		Y:   float32(cAdapted.y),
		A11: float32(cAdapted.a11),
		A12: float32(cAdapted.a12),
		A21: float32(cAdapted.a21),
		A22: float32(cAdapted.a22),
	}
	return VlErrorType(vlerr), adapted
}

// https://www.vlfeat.org/api/covdet_8c.html#a5332ef1f0e09654f5787c19e156404a9
func (covdet *CovDet) ExtractPatchForFrame(patchSize int, resolution uint, extent, sigma float64, frame CovDetFrameOrientedEllipse) (bool, []float32) {
	cPatch := make([]C.float, patchSize)
	cFrame := C.VlFrameOrientedEllipse{
		x:   C.float(frame.X),
		y:   C.float(frame.Y),
		a11: C.float(frame.A11),
		a12: C.float(frame.A12),
		a21: C.float(frame.A21),
		a22: C.float(frame.A22),
	}
	result := int(C.vl_covdet_extract_patch_for_frame(covdet.p, &cPatch[0], C.uint(resolution), C.double(extent), C.double(sigma), cFrame))
	patch := make([]float32, patchSize)
	for i, data := range cPatch {
		patch[i] = float32(data)
	}
	return result != 0, patch
}

// https://www.vlfeat.org/api/covdet_8c.html#a614fd4d42d8d2c938c945c76a544680f
func (covdet *CovDet) DropFeaturesOutside(margin float64) {
	C.vl_covdet_drop_features_outside(covdet.p, C.double(margin))
}

/* Retrieve data and parameters */

// https://www.vlfeat.org/api/covdet_8c.html#a463ed222516fe046e8b170e3d3c06585
func (covdet *CovDet) GetFeaturesNum() uint {
	return uint(C.vl_covdet_get_num_features(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#a8a703daa248e50accd2c09430caeeadd
func (covdet *CovDet) Features() []CovDetFeature {
	length := covdet.GetFeaturesNum()
	features := (*C.VlCovDetFeature)(C.vl_covdet_get_features(covdet.p))
	return getCovDetFeatures(features, int(length))
}

// https://www.vlfeat.org/api/covdet_8c.html#ab06c5d750facf9105d6e3036c1ba40bb
func (covdet *CovDet) GetFirstOctave() int {
	return int(C.vl_covdet_get_first_octave(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#ad3247e72c38a9c713aecc0752468f954
func (covdet *CovDet) GetOctavesNum() uint {
	return uint(C.vl_covdet_get_num_features(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#af57b014275578f5d4513a65b0523ac11
func (covdet *CovDet) GetBaseScale() float64 {
	return float64(C.vl_covdet_get_base_scale(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#a9ccdfbf6b9fbfa01429ed2713c555846
func (covdet *CovDet) GetOctaveResolution() uint {
	return uint(C.vl_covdet_get_octave_resolution(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#a4fd154b4600ca2bf32b87adf77f62825
func (covdet *CovDet) GetPeakThreshold() float64 {
	return float64(C.vl_covdet_get_peak_threshold(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#acd380e24e7525b6c5bc2dfa656d7a9b2
func (covdet *CovDet) GetEdgeThreshold() float64 {
	return float64(C.vl_covdet_get_edge_threshold(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#ab921cc0292adf7c88b17ec7426d2bdd8
func (covdet *CovDet) GetMaxNumOrientations() uint {
	return uint(C.vl_covdet_get_max_num_orientations(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#a55d10412a55bc52e6468f347daa84c2d
func (covdet *CovDet) GetTransposed() bool {
	return C.vl_covdet_get_transposed(covdet.p) != 0
}

// skiped because the return struct in other algorithm
func (covdet *CovDet) GetGss() {
}

func (covdet *CovDet) GetCss() {
}

// https://www.vlfeat.org/api/covdet_8c.html#a997689d5bf1078028223ebad6f1f0626
// because libvl not have vl_covdet_get_laplacian_peak_threshold function
// and VlCovDet strtuct not have this field
func (covdet *CovDet) GetLaplacianPeakThreshold() {
}

// https://www.vlfeat.org/api/covdet_8c.html#ad91ad31d5ca752e14c351e9cffbde68c
func (covdet *CovDet) GetAaAccurateSmoothing() bool {
	return C.vl_covdet_get_aa_accurate_smoothing(covdet.p) != 0
}

// https://www.vlfeat.org/api/covdet_8c.html#a2fce4c4f82f50598abe75f756b631366
func (covdet *CovDet) GetLaplacianScalesStatistics(numScales int) []uint {
	cNumScales := C.uint(numScales)
	cHistogram := C.vl_covdet_get_laplacian_scales_statistics(covdet.p, &cNumScales)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cHistogram)),
		Len:  numScales,
		Cap:  numScales,
	}
	cHistogramSlice := *(*[]C.uint)(unsafe.Pointer(&hdr))
	histogram := make([]uint, numScales)
	for i, data := range cHistogramSlice {
		histogram[i] = uint(data)
	}
	return histogram
}

// https://www.vlfeat.org/api/covdet_8c.html#a68dbea737c1c5ced8590b0b6ae66a3ad
func (covdet *CovDet) GetNonExtremaSuppressionThreshold() float64 {
	return float64(C.vl_covdet_get_non_extrema_suppression_threshold(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#a6e2c647eb87bea3d0935734acd2ba8f5
func (covdet *CovDet) GetNumNonExtremaSuppressed() uint {
	return uint(C.vl_covdet_get_num_non_extrema_suppressed(covdet.p))
}

// https://www.vlfeat.org/api/covdet_8c.html#a1dff0f21a42e510718b35c5020204ca0
func (covdet *CovDet) GetAllowPaddedWarping() bool {
	return C.vl_covdet_get_allow_padded_warping(covdet.p) != 0
}

/*  Set parameters */

// https://www.vlfeat.org/api/covdet_8c.html#af51d1a729e1611201ae7208a1144b8d8
func (covdet *CovDet) SetFirstOctave(o int) {
	C.vl_covdet_set_first_octave(covdet.p, C.int(o))
}

// https://www.vlfeat.org/api/covdet_8c.html#a234342d7c689f6d16e30e5b8433f3170
func (covdet *CovDet) SetNumOctaves(o uint) {
	C.vl_covdet_set_num_octaves(covdet.p, C.uint(o))
}

// https://www.vlfeat.org/api/covdet_8c.html#a83fcd4f56f9ced56232b58d13f869c02
func (covdet *CovDet) SetBaseScale(s float64) {
	C.vl_covdet_set_base_scale(covdet.p, C.double(s))
}

// https://www.vlfeat.org/api/covdet_8c.html#a3ab247f8ce822636bb32ccd48c753a1d
func (covdet *CovDet) SetOctaveResolution(r uint) {
	C.vl_covdet_set_octave_resolution(covdet.p, C.uint(r))
}

// https://www.vlfeat.org/api/covdet_8c.html#afcbfb9f6cade3f20bf429cd520b9cee8
func (covdet *CovDet) SetPeakThreshold(peakThreshold float64) {
	C.vl_covdet_set_peak_threshold(covdet.p, C.double(peakThreshold))
}

// https://www.vlfeat.org/api/covdet_8c.html#a16e0f9b96b16727a1bcec2df7bff0017
func (covdet *CovDet) SetEdgeThreshold(edgeThreshold float64) {
	C.vl_covdet_set_edge_threshold(covdet.p, C.double(edgeThreshold))
}

// https://www.vlfeat.org/api/covdet_8c.html#a214c8968e436f281f70532fa85b54043
func (covdet *CovDet) SetLaplacianPeakThreshold(peakThreshold float64) {
	C.vl_covdet_set_laplacian_peak_threshold(covdet.p, C.double(peakThreshold))
}

// https://www.vlfeat.org/api/covdet_8c.html#a1064241f872090d4ba32c8a7de9dc89b
func (covdet *CovDet) SetMaxNumOrientations(m uint) {
	C.vl_covdet_set_max_num_orientations(covdet.p, C.uint(m))
}

// https://www.vlfeat.org/api/covdet_8c.html#aca7fefc61bf17249380bef9410e0b2b2
func (covdet *CovDet) SetTransposed(t bool) {
	cT := 0
	if t {
		cT = 1
	}
	C.vl_covdet_set_transposed(covdet.p, C.int(cT))
}

// https://www.vlfeat.org/api/covdet_8c.html#a8573f2724701397d5684e0dfbc987423
func (covdet *CovDet) SetAaAccurateSmoothing(x bool) {
	cX := 0
	if x {
		cX = 1
	}
	C.vl_covdet_set_aa_accurate_smoothing(covdet.p, C.int(cX))
}

// https://www.vlfeat.org/api/covdet_8c.html#a50bd5cc5eb6584f585d0bf7af8feff2f
func (covdet *CovDet) SetNonExtremaSuppressionThreshold(x float64) {
	C.vl_covdet_set_non_extrema_suppression_threshold(covdet.p, C.double(x))
}

// https://www.vlfeat.org/api/covdet_8c.html#a1d121b6a8ab2f8a6bda93336731b3b7d
func (covdet *CovDet) SetAllowPaddedWarping(x bool) {
	cX := 0
	if x {
		cX = 1
	}
	C.vl_covdet_set_allow_padded_warping(covdet.p, C.int(cX))
}
