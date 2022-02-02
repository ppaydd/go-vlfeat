package vlfeat

/*
#include <stdlib.h>
#include <slic.h>
*/
import "C"

// https://www.vlfeat.org/api/slic_8c.html#adb6a4c91f40fc32528ba88cffba756ab
func SlicSegment(image []float32, width, height, numChannels, regionSize uint, regularization float32, minRegionSize uint) []uint {
	imgPtr := toCFloatArrayPtr(image)
	segmentationLength := width * height
	cSegmentation := make([]C.uint, segmentationLength)
	C.vl_slic_segment(&cSegmentation[0], imgPtr, C.uint(width), C.uint(height), C.uint(numChannels), C.uint(regionSize), C.float(regularization), C.uint(minRegionSize))
	segmentation := make([]uint, segmentationLength)
	for i, data := range cSegmentation {
		segmentation[i] = uint(data)
	}
	return segmentation
}
