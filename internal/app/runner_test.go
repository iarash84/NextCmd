package app

import (
	"context"
	"io"
	"testing"

	"nextcmd/sdk"
)

type fakeStreamingRunner struct{}

func (fakeStreamingRunner) Run(_ context.Context, command sdk.Command) sdk.ExecutionResult {
	return sdk.ExecutionResult{Command: command, ExitCode: 0}
}

func (fakeStreamingRunner) RunStreaming(_ context.Context, command sdk.Command, _, _ io.Writer) sdk.ExecutionResult {
	return sdk.ExecutionResult{Command: command, ExitCode: 0}
}

func TestNewWithRunnerInjectsRunner(t *testing.T) {
	runner := fakeStreamingRunner{}
	application := NewWithRunner(nil, nil, nil, nil, "work", RuntimeSettings{}, runner)
	if application.executor != runner {
		t.Fatalf("runner was not injected: %#v", application.executor)
	}
}
