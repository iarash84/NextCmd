package git

import (
	"nextcmd/internal/execution"
	"nextcmd/sdk"
)

type Plugin struct{ runner sdk.Runner }

func New() *Plugin { return NewWithRunner(execution.Executor{}) }
func NewWithRunner(runner sdk.Runner) *Plugin {
	if runner == nil {
		runner = execution.Executor{}
	}
	return &Plugin{runner: runner}
}
func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "git", Name: "Git", Version: "1.0.0", Description: "Context-aware Git commands"}
}
