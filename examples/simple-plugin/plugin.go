// Package simpleplugin demonstrates a plugin that depends only on the public SDK.
package simpleplugin

import (
	"context"
	"nextcmd/sdk"
	"strings"
)

type Plugin struct{}

func New() Plugin { return Plugin{} }
func (Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "simple", Name: "Simple", Version: "1.0.0"}
}
func (Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	if input.Input != "" && !strings.HasPrefix("hello", input.Input) {
		return nil, nil
	}
	return []sdk.Suggestion{{Command: sdk.Command{Executable: "hello"}, Title: "Say hello", Kind: sdk.Completion, Risk: sdk.Safe, Priority: 50}}, nil
}
