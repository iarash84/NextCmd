// Package npm implements npm suggestions without exposing Node.js concepts to Core.
package npm

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin { return &Plugin{} }
func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "npm", Name: "npm", Version: "1.0.0", Description: "Context-aware npm project and workspace commands"}
}
