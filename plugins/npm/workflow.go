package npm

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "npm") || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	switch input.Result.Command.Args[0] {
	case "install", "ci", "update", "uninstall":
		return []sdk.Suggestion{suggest([]string{"test"}, "Test after dependency changes", sdk.NextAction, sdk.Mutating, 90, "Dependencies were installed or changed"), suggest([]string{"audit"}, "Audit installed dependencies", sdk.NextAction, sdk.Safe, 84, "Dependencies were installed or changed")}, nil
	case "test":
		return []sdk.Suggestion{suggest([]string{"run", "build"}, "Build the tested package", sdk.NextAction, sdk.Mutating, 76, "Tests completed successfully")}, nil
	case "audit":
		return []sdk.Suggestion{suggest([]string{"outdated"}, "Review outdated dependencies", sdk.NextAction, sdk.Safe, 64, "Dependency audit completed")}, nil
	}
	return nil, nil
}
func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if state.Root == "" {
		return nil, nil
	}
	out := []sdk.Suggestion{suggest([]string{"test"}, "Run package tests", sdk.BestPractice, sdk.Mutating, 66, "A package.json was detected"), suggest([]string{"audit"}, "Audit dependencies", sdk.BestPractice, sdk.Safe, 62, "Review known dependency vulnerabilities")}
	if state.HasLock {
		out = append(out, suggest([]string{"ci"}, "Use the lockfile for a reproducible install", sdk.BestPractice, sdk.Destructive, 58, "package-lock.json was detected"))
	}
	return out, nil
}
func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "npm") {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	state, _ := input.Project.(State)
	switch {
	case state.Root == "" || strings.Contains(message, "could not read package.json"), strings.Contains(message, "enoent") && strings.Contains(message, "package.json"):
		return []sdk.Suggestion{suggest([]string{"init", "-y"}, "Initialize package.json", sdk.Recovery, sdk.Mutating, 96, "No package.json was found")}, nil
	case strings.Contains(message, "missing script"):
		out := []sdk.Suggestion{}
		for _, script := range state.Scripts {
			out = append(out, suggest([]string{"run", script}, "Run script "+script, sdk.Recovery, sdk.Mutating, 94, "Use a script declared in package.json"))
		}
		return out, nil
	case strings.Contains(message, "npm ci") && strings.Contains(message, "package-lock"):
		return []sdk.Suggestion{suggest([]string{"install"}, "Create or update package-lock.json", sdk.Recovery, sdk.Mutating, 90, "npm ci requires a compatible lockfile")}, nil
	}
	return nil, nil
}
