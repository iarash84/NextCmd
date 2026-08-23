// Package docker implements Docker CLI suggestions without exposing Docker
// concepts to Core.
package docker

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin { return &Plugin{} }

func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "docker", Name: "Docker", Version: "1.0.0", Description: "Context-aware Docker and Compose commands"}
}
