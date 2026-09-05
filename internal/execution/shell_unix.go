//go:build !windows

package execution

func platformShell(command string) (string, []string) {
	return "/bin/sh", []string{"-c", command}
}
