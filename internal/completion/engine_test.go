package completion

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
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
	for _, input := range []string{"cd", "cd ..", ":cd project", "pwd", ":pwd", ":ls", ":ls ..", ":mkdir old", ":del old", ":trash old", ":undo", ":history", ":history 5", ":plugins", ":clear", ":config", ":which go", ":version"} {
		if suggestions := engine.Complete(t.Context(), input, ".", nil); len(suggestions) != 0 {
			t.Errorf("%q received plugin suggestions: %#v", input, suggestions)
		}
	}
}

func TestColonSuggestsBuiltinCommands(t *testing.T) {
	engine := New([]sdk.Plugin{catalogPlugin{}}, 20, slog.New(slog.NewTextHandler(io.Discard, nil)))
	suggestions := engine.Complete(t.Context(), ":", ".", nil)
	wanted := map[string]bool{":q": false, ":ls": false, ":mkdir <path>": false, ":del <path>": false, ":trash <path>": false, ":undo": false, ":plugins": false, ":history": false, ":clear": false, ":config": false, ":which <command>": false, ":version": false, ":?": false}
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

func TestBuiltinPathCompletionSuggestsMatchingPaths(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "old file.txt"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	engine := New(nil, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	suggestions := engine.Complete(t.Context(), ":del old", dir, nil)
	if len(suggestions) != 1 || suggestions[0].Command.Display() != `:del "old file.txt"` {
		t.Fatalf("Complete(:del old) = %#v", suggestions)
	}
}

func TestCommandPathCompletionPreservesEarlierArguments(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatal(err)
	}
	engine := New(nil, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	suggestions := engine.Complete(t.Context(), "tool run ma", dir, nil)
	if len(suggestions) != 1 || suggestions[0].Command.Display() != "tool run main.go" {
		t.Fatalf("Complete(tool run ma) = %#v", suggestions)
	}
}

func TestCommandPathCompletionQuotesSpacesAndDescendsDirectories(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "my folder"), 0o755); err != nil {
		t.Fatal(err)
	}
	engine := New(nil, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	suggestions := engine.Complete(t.Context(), `tool "my f`, dir, nil)
	want := (sdk.Command{Executable: "tool", Args: []string{"my folder" + string(os.PathSeparator)}}).Display()
	if len(suggestions) != 1 || suggestions[0].Command.Display() != want {
		t.Fatalf("Complete quoted path = %#v, want %q", suggestions, want)
	}
}

func TestCommandPathCompletionSupportsNestedPaths(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "src", "main.go"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	engine := New(nil, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	input := "tool " + filepath.Join("src", "ma")
	want := (sdk.Command{Executable: "tool", Args: []string{filepath.Join("src", "main.go")}}).Display()
	suggestions := engine.Complete(t.Context(), input, dir, nil)
	if len(suggestions) != 1 || suggestions[0].Command.Display() != want {
		t.Fatalf("Complete(%q) = %#v, want %q", input, suggestions, want)
	}
}

func TestCommandPathCompletionPreservesDotPrefix(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	engine := New(nil, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	input := "tool ." + string(os.PathSeparator) + "ma"
	want := (sdk.Command{Executable: "tool", Args: []string{"." + string(os.PathSeparator) + "main.go"}}).Display()
	suggestions := engine.Complete(t.Context(), input, dir, nil)
	if len(suggestions) != 1 || suggestions[0].Command.Display() != want {
		t.Fatalf("Complete(%q) = %#v, want %q", input, suggestions, want)
	}
}

func TestCommandPathCompletionContinuesInsideAcceptedDirectory(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "src")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "main.go"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	engine := New(nil, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	input := "tool src" + string(os.PathSeparator)
	want := (sdk.Command{Executable: "tool", Args: []string{filepath.Join("src", "main.go")}}).Display()
	suggestions := engine.Complete(t.Context(), input, dir, nil)
	if len(suggestions) != 1 || suggestions[0].Command.Display() != want {
		t.Fatalf("Complete(%q) = %#v, want %q", input, suggestions, want)
	}
}

func TestCommandPathCompletionIgnoresOptionsExecutablesAndShell(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	engine := New(nil, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, input := range []string{"to", "tool --ma", "! tool ma"} {
		if suggestions := engine.Complete(t.Context(), input, dir, nil); len(suggestions) != 0 {
			t.Errorf("Complete(%q) = %#v", input, suggestions)
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
