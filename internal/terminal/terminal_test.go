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

// runKeystrokes feeds keystrokes to ReadCommand and ends the line with
// Enter, returning the submitted line and the caret position at submit
// time. (EOF returns an empty line, so it cannot carry caret state.)
func runKeystrokes(t *testing.T, keystrokes []byte) (string, int) {
	t.Helper()
	var output bytes.Buffer
	completer := &directoryCompleter{}
	ui := &UI{input: bytes.NewReader(keystrokes), output: &output, directory: "d"}
	line, err := ui.ReadCommand(context.Background(), completer, nil)
	if err != nil {
		t.Fatal(err)
	}
	return line, ui.caretForTest()
}

func TestLeftArrowMovesCaretInsteadOfClearingLine(t *testing.T) {
	// "ab" then Left, Left (caret=0), then Left again (no-op), Enter.
	line, caret := runKeystrokes(t, []byte("ab\x1b[D\x1b[D\x1b[D\r"))
	if line != "ab" {
		t.Fatalf("left arrow must not clear the line: %q", line)
	}
	if caret != 0 {
		t.Fatalf("caret should be at home after 3 lefts: %d", caret)
	}
}

func TestTypingAtCaretInsertsMidLine(t *testing.T) {
	// "ac", Left (caret=1), type "b" => "abc" with caret at 2.
	line, caret := runKeystrokes(t, []byte("ac\x1b[Db\r"))
	if line != "abc" {
		t.Fatalf("mid-line insert failed: %q", line)
	}
	if caret != 2 {
		t.Fatalf("caret after insert: %d", caret)
	}
}

func TestRightArrowMovesCaretWhenNoSuggestions(t *testing.T) {
	// "ab", Left, Left, Right, Right, Right (clamped at end), Enter.
	line, caret := runKeystrokes(t, []byte("ab\x1b[D\x1b[D\x1b[C\x1b[C\x1b[C\r"))
	if line != "ab" {
		t.Fatalf("right arrow must not change the line: %q", line)
	}
	if caret != 2 {
		t.Fatalf("caret should clamp at line end: %d", caret)
	}
}

func TestAcceptSelectedRejectsInvalidIndex(t *testing.T) {
	line, ok := acceptSelected("git", nil, 0)
	if ok || line != "git" {
		t.Fatalf("invalid selection changed the line: %q, %v", line, ok)
	}
}

func TestAcceptSelectedAcceptsBuiltinFromColonPrefix(t *testing.T) {
	suggestions := []sdk.Suggestion{{Command: sdk.Command{Executable: ":plugins"}, Source: "nextcmd"}}
	line, ok := acceptSelected(":pl", suggestions, 0)
	if !ok || line != ":plugins" {
		t.Fatalf("acceptSelected() = %q, %v", line, ok)
	}
	if line, ok = acceptSelected(line, suggestions, 0); ok || line != ":plugins" {
		t.Fatalf("accepted command must be ready to execute: %q, %v", line, ok)
	}
}

func TestDirectoryCommandsExecuteWithoutAcceptingPluginSuggestion(t *testing.T) {
	suggestions := []sdk.Suggestion{{Command: sdk.Command{Executable: "cargo", Args: []string{"add", "<crate>"}}}}
	for _, line := range []string{"cd ..", `cd "project with spaces"`, ":cd ..", "pwd", ":pwd", ":ls", ":ls ..", ":mkdir old", ":del old", ":history", ":history 5", ":plugins", ":clear", ":config", ":which go", ":version"} {
		if accepted, ok := acceptSelected(line, suggestions, 0); ok || accepted != line {
			t.Errorf("directory command %q accepted suggestion: %q, %v", line, accepted, ok)
		}
	}
}

func TestEscapeClearsLine(t *testing.T) {
	line, caret := runKeystrokes(t, []byte("abc\x1bde\r"))
	if line != "de" {
		t.Fatalf("escape should clear current line before more typing: %q", line)
	}
	if caret != 2 {
		t.Fatalf("caret after escape and typing: %d", caret)
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
