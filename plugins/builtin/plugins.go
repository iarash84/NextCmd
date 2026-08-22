package builtin

import (
	"nextcmd/plugins/dotnet"
	"nextcmd/plugins/git"
	"nextcmd/sdk"
)

func All(gitEnabled, dotnetEnabled bool) []sdk.Plugin {
	plugins := []sdk.Plugin{}
	if gitEnabled {
		plugins = append(plugins, git.New())
	}
	if dotnetEnabled {
		plugins = append(plugins, dotnet.New())
	}
	return plugins
}
