package vlfeat

/*
#cgo CXXFLAGS: --std=c++11
#cgo !windows CPPFLAGS: -I/usr/local/include -I/usr/local/include/vlfeat
#cgo !windows LDFLAGS: -L/usr/local/lib -lvl
#cgo windows  CPPFLAGS: -IC:/vlfeat-0.9.21/vl
#cgo windows  LDFLAGS: -LC:/vlfeat-0.9.21/bin/win64 -lvl
*/
import "C"
