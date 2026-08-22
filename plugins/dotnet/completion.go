package dotnet

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
	{[]string{"--info"}, "Show .NET SDK information", sdk.Safe, 45},
	{[]string{"sdk", "check"}, "Check installed SDK and runtime status", sdk.Safe, 48},
	{[]string{"new", "list"}, "List installed project templates", sdk.Safe, 55},
	{[]string{"new", "console", "-n", "<name>"}, "Create a console application", sdk.Mutating, 72},
	{[]string{"new", "webapi", "-n", "<name>"}, "Create an ASP.NET Core Web API", sdk.Mutating, 72},
	{[]string{"new", "webapp", "-n", "<name>"}, "Create an ASP.NET Core web application", sdk.Mutating, 65},
	{[]string{"new", "classlib", "-n", "<name>"}, "Create a class library", sdk.Mutating, 66},
	{[]string{"new", "xunit", "-n", "<name>.Tests"}, "Create an xUnit test project", sdk.Mutating, 68},
	{[]string{"new", "sln", "-n", "<name>"}, "Create a solution", sdk.Mutating, 70},
	{[]string{"restore"}, "Restore project dependencies", sdk.Mutating, 82},
	{[]string{"build"}, "Build the project or solution", sdk.Mutating, 90},
	{[]string{"build", "--no-restore"}, "Build without restoring again", sdk.Mutating, 72},
	{[]string{"run"}, "Run the application", sdk.Mutating, 86},
	{[]string{"watch", "run"}, "Run with hot reload", sdk.Mutating, 76},
	{[]string{"test"}, "Run all tests", sdk.Mutating, 88},
	{[]string{"test", "--no-build"}, "Run tests without rebuilding", sdk.Mutating, 70},
	{[]string{"clean"}, "Clean build outputs", sdk.Destructive, 45},
	{[]string{"publish", "-c", "Release"}, "Publish a release build", sdk.Mutating, 62},
	{[]string{"pack", "-c", "Release"}, "Create a NuGet package", sdk.Mutating, 55},
	{[]string{"format"}, "Format and analyze the solution", sdk.Mutating, 64},
	{[]string{"sln", "list"}, "List projects in the solution", sdk.Safe, 58},
	{[]string{"sln", "<solution>", "add", "<project>"}, "Add a project to the solution", sdk.Mutating, 60},
	{[]string{"add", "<project>", "package", "<package>"}, "Add a NuGet package (compatible form)", sdk.Mutating, 58},
	{[]string{"package", "add", "<package>", "--project", "<project>"}, "Add a NuGet package (.NET 10+)", sdk.Mutating, 52},
	{[]string{"list", "<project>", "package"}, "List project packages", sdk.Safe, 52},
	{[]string{"add", "<project>", "reference", "<reference>"}, "Add a project reference", sdk.Mutating, 56},
	{[]string{"list", "<project>", "reference"}, "List project references", sdk.Safe, 48},
	{[]string{"tool", "restore"}, "Restore local .NET tools", sdk.Mutating, 48},
	{[]string{"tool", "list", "--local"}, "List local .NET tools", sdk.Safe, 42},
	{[]string{"ef", "migrations", "add", "<name>"}, "Create an EF Core migration", sdk.Mutating, 48},
	{[]string{"ef", "migrations", "list"}, "List EF Core migrations", sdk.Safe, 44},
	{[]string{"ef", "database", "update"}, "Apply EF Core migrations", sdk.Destructive, 34},
}

func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	trimmed := strings.TrimSpace(input.Input)
	if trimmed != "" && !strings.HasPrefix(strings.ToLower(trimmed), "dotnet") {
		return nil, nil
	}
	state, _ := input.Project.(State)
	out := make([]sdk.Suggestion, 0, len(commands)+len(state.Projects))
	for _, spec := range commands {
		projectRequired := requiresProject(spec.args)
		if projectRequired && len(state.Projects) == 0 {
			continue
		}
		priority := spec.priority
		if len(state.Projects) > 0 && (spec.args[0] == "build" || spec.args[0] == "restore") {
			priority += 15
		}
		if hasTests(state) && spec.args[0] == "test" {
			priority += 18
		}
		out = append(out, suggestion(spec.args, spec.title, sdk.Completion, spec.risk, priority, "Matches the current input and .NET workspace"))
	}
	return append(out, dynamic(input.Input, state)...), nil
}

func dynamic(input string, state State) []sdk.Suggestion {
	fields := strings.Fields(strings.ToLower(input))
	if len(fields) < 2 || fields[0] != "dotnet" {
		return nil
	}
	verb := fields[1]
	base := []string{verb}
	projects := state.Projects
	risk := sdk.Mutating
	switch verb {
	case "build", "clean", "publish", "pack", "restore":
	case "run":
		base = []string{"run", "--project"}
	case "test":
		projects = testProjects(state)
	default:
		return nil
	}
	out := make([]sdk.Suggestion, 0, len(projects))
	for _, project := range projects {
		out = append(out, suggestion(append(base, project.Path), verb+" "+project.Name, sdk.Completion, risk, 96, "Discovered from the current .NET workspace"))
	}
	return out
}

func requiresProject(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "restore", "build", "run", "watch", "test", "clean", "publish", "pack", "format", "sln", "add", "package", "list", "ef":
		return true
	default:
		return false
	}
}
func hasTests(state State) bool { return len(testProjects(state)) > 0 }
func testProjects(state State) []Project {
	out := []Project{}
	for _, project := range state.Projects {
		if project.Test {
			out = append(out, project)
		}
	}
	return out
}
func suggestion(args []string, title string, kind sdk.SuggestionKind, risk sdk.Risk, priority int, reason string) sdk.Suggestion {
	copied := append([]string(nil), args...)
	placeholders := []sdk.Placeholder{}
	for i, arg := range copied {
		start, end := strings.IndexByte(arg, '<'), strings.IndexByte(arg, '>')
		if start >= 0 && end > start {
			placeholders = append(placeholders, sdk.Placeholder{Name: arg[start+1 : end], ArgIndex: i, Start: start, End: end + 1})
		}
	}
	return sdk.Suggestion{Command: sdk.Command{Executable: "dotnet", Args: copied}, Title: title, Description: title, Reason: reason, Kind: kind, Priority: priority, Risk: risk, Source: "dotnet", Placeholders: placeholders}
}
