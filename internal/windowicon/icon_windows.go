//go:build windows

package windowicon

import (
	"syscall"
	"unsafe"
)

const (
	iconResourceID = 1
	imageIcon      = 1
	lrDefaultSize  = 0x0040
	lrShared       = 0x8000
	wmSetIcon      = 0x0080
	iconSmall      = 0
	iconBig        = 1
)

var (
	kernel32                           = syscall.NewLazyDLL("kernel32.dll")
	user32                             = syscall.NewLazyDLL("user32.dll")
	shell32                            = syscall.NewLazyDLL("shell32.dll")
	getConsoleWindow                   = kernel32.NewProc("GetConsoleWindow")
	getModuleHandle                    = kernel32.NewProc("GetModuleHandleW")
	setConsoleIcon                     = kernel32.NewProc("SetConsoleIcon")
	loadImage                          = user32.NewProc("LoadImageW")
	sendMessage                        = user32.NewProc("SendMessageW")
	setCurrentProcessExplicitAppUserID = shell32.NewProc("SetCurrentProcessExplicitAppUserModelID")
)

// Apply gives the console window and its taskbar group the icon embedded in
// the executable. Failures are intentionally non-fatal on older Windows hosts
// and pseudoconsole terminals that own their window chrome themselves.
func Apply() {
	appID, err := syscall.UTF16PtrFromString("NextCmd.TerminalAssistant")
	if err == nil {
		_, _, _ = setCurrentProcessExplicitAppUserID.Call(uintptr(unsafe.Pointer(appID)))
	}

	module, _, _ := getModuleHandle.Call(0)
	if module == 0 {
		return
	}
	icon, _, _ := loadImage.Call(module, iconResourceID, imageIcon, 0, 0, lrDefaultSize|lrShared)
	if icon == 0 {
		return
	}
	_, _, _ = setConsoleIcon.Call(icon)
	window, _, _ := getConsoleWindow.Call()
	if window != 0 {
		_, _, _ = sendMessage.Call(window, wmSetIcon, iconBig, icon)
		_, _, _ = sendMessage.Call(window, wmSetIcon, iconSmall, icon)
	}
}
