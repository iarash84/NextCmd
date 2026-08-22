package dotnet

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "dotnet") || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	state, _ := input.Project.(State)
	out := []sdk.Suggestion{}
	switch input.Result.Command.Args[0] {
	case "new":
		out = append(out, suggestion([]string{"restore"}, "Restore generated dependencies", sdk.NextAction, sdk.Mutating, 86, "A project or solution was created"), suggestion([]string{"build"}, "Build the generated project", sdk.NextAction, sdk.Mutating, 80, "Verify the generated project"))
	case "restore":
		out = append(out, suggestion([]string{"build", "--no-restore"}, "Build without restoring again", sdk.NextAction, sdk.Mutating, 92, "Dependencies were restored"))
	case "build":
		if hasTests(state) {
			out = append(out, suggestion([]string{"test", "--no-build"}, "Run tests without rebuilding", sdk.NextAction, sdk.Mutating, 94, "The solution built successfully"))
		}
		out = append(out, suggestion([]string{"run", "--no-build"}, "Run without rebuilding", sdk.NextAction, sdk.Mutating, 76, "The application built successfully"))
	case "test":
		out = append(out, suggestion([]string{"publish", "-c", "Release"}, "Publish a release build", sdk.NextAction, sdk.Mutating, 72, "Tests completed successfully"), suggestion([]string{"format", "--verify-no-changes"}, "Verify formatting", sdk.NextAction, sdk.Safe, 78, "Check source formatting before release"))
	case "package", "add":
		out = append(out, suggestion([]string{"restore"}, "Restore changed dependencies", sdk.NextAction, sdk.Mutating, 86, "Project dependencies changed"), suggestion([]string{"build"}, "Verify dependency compatibility", sdk.NextAction, sdk.Mutating, 80, "Build after changing dependencies"))
	case "ef":
		out = append(out, suggestion([]string{"ef", "migrations", "list"}, "Inspect EF Core migrations", sdk.NextAction, sdk.Safe, 72, "Review migration state"), suggestion([]string{"test"}, "Run tests after schema changes", sdk.NextAction, sdk.Mutating, 78, "Verify schema-dependent behavior"))
	}
	return out, nil
}

func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if len(state.Projects) == 0 {
		return nil, nil
	}
	out := []sdk.Suggestion{suggestion([]string{"format", "--verify-no-changes"}, "Verify .NET formatting", sdk.BestPractice, sdk.Safe, 58, "Keep formatting deterministic in CI")}
	if hasTests(state) {
		out = append(out, suggestion([]string{"test"}, "Run the test projects", sdk.BestPractice, sdk.Mutating, 66, "Test projects were detected"))
	}
	return out, nil
}

func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "dotnet") {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	state, _ := input.Project.(State)
	if strings.Contains(message, "no executable found matching command") || strings.Contains(message, "could not execute because") {
		return []sdk.Suggestion{suggestion([]string{"--info"}, "Check the installed .NET SDK", sdk.Recovery, sdk.Safe, 95, "The requested SDK command or tool may be unavailable"), suggestion([]string{"tool", "restore"}, "Restore local .NET tools", sdk.Recovery, sdk.Mutating, 82, "A local tool may be missing")}, nil
	}
	if len(state.Projects) == 0 {
		return []sdk.Suggestion{suggestion([]string{"new", "sln", "-n", "<name>"}, "Create a solution", sdk.Recovery, sdk.Mutating, 85, "No .NET project or solution was detected"), suggestion([]string{"new", "console", "-n", "<name>"}, "Create a console project", sdk.Recovery, sdk.Mutating, 80, "Start a new .NET project")}, nil
	}
	if strings.Contains(message, "assets file") || strings.Contains(message, "restore") {
		return []sdk.Suggestion{suggestion([]string{"restore"}, "Restore project dependencies", sdk.Recovery, sdk.Mutating, 94, "Build assets or packages appear to be missing")}, nil
	}
	return nil, nil
}
