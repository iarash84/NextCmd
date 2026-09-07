//go:build !windows

package terminal

func readConsoleKey(bool) (keyEvent, bool, error) {
	return keyEvent{}, false, nil
}
