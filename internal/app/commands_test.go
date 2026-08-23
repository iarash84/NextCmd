package app

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"nextcmd/internal/completion"
	"nextcmd/internal/history"
	"nextcmd/sdk"
)

func TestParseUtilityCommand(t *testing.T) {
	tests := []struct {
		input, name, arg string
		count            int
		handled          bool
	}{
		{input: ":history", name: ":history", count: 20, handled: true},
		{input: ":history 5", name: ":history", count: 5, handled: true},
		{input: ":PLUGINS", name: ":plugins", handled: true},
		{input: ":clear", name: ":clear", handled: true},
		{input: ":config", name: ":config", handled: true},
		{input: ":which go", name: ":which", arg: "go", handled: true},
		{input: ":version", name: ":version", handled: true},
		{input: "git status", handled: false},
	}
	for _, test := range tests {
		got, handled, err := parseUtilityCommand(test.input)
		if err != nil || handled != test.handled || got.name != test.name || got.arg != test.arg || got.count != test.count {
			t.Errorf("parseUtilityCommand(%q) = %#v, %v, %v", test.input, got, handled, err)
		}
	}
}

func TestParseUtilityCommandRejectsInvalidArguments(t *testing.T) {
	for _, input := range []string{":history 0", ":history 1001", ":history many", ":which", ":which go extra", ":clear now", ":version extra"} {
		command, handled, err := parseUtilityCommand(input)
		if !handled || err == nil || command.name == "" {
			t.Errorf("parseUtilityCommand(%q) = %#v, %v, %v", input, command, handled, err)
		}
	}
}

func TestPrintHistoryUsesRedactedStoredCommands(t *testing.T) {
	store := history.New(t.TempDir()+"/history.jsonl", true)
	entry := sdk.HistoryEntry{
		Command:          sdk.Command{Executable: "tool", Args: []string{"--token", "secret-value"}},
		WorkingDirectory: "project",
		Timestamp:        time.Date(2026, 8, 23, 10, 30, 0, 0, time.Local),
		ExitCode:         0,
		Duration:         time.Second,
		Plugin:           "example",
	}
	if err := store.Append(entry); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := printHistory(&output, store, 1); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, want := range []string{"TIME", "EXIT", "DURATION", "DIRECTORY", "tool --token <redacted>", "project", "example"} {
		if !strings.Contains(text, want) {
			t.Errorf("history does not contain %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "secret-value") {
		t.Fatalf("history output leaked a secret: %s", text)
	}
}

type utilityPlugin struct{ id string }

func (p utilityPlugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: p.id, Name: strings.ToUpper(p.id), Version: "1.0.0", Description: p.id + " commands"}
}

func TestPrintPluginsIsSorted(t *testing.T) {
	engine := completion.New([]sdk.Plugin{utilityPlugin{"zeta"}, utilityPlugin{"alpha"}}, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	var output bytes.Buffer
	printPlugins(&output, engine)
	text := output.String()
	if strings.Index(text, "alpha") > strings.Index(text, "zeta") {
		t.Fatalf("plugins are not sorted:\n%s", text)
	}
}

func TestPrintConfigSortsPluginOverrides(t *testing.T) {
	settings := RuntimeSettings{ConfigPath: "config.json", HistoryEnabled: true, MaxSuggestions: 8, Debug: false, PluginOverrides: map[string]bool{"zeta": false, "alpha": true}}
	var output bytes.Buffer
	printConfig(&output, settings, "history.jsonl")
	text := output.String()
	if strings.Index(text, "alpha") > strings.Index(text, "zeta") {
		t.Fatalf("plugin overrides are not sorted:\n%s", text)
	}
}

func TestFindExecutableAcceptsAnExistingExecutablePath(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	got, err := findExecutable(executable)
	if err != nil || got == "" {
		t.Fatalf("findExecutable() = %q, %v", got, err)
	}
}
