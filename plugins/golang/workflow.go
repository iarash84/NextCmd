package golang

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "go") || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	args := input.Result.Command.Args
	switch args[0] {
	case "build":
		return []sdk.Suggestion{
			suggestion([]string{"test", "./..."}, "Test all packages after the build", sdk.NextAction, sdk.Mutating, 94, "The Go build completed successfully"),
			suggestion([]string{"run", "."}, "Run the current main package", sdk.NextAction, sdk.Mutating, 72, "The Go build completed successfully"),
		}, nil
	case "test":
		return []sdk.Suggestion{
			suggestion([]string{"vet", "./..."}, "Analyze packages after tests", sdk.NextAction, sdk.Safe, 84, "The test suite completed successfully"),
			suggestion([]string{"build", "./..."}, "Build all tested packages", sdk.NextAction, sdk.Mutating, 78, "The test suite completed successfully"),
		}, nil
	case "fmt":
		return []sdk.Suggestion{
			suggestion([]string{"vet", "./..."}, "Analyze formatted packages", sdk.NextAction, sdk.Safe, 88, "Formatting completed successfully"),
			suggestion([]string{"test", "./..."}, "Test formatted packages", sdk.NextAction, sdk.Mutating, 86, "Formatting completed successfully"),
		}, nil
	case "vet":
		return []sdk.Suggestion{suggestion([]string{"test", "./..."}, "Run all tests", sdk.NextAction, sdk.Mutating, 90, "Static analysis completed successfully")}, nil
	case "generate":
		return []sdk.Suggestion{
			suggestion([]string{"fmt", "./..."}, "Format generated source", sdk.NextAction, sdk.Mutating, 90, "Source generation completed successfully"),
			suggestion([]string{"test", "./..."}, "Test generated source", sdk.NextAction, sdk.Mutating, 84, "Source generation completed successfully"),
		}, nil
	case "get":
		return dependencyNextActions("A module dependency changed"), nil
	case "mod":
		if len(args) > 1 && args[1] == "init" {
			return []sdk.Suggestion{
				suggestion([]string{"mod", "tidy"}, "Resolve module dependencies", sdk.NextAction, sdk.Mutating, 92, "A Go module was initialized"),
				suggestion([]string{"test", "./..."}, "Test the new module", sdk.NextAction, sdk.Mutating, 82, "A Go module was initialized"),
			}, nil
		}
		if len(args) > 1 && args[1] == "tidy" {
			return dependencyNextActions("Module files were synchronized"), nil
		}
	}
	return nil, nil
}

func dependencyNextActions(reason string) []sdk.Suggestion {
	return []sdk.Suggestion{
		suggestion([]string{"mod", "verify"}, "Verify downloaded modules", sdk.NextAction, sdk.Safe, 88, reason),
		suggestion([]string{"test", "./..."}, "Test after dependency changes", sdk.NextAction, sdk.Mutating, 90, reason),
	}
}

func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if state.Root == "" {
		return nil, nil
	}
	return []sdk.Suggestion{
		suggestion([]string{"fmt", "./..."}, "Format all Go packages", sdk.BestPractice, sdk.Mutating, 62, "Keep Go source formatting deterministic"),
		suggestion([]string{"vet", "./..."}, "Analyze all Go packages", sdk.BestPractice, sdk.Safe, 68, "Catch suspicious constructs before committing"),
		suggestion([]string{"test", "./..."}, "Run the complete Go test suite", sdk.BestPractice, sdk.Mutating, 72, "A Go module or workspace was detected"),
		suggestion([]string{"mod", "tidy"}, "Keep module files synchronized", sdk.BestPractice, sdk.Mutating, 58, "Remove unused requirements and add missing ones"),
	}, nil
}

func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "go") {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	state, _ := input.Project.(State)
	if state.Root == "" || strings.Contains(message, "go.mod file not found") || strings.Contains(message, "cannot find main module") {
		return []sdk.Suggestion{suggestion([]string{"mod", "init", "<module>"}, "Initialize a Go module", sdk.Recovery, sdk.Mutating, 98, "No go.mod or go.work file was found")}, nil
	}
	if strings.Contains(message, "updates to go.mod needed") || strings.Contains(message, "go mod tidy") {
		return []sdk.Suggestion{suggestion([]string{"mod", "tidy"}, "Synchronize module dependencies", sdk.Recovery, sdk.Mutating, 96, "The Go tool reported that go.mod needs updates")}, nil
	}
	if strings.Contains(message, "no required module provides package") || strings.Contains(message, "is not in std") {
		return []sdk.Suggestion{suggestion([]string{"get", "<module>"}, "Add the missing module dependency", sdk.Recovery, sdk.Mutating, 90, "A required package could not be resolved")}, nil
	}
	if strings.Contains(message, "build constraints exclude all go files") {
		return []sdk.Suggestion{suggestion([]string{"env", "GOOS", "GOARCH", "CGO_ENABLED"}, "Inspect build environment", sdk.Recovery, sdk.Safe, 88, "Build constraints excluded the available Go files")}, nil
	}
	return nil, nil
}
