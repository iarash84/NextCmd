package npm

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
	{[]string{"--version"}, "Show npm version", sdk.Safe, 40},
	{[]string{"init", "-y"}, "Initialize package.json", sdk.Mutating, 58},
	{[]string{"install"}, "Install declared dependencies", sdk.Mutating, 86},
	{[]string{"ci"}, "Perform a clean lockfile install", sdk.Destructive, 88},
	{[]string{"run", "<script>"}, "Run a package script", sdk.Mutating, 90},
	{[]string{"test"}, "Run the test script", sdk.Mutating, 84},
	{[]string{"start"}, "Run the start script", sdk.Mutating, 74},
	{[]string{"install", "<package>"}, "Install a dependency", sdk.Mutating, 62},
	{[]string{"install", "--save-dev", "<package>"}, "Install a development dependency", sdk.Mutating, 58},
	{[]string{"uninstall", "<package>"}, "Remove a dependency", sdk.Mutating, 46},
	{[]string{"update"}, "Update dependencies", sdk.Mutating, 44},
	{[]string{"outdated"}, "List outdated dependencies", sdk.Safe, 66},
	{[]string{"audit"}, "Audit dependency vulnerabilities", sdk.Safe, 68},
	{[]string{"audit", "fix"}, "Apply compatible audit fixes", sdk.Mutating, 42},
	{[]string{"list", "--depth=0"}, "List direct dependencies", sdk.Safe, 56},
	{[]string{"view", "<package>", "version"}, "Show a package version", sdk.Safe, 40},
	{[]string{"exec", "<command>"}, "Run a package command", sdk.Mutating, 38},
	{[]string{"cache", "verify"}, "Verify npm cache", sdk.Safe, 36},
	{[]string{"cache", "clean", "--force"}, "Clear npm cache", sdk.Destructive, 18},
	{[]string{"publish", "--dry-run"}, "Preview package publication", sdk.Mutating, 30},
	{[]string{"publish"}, "Publish the package", sdk.Dangerous, 12},
}

func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	trimmed := strings.TrimSpace(input.Input)
	if !matches(trimmed, "npm") {
		return nil, nil
	}
	state, _ := input.Project.(State)
	out := make([]sdk.Suggestion, 0, len(commands)+len(state.Scripts))
	for _, spec := range commands {
		priority, reason := spec.priority, "Matches the current npm input"
		if requiresProject(spec.args) && state.Root == "" {
			priority -= 20
			reason = "No package.json was detected"
		}
		out = append(out, suggest(spec.args, spec.title, sdk.Completion, spec.risk, priority, reason))
	}
	return append(out, dynamic(input.Input, state)...), nil
}
func (*Plugin) Help() []sdk.CommandHelp {
	out := make([]sdk.CommandHelp, 0, len(commands))
	for _, spec := range commands {
		out = append(out, sdk.CommandHelp{Command: sdk.Command{Executable: "npm", Args: append([]string(nil), spec.args...)}, Description: spec.title, Risk: spec.risk})
	}
	return out
}
func dynamic(input string, state State) []sdk.Suggestion {
	fields := strings.Fields(input)
	if len(fields) < 2 || !strings.EqualFold(fields[0], "npm") {
		return nil
	}
	out := []sdk.Suggestion{}
	if strings.EqualFold(fields[1], "run") {
		for _, script := range state.Scripts {
			out = append(out, suggest([]string{"run", script}, "Run script "+script, sdk.Completion, sdk.Mutating, 99, "Script declared in package.json"))
		}
	}
	if strings.EqualFold(fields[1], "uninstall") {
		for _, dependency := range state.Dependencies {
			out = append(out, suggest([]string{"uninstall", dependency}, "Remove "+dependency, sdk.Completion, sdk.Mutating, 96, "Dependency declared in package.json"))
		}
	}
	if strings.EqualFold(fields[1], "--workspace") || strings.EqualFold(fields[1], "-w") {
		for _, workspace := range state.Workspaces {
			out = append(out, suggest([]string{"--workspace", workspace, "run", "<script>"}, "Run a workspace script in "+workspace, sdk.Completion, sdk.Mutating, 96, "Workspace package discovered locally"))
		}
	}
	return out
}
func requiresProject(args []string) bool {
	if len(args) == 0 {
		return false
	}
	return args[0] != "--version" && args[0] != "init" && args[0] != "view"
}
func matches(input, executable string) bool {
	if input == "" {
		return true
	}
	first := strings.ToLower(strings.Fields(input)[0])
	return strings.HasPrefix(executable, first) || first == executable
}
func suggest(args []string, title string, kind sdk.SuggestionKind, risk sdk.Risk, priority int, reason string) sdk.Suggestion {
	args = append([]string(nil), args...)
	placeholders := []sdk.Placeholder{}
	for i, arg := range args {
		s, e := strings.IndexByte(arg, '<'), strings.IndexByte(arg, '>')
		if s >= 0 && e > s {
			placeholders = append(placeholders, sdk.Placeholder{Name: arg[s+1 : e], ArgIndex: i, Start: s, End: e + 1})
		}
	}
	return sdk.Suggestion{Command: sdk.Command{Executable: "npm", Args: args}, Title: title, Description: title, Reason: reason, Kind: kind, Risk: risk, Priority: priority, Source: "npm", Placeholders: placeholders}
}
