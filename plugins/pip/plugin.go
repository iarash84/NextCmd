// Package pip implements pip and pip3 suggestions without exposing Python concepts to Core.
package pip

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin { return &Plugin{} }
func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "pip", Name: "pip / pip3", Version: "1.0.0", Description: "Context-aware Python package management commands"}
}
