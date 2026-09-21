package dart

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	executable := strings.ToLower(input.Result.Command.Executable)
	if executable != "dart" && executable != "flutter" || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	switch input.Result.Command.Args[0] {
	case "pub":
		if len(input.Result.Command.Args) > 1 && (input.Result.Command.Args[1] == "get" || input.Result.Command.Args[1] == "upgrade") {
			return []sdk.Suggestion{suggest(executable, []string{"analyze"}, "Analyze after dependency changes", sdk.Safe, 86, "Dart dependencies changed"), suggest(executable, []string{"test"}, "Run tests after dependency changes", sdk.Mutating, 82, "Dart dependencies changed")}, nil
		}
	case "format":
		return []sdk.Suggestion{suggest(executable, []string{"analyze"}, "Analyze formatted Dart source", sdk.Safe, 88, "Formatting completed successfully")}, nil
	case "analyze":
		return []sdk.Suggestion{suggest(executable, []string{"test"}, "Run the Dart test suite", sdk.Mutating, 86, "Static analysis completed successfully")}, nil
	case "test":
		return []sdk.Suggestion{suggest(executable, []string{"format", "."}, "Format tested Dart source", sdk.Mutating, 72, "Tests completed successfully")}, nil
	}
	return nil, nil
}

func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if state.Root == "" {
		return nil, nil
	}
	executable := "dart"
	if state.IsFlutter {
		executable = "flutter"
	}
	return []sdk.Suggestion{suggest(executable, []string{"format", "."}, "Format Dart source", sdk.Mutating, 68, "Keep Dart formatting deterministic"), suggest(executable, []string{"analyze"}, "Analyze Dart source", sdk.Safe, 72, "Catch issues before committing"), suggest(executable, []string{"test"}, "Run Dart tests", sdk.Mutating, 76, "Verify project behavior")}, nil
}

func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	executable := strings.ToLower(input.Result.Command.Executable)
	if executable != "dart" && executable != "flutter" {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	if strings.Contains(message, "pubspec.yaml") || strings.Contains(message, "could not find a file named") {
		return []sdk.Suggestion{suggest(executable, []string{"pub", "get"}, "Resolve Dart dependencies", sdk.Mutating, 94, "Project metadata or dependencies are unavailable")}, nil
	}
	if strings.Contains(message, "version solving failed") || strings.Contains(message, "because") && strings.Contains(message, "depends on") {
		return []sdk.Suggestion{suggest(executable, []string{"pub", "outdated"}, "Inspect dependency constraints", sdk.Safe, 90, "Pub dependency resolution failed")}, nil
	}
	return nil, nil
}
