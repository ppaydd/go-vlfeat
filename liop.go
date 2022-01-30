package vlfeat

/*
#include <stdlib.h>
#include <liop.h>
*/
import "C"

type LiopDesc struct {
	p *C.VlLiopDesc
}

// https://www.vlfeat.org/api/liop_8c.html#a58f0187de91697299f22036831453a9e
func NewLiopDesc(numNeighbours, numSpatialBins int, radius float32, sideLength uint) LiopDesc {
	p := C.vl_liopdesc_new(C.int(numNeighbours), C.int(numSpatialBins), C.float(radius), C.uint(sideLength))
	return LiopDesc{p: p}
}

// https://www.vlfeat.org/api/liop_8c.html#a62182ec2c1ee31d0eba5b7203f64f8bc
func newLiopDescBasic(sideLength uint) LiopDesc {
	p := C.vl_liopdesc_new_basic(C.uint(sideLength))
	return LiopDesc{p: p}
}

// https://www.vlfeat.org/api/liop_8c.html#a2e765f1f59a64454f05999de06e48d28
func (ld *LiopDesc) Delete() {
	C.vl_liopdesc_delete(ld.p)
}

// https://www.vlfeat.org/api/liop_8c.html#aea42b7981535254d18c8d55a4c32f606
func (ld *LiopDesc) GetDimension() uint {
	return uint(C.vl_liopdesc_get_dimension(ld.p))
}

// https://www.vlfeat.org/api/liop_8c.html#ac90f67676f91f17d9f3b51df4dab4730
func (ld *LiopDesc) GetNumNeighbours() uint {
	return uint(C.vl_liopdesc_get_num_neighbours(ld.p))
}

// https://www.vlfeat.org/api/liop_8c.html#ad861a14df3f0236ca908726c9fbe269a
func (ld *LiopDesc) GetIntensityThreshold() float32 {
	return float32(C.vl_liopdesc_get_intensity_threshold(ld.p))
}

// https://www.vlfeat.org/api/liop_8c.html#afbb31827623d0b7f7957ba76e6a705cb
func (ld *LiopDesc) GetNumSpatialBins() uint {
	return uint(C.vl_liopdesc_get_num_spatial_bins(ld.p))
}

// https://www.vlfeat.org/api/liop_8c.html#a409893390d672968c871861a37c3dc0c
func (ld *LiopDesc) GetNeighbourhoodRadius() float64 {
	return float64(C.vl_liopdesc_get_neighbourhood_radius(ld.p))
}

// https://www.vlfeat.org/api/liop_8c.html#a1d385f07442b954b658ae419abdd31eb
func (ld *LiopDesc) SetIntensityThreshold(x float32) {
	C.vl_liopdesc_set_intensity_threshold(ld.p, C.float(x))
}

// https://www.vlfeat.org/api/liop_8c.html#a33e638228068ce2674b23eee7af5ce43
func (ld *LiopDesc) Process(patch []float32) []float32 {
	descLength := ld.GetDimension()
	patchPtr := toCFloatArrayPtr(patch)
	cDesc := make([]C.float, descLength)
	C.vl_liopdesc_process(ld.p, &cDesc[0], patchPtr)
	desc := make([]float32, descLength)
	for i, data := range cDesc {
		desc[i] = float32(data)
	}
	return desc
}
