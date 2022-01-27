package vlfeat

/*
#include <stdlib.h>
#include <scalespace.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type ScaleSpaceGeometry struct {
	Width                  uint    `json:"width"`
	Height                 uint    `json:"height"`
	FirstOctave            int     `json:"firstOctave"`
	LastOctave             int     `json:"lastOctave"`
	OctaveResolution       uint    `json:"octaveResolution"`
	OctaveFirstSubdivision int     `json:"octaveFirstSubdivision"`
	OctaveLastSubdivision  int     `json:"octaveLastSubdivision"`
	BaseScale              float64 `json:"baseScale"`
	NominalScale           float64 `json:"nominalScale"`
}

type ScaleSpaceOctaveGeometry struct {
	Width  uint    `json:"width"`
	Height uint    `json:"height"`
	Step   float64 `json:"step"`
}

type ScaleSpace struct {
	p *C.VlScaleSpace
}

/* ScaleSpace  Create and destroy */

// https://www.vlfeat.org/api/scalespace_8c.html#ab9a838a41d9fe1e04f1477061055b6d9
func ScaleSpaceGetDefaultGeometry(width, height uint) ScaleSpaceGeometry {
	cGeom := C.vl_scalespace_get_default_geometry(C.uint(width), C.uint(height))
	return ScaleSpaceGeometry{
		Width:                  uint(cGeom.width),
		Height:                 uint(cGeom.height),
		FirstOctave:            int(cGeom.firstOctave),
		LastOctave:             int(cGeom.lastOctave),
		OctaveResolution:       uint(cGeom.octaveResolution),
		OctaveFirstSubdivision: int(cGeom.octaveFirstSubdivision),
		OctaveLastSubdivision:  int(cGeom.octaveLastSubdivision),
		BaseScale:              float64(cGeom.baseScale),
		NominalScale:           float64(cGeom.nominalScale),
	}
}

// https://www.vlfeat.org/api/scalespace_8c.html#ad515a7ec3d620d6b55f9869c5318b546
func vl_scalespacegeometry_is_equal(a, b ScaleSpaceGeometry) bool {
	cGeomA := C.VlScaleSpaceGeometry{
		width:                  C.uint(a.Width),
		height:                 C.uint(a.Height),
		firstOctave:            C.int(a.FirstOctave),
		lastOctave:             C.int(a.LastOctave),
		octaveResolution:       C.uint(a.OctaveResolution),
		octaveFirstSubdivision: C.int(a.OctaveFirstSubdivision),
		octaveLastSubdivision:  C.int(a.OctaveLastSubdivision),
		baseScale:              C.double(a.BaseScale),
		nominalScale:           C.double(a.NominalScale),
	}
	cGeomB := C.VlScaleSpaceGeometry{
		width:                  C.uint(b.Width),
		height:                 C.uint(b.Height),
		firstOctave:            C.int(b.FirstOctave),
		lastOctave:             C.int(b.LastOctave),
		octaveResolution:       C.uint(b.OctaveResolution),
		octaveFirstSubdivision: C.int(b.OctaveFirstSubdivision),
		octaveLastSubdivision:  C.int(b.OctaveLastSubdivision),
		baseScale:              C.double(b.BaseScale),
		nominalScale:           C.double(b.NominalScale),
	}
	result := C.vl_scalespacegeometry_is_equal(cGeomA, cGeomB)
	return result != 0
}

// https://www.vlfeat.org/api/scalespace_8c.html#a42931ebe6e7e762b7d0ccd7856d97834
func NewScaleSpace(width, height uint) ScaleSpace {
	p := C.vl_scalespace_new(C.uint(width), C.uint(height))
	return ScaleSpace{p: p}
}

// https://www.vlfeat.org/api/scalespace_8c.html#a6e997ce733e67a186769e2877a5aaa6a
func NewScaleSpaceWithGeometry(geom ScaleSpaceGeometry) ScaleSpace {
	cGeom := C.VlScaleSpaceGeometry{
		width:                  C.uint(geom.Width),
		height:                 C.uint(geom.Height),
		firstOctave:            C.int(geom.FirstOctave),
		lastOctave:             C.int(geom.LastOctave),
		octaveResolution:       C.uint(geom.OctaveResolution),
		octaveFirstSubdivision: C.int(geom.OctaveFirstSubdivision),
		octaveLastSubdivision:  C.int(geom.OctaveLastSubdivision),
		baseScale:              C.double(geom.BaseScale),
		nominalScale:           C.double(geom.NominalScale),
	}
	p := C.vl_scalespace_new_with_geometry(cGeom)
	return ScaleSpace{p: p}
}

// https://www.vlfeat.org/api/scalespace_8c.html#aaea616bc66aee097260035536befd845
func (ss *ScaleSpace) Copy() ScaleSpace {
	p := C.vl_scalespace_new_copy(ss.p)
	return ScaleSpace{p: p}
}

// https://www.vlfeat.org/api/scalespace_8c.html#a72281dc130acc23ebbcfed24aafc1ffd
func (ss *ScaleSpace) ShallowCopy() ScaleSpace {
	p := C.vl_scalespace_new_shallow_copy(ss.p)
	return ScaleSpace{p: p}
}

// https://www.vlfeat.org/api/scalespace_8c.html#af4d66bec63bb5a3670f55a2fa7786830
func (ss *ScaleSpace) Delete() {
	C.vl_scalespace_delete(ss.p)
}

/* Process data */

// https://www.vlfeat.org/api/scalespace_8c.html#a2e2be28c7c1017222fec22a8c0df02f5
func (ss *ScaleSpace) PutImage(img []float32) {
	imgPtr := toCFloatArrayPtr(img)
	C.vl_scalespace_put_image(ss.p, imgPtr)
}

/* Retrieve data and parameters */
// https://www.vlfeat.org/api/scalespace_8c.html#a13e0136527672f35a76ffb60fcea2bba
func (ss *ScaleSpace) GetGeometry() ScaleSpaceGeometry {
	cGeom := C.vl_scalespace_get_geometry(ss.p)
	return ScaleSpaceGeometry{
		Width:                  uint(cGeom.width),
		Height:                 uint(cGeom.height),
		FirstOctave:            int(cGeom.firstOctave),
		LastOctave:             int(cGeom.lastOctave),
		OctaveResolution:       uint(cGeom.octaveResolution),
		OctaveFirstSubdivision: int(cGeom.octaveFirstSubdivision),
		OctaveLastSubdivision:  int(cGeom.octaveLastSubdivision),
		BaseScale:              float64(cGeom.baseScale),
		NominalScale:           float64(cGeom.nominalScale),
	}
}

// https://www.vlfeat.org/api/scalespace_8c.html#a5c63aee8d9c6e9393308eb3306264108
func (ss *ScaleSpace) GetOctaveGeometry(o int) ScaleSpaceOctaveGeometry {
	cGeom := C.vl_scalespace_get_octave_geometry(ss.p, C.int(o))
	return ScaleSpaceOctaveGeometry{
		Width:  uint(cGeom.width),
		Height: uint(cGeom.height),
		Step:   float64(cGeom.step),
	}
}

// https://www.vlfeat.org/api/scalespace_8c.html#a5010e154df5b981f31afff3af704ae30
func (ss *ScaleSpace) GetLevel(o, s int) []float32 {
	cLevel := C.vl_scalespace_get_level(ss.p, C.int(o), C.int(s))
	ogeo := ss.GetOctaveGeometry(o)
	geom := ss.GetGeometry()
	length := int(ogeo.Width * ogeo.Height * uint(geom.OctaveLastSubdivision-geom.OctaveFirstSubdivision+1))
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cLevel)),
		Len:  length,
		Cap:  length,
	}
	cLevelSlice := *(*[]C.float)(unsafe.Pointer(&hdr))
	levels := make([]float32, length)
	for i, des := range cLevelSlice {
		levels[i] = float32(des)
	}
	return levels
}

// https://www.vlfeat.org/api/scalespace_8c.html#a105e697419bf9ddfaef70e498353a968
func (ss *ScaleSpace) GetLevelSigma(o, s int) float64 {
	return float64(C.vl_scalespace_get_level_sigma(ss.p, C.int(o), C.int(s)))
}
