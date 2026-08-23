// Package cargo implements Cargo suggestions without exposing Rust concepts to Core.
package cargo

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin { return &Plugin{} }

func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "cargo", Name: "Cargo", Version: "1.0.0", Description: "Context-aware Rust and Cargo commands"}
}
