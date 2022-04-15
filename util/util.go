// util.go
// 需要导入 unsafe 包

package util

import _ "unsafe"

//go:linkname hello awesomeGolang/third_party.byteToString
func hello(d []byte) string

func CallHello(d []byte) string {
	return hello(d)
}
