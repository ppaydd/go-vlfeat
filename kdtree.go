package vlfeat

/*
#include <stdlib.h>
#include <kdtree.h>
*/
import "C"
import (
	"errors"
)

type VlKDTreeThresholdingMethod int

const (
	KDTreeMedian VlKDTreeThresholdingMethod = 0
	KDTreeMean   VlKDTreeThresholdingMethod = 1
)

type KDForest struct {
	p *C.VlKDForest
}

type KDForestSearcher struct {
	p *C.VlKDForestSearcher
}

type KDForestNeighbor struct {
	Distance float64
	Index    uint
}

/* Creating, copying and disposing */

// https://www.vlfeat.org/api/kdtree_8c.html#a52564e86ef0d9294a9bc9b13c5d44427
func NewKDForest(dataType VlType, dimension, numTress uint, normType VlVectorComparisonType) (KDForest, error) {
	if dataType != VlTypeFloat && dataType != VlTypeDouble {
		return KDForest{}, errors.New("Kmeans just support VlTypeFloat and VlTypeDouble")
	}
	p := C.vl_kdforest_new(C.vl_type(dataType), C.uint(dimension), C.uint(numTress), C.VlVectorComparisonType(normType))
	return KDForest{p: p}, nil
}

// https://www.vlfeat.org/api/kdtree_8c.html#a9d909b0b42489ce438b03e99be9fd5d1
func (kdforest *KDForest) NewSearcher() KDForestSearcher {
	p := C.vl_kdforest_new_searcher(kdforest.p)
	return KDForestSearcher{p: p}
}

// https://www.vlfeat.org/api/kdtree_8c.html#a68bcecfea6e63a41aafbd5ea9ca1a418
func (kdforest *KDForest) Delete() {
	C.vl_kdforest_delete(kdforest.p)
}

// https://www.vlfeat.org/api/kdtree_8c.html#aaf7bb0d93fffba8cc0b1967b6a94293a
func (kdfs *KDForestSearcher) Delete() {
	C.vl_kdforestsearcher_delete(kdfs.p)
}

/* Building and querying */

// https://www.vlfeat.org/api/kdtree_8c.html#ac886f1fd6024a74e9e4a5d7566b2125f
func (kdforest *KDForest) Build(data interface{}) error {
	vltype := kdforest.GetDataType()
	dataPtr, numData, err := ToCVlTypeArrayPtr(data, vltype)
	if err != nil {
		return err
	}
	C.vl_kdforest_build(kdforest.p, C.uint(numData), dataPtr)
	return nil
}

// https://www.vlfeat.org/api/kdtree_8c.html#a2af87b58193ea0314fa971f8678e4e8c
func (kdforest *KDForest) Query(numNeighbors uint, query interface{}) (uint, []KDForestNeighbor, error) {
	vltype := kdforest.GetDataType()
	queryPtr, _, err := ToCVlTypeArrayPtr(query, vltype)
	if err != nil {
		return 0, []KDForestNeighbor{}, err
	}
	cNeighbors := make([]C.VlKDForestNeighbor, numNeighbors)
	neighbors := make([]KDForestNeighbor, numNeighbors)
	result := C.vl_kdforest_query(kdforest.p, &cNeighbors[0], C.uint(numNeighbors), queryPtr)
	for i, neighbor := range cNeighbors {
		neighbors[i] = KDForestNeighbor{
			float64(neighbor.distance),
			uint(neighbor.index),
		}
	}
	return uint(result), neighbors, nil
}

/* Retrieving and setting parameters */

// https://www.vlfeat.org/api/kdtree_8c.html#a55728e3e3e7a3619ed24e8b016dbf2a4
func (kdforest *KDForest) GetDepthOfTree(treeIndex uint) uint {
	return uint(C.vl_kdforest_get_depth_of_tree(kdforest.p, C.uint(treeIndex)))
}

// https://www.vlfeat.org/api/kdtree_8c.html#a766017f8aefa345762c3503f6a2b0b75
func (kdforest *KDForest) GetNumNodesOfTree(treeIndex uint) uint {
	return uint(C.vl_kdforest_get_num_nodes_of_tree(kdforest.p, C.uint(treeIndex)))
}

// https://www.vlfeat.org/api/kdtree_8c.html#a054701571177903a5369bb73da3139ef
func (kdforest *KDForest) GetNumTrees() uint {
	return uint(C.vl_kdforest_get_num_trees(kdforest.p))
}

// https://www.vlfeat.org/api/kdtree_8c.html#a548f72a02684f50ed2fc5d54e256b752
func (kdforest *KDForest) GetDataDimension() uint {
	return uint(C.vl_kdforest_get_data_dimension(kdforest.p))
}

// https://www.vlfeat.org/api/kdtree_8c.html#abba39fcd9a3693e9c8851b0b008c2db0
func (kdforest *KDForest) GetDataType() VlType {
	return VlType(C.vl_kdforest_get_data_type(kdforest.p))
}

// https://www.vlfeat.org/api/kdtree_8c.html#af0d6436d4b42826cf4aaf50b2ec59b53
func (kdforest *KDForest) GetMaxNumComparisons() uint {
	return uint(C.vl_kdforest_get_max_num_comparisons(kdforest.p))
}

// https://www.vlfeat.org/api/kdtree_8c.html#ac153ede5b0d8716f2ae7887ef7060884
func (kdforest *KDForest) GetThresholdingMethod() VlKDTreeThresholdingMethod {
	return VlKDTreeThresholdingMethod(C.vl_kdforest_get_thresholding_method(kdforest.p))
}

// https://www.vlfeat.org/api/kdtree_8c.html#adf66cf1b6be55d82a2ad4d568271f157
func (kdforest *KDForest) SetThresholdingMethod(method VlKDTreeThresholdingMethod) {
	C.vl_kdforest_set_thresholding_method(kdforest.p, C.VlKDTreeThresholdingMethod(method))
}

// https://www.vlfeat.org/api/kdtree_8c.html#a4bf926f9406cec740d564b703236c68d
func (kdforest *KDForest) SetMaxNumComparisons(n uint) {
	C.vl_kdforest_set_max_num_comparisons(kdforest.p, C.uint(n))
}
