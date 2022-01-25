package vlfeat

/*
#cgo CXXFLAGS: --std=c++11
#cgo !windows CPPFLAGS: -I./vl
#cgo !windows LDFLAGS: -L. -lvl
#cgo windows  CPPFLAGS: -I ./vl
#cgo windows  LDFLAGS: -L. -lvl
*/
import "C"
