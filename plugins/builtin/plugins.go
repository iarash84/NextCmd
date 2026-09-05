package builtin

import (
	"nextcmd/plugins/cargo"
	"nextcmd/plugins/curl"
	"nextcmd/plugins/docker"
	"nextcmd/plugins/dotnet"
	"nextcmd/plugins/git"
	"nextcmd/plugins/golang"
	"nextcmd/plugins/kubernetes"
	"nextcmd/plugins/npm"
	"nextcmd/plugins/pip"
	"nextcmd/plugins/terraform"
	"nextcmd/sdk"
)

// All is the single explicit composition point for built-in plugins. Adding a
// plugin here requires no matching flag, constructor argument, or Core change.
func All() []sdk.Plugin {
	return []sdk.Plugin{
		git.New(),
		dotnet.New(),
		cargo.New(),
		curl.New(),
		golang.New(),
		docker.New(),
		npm.New(),
		pip.New(),
		kubernetes.New(),
		terraform.New(),
	}
}
