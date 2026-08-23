// Package curl implements curl suggestions without exposing HTTP concepts to Core.
package curl

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin { return &Plugin{} }

func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "curl", Name: "Curl", Version: "1.0.0", Description: "HTTP transfer and diagnostics commands"}
}
