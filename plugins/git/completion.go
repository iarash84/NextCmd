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
	{[]string{"status", "--short", "--branch"}, "Show concise branch and file status", sdk.Safe, 95},
	{[]string{"status"}, "Show working tree status", sdk.Safe, 90},
	{[]string{"add", "."}, "Stage all changes", sdk.Mutating, 70},
	{[]string{"add", "-p"}, "Interactively stage selected changes", sdk.Mutating, 78},
	{[]string{"commit", "-m", "<message>"}, "Commit staged changes", sdk.Mutating, 70},
	{[]string{"commit", "--amend"}, "Amend the latest commit", sdk.Destructive, 35},
	{[]string{"diff"}, "Review unstaged changes", sdk.Safe, 80},
	{[]string{"diff", "--cached"}, "Review staged changes", sdk.Safe, 85},
	{[]string{"diff", "HEAD~1", "HEAD"}, "Compare the latest commit", sdk.Safe, 55},
	{[]string{"log", "--oneline", "--graph", "--decorate", "--all"}, "Show branch graph", sdk.Safe, 72},
	{[]string{"log", "--oneline"}, "Show concise history", sdk.Safe, 65},
	{[]string{"show", "--stat", "HEAD"}, "Show the latest commit summary", sdk.Safe, 62},
	{[]string{"reflog"}, "Show local reference history", sdk.Safe, 48},
	{[]string{"blame", "<file>"}, "Show line-by-line authorship", sdk.Safe, 42},
	{[]string{"grep", "<pattern>"}, "Search tracked files", sdk.Safe, 48},
	{[]string{"branch"}, "List local branches", sdk.Safe, 60},
	{[]string{"branch", "-a"}, "List local and remote branches", sdk.Safe, 58},
	{[]string{"switch", "<branch>"}, "Switch branch", sdk.Mutating, 55},
	{[]string{"switch", "-c", "feature/<name>"}, "Create a feature branch", sdk.Mutating, 82},
	{[]string{"switch", "-c", "bugfix/<name>"}, "Create a bug-fix branch", sdk.Mutating, 82},
	{[]string{"switch", "-c", "hotfix/<name>"}, "Create a hot-fix branch", sdk.Mutating, 76},
	{[]string{"switch", "-c", "release/<version>"}, "Create a release branch", sdk.Mutating, 68},
	{[]string{"switch", "-c", "chore/<name>"}, "Create a maintenance branch", sdk.Mutating, 66},
	{[]string{"switch", "-c", "docs/<name>"}, "Create a documentation branch", sdk.Mutating, 62},
	{[]string{"switch", "-c", "refactor/<name>"}, "Create a refactoring branch", sdk.Mutating, 62},
	{[]string{"checkout", "<branch>"}, "Checkout branch", sdk.Mutating, 45},
	{[]string{"restore", "<file>"}, "Restore a file", sdk.Destructive, 40},
	{[]string{"restore", "--staged", "<file>"}, "Unstage a file", sdk.Mutating, 58},
	{[]string{"stash", "push", "-m", "<message>"}, "Stash changes with a message", sdk.Mutating, 52},
	{[]string{"stash", "list"}, "List stashes", sdk.Safe, 50},
	{[]string{"stash", "pop"}, "Apply and remove the latest stash", sdk.Destructive, 38},
	{[]string{"pull", "--rebase"}, "Pull and rebase local commits", sdk.Mutating, 62},
	{[]string{"pull"}, "Pull upstream changes", sdk.Mutating, 55},
	{[]string{"push"}, "Push commits", sdk.Mutating, 55},
	{[]string{"push", "--force-with-lease"}, "Safely force-update the remote branch", sdk.Destructive, 22},
	{[]string{"fetch", "--all", "--prune"}, "Fetch all remotes and prune stale refs", sdk.Safe, 60},
	{[]string{"fetch"}, "Fetch remotes", sdk.Safe, 50},
	{[]string{"remote", "-v"}, "List remotes", sdk.Safe, 45},
	{[]string{"remote", "add", "<name>", "<url>"}, "Add a remote", sdk.Mutating, 35},
	{[]string{"tag", "--list"}, "List tags", sdk.Safe, 45},
	{[]string{"tag", "-a", "<version>", "-m", "<message>"}, "Create an annotated tag", sdk.Mutating, 42},
	{[]string{"merge", "<branch>"}, "Merge a branch", sdk.Mutating, 35},
	{[]string{"merge", "--abort"}, "Abort the current merge", sdk.Mutating, 48},
	{[]string{"rebase", "<branch>"}, "Rebase onto a branch", sdk.Destructive, 25},
	{[]string{"rebase", "--abort"}, "Abort the current rebase", sdk.Mutating, 48},
	{[]string{"cherry-pick", "<commit>"}, "Apply a specific commit", sdk.Mutating, 38},
	{[]string{"revert", "<commit>"}, "Revert a commit safely", sdk.Mutating, 44},
	{[]string{"reset", "--soft", "HEAD~1"}, "Undo the latest commit but keep changes staged", sdk.Destructive, 28},
	{[]string{"clean", "-n"}, "Preview untracked files that would be removed", sdk.Safe, 52},
	{[]string{"worktree", "list"}, "List linked worktrees", sdk.Safe, 38},
	{[]string{"submodule", "status"}, "Show submodule status", sdk.Safe, 35},
	{[]string{"config", "--list", "--show-origin"}, "List Git configuration sources", sdk.Safe, 32},
	{[]string{"init"}, "Initialize repository", sdk.Mutating, 75},
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
	fields := strings.Fields(strings.ToLower(input))
	if len(fields) < 2 || fields[0] != "git" {
		return nil
	}
	base, values, risk := dynamicValues(fields, state)
	if len(base) == 0 || len(values) == 0 {
		return nil
	}
	out := make([]sdk.Suggestion, 0, len(values))
	for _, value := range values {
		out = append(out, makeSuggestion(append(base, value), "Use "+value, sdk.Completion, risk, 95, "Discovered from the current repository"))
	}
	return out
}

func dynamicValues(fields []string, state State) ([]string, []string, sdk.Risk) {
	switch fields[1] {
	case "switch", "checkout", "merge", "rebase":
		if len(fields) >= 3 && strings.HasPrefix(fields[2], "-") {
			return nil, nil, sdk.Safe
		}
		return []string{fields[1]}, state.Branches, sdk.Mutating
	case "branch":
		if len(fields) >= 3 && (fields[2] == "-d" || fields[2] == "-D") {
			return []string{"branch", fields[2]}, state.Branches, sdk.Destructive
		}
	case "add":
		return []string{"add"}, state.Files, sdk.Mutating
	case "restore":
		if len(fields) >= 3 && fields[2] == "--staged" {
			return []string{"restore", "--staged"}, state.Files, sdk.Mutating
		}
		return []string{"restore"}, state.Files, sdk.Destructive
	case "blame":
		return []string{"blame"}, state.Files, sdk.Safe
	case "push", "pull", "fetch":
		risk := sdk.Mutating
		if fields[1] == "fetch" {
			risk = sdk.Safe
		}
		return []string{fields[1]}, state.Remotes, risk
	}
	return nil, nil, sdk.Safe
}
func makeSuggestion(args []string, title string, kind sdk.SuggestionKind, risk sdk.Risk, priority int, reason string) sdk.Suggestion {
	copied := append([]string(nil), args...)
	placeholders := []sdk.Placeholder{}
	for i, arg := range copied {
		start := strings.IndexByte(arg, '<')
		end := strings.IndexByte(arg, '>')
		if start >= 0 && end > start {
			placeholders = append(placeholders, sdk.Placeholder{Name: arg[start+1 : end], ArgIndex: i, Start: start, End: end + 1})
		}
	}
	return sdk.Suggestion{Command: sdk.Command{Executable: "git", Args: copied}, Title: title, Description: title, Reason: reason, Kind: kind, Risk: risk, Priority: priority, Source: "git", Placeholders: placeholders}
}
