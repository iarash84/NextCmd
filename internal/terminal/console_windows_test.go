//go:build windows

package terminal

import "testing"

func TestTranslateConsoleEnterPreservesShift(t *testing.T) {
	plain, ok := translateConsoleKey(consoleKeyEvent{VirtualKey: vkReturn})
	if !ok || plain.kind != KeyEnter {
		t.Fatalf("plain Enter = %#v, %v", plain, ok)
	}
	shifted, ok := translateConsoleKey(consoleKeyEvent{VirtualKey: vkReturn, Control: shiftPressed})
	if !ok || shifted.kind != KeyNewline {
		t.Fatalf("Shift+Enter = %#v, %v", shifted, ok)
	}
}

func TestTranslateConsoleUnicodeCharacter(t *testing.T) {
	event, ok := translateConsoleKey(consoleKeyEvent{UnicodeChar: 'ش'})
	if !ok || event.kind != KeyRune || event.text != "ش" {
		t.Fatalf("unicode key = %#v, %v", event, ok)
	}
}
