package golang

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
	{[]string{"version"}, "Show the installed Go version", sdk.Safe, 45},
	{[]string{"env"}, "Show Go environment settings", sdk.Safe, 42},
	{[]string{"env", "GOMOD"}, "Show the active module file", sdk.Safe, 48},
	{[]string{"mod", "init", "<module>"}, "Initialize a Go module", sdk.Mutating, 72},
	{[]string{"mod", "tidy"}, "Synchronize module dependencies", sdk.Mutating, 88},
	{[]string{"mod", "download"}, "Download module dependencies", sdk.Mutating, 58},
	{[]string{"mod", "verify"}, "Verify downloaded module content", sdk.Safe, 70},
	{[]string{"mod", "graph"}, "Print the module dependency graph", sdk.Safe, 48},
	{[]string{"mod", "vendor"}, "Create a vendor directory", sdk.Mutating, 46},
	{[]string{"work", "init"}, "Initialize a Go workspace", sdk.Mutating, 48},
	{[]string{"work", "use", "<directory>"}, "Add a module to the workspace", sdk.Mutating, 44},
	{[]string{"work", "sync"}, "Synchronize workspace modules", sdk.Mutating, 52},
	{[]string{"build", "./..."}, "Build every package", sdk.Mutating, 90},
	{[]string{"run", "."}, "Run the current main package", sdk.Mutating, 82},
	{[]string{"test", "./..."}, "Run all tests", sdk.Mutating, 94},
	{[]string{"test", "-race", "./..."}, "Run all tests with race detection", sdk.Mutating, 76},
	{[]string{"test", "-cover", "./..."}, "Run all tests with coverage", sdk.Mutating, 74},
	{[]string{"vet", "./..."}, "Analyze all packages", sdk.Safe, 86},
	{[]string{"fmt", "./..."}, "Format all Go packages", sdk.Mutating, 84},
	{[]string{"generate", "./..."}, "Run source generators", sdk.Mutating, 56},
	{[]string{"list", "./..."}, "List all packages", sdk.Safe, 60},
	{[]string{"doc", "<symbol>"}, "Show documentation for a package or symbol", sdk.Safe, 42},
	{[]string{"get", "<module>"}, "Add or update a module dependency", sdk.Mutating, 52},
	{[]string{"install", "<package>@latest"}, "Install a Go command", sdk.Mutating, 38},
	{[]string{"clean", "-cache"}, "Remove the Go build cache", sdk.Destructive, 24},
}

func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	trimmed := strings.TrimSpace(input.Input)
	if !matchesExecutable(trimmed, "go") {
		return nil, nil
	}
	state, _ := input.Project.(State)
	out := make([]sdk.Suggestion, 0, len(commands)+len(state.Packages)+len(state.GoFiles))
	for _, spec := range commands {
		priority := spec.priority
		reason := "Matches the current Go command input"
		if requiresProject(spec.args) && state.Root == "" {
			priority -= 18
			reason = "No go.mod or go.work file was detected"
		}
		if state.Workspace && len(spec.args) > 0 && (spec.args[0] == "build" || spec.args[0] == "test" || spec.args[0] == "list") {
			priority += 8
			reason = "A Go workspace was detected"
		}
		out = append(out, suggestion(spec.args, spec.title, sdk.Completion, spec.risk, priority, reason))
	}
	return append(out, dynamic(input.Input, state)...), nil
}

func (*Plugin) Help() []sdk.CommandHelp {
	out := make([]sdk.CommandHelp, 0, len(commands))
	for _, spec := range commands {
		out = append(out, sdk.CommandHelp{Command: sdk.Command{Executable: "go", Args: append([]string(nil), spec.args...)}, Description: spec.title, Risk: spec.risk})
	}
	return out
}

func matchesExecutable(input, executable string) bool {
	if input == "" {
		return true
	}
	first := strings.ToLower(strings.Fields(input)[0])
	return strings.HasPrefix(executable, first) || first == executable
}

func dynamic(input string, state State) []sdk.Suggestion {
	fields := strings.Fields(input)
	if len(fields) < 2 || !strings.EqualFold(fields[0], "go") {
		return nil
	}
	verb := strings.ToLower(fields[1])
	out := []sdk.Suggestion{}
	if supportsPackage(verb) {
		for _, pkg := range state.Packages {
			out = append(out, suggestion(dynamicArgs(fields, input, pkg), verb+" package "+pkg, sdk.Completion, dynamicRisk(verb), 98, "Package discovered in the current Go project"))
		}
		if len(state.Packages) > 1 {
			out = append(out, suggestion(dynamicArgs(fields, input, "./..."), verb+" all packages", sdk.Completion, dynamicRisk(verb), 99, "Multiple Go packages were discovered"))
		}
	}
	if verb == "run" {
		for _, file := range state.GoFiles {
			if !strings.HasSuffix(strings.ToLower(file), "_test.go") {
				out = append(out, suggestion(dynamicArgs(fields, input, file), "Run "+file, sdk.Completion, sdk.Mutating, 98, "Go source file discovered in the current project"))
			}
		}
	}
	return out
}

func dynamicArgs(fields []string, input, value string) []string {
	args := append([]string(nil), fields[1:]...)
	if len(args) == 1 || strings.HasSuffix(input, " ") || strings.HasSuffix(input, "\t") || strings.HasPrefix(args[len(args)-1], "-") {
		return append(args, value)
	}
	args[len(args)-1] = value
	return args
}

func supportsPackage(verb string) bool {
	switch verb {
	case "build", "test", "vet", "fmt", "generate", "list":
		return true
	default:
		return false
	}
}

func dynamicRisk(verb string) sdk.Risk {
	switch verb {
	case "vet", "list":
		return sdk.Safe
	default:
		return sdk.Mutating
	}
}

func requiresProject(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "version", "env", "install":
		return false
	case "mod":
		return len(args) < 2 || args[1] != "init"
	case "work":
		return len(args) < 2 || args[1] != "init"
	default:
		return true
	}
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
	return sdk.Suggestion{Command: sdk.Command{Executable: "go", Args: copied}, Title: title, Description: title, Reason: reason, Kind: kind, Priority: priority, Risk: risk, Source: "go", Placeholders: placeholders}
}
