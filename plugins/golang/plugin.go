// Package golang implements Go toolchain suggestions without exposing
// Go-specific concepts to Core.
package golang

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin { return &Plugin{} }

func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "go", Name: "Go", Version: "1.0.0", Description: "Context-aware Go toolchain commands"}
}
