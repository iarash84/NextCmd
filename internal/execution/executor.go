package execution

import (
	"bytes"
	"context"
	"errors"
	"nextcmd/sdk"
	"os"
	"os/exec"
	"time"
)

type Executor struct{}

func (Executor) Run(ctx context.Context, command sdk.Command) sdk.ExecutionResult {
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
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	result.Stdout, result.Stderr, result.Duration, result.Err = stdout.String(), stderr.String(), time.Since(started), err
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
