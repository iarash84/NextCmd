//go:build darwin

package terminal

import (
	"os"
	"syscall"
	"unsafe"
)

func makeRaw() (func(), bool) {
	fd := os.Stdin.Fd()
	var old syscall.Termios
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd, 0x40487413, uintptr(unsafe.Pointer(&old)), 0, 0, 0); errno != 0 {
		return func() {}, false
	}
	raw := old
	raw.Lflag &^= syscall.ICANON | syscall.ECHO | syscall.ISIG
	raw.Iflag &^= syscall.ICRNL | syscall.IXON
	raw.Cc[syscall.VMIN], raw.Cc[syscall.VTIME] = 1, 0
	if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, fd, 0x80487414, uintptr(unsafe.Pointer(&raw)), 0, 0, 0); errno != 0 {
		return func() {}, false
	}
	return func() { syscall.Syscall6(syscall.SYS_IOCTL, fd, 0x80487414, uintptr(unsafe.Pointer(&old)), 0, 0, 0) }, true
}
