//go:build windows

package terminal

import (
	"os"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var getConsoleMode = kernel32.NewProc("GetConsoleMode")
var setConsoleMode = kernel32.NewProc("SetConsoleMode")

func makeRaw() (func(), bool) {
	handle := os.Stdin.Fd()
	var old uint32
	ok, _, _ := getConsoleMode.Call(handle, uintptr(unsafe.Pointer(&old)))
	if ok == 0 {
		return func() {}, false
	}
	raw := old &^ uint32(0x0002|0x0004)
	raw |= 0x0200
	ok, _, _ = setConsoleMode.Call(handle, uintptr(raw))
	if ok == 0 {
		return func() {}, false
	}
	return func() { setConsoleMode.Call(handle, uintptr(old)) }, true
}
