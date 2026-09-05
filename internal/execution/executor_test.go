package execution

import (
	"bytes"
	"context"
	"fmt"
	"nextcmd/sdk"
	"os"
	"strings"
	"testing"
	"time"
)

func TestExecutorCapturesOutputAndExitCode(t *testing.T) {
	if os.Getenv("NEXTCMD_HELPER") == "1" {
		fmt.Print("out")
		fmt.Fprint(os.Stderr, "err")
		os.Exit(7)
	}
	var streamedStdout, streamedStderr bytes.Buffer
	result := (Executor{}).RunStreaming(context.Background(), sdk.Command{Executable: os.Args[0], Args: []string{"-test.run=TestExecutorCapturesOutputAndExitCode"}, Environment: map[string]string{"NEXTCMD_HELPER": "1"}}, &streamedStdout, &streamedStderr)
	if result.ExitCode != 7 || result.Stdout != "out" || result.Stderr != "err" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if streamedStdout.String() != "out" || streamedStderr.String() != "err" {
		t.Fatalf("unexpected streamed output: stdout=%q stderr=%q", streamedStdout.String(), streamedStderr.String())
	}
}

type cancelWriter struct{ cancel context.CancelFunc }

func (w cancelWriter) Write(data []byte) (int, error) {
	w.cancel()
	return len(data), nil
}

func TestExecutorCanCancelWhileStreaming(t *testing.T) {
	if os.Getenv("NEXTCMD_STREAM_HELPER") == "1" {
		fmt.Print("ready")
		time.Sleep(30 * time.Second)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := (Executor{}).RunStreaming(ctx, sdk.Command{
		Executable:  os.Args[0],
		Args:        []string{"-test.run=TestExecutorCanCancelWhileStreaming"},
		Environment: map[string]string{"NEXTCMD_STREAM_HELPER": "1"},
	}, cancelWriter{cancel: cancel}, nil)
	if !result.Canceled || result.Err == nil {
		t.Fatalf("expected canceled execution: %#v", result)
	}
	if !strings.Contains(result.Stdout, "ready") {
		t.Fatalf("streamed output was not captured: %q", result.Stdout)
	}
}

func TestExecutorRunsPlatformShell(t *testing.T) {
	result := (Executor{}).Run(context.Background(), sdk.Command{ShellCommand: "echo nextcmd-shell"})
	if result.Err != nil || result.ExitCode != 0 || strings.TrimSpace(result.Stdout) != "nextcmd-shell" {
		t.Fatalf("unexpected shell result: %#v", result)
	}
}
