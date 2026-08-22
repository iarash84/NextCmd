// Package dotnet implements .NET CLI suggestions without exposing .NET concepts to Core.
package dotnet

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin { return &Plugin{} }

func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "dotnet", Name: ".NET CLI", Version: "1.0.0", Description: "Context-aware .NET project commands"}
}
