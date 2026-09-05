//go:build windows

package execution

import "os"

func platformShell(command string) (string, []string) {
	shell := os.Getenv("COMSPEC")
	if shell == "" {
		shell = "cmd.exe"
	}
	return shell, []string{"/D", "/S", "/C", command}
}
