package builtin

import (
	"nextcmd/plugins/git"
	"nextcmd/sdk"
)

func All(gitEnabled bool) []sdk.Plugin {
	if !gitEnabled {
		return nil
	}
	return []sdk.Plugin{git.New()}
}
