package builtin

import (
	"nextcmd/plugins/cargo"
	"nextcmd/plugins/dotnet"
	"nextcmd/plugins/git"
	"nextcmd/sdk"
)

func All(gitEnabled, dotnetEnabled, cargoEnabled bool) []sdk.Plugin {
	plugins := []sdk.Plugin{}
	if gitEnabled {
		plugins = append(plugins, git.New())
	}
	if dotnetEnabled {
		plugins = append(plugins, dotnet.New())
	}
	if cargoEnabled {
		plugins = append(plugins, cargo.New())
	}
	return plugins
}
