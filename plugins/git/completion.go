package git

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

type commandSpec struct {
	args     []string
	title    string
	risk     sdk.Risk
	priority int
}

var commands = []commandSpec{
	{[]string{"status"}, "Show working tree status", sdk.Safe, 90}, {[]string{"add", "."}, "Stage all changes", sdk.Mutating, 70},
	{[]string{"commit", "-m", "<message>"}, "Commit staged changes", sdk.Mutating, 70}, {[]string{"diff"}, "Review unstaged changes", sdk.Safe, 80},
	{[]string{"diff", "--cached"}, "Review staged changes", sdk.Safe, 85}, {[]string{"log", "--oneline"}, "Show concise history", sdk.Safe, 65},
	{[]string{"branch"}, "List branches", sdk.Safe, 60}, {[]string{"switch", "<branch>"}, "Switch branch", sdk.Mutating, 55},
	{[]string{"checkout", "<branch>"}, "Checkout branch", sdk.Mutating, 45}, {[]string{"restore", "<file>"}, "Restore a file", sdk.Destructive, 40},
	{[]string{"stash"}, "Stash changes", sdk.Mutating, 45}, {[]string{"pull"}, "Pull upstream changes", sdk.Mutating, 55},
	{[]string{"push"}, "Push commits", sdk.Mutating, 55}, {[]string{"fetch"}, "Fetch remotes", sdk.Safe, 50},
	{[]string{"remote", "-v"}, "List remotes", sdk.Safe, 45}, {[]string{"merge", "<branch>"}, "Merge a branch", sdk.Mutating, 35},
	{[]string{"rebase", "<branch>"}, "Rebase onto a branch", sdk.Destructive, 25}, {[]string{"init"}, "Initialize repository", sdk.Mutating, 75},
	{[]string{"clone", "<repository>"}, "Clone repository", sdk.Mutating, 70},
}

func (p *Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	trimmed := strings.TrimSpace(input.Input)
	if trimmed != "" && !strings.HasPrefix(strings.ToLower(trimmed), "git") {
		return nil, nil
	}
	state, _ := input.Project.(State)
	suggestions := []sdk.Suggestion{}
	for _, spec := range commands {
		repoSpecific := spec.args[0] != "init" && spec.args[0] != "clone"
		if !state.InRepository && repoSpecific {
			continue
		}
		priority := spec.priority
		if state.Modified && (spec.args[0] == "diff" || spec.args[0] == "add") {
			priority += 30
		}
		if state.Staged && (spec.args[0] == "commit" || (len(spec.args) > 1 && spec.args[1] == "--cached")) {
			priority += 40
		}
		if state.Ahead && spec.args[0] == "push" {
			priority += 40
		}
		suggestions = append(suggestions, makeSuggestion(spec.args, spec.title, sdk.Completion, spec.risk, priority, "Matches current input and repository state"))
	}
	suggestions = append(suggestions, dynamic(input.Input, state)...)
	return suggestions, nil
}
func dynamic(input string, state State) []sdk.Suggestion {
	lower := strings.ToLower(strings.TrimSpace(input))
	var values []string
	var prefix []string
	risk := sdk.Mutating
	switch {
	case strings.HasPrefix(lower, "git switch"), strings.HasPrefix(lower, "git checkout"), strings.HasPrefix(lower, "git branch -d"), strings.HasPrefix(lower, "git merge"), strings.HasPrefix(lower, "git rebase"):
		values = state.Branches
		fields := strings.Fields(input)
		if len(fields) >= 2 {
			prefix = fields[1:]
		}
	case strings.HasPrefix(lower, "git add"), strings.HasPrefix(lower, "git restore"):
		values = state.Files
		fields := strings.Fields(input)
		if len(fields) >= 2 {
			prefix = fields[1:]
		}
		if strings.HasPrefix(lower, "git restore") {
			risk = sdk.Destructive
		}
	case strings.HasPrefix(lower, "git push"), strings.HasPrefix(lower, "git pull"), strings.HasPrefix(lower, "git fetch"):
		values = state.Remotes
		fields := strings.Fields(input)
		if len(fields) >= 2 {
			prefix = fields[1:]
		}
	default:
		return nil
	}
	if len(prefix) == 0 {
		return nil
	}
	base := prefix[:1]
	out := make([]sdk.Suggestion, 0, len(values))
	for _, value := range values {
		out = append(out, makeSuggestion(append(base, value), "Use "+value, sdk.Completion, risk, 95, "Discovered from the current repository"))
	}
	return out
}
func makeSuggestion(args []string, title string, kind sdk.SuggestionKind, risk sdk.Risk, priority int, reason string) sdk.Suggestion {
	copied := append([]string(nil), args...)
	placeholders := []sdk.Placeholder{}
	for i, arg := range copied {
		if strings.HasPrefix(arg, "<") {
			placeholders = append(placeholders, sdk.Placeholder{Name: strings.Trim(arg, "<>"), ArgIndex: i})
		}
	}
	return sdk.Suggestion{Command: sdk.Command{Executable: "git", Args: copied}, Title: title, Description: title, Reason: reason, Kind: kind, Risk: risk, Priority: priority, Source: "git", Placeholders: placeholders}
}
