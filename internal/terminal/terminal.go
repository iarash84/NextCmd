package terminal

import (
	"bufio"
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
	KeyRight
	KeyLeft
	KeyBackspace
	KeyEscape
	KeyHome
	KeyEnd
	KeyHistoryPrevious
	KeyHistoryNext
	KeyHistorySearch
	KeyClearLine
	KeyEOF
	KeyIgnored
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
	caret              int
	color              bool
	history            []string
}

// caretForTest exposes the caret position for regression tests.
func (u *UI) caretForTest() int { return u.caret }

func New(directory string) *UI {
	return &UI{input: os.Stdin, output: os.Stdout, directory: directory, color: supportsColor(os.Stdout)}
}
func (u *UI) SetDirectory(directory string) { u.directory = directory }
func (u *UI) SetHistory(entries []sdk.HistoryEntry) {
	u.history = u.history[:0]
	for _, entry := range entries {
		u.AddHistory(entry.Command.Display())
	}
}
func (u *UI) AddHistory(command string) {
	command = strings.TrimSpace(command)
	if command == "" || (len(u.history) > 0 && u.history[len(u.history)-1] == command) {
		return
	}
	u.history = append(u.history, command)
}
func (u *UI) Clear() {
	u.clearSuggestions()
	fmt.Fprint(u.output, "\x1b[2J\x1b[H")
}
func (u *UI) ReadCommand(ctx context.Context, completer Completer, previous *sdk.ExecutionResult) (string, error) {
	restore, raw := makeRaw()
	if raw {
		defer restore()
	}
	line := ""
	historyIndex := len(u.history)
	draft := ""
	searching := false
	searchQuery := ""
	searchMatch := ""
	searchIndex := -1
	searchDraft := ""
	// caret is a byte offset into line. Left/Right move it so an accepted
	// command can be edited in place; appending at the end (the common
	// case) behaves exactly like the old line += ... code path.
	caret := 0
	u.caret = 0
	fmt.Fprintf(u.output, "%s %s\n", paint(u.color, ansiDim, "cwd"), paint(u.color, ansiCyan, u.directory))
	suggestions := completer.Complete(ctx, line, u.directory, previous)
	u.render(line, suggestions, caret)
	keyInput := bufio.NewReader(u.input)
	for {
		event, err := readKey(keyInput)
		if err != nil {
			return "", err
		}
		if searching {
			switch event.kind {
			case KeyRune:
				searchQuery += string(event.value)
				searchMatch, searchIndex = reverseHistorySearch(u.history, searchQuery, len(u.history))
			case KeyBackspace:
				if len(searchQuery) > 0 {
					searchQuery = searchQuery[:len(searchQuery)-1]
				}
				searchMatch, searchIndex = reverseHistorySearch(u.history, searchQuery, len(u.history))
			case KeyHistorySearch:
				before := searchIndex
				if before < 0 {
					before = len(u.history)
				}
				if match, index := reverseHistorySearch(u.history, searchQuery, before); index >= 0 {
					searchMatch, searchIndex = match, index
				}
			case KeyClearLine:
				searchQuery = ""
				searchMatch, searchIndex = reverseHistorySearch(u.history, searchQuery, len(u.history))
			case KeyEnter, KeyTab, KeyRight:
				if searchMatch != "" {
					line = searchMatch
					caret = len(line)
					historyIndex = len(u.history)
					searching = false
				}
			case KeyEscape:
				line = searchDraft
				caret = len(line)
				historyIndex = len(u.history)
				searching = false
			case KeyEOF:
				u.clearSuggestions()
				fmt.Fprint(u.output, "\r\x1b[2K\n")
				return "", io.EOF
			}
			if searching {
				u.renderHistorySearch(searchQuery, searchMatch)
				u.caret = len(searchQuery)
				continue
			}
			suggestions = completer.Complete(ctx, line, u.directory, previous)
			u.render(line, suggestions, caret)
			u.caret = caret
			continue
		}
		switch event.kind {
		case KeyRune:
			line = line[:caret] + string(event.value) + line[caret:]
			caret++
			u.selected = 0
			historyIndex = len(u.history)
		case KeyBackspace:
			if caret > 0 {
				line = line[:caret-1] + line[caret:]
				caret--
			}
			u.selected = 0
			historyIndex = len(u.history)
		case KeyUp:
			if len(suggestions) > 0 {
				u.selected = (u.selected - 1 + len(suggestions)) % len(suggestions)
			}
		case KeyDown:
			if len(suggestions) > 0 {
				u.selected = (u.selected + 1) % len(suggestions)
			}
		case KeyHistoryPrevious:
			if len(u.history) > 0 {
				if historyIndex == len(u.history) {
					draft = line
				}
				if historyIndex > 0 {
					historyIndex--
					line = u.history[historyIndex]
					caret = len(line)
					u.selected = 0
				}
			}
		case KeyHistoryNext:
			if historyIndex < len(u.history) {
				historyIndex++
				if historyIndex == len(u.history) {
					line = draft
				} else {
					line = u.history[historyIndex]
				}
				caret = len(line)
				u.selected = 0
			}
		case KeyHistorySearch:
			searching = true
			searchDraft = line
			searchQuery = line
			searchMatch, searchIndex = reverseHistorySearch(u.history, searchQuery, len(u.history))
			u.renderHistorySearch(searchQuery, searchMatch)
			u.caret = len(searchQuery)
			continue
		case KeyTab, KeyRight:
			if len(suggestions) > 0 {
				line = suggestions[u.selected].Command.Display()
				caret = len(line)
				historyIndex = len(u.history)
			} else if event.kind == KeyRight && caret < len(line) {
				caret++
			}
		case KeyLeft:
			// Move the caret instead of clearing the line so part of an
			// accepted command can be rewritten. Home (column 0) stops.
			if caret > 0 {
				caret--
			}
		case KeyEscape:
			line = ""
			caret = 0
			u.selected = 0
			historyIndex = len(u.history)
		case KeyHome:
			caret = 0
		case KeyEnd:
			caret = len(line)
		case KeyClearLine:
			line = ""
			caret = 0
			u.selected = 0
			historyIndex = len(u.history)
		case KeyEnter:
			if accepted, ok := acceptSelected(line, suggestions, u.selected); ok {
				line = accepted
				caret = len(line)
				u.selected = 0
				break
			}
			u.caret = caret
			u.clearSuggestions()
			fmt.Fprint(u.output, "\r\x1b[2K", paint(u.color, ansiBold+ansiCyan, "❯ "), line, "\n")
			return strings.TrimSpace(line), nil
		case KeyEOF:
			u.clearSuggestions()
			fmt.Fprint(u.output, "\r\x1b[2K\n")
			return "", io.EOF
		case KeyIgnored:
			// Unsupported terminal sequences must not terminate the session.
		}
		suggestions = completer.Complete(ctx, line, u.directory, previous)
		if u.selected >= len(suggestions) {
			u.selected = 0
		}
		u.render(line, suggestions, caret)
		u.caret = caret
	}
}

func reverseHistorySearch(history []string, query string, before int) (string, int) {
	if before > len(history) {
		before = len(history)
	}
	for index := before - 1; index >= 0; index-- {
		if fuzzyHistoryMatch(history[index], query) {
			return history[index], index
		}
	}
	return "", -1
}

func fuzzyHistoryMatch(command, query string) bool {
	queryRunes := []rune(strings.ToLower(strings.TrimSpace(query)))
	if len(queryRunes) == 0 {
		return true
	}
	matched := 0
	for _, candidate := range []rune(strings.ToLower(command)) {
		if candidate == queryRunes[matched] {
			matched++
			if matched == len(queryRunes) {
				return true
			}
		}
	}
	return false
}

// acceptSelected keeps suggestion acceptance separate from command execution.
// Enter accepts a highlighted suggestion first; a subsequent Enter executes it.
func acceptSelected(line string, suggestions []sdk.Suggestion, selected int) (string, bool) {
	if isInternalCommand(line) {
		return line, false
	}
	if selected < 0 || selected >= len(suggestions) {
		return line, false
	}
	command := suggestions[selected].Command.Display()
	if strings.TrimSpace(line) == command {
		return line, false
	}
	return command, true
}

func isInternalCommand(line string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(line))
	if trimmed == "cd" || trimmed == ":cd" || strings.HasPrefix(trimmed, "cd ") || strings.HasPrefix(trimmed, ":cd ") || trimmed == "pwd" || trimmed == ":pwd" || trimmed == ":ls" || strings.HasPrefix(trimmed, ":ls ") || trimmed == ":mkdir" || strings.HasPrefix(trimmed, ":mkdir ") || trimmed == ":del" || strings.HasPrefix(trimmed, ":del ") || trimmed == ":trash" || strings.HasPrefix(trimmed, ":trash ") || trimmed == ":undo" {
		return true
	}
	for _, name := range []string{":history", ":plugins", ":clear", ":config", ":which", ":version"} {
		if trimmed == name || strings.HasPrefix(trimmed, name+" ") {
			return true
		}
	}
	return false
}

func (u *UI) render(line string, suggestions []sdk.Suggestion, caret int) {
	u.clearSuggestions()
	fmt.Fprint(u.output, "\r\x1b[2K", paint(u.color, ansiBold+ansiCyan, "❯ "), line)
	for i, suggestion := range suggestions {
		marker := "  "
		if i == u.selected {
			marker = paint(u.color, ansiBold+ansiCyan, "❯ ")
		}
		commandColor := ansiWhite
		if i == u.selected {
			commandColor = ansiBold + ansiWhite
		}
		fmt.Fprintf(u.output, "\n\x1b[2K%s%s  %s %s %s",
			marker,
			paint(u.color, commandColor, suggestion.Command.Display()),
			paint(u.color, kindColor(suggestion.Kind), suggestionBadge(suggestion.Kind)),
			paint(u.color, riskColor(suggestion.Risk), strings.ToUpper(string(suggestion.Risk))),
			paint(u.color, ansiDim, "· "+suggestion.Source))
	}
	if len(suggestions) > 0 {
		fmt.Fprintf(u.output, "\x1b[%dA", len(suggestions))
	}
	// Park the terminal cursor on the editing caret, not the line end.
	fmt.Fprintf(u.output, "\r\x1b[%dC", caret+2)
	u.rendered = len(suggestions)
}
func (u *UI) renderHistorySearch(query, match string) {
	u.clearSuggestions()
	if match == "" {
		match = "no match"
	}
	fmt.Fprintf(u.output, "\r\x1b[2K%s %q: %s", paint(u.color, ansiBold+ansiCyan, "reverse-i-search"), query, match)
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
