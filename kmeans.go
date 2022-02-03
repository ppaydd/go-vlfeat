package vlfeat

/*
#include <stdlib.h>
#include <kmeans.h>
*/
import "C"
import (
	"errors"
	"reflect"
	"unsafe"
)

type VlKMeansAlgorithm int

const (
	VlKMeansLloyd VlKMeansAlgorithm = 0
	VlKMeansElkan VlKMeansAlgorithm = 1
	VlKMeansANN   VlKMeansAlgorithm = 2
)

type VlKMeansInitialization int

const (
	VlKMeansRandomSelection VlKMeansInitialization = 0
	VlKMeansPlusPlus        VlKMeansInitialization = 1
)

type Kmeans struct {
	p *C.VlKMeans
}

// https://www.vlfeat.org/api/kmeans_8c.html#a868a729d2ea5b9f9fec15a18e0a27a76
func NewKeans(dataType VlType, distance VlVectorComparisonType) (Kmeans, error) {
	if dataType != VlTypeFloat && dataType != VlTypeDouble {
		return Kmeans{}, errors.New("Kmeans just support VlTypeFloat and VlTypeDouble")
	}
	p := C.vl_kmeans_new(C.vl_type(dataType), C.VlVectorComparisonType(distance))
	return Kmeans{p: p}, nil
}

// https://www.vlfeat.org/api/kmeans_8c.html#ae251eb379788d26613057f0014bb15bd
func (kmeans *Kmeans) Copy() Kmeans {
	p := C.vl_kmeans_new_copy(kmeans.p)
	return Kmeans{p: p}
}

// https://www.vlfeat.org/api/kmeans_8c.html#a55a50b06dfd493861651200c61458609
func (kmeans *Kmeans) Delete() {
	C.vl_kmeans_delete(kmeans.p)
}

/* Basic data processing*/

// https://www.vlfeat.org/api/kmeans_8c.html#a77b5f58050110584e188534ad15c1dd0
func (kmeans *Kmeans) Reset() {
	C.vl_kmeans_reset(kmeans.p)
}

// https://www.vlfeat.org/api/kmeans_8c.html#a3f35fc9b75799b10a6e32ec81a9cd54d
func (kmeans *Kmeans) Cluster(data interface{}, dimension, numData, numCenters uint) (float64, error) {
	vltype := kmeans.GetDataType()
	dataPtr, _, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return 0, err
	}
	return float64(C.vl_kmeans_cluster(kmeans.p, dataPtr, C.uint(dimension), C.uint(numData), C.uint(numCenters))), nil
}

// https://www.vlfeat.org/api/kmeans_8c.html#a3649fe42a94e9b4945511b5665f82355
// because distances is float or double,so return double
func (kmeans *Kmeans) Quantize(data interface{}, numData uint) ([]uint, []float64, error) {
	distances := make([]float64, numData)
	cDistances := make([]C.double, numData)
	distancesPtr := unsafe.Pointer(&cDistances)

	assignments := make([]uint, numData)
	cAssignments := make([]C.uint, numData)

	vltype := kmeans.GetDataType()
	dataPtr, _, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return assignments, distances, err
	}
	C.vl_kmeans_quantize(kmeans.p, &cAssignments[0], distancesPtr, dataPtr, C.uint(numData))
	for i := 0; i < int(numData); i++ {
		distances[i] = float64(cDistances[i])
		assignments[i] = uint(cAssignments[i])
	}
	return assignments, distances, nil
}

/*
 this function of kemans header file is error
// https://www.vlfeat.org/api/kmeans_8c.html#aa1b270ad6b6e303994629b74b4862f8e
func (kmeans *Kmeans) QuantizeAnn(data interface{}, numData uint, update bool) ([]uint, []float64, error) {
	distances := make([]float64, numData)
	cDistances := make([]C.double, numData)
	distancesPtr := unsafe.Pointer(&cDistances)

	assignments := make([]uint, numData)
	cAssignments := make([]C.uint, numData)

	vltype := kmeans.GetDataType()
	dataPtr, _, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return assignments, distances, err
	}
	cUpdate := 0
	if update {
		cUpdate = 1
	}
	C.vl_kmeans_quantize_ann(kmeans.p, &cAssignments[0], distancesPtr, dataPtr, C.uint(numData), C.int(cUpdate))
	for i := 0; i < int(numData); i++ {
		distances[i] = float64(cDistances[i])
		assignments[i] = uint(cAssignments[i])
	}
	return assignments, distances, nil
}
*/
/* Advanced data processing */

// https://www.vlfeat.org/api/kmeans_8c.html#ac86bd2fa181f6e23e22a6ad92f25288c
func (kmeans *Kmeans) SetCenters(centers interface{}, dimension, numCenters uint) error {
	vltype := kmeans.GetDataType()
	centersPtr, _, err := ToCVlTypeArrayPtr(centers, vltype)
	if err != nil {
		return err
	}
	C.vl_kmeans_set_centers(kmeans.p, centersPtr, C.uint(dimension), C.uint(numCenters))
	return nil
}

// https://www.vlfeat.org/api/kmeans_8c.html#ae32387a856746fe4c39ae10fd533c8d3
func (kmeans *Kmeans) InitCentersWithRandData(data interface{}, dimension, numData, numCenters uint) error {
	vltype := kmeans.GetDataType()
	dataPtr, _, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return err
	}
	C.vl_kmeans_init_centers_with_rand_data(kmeans.p, dataPtr, C.uint(dimension), C.uint(numData), C.uint(numCenters))
	return nil
}

// https://www.vlfeat.org/api/kmeans_8c.html#a5867c89e2916d933ecbc383c4da348c9
func (kmeans *Kmeans) InitCentersPlusPlus(data interface{}, dimension, numData, numCenters uint) error {
	vltype := kmeans.GetDataType()
	dataPtr, _, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return err
	}
	C.vl_kmeans_init_centers_plus_plus(kmeans.p, dataPtr, C.uint(dimension), C.uint(numData), C.uint(numCenters))
	return nil
}

// https://www.vlfeat.org/api/kmeans_8c.html#a9fd1885e6b4742a93b4672f56eb9f2ce
func (kmeans *Kmeans) RefineCenters(data interface{}, numData uint) (float64, error) {
	vltype := kmeans.GetDataType()
	dataPtr, _, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return 0, err
	}
	result := C.vl_kmeans_refine_centers(kmeans.p, dataPtr, C.uint(numData))
	return float64(result), nil
}

/* Retrieve data and parameters */
// https://www.vlfeat.org/api/kmeans_8h.html#abc797cd0e7228d096313fd97b412a21c
func (kmeans *Kmeans) GetDataType() VlType {
	return VlType(C.vl_kmeans_get_data_type(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a61bd032e28960a20cf86e8c77577b67e
func (kmeans *Kmeans) GetDistance() VlVectorComparisonType {
	return VlVectorComparisonType(C.vl_kmeans_get_distance(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#ae6d9f44be8b5d6fc71f11554de742a63
func (kmeans *Kmeans) GetAlgorithm() VlKMeansAlgorithm {
	return VlKMeansAlgorithm(C.vl_kmeans_get_algorithm(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#af920d0fc7e802901fd071eec3def47d2
func (kmeans *Kmeans) GetInitialization() VlKMeansInitialization {
	return VlKMeansInitialization(C.vl_kmeans_get_initialization(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#ae84b1f054eacec5dd7cf2d54668ba181
func (kmeans *Kmeans) GetNumRepetitions() uint {
	return uint(C.vl_kmeans_get_num_repetitions(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a16d8a8bb5d3d7484c4645561dc42ee53
func (kmeans *Kmeans) GetDimension() uint {
	return uint(C.vl_kmeans_get_dimension(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#acba49a529f1393cd04c02a6ef8e2cacd
func (kmeans *Kmeans) GetNumCenters() uint {
	return uint(C.vl_kmeans_get_num_centers(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a7747422051ced08941d2306951445d79
func (kmeans *Kmeans) GetVerbosity() int {
	return int(C.vl_kmeans_get_verbosity(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a6312397c35e56ccd50ac1fa8dbc6bcc2
func (kmeans *Kmeans) GetMaxNumIterations() uint {
	return uint(C.vl_kmeans_get_max_num_iterations(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a49bb2db622503ca24038be22e76aae24
func (kmeans *Kmeans) GetMinEnergyVariation() float64 {
	return float64(C.vl_kmeans_get_min_energy_variation(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a6312397c35e56ccd50ac1fa8dbc6bcc2
func (kmeans *Kmeans) GetMaxNumComparisons() uint {
	return uint(C.vl_kmeans_get_max_num_comparisons(kmeans.p))
}

func (kmeans *Kmeans) GetNumTrees() uint {
	return uint(C.vl_kmeans_get_num_trees(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a02365a4146fabad03f76bb4cdad7cb77
func (kmeans *Kmeans) GetEnergy() float64 {
	return float64(C.vl_kmeans_get_energy(kmeans.p))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a05bd0ac3529beeb460c3dcef4f1a594f
func (kmeans *Kmeans) GetCenters() []float64 {
	dimension := kmeans.GetDimension()
	numCenters := kmeans.GetNumCenters()
	length := dimension * numCenters
	cCenterPtr := C.vl_kmeans_get_centers(kmeans.p)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cCenterPtr)),
		Len:  int(length),
		Cap:  int(length),
	}
	cCenterSlice := *(*[]C.double)(unsafe.Pointer(&hdr))
	centers := make([]float64, length)
	for i, center := range cCenterSlice {
		centers[i] = float64(center)
	}
	return centers
}

/* Set parameters */

// https://www.vlfeat.org/api/kmeans_8h.html#ae2ab27c25bb4730d854219af20f37804
func (kmeans *Kmeans) SetAlgorithm(algorithm VlKMeansAlgorithm) {
	C.vl_kmeans_set_algorithm(kmeans.p, C.VlKMeansAlgorithm(algorithm))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a7c61513feb74bb5393326ed3cced2650
func (kmeans *Kmeans) SetInitialization(initialization VlKMeansInitialization) {
	C.vl_kmeans_set_initialization(kmeans.p, C.VlKMeansInitialization(initialization))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a7a1095cdba2192ee48e15d19144315fb
func (kmeans *Kmeans) SetNumRepetitions(numRepetitions uint) {
	C.vl_kmeans_set_num_repetitions(kmeans.p, C.uint(numRepetitions))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a34f80e7e3f4c7213366b88169cc2f70f
func (kmeans *Kmeans) SetMaxNumIterations(maxNumIterations uint) {
	C.vl_kmeans_set_max_num_iterations(kmeans.p, C.uint(maxNumIterations))
}

// https://www.vlfeat.org/api/kmeans_8h.html#aa807a9a807f80ad1dc5359a81e06566b
func (kmeans *Kmeans) SetMinEnergyVariation(minEnergyVariation float64) {
	C.vl_kmeans_set_min_energy_variation(kmeans.p, C.double(minEnergyVariation))
}

// https://www.vlfeat.org/api/kmeans_8h.html#af2411b97e440ef5419657c41273d38cb
func (kmeans *Kmeans) SetVerbosity(verbosity int) {
	C.vl_kmeans_set_verbosity(kmeans.p, C.int(verbosity))
}

// https://www.vlfeat.org/api/kmeans_8h.html#ac19079ea46c3b5d719dc66f81b38bb5a
func (kmeans *Kmeans) SetMaxNumComparisons(maxNumComparisons uint) {
	C.vl_kmeans_set_max_num_comparisons(kmeans.p, C.uint(maxNumComparisons))
}

// https://www.vlfeat.org/api/kmeans_8h.html#a1a607e91823a83cbcb9c5eaecab99843
func (kmeans *Kmeans) SetNumTrees(numTrees uint) {
	C.vl_kmeans_set_num_trees(kmeans.p, C.uint(numTrees))
}
