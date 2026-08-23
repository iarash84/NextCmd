package terminal

import (
	"bytes"
	"context"
	"io"
	"testing"

	"nextcmd/sdk"
)

type directoryCompleter struct{ directory string }

func (c *directoryCompleter) Complete(_ context.Context, _ string, directory string, _ *sdk.ExecutionResult) []sdk.Suggestion {
	c.directory = directory
	return nil
}

func TestSetDirectoryUpdatesCompletionContextAndDisplay(t *testing.T) {
	var output bytes.Buffer
	completer := &directoryCompleter{}
	ui := &UI{input: bytes.NewReader([]byte{'\r'}), output: &output, directory: "old"}
	ui.SetDirectory("new-project")
	command, err := ui.ReadCommand(context.Background(), completer, nil)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if command != "" || completer.directory != "new-project" {
		t.Fatalf("command=%q directory=%q", command, completer.directory)
	}
	if !bytes.Contains(output.Bytes(), []byte("cwd new-project")) {
		t.Fatalf("working directory is not visible: %q", output.String())
	}
}

func TestAcceptSelected(t *testing.T) {
	suggestions := []sdk.Suggestion{{
		Command: sdk.Command{Executable: "git", Args: []string{"status"}},
	}}

	accepted, ok := acceptSelected("git sta", suggestions, 0)
	if !ok || accepted != "git status" {
		t.Fatalf("acceptSelected() = %q, %v", accepted, ok)
	}

	accepted, ok = acceptSelected("git status", suggestions, 0)
	if ok || accepted != "git status" {
		t.Fatalf("an already accepted command must be ready to execute: %q, %v", accepted, ok)
	}
}

func TestReadKeyRecognizesRightArrow(t *testing.T) {
	tests := []struct {
		name     string
		sequence []byte
	}{
		{"ANSI", []byte{27, '[', 'C'}},
		{"Windows", []byte{224, 77}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event, err := readKey(bytes.NewReader(test.sequence))
			if err != nil || event.kind != KeyRight {
				t.Fatalf("readKey() = %#v, %v", event, err)
			}
		})
	}
}

func TestReadKeyIgnoresUnsupportedArrowSequence(t *testing.T) {
	event, err := readKey(bytes.NewReader([]byte{27, '[', 'H'}))
	if err != nil || event.kind != KeyIgnored {
		t.Fatalf("readKey() = %#v, %v", event, err)
	}
}

func TestReadKeyRecognizesLeftArrow(t *testing.T) {
	tests := []struct {
		name     string
		sequence []byte
	}{
		{"ANSI", []byte{27, '[', 'D'}},
		{"Windows", []byte{224, 75}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event, err := readKey(bytes.NewReader(test.sequence))
			if err != nil || event.kind != KeyLeft {
				t.Fatalf("readKey() = %#v, %v", event, err)
			}
		})
	}
}

func TestAcceptSelectedRejectsInvalidIndex(t *testing.T) {
	line, ok := acceptSelected("git", nil, 0)
	if ok || line != "git" {
		t.Fatalf("invalid selection changed the line: %q, %v", line, ok)
	}
}

func TestDirectoryCommandsExecuteWithoutAcceptingPluginSuggestion(t *testing.T) {
	suggestions := []sdk.Suggestion{{Command: sdk.Command{Executable: "cargo", Args: []string{"add", "<crate>"}}}}
	for _, line := range []string{"cd ..", `cd "project with spaces"`, ":cd ..", "pwd", ":pwd", ":ls", ":ls ..", ":history", ":history 5", ":plugins", ":clear", ":config", ":which go", ":version"} {
		if accepted, ok := acceptSelected(line, suggestions, 0); ok || accepted != line {
			t.Errorf("directory command %q accepted suggestion: %q, %v", line, accepted, ok)
		}
	}
}

func TestClearResetsScreenAndRenderedSuggestions(t *testing.T) {
	var output bytes.Buffer
	ui := &UI{output: &output, rendered: 0}
	ui.Clear()
	if output.String() != "\x1b[2J\x1b[H" || ui.rendered != 0 {
		t.Fatalf("Clear() output=%q rendered=%d", output.String(), ui.rendered)
	}
}
