package execution

import (
	"context"
	"fmt"
	"nextcmd/sdk"
	"os"
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
