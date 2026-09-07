//go:build windows

package terminal

import (
	"io"
	"os"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	keyEventType = 0x0001
	shiftPressed = 0x0010
	vkBack       = 0x08
	vkTab        = 0x09
	vkReturn     = 0x0D
	vkEscape     = 0x1B
	vkEnd        = 0x23
	vkHome       = 0x24
	vkLeft       = 0x25
	vkUp         = 0x26
	vkRight      = 0x27
	vkDown       = 0x28
)

type consoleKeyEvent struct {
	KeyDown     int32
	RepeatCount uint16
	VirtualKey  uint16
	ScanCode    uint16
	UnicodeChar uint16
	Control     uint32
}

type consoleInputRecord struct {
	EventType uint16
	_         uint16
	Key       consoleKeyEvent
}

var readConsoleInputW = kernel32.NewProc("ReadConsoleInputW")

// readConsoleKey uses native console records because the byte stream collapses
// Shift+Enter and Enter to the same carriage-return byte on Windows.
func readConsoleKey(enabled bool) (keyEvent, bool, error) {
	if !enabled {
		return keyEvent{}, false, nil
	}
	for {
		var record consoleInputRecord
		var count uint32
		ok, _, callErr := readConsoleInputW.Call(os.Stdin.Fd(), uintptr(unsafe.Pointer(&record)), 1, uintptr(unsafe.Pointer(&count)))
		if ok == 0 {
			if callErr != syscall.Errno(0) {
				return keyEvent{}, true, callErr
			}
			return keyEvent{}, true, io.ErrUnexpectedEOF
		}
		if count == 0 || record.EventType != keyEventType || record.Key.KeyDown == 0 {
			continue
		}
		if event, ok := translateConsoleKey(record.Key); ok {
			return event, true, nil
		}
	}
}

func translateConsoleKey(key consoleKeyEvent) (keyEvent, bool) {
	if key.VirtualKey == vkReturn {
		if key.Control&shiftPressed != 0 {
			return keyEvent{kind: KeyNewline}, true
		}
		return keyEvent{kind: KeyEnter}, true
	}
	switch key.VirtualKey {
	case vkTab:
		return keyEvent{kind: KeyTab}, true
	case vkBack:
		return keyEvent{kind: KeyBackspace}, true
	case vkEscape:
		return keyEvent{kind: KeyEscape}, true
	case vkUp:
		return keyEvent{kind: KeyUp}, true
	case vkDown:
		return keyEvent{kind: KeyDown}, true
	case vkLeft:
		return keyEvent{kind: KeyLeft}, true
	case vkRight:
		return keyEvent{kind: KeyRight}, true
	case vkHome:
		return keyEvent{kind: KeyHome}, true
	case vkEnd:
		return keyEvent{kind: KeyEnd}, true
	}
	if event, ok := consoleControlKey(key.UnicodeChar); ok {
		return event, true
	}
	if key.UnicodeChar != 0 {
		runes := utf16.Decode([]uint16{key.UnicodeChar})
		if len(runes) > 0 {
			return keyEvent{kind: KeyRune, text: string(runes)}, true
		}
	}
	return keyEvent{}, false
}

func consoleControlKey(char uint16) (keyEvent, bool) {
	switch char {
	case 1:
		return keyEvent{kind: KeyHome}, true
	case 3, 4:
		return keyEvent{kind: KeyEOF}, true
	case 5:
		return keyEvent{kind: KeyEnd}, true
	case 14:
		return keyEvent{kind: KeyHistoryNext}, true
	case 16:
		return keyEvent{kind: KeyHistoryPrevious}, true
	case 18:
		return keyEvent{kind: KeyHistorySearch}, true
	case 21:
		return keyEvent{kind: KeyClearLine}, true
	}
	return keyEvent{}, false
}
