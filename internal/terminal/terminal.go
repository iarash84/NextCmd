package terminal

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"nextcmd/sdk"
)

type Completer interface {
	Complete(context.Context, string, string, *sdk.ExecutionResult) []sdk.Suggestion
}
type Key byte

const (
	KeyRune Key = iota
	KeyEnter
	KeyTab
	KeyUp
	KeyDown
	KeyBackspace
	KeyEscape
	KeyEOF
)

type keyEvent struct {
	kind  Key
	value byte
}
type UI struct {
	input              io.Reader
	output             io.Writer
	directory          string
	selected, rendered int
}

func New(directory string) *UI { return &UI{input: os.Stdin, output: os.Stdout, directory: directory} }
func (u *UI) ReadCommand(ctx context.Context, completer Completer, previous *sdk.ExecutionResult) (string, error) {
	restore, raw := makeRaw()
	if raw {
		defer restore()
	}
	line := ""
	suggestions := completer.Complete(ctx, line, u.directory, previous)
	u.render(line, suggestions)
	for {
		event, err := readKey(u.input)
		if err != nil {
			return "", err
		}
		switch event.kind {
		case KeyRune:
			line += string(event.value)
			u.selected = 0
		case KeyBackspace:
			if len(line) > 0 {
				line = line[:len(line)-1]
			}
			u.selected = 0
		case KeyUp:
			if len(suggestions) > 0 {
				u.selected = (u.selected - 1 + len(suggestions)) % len(suggestions)
			}
		case KeyDown:
			if len(suggestions) > 0 {
				u.selected = (u.selected + 1) % len(suggestions)
			}
		case KeyTab:
			if len(suggestions) > 0 {
				line = suggestions[u.selected].Command.Display()
			}
		case KeyEnter:
			u.clearSuggestions()
			fmt.Fprint(u.output, "\r\x1b[2K> ", line, "\n")
			return strings.TrimSpace(line), nil
		case KeyEscape, KeyEOF:
			return "", io.EOF
		}
		suggestions = completer.Complete(ctx, line, u.directory, previous)
		if u.selected >= len(suggestions) {
			u.selected = 0
		}
		u.render(line, suggestions)
	}
}
func (u *UI) render(line string, suggestions []sdk.Suggestion) {
	u.clearSuggestions()
	fmt.Fprint(u.output, "\r\x1b[2K> ", line)
	for i, suggestion := range suggestions {
		marker := "  "
		if i == u.selected {
			marker = "❯ "
		}
		fmt.Fprintf(u.output, "\n\x1b[2K%s%s  [%s/%s]", marker, suggestion.Command.Display(), suggestion.Kind, suggestion.Risk)
	}
	if len(suggestions) > 0 {
		fmt.Fprintf(u.output, "\x1b[%dA\r\x1b[%dC", len(suggestions), len(line)+2)
	}
	u.rendered = len(suggestions)
}
func (u *UI) clearSuggestions() {
	if u.rendered == 0 {
		return
	}
	for i := 0; i < u.rendered; i++ {
		fmt.Fprint(u.output, "\x1b[1B\r\x1b[2K")
	}
	fmt.Fprintf(u.output, "\x1b[%dA", u.rendered)
	u.rendered = 0
}
