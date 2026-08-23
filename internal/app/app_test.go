package app

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"

	"nextcmd/internal/completion"
	"nextcmd/sdk"
)

func TestIsExitCommand(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"exit", true},
		{" QUIT ", true},
		{":q", true},
		{"git status", false},
		{"exit now", false},
	}

	for _, test := range tests {
		if got := isExitCommand(test.input); got != test.want {
			t.Errorf("isExitCommand(%q) = %v, want %v", test.input, got, test.want)
		}
	}
}

type helpPlugin struct{}

func (helpPlugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "demo", Name: "Demo", Description: "Demo commands"}
}
func (helpPlugin) Help() []sdk.CommandHelp {
	return []sdk.CommandHelp{{Command: sdk.Command{Executable: "demo", Args: []string{"run"}}, Description: "Run demo", Risk: sdk.Safe}}
}

func TestParseHelpCommand(t *testing.T) {
	tests := []struct {
		input, plugin string
		help          bool
	}{{":?", "", true}, {":؟", "", true}, {":? git", "git", true}, {":؟ dotnet", "dotnet", true}, {"git status", "", false}}
	for _, test := range tests {
		plugin, help := parseHelpCommand(test.input)
		if plugin != test.plugin || help != test.help {
			t.Errorf("parseHelpCommand(%q) = %q, %v", test.input, plugin, help)
		}
	}
}

func TestPrintPluginHelp(t *testing.T) {
	engine := completion.New([]sdk.Plugin{helpPlugin{}}, 8, slog.New(slog.NewTextHandler(io.Discard, nil)))
	var output bytes.Buffer
	printHelp(&output, engine, "demo")
	if text := output.String(); !strings.Contains(text, "demo run") || !strings.Contains(text, "Run demo") {
		t.Fatalf("unexpected help: %s", text)
	}
}
