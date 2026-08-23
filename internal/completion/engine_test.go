package completion

import (
	"io"
	"log/slog"
	"testing"

	"nextcmd/sdk"
)

type catalogPlugin struct{}

func (catalogPlugin) Info() sdk.PluginInfo { return sdk.PluginInfo{ID: "tool"} }
func (catalogPlugin) Help() []sdk.CommandHelp {
	return []sdk.CommandHelp{{Command: sdk.Command{Executable: "tool", Args: []string{"run"}}}}
}

func TestPluginForExecutableUsesPublicCatalog(t *testing.T) {
	engine := New([]sdk.Plugin{catalogPlugin{}}, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if got := engine.PluginForExecutable("TOOL"); got != "tool" {
		t.Fatalf("PluginForExecutable() = %q", got)
	}
	if got := engine.PluginForExecutable("unknown"); got != "" {
		t.Fatalf("unknown executable matched %q", got)
	}
}

func TestDirectoryCommandsDoNotReceivePluginSuggestions(t *testing.T) {
	engine := New([]sdk.Plugin{catalogPlugin{}}, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, input := range []string{"cd", "cd ..", ":cd project", "pwd", ":pwd"} {
		if suggestions := engine.Complete(t.Context(), input, ".", nil); len(suggestions) != 0 {
			t.Errorf("%q received plugin suggestions: %#v", input, suggestions)
		}
	}
}
