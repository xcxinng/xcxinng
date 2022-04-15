package third_party

import (
	_ "awesomeProject/util"
	"unsafe"
)

// Assume that byteToString is not exported at the beginning, and
// another package (package util for this example) want to call it,
// however, this need should not change the exported quality of byteToString.
//
// To resolve such a dilemma, generic way is to fork the 3rd party codebase
// and add into your project codebase, (if it is in the same project that
// can't be more convenient)
// and the last but not least, use go:linkname to replace the bodiless function
// in other package.
//
//go:linkname byteToString awesomeGolang/util/
func byteToString(d []byte) string {
	return *(*string)(unsafe.Pointer(&d))
}
