package git

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

func (p *Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if strings.ToLower(input.Result.Command.Executable) != "git" || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	state, _ := input.Project.(State)
	var out []sdk.Suggestion
	switch input.Result.Command.Args[0] {
	case "add":
		out = append(out, makeSuggestion([]string{"diff", "--cached"}, "Review staged changes", sdk.NextAction, sdk.Safe, 100, "Review what will be committed"), makeSuggestion([]string{"commit", "-m", "<message>"}, "Commit staged changes", sdk.NextAction, sdk.Mutating, 85, "Changes were staged"), makeSuggestion([]string{"status"}, "Check status", sdk.NextAction, sdk.Safe, 70, "Confirm repository state"))
	case "commit":
		out = append(out, makeSuggestion([]string{"status"}, "Check status", sdk.NextAction, sdk.Safe, 80, "Confirm a clean working tree"), makeSuggestion([]string{"log", "-1"}, "Inspect the commit", sdk.NextAction, sdk.Safe, 75, "Verify the new commit"))
		if state.HasUpstream {
			out = append(out, makeSuggestion([]string{"push"}, "Push the commit", sdk.NextAction, sdk.Mutating, 95, "Publish the new commit"))
		} else if state.Branch != "" && len(state.Remotes) > 0 {
			out = append(out, makeSuggestion([]string{"push", "-u", state.Remotes[0], state.Branch}, "Set upstream and push", sdk.NextAction, sdk.Mutating, 95, "No upstream is configured"))
		}
	case "pull":
		out = append(out, makeSuggestion([]string{"status"}, "Check merged state", sdk.NextAction, sdk.Safe, 80, "Inspect repository after pulling"))
		if state.Modified {
			out = append(out, makeSuggestion([]string{"diff"}, "Review changes", sdk.NextAction, sdk.Safe, 75, "Working tree contains changes"))
		}
	case "status":
		if state.Modified || state.Untracked {
			out = append(out, makeSuggestion([]string{"diff"}, "Review changes", sdk.NextAction, sdk.Safe, 90, "Working tree contains changes"), makeSuggestion([]string{"add", "."}, "Stage changes", sdk.NextAction, sdk.Mutating, 75, "Prepare changes for commit"))
		}
		if state.Staged {
			out = append(out, makeSuggestion([]string{"diff", "--cached"}, "Review staged changes", sdk.NextAction, sdk.Safe, 100, "Review what will be committed"))
		}
	}
	return out, nil
}
func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if !state.Staged {
		return nil, nil
	}
	return []sdk.Suggestion{makeSuggestion([]string{"diff", "--cached"}, "Review staged changes", sdk.BestPractice, sdk.Safe, 98, "Recommended before committing")}, nil
}
func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if strings.ToLower(input.Result.Command.Executable) != "git" {
		return nil, nil
	}
	state, _ := input.Project.(State)
	if !state.InRepository {
		return []sdk.Suggestion{makeSuggestion([]string{"init"}, "Initialize a repository", sdk.Recovery, sdk.Mutating, 100, "The previous Git command requires a repository")}, nil
	}
	message := strings.ToLower(input.Result.Stderr + input.Result.Stdout)
	if strings.Contains(message, "pathspec") || strings.Contains(message, "unknown revision") {
		out := []sdk.Suggestion{}
		for _, branch := range state.Branches {
			out = append(out, makeSuggestion([]string{"switch", branch}, "Switch to "+branch, sdk.Recovery, sdk.Mutating, 85, "Choose an existing branch"))
		}
		return out, nil
	}
	return nil, nil
}
