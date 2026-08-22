package git

import (
	"bytes"
	"context"
	"os/exec"

	"nextcmd/sdk"
)

type commandRunner struct{}

func (commandRunner) Run(ctx context.Context, directory string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, "git", args...)
	command.Dir = directory
	var output bytes.Buffer
	command.Stdout, command.Stderr = &output, &output
	err := command.Run()
	return output.String(), err
}

type Plugin struct{ runner Runner }

func New() *Plugin                        { return &Plugin{runner: commandRunner{}} }
func NewWithRunner(runner Runner) *Plugin { return &Plugin{runner: runner} }
func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "git", Name: "Git", Version: "1.0.0", Description: "Context-aware Git commands"}
}
