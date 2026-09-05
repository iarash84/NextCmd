package execution

import (
	"context"
	"fmt"
	"nextcmd/sdk"
	"os"
	"strings"
	"testing"
)

func TestExecutorCapturesOutputAndExitCode(t *testing.T) {
	if os.Getenv("NEXTCMD_HELPER") == "1" {
		fmt.Print("out")
		fmt.Fprint(os.Stderr, "err")
		os.Exit(7)
	}
	result := (Executor{}).Run(context.Background(), sdk.Command{Executable: os.Args[0], Args: []string{"-test.run=TestExecutorCapturesOutputAndExitCode"}, Environment: map[string]string{"NEXTCMD_HELPER": "1"}})
	if result.ExitCode != 7 || result.Stdout != "out" || result.Stderr != "err" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestExecutorRunsPlatformShell(t *testing.T) {
	result := (Executor{}).Run(context.Background(), sdk.Command{ShellCommand: "echo nextcmd-shell"})
	if result.Err != nil || result.ExitCode != 0 || strings.TrimSpace(result.Stdout) != "nextcmd-shell" {
		t.Fatalf("unexpected shell result: %#v", result)
	}
}
