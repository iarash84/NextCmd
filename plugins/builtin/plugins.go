package builtin

import (
	"nextcmd/plugins/cargo"
	"nextcmd/plugins/dotnet"
	"nextcmd/plugins/git"
	"nextcmd/sdk"
)

// All is the single explicit composition point for built-in plugins. Adding a
// plugin here requires no matching flag, constructor argument, or Core change.
func All() []sdk.Plugin {
	return []sdk.Plugin{
		git.New(),
		dotnet.New(),
		cargo.New(),
	}
}
