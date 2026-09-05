package terminal

import (
	"bufio"
	"io"
)

type byteScanner interface {
	io.ByteReader
	UnreadByte() error
}

func readKey(reader io.Reader) (keyEvent, error) {
	scanner, ok := reader.(byteScanner)
	if !ok {
		scanner = bufio.NewReader(reader)
	}
	first, err := scanner.ReadByte()
	if err != nil {
		return keyEvent{}, err
	}
	switch first {
	case 1:
		return keyEvent{kind: KeyHome}, nil
	case 3, 4:
		return keyEvent{kind: KeyEOF}, nil
	case 5:
		return keyEvent{kind: KeyEnd}, nil
	case 14:
		return keyEvent{kind: KeyHistoryNext}, nil
	case 16:
		return keyEvent{kind: KeyHistoryPrevious}, nil
	case 18:
		return keyEvent{kind: KeyHistorySearch}, nil
	case 21:
		return keyEvent{kind: KeyClearLine}, nil
	case 13, 10:
		return keyEvent{kind: KeyEnter}, nil
	case 9:
		return keyEvent{kind: KeyTab}, nil
	case 8, 127:
		return keyEvent{kind: KeyBackspace}, nil
	case 0, 224:
		next, err := scanner.ReadByte()
		if err != nil {
			return keyEvent{}, err
		}
		if next == 72 {
			return keyEvent{kind: KeyUp}, nil
		}
		if next == 80 {
			return keyEvent{kind: KeyDown}, nil
		}
		if next == 77 {
			return keyEvent{kind: KeyRight}, nil
		}
		if next == 75 {
			return keyEvent{kind: KeyLeft}, nil
		}
		if next == 71 {
			return keyEvent{kind: KeyHome}, nil
		}
		if next == 79 {
			return keyEvent{kind: KeyEnd}, nil
		}
		return keyEvent{kind: KeyIgnored}, nil
	case 27:
		next, err := scanner.ReadByte()
		if err != nil {
			return keyEvent{kind: KeyEscape}, nil
		}
		if next != '[' {
			_ = scanner.UnreadByte()
			return keyEvent{kind: KeyEscape}, nil
		}
		final, err := scanner.ReadByte()
		if err != nil {
			return keyEvent{kind: KeyEscape}, nil
		}
		if final == 'A' {
			return keyEvent{kind: KeyUp}, nil
		}
		if final == 'B' {
			return keyEvent{kind: KeyDown}, nil
		}
		if final == 'C' {
			return keyEvent{kind: KeyRight}, nil
		}
		if final == 'D' {
			return keyEvent{kind: KeyLeft}, nil
		}
		if final == 'H' {
			return keyEvent{kind: KeyHome}, nil
		}
		if final == 'F' {
			return keyEvent{kind: KeyEnd}, nil
		}
		return keyEvent{kind: KeyIgnored}, nil
	default:
		return keyEvent{kind: KeyRune, value: first}, nil
	}
}
