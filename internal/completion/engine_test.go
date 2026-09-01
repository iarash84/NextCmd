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
	for _, input := range []string{"cd", "cd ..", ":cd project", "pwd", ":pwd", ":ls", ":ls ..", ":mkdir old", ":del old", ":history", ":history 5", ":plugins", ":clear", ":config", ":which go", ":version"} {
		if suggestions := engine.Complete(t.Context(), input, ".", nil); len(suggestions) != 0 {
			t.Errorf("%q received plugin suggestions: %#v", input, suggestions)
		}
	}
}

func TestColonSuggestsBuiltinCommands(t *testing.T) {
	engine := New([]sdk.Plugin{catalogPlugin{}}, 20, slog.New(slog.NewTextHandler(io.Discard, nil)))
	suggestions := engine.Complete(t.Context(), ":", ".", nil)
	wanted := map[string]bool{":q": false, ":ls": false, ":mkdir <path>": false, ":del <path>": false, ":plugins": false, ":history": false, ":clear": false, ":config": false, ":which <command>": false, ":version": false, ":?": false}
	for _, suggestion := range suggestions {
		if _, ok := wanted[suggestion.Command.Display()]; ok {
			wanted[suggestion.Command.Display()] = true
		}
		if suggestion.Source != "nextcmd" {
			t.Errorf("built-in suggestion source = %q", suggestion.Source)
		}
	}
	for command, found := range wanted {
		if !found {
			t.Errorf("%s was not suggested for colon input", command)
		}
	}
}

func TestBuiltinSuggestionsSupportPrefixFiltering(t *testing.T) {
	engine := New(nil, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	suggestions := engine.Complete(t.Context(), ":pl", ".", nil)
	if len(suggestions) != 1 || suggestions[0].Command.Display() != ":plugins" {
		t.Fatalf("Complete(:pl) = %#v", suggestions)
	}
}

func TestAcceptedBuiltinCommandHasNoSuggestions(t *testing.T) {
	engine := New(nil, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if suggestions := engine.Complete(t.Context(), ":plugins", ".", nil); len(suggestions) != 0 {
		t.Fatalf("accepted built-in received suggestions: %#v", suggestions)
	}
}
