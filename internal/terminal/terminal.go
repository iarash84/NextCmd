package terminal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"nextcmd/sdk"
)

type Completer interface {
	Complete(context.Context, string, string, *sdk.ExecutionResult) []sdk.Suggestion
}
type Key byte

const (
	KeyRune Key = iota
	KeyEnter
	KeyNewline
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
	KeyPaste
	KeyEOF
	KeyIgnored
)

type keyEvent struct {
	kind  Key
	value byte
	text  string
}
type placeholderRange struct {
	start, end int
}
type UI struct {
	input              io.Reader
	output             io.Writer
	directory          string
	selected, rendered int
	renderedCursorRow  int
	caret              int
	color              bool
	history            []string
	keyInput           *bufio.Reader
	pending            []string
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
	if len(u.pending) > 0 {
		line := u.pending[0]
		u.pending = u.pending[1:]
		fmt.Fprintf(u.output, "%s %s\n", paint(u.color, ansiDim, "cwd"), paint(u.color, ansiCyan, u.directory))
		fmt.Fprint(u.output, paint(u.color, ansiBold+ansiCyan, "❯ "), line, "\n")
		return line, nil
	}
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
	placeholderMode := false
	placeholders := []placeholderRange(nil)
	activePlaceholder := -1
	editingAccepted := false
	// caret is a byte offset into line. Left/Right move it so an accepted
	// command can be edited in place; appending at the end (the common
	// case) behaves exactly like the old line += ... code path.
	caret := 0
	u.caret = 0
	fmt.Fprintf(u.output, "%s %s\n", paint(u.color, ansiDim, "cwd"), paint(u.color, ansiCyan, u.directory))
	suggestions := completeEditorLine(ctx, completer, line, u.directory, previous)
	u.render(line, suggestions, caret)
	if u.keyInput == nil {
		u.keyInput = bufio.NewReader(u.input)
	}
	for {
		event, handled, err := readConsoleKey(raw && u.input == os.Stdin)
		if err == nil && !handled {
			event, err = readKeyMode(u.keyInput, raw)
		}
		if err != nil {
			return "", err
		}
		if searching {
			switch event.kind {
			case KeyRune:
				searchQuery += keyText(event)
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
					placeholderMode = false
					placeholders = nil
					activePlaceholder = -1
					editingAccepted = false
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
			suggestions = completeEditorLine(ctx, completer, line, u.directory, previous)
			u.render(line, suggestions, caret)
			u.caret = caret
			continue
		}
		switch event.kind {
		case KeyRune:
			inserted := keyText(event)
			if placeholderMode && activePlaceholder >= 0 && activePlaceholder < len(placeholders) {
				placeholder := placeholders[activePlaceholder]
				line = line[:placeholder.start] + inserted + line[placeholder.end:]
				caret = placeholder.start
				activePlaceholder = -1
			} else {
				line = line[:caret] + inserted + line[caret:]
			}
			caret += len(inserted)
			if placeholderMode {
				placeholders = findPlaceholderRanges(line)
				placeholderMode = len(placeholders) > 0
			}
			u.selected = 0
			historyIndex = len(u.history)
		case KeyPaste:
			if placeholderMode && activePlaceholder >= 0 && activePlaceholder < len(placeholders) {
				placeholder := placeholders[activePlaceholder]
				line = line[:placeholder.start] + event.text + line[placeholder.end:]
				caret = placeholder.start + len(event.text)
			} else {
				line = line[:caret] + event.text + line[caret:]
				caret += len(event.text)
			}
			placeholderMode = false
			placeholders = nil
			activePlaceholder = -1
			editingAccepted = true
			u.selected = 0
			historyIndex = len(u.history)
		case KeyNewline:
			line = line[:caret] + "\n" + line[caret:]
			caret++
			placeholderMode = false
			placeholders = nil
			activePlaceholder = -1
			editingAccepted = true
			u.selected = 0
			historyIndex = len(u.history)
		case KeyBackspace:
			if placeholderMode && activePlaceholder >= 0 && activePlaceholder < len(placeholders) {
				placeholder := placeholders[activePlaceholder]
				line = line[:placeholder.start] + line[placeholder.end:]
				caret = placeholder.start
				activePlaceholder = -1
			} else if caret > 0 {
				line = line[:caret-1] + line[caret:]
				caret--
			}
			if placeholderMode {
				placeholders = findPlaceholderRanges(line)
				placeholderMode = len(placeholders) > 0
			}
			u.selected = 0
			historyIndex = len(u.history)
		case KeyUp:
			if len(suggestions) > 0 {
				u.selected = (u.selected - 1 + len(suggestions)) % len(suggestions)
				editingAccepted = false
			}
		case KeyDown:
			if len(suggestions) > 0 {
				u.selected = (u.selected + 1) % len(suggestions)
				editingAccepted = false
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
					placeholderMode = false
					editingAccepted = false
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
				placeholderMode = false
				editingAccepted = false
			}
		case KeyHistorySearch:
			searching = true
			searchDraft = line
			searchQuery = line
			searchMatch, searchIndex = reverseHistorySearch(u.history, searchQuery, len(u.history))
			u.renderHistorySearch(searchQuery, searchMatch)
			u.caret = len(searchQuery)
			continue
		case KeyTab:
			if placeholderMode && len(placeholders) > 0 {
				activePlaceholder = nextPlaceholder(placeholders, activePlaceholder, caret)
				caret = placeholders[activePlaceholder].start
				break
			}
			if editingAccepted {
				caret = len(line)
				break
			}
			if len(suggestions) > 0 {
				line = suggestions[u.selected].Command.Display()
				placeholders = placeholderRangesForSuggestion(suggestions[u.selected])
				placeholderMode = len(placeholders) > 0
				activePlaceholder = -1
				caret = len(line)
				if placeholderMode {
					activePlaceholder = 0
					caret = placeholders[0].start
				}
				editingAccepted = true
				historyIndex = len(u.history)
			}
		case KeyRight:
			if !editingAccepted && len(suggestions) > 0 {
				line = suggestions[u.selected].Command.Display()
				placeholders = placeholderRangesForSuggestion(suggestions[u.selected])
				placeholderMode = len(placeholders) > 0
				activePlaceholder = -1
				caret = len(line)
				if placeholderMode {
					activePlaceholder = 0
					caret = placeholders[0].start
				}
				editingAccepted = true
				historyIndex = len(u.history)
			} else if caret < len(line) {
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
			placeholderMode = false
			placeholders = nil
			activePlaceholder = -1
			editingAccepted = false
		case KeyHome:
			caret = 0
		case KeyEnd:
			caret = len(line)
		case KeyClearLine:
			line = ""
			caret = 0
			u.selected = 0
			historyIndex = len(u.history)
			placeholderMode = false
			placeholders = nil
			activePlaceholder = -1
			editingAccepted = false
		case KeyEnter:
			if placeholderMode && len(placeholders) > 0 {
				if activePlaceholder < 0 {
					activePlaceholder = nextPlaceholder(placeholders, -1, caret)
				}
				caret = placeholders[activePlaceholder].start
				break
			}
			commands := splitCommandLines(line)
			if len(commands) == 0 {
				line = ""
				caret = 0
				break
			}
			line = commands[0]
			u.pending = append(u.pending, commands[1:]...)
			if len(commands) > 1 {
				u.caret = len(line)
			} else {
				u.caret = caret
			}
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
		suggestions = completeEditorLine(ctx, completer, line, u.directory, previous)
		if u.selected >= len(suggestions) {
			u.selected = 0
		}
		u.render(line, suggestions, caret)
		u.caret = caret
	}
}

func keyText(event keyEvent) string {
	if event.text != "" {
		return event.text
	}
	return string(event.value)
}

func completeEditorLine(ctx context.Context, completer Completer, line, directory string, previous *sdk.ExecutionResult) []sdk.Suggestion {
	// Suggestion acceptance currently replaces the complete editor value. Hide
	// suggestions in multi-line mode so Tab/Right cannot discard earlier lines.
	if strings.ContainsRune(line, '\n') {
		return nil
	}
	return completer.Complete(ctx, line, directory, previous)
}

func findPlaceholderRanges(line string) []placeholderRange {
	var ranges []placeholderRange
	for offset := 0; offset < len(line); {
		start := strings.IndexByte(line[offset:], '<')
		if start < 0 {
			break
		}
		start += offset
		end := strings.IndexByte(line[start+1:], '>')
		if end < 0 {
			break
		}
		end += start + 2
		if end > start+2 {
			ranges = append(ranges, placeholderRange{start: start, end: end})
		}
		offset = end
	}
	return ranges
}

func placeholderRangesForSuggestion(suggestion sdk.Suggestion) []placeholderRange {
	if len(suggestion.Placeholders) == 0 {
		return nil
	}
	return findPlaceholderRanges(suggestion.Command.Display())
}

func nextPlaceholder(placeholders []placeholderRange, active, caret int) int {
	if len(placeholders) == 0 {
		return -1
	}
	if active >= 0 && active < len(placeholders) {
		return (active + 1) % len(placeholders)
	}
	for index, placeholder := range placeholders {
		if placeholder.start >= caret {
			return index
		}
	}
	return 0
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

func (u *UI) render(line string, suggestions []sdk.Suggestion, caret int) {
	u.clearSuggestions()
	lines := strings.Split(line, "\n")
	beforeCaret := strings.Split(line[:caret], "\n")
	caretRow := len(beforeCaret) - 1
	caretColumn := utf8.RuneCountInString(beforeCaret[len(beforeCaret)-1])

	for index, editorLine := range lines {
		if index > 0 {
			fmt.Fprint(u.output, "\n")
		}
		prompt := paint(u.color, ansiBold+ansiCyan, "❯ ")
		if index > 0 {
			prompt = paint(u.color, ansiDim+ansiCyan, "│ ")
		}
		fmt.Fprint(u.output, "\r\x1b[2K", prompt, editorLine)
	}
	for i, suggestion := range suggestions {
		marker := "  "
		if i == u.selected {
			marker = paint(u.color, ansiBold+ansiCyan, "❯ ")
		}
		commandColor := ansiWhite
		if i == u.selected {
			commandColor = ansiBold + ansiWhite
		}
		fmt.Fprintf(u.output, "\n\r\x1b[2K%s%s  %s %s %s",
			marker,
			paint(u.color, commandColor, suggestion.Command.Display()),
			paint(u.color, kindColor(suggestion.Kind), suggestionBadge(suggestion.Kind)),
			paint(u.color, riskColor(suggestion.Risk), strings.ToUpper(string(suggestion.Risk))),
			paint(u.color, ansiDim, "· "+suggestion.Source))
	}
	totalRows := len(lines) + len(suggestions)
	lastRow := totalRows - 1
	if rowsUp := lastRow - caretRow; rowsUp > 0 {
		fmt.Fprintf(u.output, "\x1b[%dA", rowsUp)
	}
	fmt.Fprintf(u.output, "\r\x1b[%dC", caretColumn+2)
	u.rendered = totalRows
	u.renderedCursorRow = caretRow
}

func splitCommandLines(input string) []string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")
	lines := strings.Split(input, "\n")
	commands := make([]string, 0, len(lines))
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			commands = append(commands, line)
		}
	}
	return commands
}
func (u *UI) renderHistorySearch(query, match string) {
	u.clearSuggestions()
	if match == "" {
		match = "no match"
	}
	fmt.Fprintf(u.output, "\r\x1b[2K%s %q: %s", paint(u.color, ansiBold+ansiCyan, "reverse-i-search"), query, match)
	u.rendered = 1
	u.renderedCursorRow = 0
}
func (u *UI) clearSuggestions() {
	if u.rendered == 0 {
		return
	}
	if u.renderedCursorRow > 0 {
		fmt.Fprintf(u.output, "\x1b[%dA", u.renderedCursorRow)
	}
	fmt.Fprint(u.output, "\r")
	for row := 0; row < u.rendered; row++ {
		fmt.Fprint(u.output, "\x1b[2K")
		if row < u.rendered-1 {
			fmt.Fprint(u.output, "\x1b[1B\r")
		}
	}
	if u.rendered > 1 {
		fmt.Fprintf(u.output, "\x1b[%dA", u.rendered-1)
	}
	fmt.Fprint(u.output, "\r")
	u.rendered = 0
	u.renderedCursorRow = 0
}
