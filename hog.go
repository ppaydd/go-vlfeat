package vlfeat

/*
#include <stdlib.h>
#include <hog.h>
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type VlHogVariant int

const (
	VlHogVariantDalalTriggs VlHogVariant = 0
	VlHogVariantUoctti      VlHogVariant = 1
)

type Hog struct {
	p *C.VlHog
}

// https://www.vlfeat.org/api/hog_8h.html#adb99ad366dbd4ea539a76f48df1dff9c
func NewHog(variant VlHogVariant, numOrientations uint, transposed bool) Hog {
	cTransposed := 0
	if transposed {
		cTransposed = 1
	}
	p := C.vl_hog_new(C.VlHogVariant(variant), C.uint(numOrientations), C.int(cTransposed))
	return Hog{p: p}
}

// https://www.vlfeat.org/api/hog_8h.html#a31692138ce8b6c925bf9cf4761f9dd71
func (hog *Hog) Delete() {
	C.vl_hog_delete(hog.p)
}

// https://www.vlfeat.org/api/hog_8h.html#a86e1faec74ae8163db8dd1e0d292c305
func (hog *Hog) PutImage(img []float32, width, height, numChannels, cellSize uint) {
	imgPtr := toCFloatArrayPtr(img)
	C.vl_hog_put_image(hog.p, imgPtr, C.uint(width), C.uint(height), C.uint(numChannels), C.uint(cellSize))
}

// https://www.vlfeat.org/api/hog_8h.html#a2da0444b21261c3db0309ca25ad5895b
func (hog *Hog) PutPolarField(modulus, angle []float32, directed bool, width, height, cellSize uint) {
	modulusPtr := toCFloatArrayPtr(modulus)
	anglePtr := toCFloatArrayPtr(angle)
	cDirected := 0
	if directed {
		cDirected = 1
	}
	C.vl_hog_put_polar_field(hog.p, modulusPtr, anglePtr, C.int(cDirected), C.uint(width), C.uint(height), C.uint(cellSize))
}

// https://www.vlfeat.org/api/hog_8h.html#ab66448d416e661344327f969aebb9a42
func (hog *Hog) Extract() []float32 {
	height := hog.GetHeight()
	width := hog.GetWidth()
	dim := hog.GetDimension()
	featuresLength := int(height * width * dim)
	cFeatures := make([]C.float, featuresLength)
	C.vl_hog_extract(hog.p, &cFeatures[0])
	features := make([]float32, featuresLength)
	for i, feature := range cFeatures {
		features[i] = float32(feature)
	}
	return features
}

// https://www.vlfeat.org/api/hog_8h.html#ad7037e000f578c2abeeb4e8a3d6bbaf6
func (hog *Hog) GetHeight() uint {
	return uint(C.vl_hog_get_height(hog.p))
}

// https://www.vlfeat.org/api/hog_8h.html#a8e899624c92435e77f3a1916b7ce5ed4
func (hog *Hog) GetWidth() uint {
	return uint(C.vl_hog_get_width(hog.p))
}

// https://www.vlfeat.org/api/hog_8h.html#acce19086c37f34edc0078933f224ebcc
func (hog *Hog) GetDimension() uint {
	return uint(C.vl_hog_get_dimension(hog.p))
}

// https://www.vlfeat.org/api/hog_8h.html#a807b3e1a41f403eab4b2c22aba5cbbc2
func (hog *Hog) Render(descriptor []float32, width, height uint) []float32 {
	descriptorPtr := toCFloatArrayPtr(descriptor)
	glyphSize := hog.GetGlyphSize()
	imgWidth := width * glyphSize
	imgHeight := height * glyphSize
	imageLength := int(imgHeight * imgWidth)
	cImage := make([]C.float, imageLength)
	C.vl_hog_render(hog.p, &cImage[0], descriptorPtr, C.uint(width), C.uint(height))
	image := make([]float32, imageLength)
	for i, pixel := range cImage {
		image[i] = float32(pixel)
	}
	return image
}

// https://www.vlfeat.org/api/hog_8h.html#a61ae53dfac6a9ed20a8865e531c23d5a
func (hog *Hog) GetGlyphSize() uint {
	return uint(C.vl_hog_get_glyph_size(hog.p))
}

// https://www.vlfeat.org/api/hog_8h.html#ac700d5506dd453b9a41d3a17b9e3cca7
func (hog *Hog) GetPermutation() []int {
	dims := C.vl_hog_get_dimension(hog.p)
	cPermutation := C.vl_hog_get_permutation(hog.p)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cPermutation)),
		Len:  int(dims),
		Cap:  int(dims),
	}
	cPermutationSlice := *(*[]C.int)(unsafe.Pointer(&hdr))
	permutation := make([]int, dims)
	for i, data := range cPermutationSlice {
		permutation[i] = int(data)
	}
	return permutation
}

// https://www.vlfeat.org/api/hog_8h.html#a901062389718b5584a194d15bbdaaad2
func (hog *Hog) GetUseBilinearOrientationAssignments() bool {
	result := C.vl_hog_get_use_bilinear_orientation_assignments(hog.p)
	return result != 0
}

// https://www.vlfeat.org/api/hog_8h.html#a1d3c6e1ebef79141b2ddbf8208f3e25e
func (hog *Hog) SetUseBilinearOrientationAssignments(x bool) {
	cX := 0
	if x {
		cX = 1
	}
	C.vl_hog_set_use_bilinear_orientation_assignments(hog.p, C.int(cX))
}
