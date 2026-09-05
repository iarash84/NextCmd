package sdk

import (
	"context"
	"io"
	"testing"
)

type testRunner struct{}

func (testRunner) Run(_ context.Context, command Command) ExecutionResult {
	return ExecutionResult{Command: command, Stdout: "controlled", ExitCode: 0}
}

type testStreamingRunner struct{ testRunner }

func (testStreamingRunner) RunStreaming(_ context.Context, command Command, stdout, _ io.Writer) ExecutionResult {
	_, _ = io.WriteString(stdout, "controlled")
	return ExecutionResult{Command: command, Stdout: "controlled", ExitCode: 0}
}

func TestRunnerContractsAcceptTestDoubles(t *testing.T) {
	var runner Runner = testRunner{}
	var streaming StreamingRunner = testStreamingRunner{}
	command := Command{Executable: "tool"}
	if got := runner.Run(context.Background(), command); got.Stdout != "controlled" {
		t.Fatalf("runner result=%#v", got)
	}
	if got := streaming.Run(context.Background(), command); got.Stdout != "controlled" {
		t.Fatalf("streaming runner result=%#v", got)
	}
}
