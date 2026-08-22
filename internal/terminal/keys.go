package terminal

import "io"

func readKey(reader io.Reader) (keyEvent, error) {
	var b [1]byte
	if _, err := reader.Read(b[:]); err != nil {
		return keyEvent{}, err
	}
	switch b[0] {
	case 3, 4:
		return keyEvent{kind: KeyEOF}, nil
	case 13, 10:
		return keyEvent{kind: KeyEnter}, nil
	case 9:
		return keyEvent{kind: KeyTab}, nil
	case 8, 127:
		return keyEvent{kind: KeyBackspace}, nil
	case 0, 224:
		if _, err := reader.Read(b[:]); err != nil {
			return keyEvent{}, err
		}
		if b[0] == 72 {
			return keyEvent{kind: KeyUp}, nil
		}
		if b[0] == 80 {
			return keyEvent{kind: KeyDown}, nil
		}
		if b[0] == 77 {
			return keyEvent{kind: KeyRight}, nil
		}
		if b[0] == 75 {
			return keyEvent{kind: KeyLeft}, nil
		}
		return keyEvent{kind: KeyIgnored}, nil
	case 27:
		var sequence [2]byte
		if _, err := io.ReadFull(reader, sequence[:]); err != nil {
			return keyEvent{kind: KeyEscape}, nil
		}
		if sequence[0] == '[' && sequence[1] == 'A' {
			return keyEvent{kind: KeyUp}, nil
		}
		if sequence[0] == '[' && sequence[1] == 'B' {
			return keyEvent{kind: KeyDown}, nil
		}
		if sequence[0] == '[' && sequence[1] == 'C' {
			return keyEvent{kind: KeyRight}, nil
		}
		if sequence[0] == '[' && sequence[1] == 'D' {
			return keyEvent{kind: KeyLeft}, nil
		}
		return keyEvent{kind: KeyIgnored}, nil
	default:
		return keyEvent{kind: KeyRune, value: b[0]}, nil
	}
}
