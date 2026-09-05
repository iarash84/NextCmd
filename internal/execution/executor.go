package execution

import (
	"bytes"
	"context"
	"errors"
	"io"
	"nextcmd/sdk"
	"os"
	"os/exec"
	"time"
)

type Executor struct{}

func (Executor) Run(ctx context.Context, command sdk.Command) sdk.ExecutionResult {
	return (Executor{}).RunStreaming(ctx, command, nil, nil)
}

// RunStreaming writes process output as it arrives while retaining the same
// output in the returned result for history, recovery, and next-action logic.
func (Executor) RunStreaming(ctx context.Context, command sdk.Command, stdoutWriter, stderrWriter io.Writer) sdk.ExecutionResult {
	started := time.Now()
	result := sdk.ExecutionResult{Command: command, ExitCode: -1}
	var cmd *exec.Cmd
	if command.ShellCommand != "" {
		name, args := platformShell(command.ShellCommand)
		cmd = exec.CommandContext(ctx, name, args...)
	} else {
		cmd = exec.CommandContext(ctx, command.Executable, command.Args...)
	}
	cmd.Dir = command.WorkingDirectory
	cmd.Env = os.Environ()
	for key, value := range command.Environment {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = captureWriter(&stdout, stdoutWriter)
	cmd.Stderr = captureWriter(&stderr, stderrWriter)
	err := cmd.Run()
	result.Stdout, result.Stderr, result.Duration, result.Err = stdout.String(), stderr.String(), time.Since(started), err
	result.Canceled = ctx.Err() != nil
	if err == nil {
		result.ExitCode = 0
	} else {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
		}
	}
	return result
}

func captureWriter(capture *bytes.Buffer, stream io.Writer) io.Writer {
	if stream == nil {
		return capture
	}
	return io.MultiWriter(stream, capture)
}
