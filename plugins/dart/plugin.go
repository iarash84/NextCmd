// Package dart implements Dart and Flutter toolchain suggestions without exposing
// Dart-specific concepts to Core.
package dart

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin { return &Plugin{} }

func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "dart", Name: "Dart / Flutter", Version: "1.0.0", Description: "Context-aware Dart and Flutter toolchain commands"}
}
