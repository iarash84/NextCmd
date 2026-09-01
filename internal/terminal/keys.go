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
	case 3, 4:
		return keyEvent{kind: KeyEOF}, nil
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
		return keyEvent{kind: KeyIgnored}, nil
	default:
		return keyEvent{kind: KeyRune, value: first}, nil
	}
}
