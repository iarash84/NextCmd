package pip

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	executable := strings.ToLower(input.Result.Command.Executable)
	if executable != "pip" && executable != "pip3" || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	switch input.Result.Command.Args[0] {
	case "install", "uninstall":
		return []sdk.Suggestion{suggest(executable, []string{"check"}, "Verify dependency compatibility", sdk.NextAction, sdk.Safe, 90, "The Python environment changed"), suggest(executable, []string{"list"}, "Review installed packages", sdk.NextAction, sdk.Safe, 70, "The Python environment changed")}, nil
	case "check":
		return []sdk.Suggestion{suggest(executable, []string{"list", "--outdated"}, "Review outdated packages", sdk.NextAction, sdk.Safe, 62, "Installed dependencies are compatible")}, nil
	}
	return nil, nil
}
func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if state.Root == "" {
		return nil, nil
	}
	executable := inputExecutable(input.Input)
	return []sdk.Suggestion{suggest(executable, []string{"check"}, "Verify dependency compatibility", sdk.BestPractice, sdk.Safe, 68, "Python project metadata was detected"), suggest(executable, []string{"list", "--outdated"}, "Review outdated dependencies", sdk.BestPractice, sdk.Safe, 58, "Keep dependency updates visible")}, nil
}
func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	executable := strings.ToLower(input.Result.Command.Executable)
	if executable != "pip" && executable != "pip3" {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	switch {
	case strings.Contains(message, "externally-managed-environment"):
		return []sdk.Suggestion{suggest(executable, []string{"--version"}, "Confirm the active pip environment", sdk.Recovery, sdk.Safe, 86, "The system Python environment is externally managed")}, nil
	case strings.Contains(message, "no matching distribution found"), strings.Contains(message, "could not find a version that satisfies"):
		return []sdk.Suggestion{suggest(executable, []string{"index", "versions", "<package>"}, "Inspect available package versions", sdk.Recovery, sdk.Safe, 92, "No compatible distribution was found")}, nil
	case strings.Contains(message, "permission denied"), strings.Contains(message, "access is denied"):
		return []sdk.Suggestion{suggest(executable, []string{"install", "--user", "<package>"}, "Install into the user package directory", sdk.Recovery, sdk.Mutating, 78, "The global environment was not writable")}, nil
	case strings.Contains(message, "certificate verify failed"):
		return []sdk.Suggestion{suggest(executable, []string{"config", "debug"}, "Inspect pip index and certificate configuration", sdk.Recovery, sdk.Safe, 90, "TLS certificate verification failed")}, nil
	}
	return nil, nil
}
func inputExecutable(input string) string {
	fields := strings.Fields(strings.ToLower(input))
	if len(fields) > 0 && fields[0] == "pip3" {
		return "pip3"
	}
	return "pip"
}
