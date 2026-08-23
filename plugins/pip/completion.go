package pip

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
	{[]string{"--version"}, "Show pip version", sdk.Safe, 42},
	{[]string{"install", "-r", "<requirements>"}, "Install from a requirements file", sdk.Mutating, 88},
	{[]string{"install", "<package>"}, "Install a package", sdk.Mutating, 76},
	{[]string{"install", "--upgrade", "<package>"}, "Upgrade a package", sdk.Mutating, 58},
	{[]string{"install", "-e", "."}, "Install the current project in editable mode", sdk.Mutating, 66},
	{[]string{"uninstall", "<package>"}, "Uninstall a package", sdk.Destructive, 42},
	{[]string{"list"}, "List installed packages", sdk.Safe, 72},
	{[]string{"list", "--outdated"}, "List outdated packages", sdk.Safe, 68},
	{[]string{"show", "<package>"}, "Show package details", sdk.Safe, 62},
	{[]string{"check"}, "Verify installed dependency compatibility", sdk.Safe, 82},
	{[]string{"freeze"}, "Print installed package versions", sdk.Safe, 64},
	{[]string{"download", "<package>"}, "Download a package archive", sdk.Mutating, 46},
	{[]string{"wheel", "-r", "<requirements>"}, "Build wheels for requirements", sdk.Mutating, 44},
	{[]string{"cache", "info"}, "Show pip cache information", sdk.Safe, 38},
	{[]string{"cache", "purge"}, "Remove pip cache content", sdk.Destructive, 20},
	{[]string{"config", "debug"}, "Show configuration sources", sdk.Safe, 36},
	{[]string{"index", "versions", "<package>"}, "List available package versions", sdk.Safe, 34},
}

func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	trimmed := strings.TrimSpace(input.Input)
	executables := matchingExecutables(trimmed)
	if len(executables) == 0 {
		return nil, nil
	}
	state, _ := input.Project.(State)
	out := []sdk.Suggestion{}
	for _, executable := range executables {
		for _, spec := range commands {
			priority, reason := spec.priority, "Matches the current "+executable+" input"
			if needsProject(spec.args) && state.Root == "" {
				priority -= 12
				reason = "No Python project metadata was detected"
			}
			out = append(out, suggest(executable, spec.args, spec.title, sdk.Completion, spec.risk, priority, reason))
		}
		out = append(out, dynamic(executable, input.Input, state)...)
	}
	return out, nil
}
func (*Plugin) Help() []sdk.CommandHelp {
	out := []sdk.CommandHelp{}
	for _, executable := range []string{"pip", "pip3"} {
		for _, spec := range commands {
			out = append(out, sdk.CommandHelp{Command: sdk.Command{Executable: executable, Args: append([]string(nil), spec.args...)}, Description: spec.title, Risk: spec.risk})
		}
	}
	return out
}
func matchingExecutables(input string) []string {
	if input == "" {
		return []string{"pip", "pip3"}
	}
	first := strings.ToLower(strings.Fields(input)[0])
	out := []string{}
	for _, name := range []string{"pip", "pip3"} {
		if strings.HasPrefix(name, first) || first == name {
			out = append(out, name)
		}
	}
	return out
}
func dynamic(executable, input string, state State) []sdk.Suggestion {
	fields := strings.Fields(input)
	if len(fields) < 2 || !strings.EqualFold(fields[0], executable) {
		return nil
	}
	out := []sdk.Suggestion{}
	if strings.EqualFold(fields[1], "install") {
		for _, file := range state.RequirementFiles {
			out = append(out, suggest(executable, []string{"install", "-r", file}, "Install "+file, sdk.Completion, sdk.Mutating, 99, "Requirements file discovered locally"))
		}
	}
	if strings.EqualFold(fields[1], "show") || strings.EqualFold(fields[1], "uninstall") {
		risk := sdk.Safe
		if strings.EqualFold(fields[1], "uninstall") {
			risk = sdk.Destructive
		}
		for _, pkg := range state.Packages {
			verb := strings.ToLower(fields[1])
			title := strings.ToUpper(verb[:1]) + verb[1:] + " " + pkg
			out = append(out, suggest(executable, []string{verb, pkg}, title, sdk.Completion, risk, 94, "Package declared in a requirements file"))
		}
	}
	return out
}
func needsProject(args []string) bool {
	return len(args) > 0 && (args[0] == "install" && len(args) > 1 && (args[1] == "-r" || args[1] == "-e") || args[0] == "wheel")
}
func suggest(executable string, args []string, title string, kind sdk.SuggestionKind, risk sdk.Risk, priority int, reason string) sdk.Suggestion {
	args = append([]string(nil), args...)
	placeholders := []sdk.Placeholder{}
	for i, arg := range args {
		s, e := strings.IndexByte(arg, '<'), strings.IndexByte(arg, '>')
		if s >= 0 && e > s {
			placeholders = append(placeholders, sdk.Placeholder{Name: arg[s+1 : e], ArgIndex: i, Start: s, End: e + 1})
		}
	}
	return sdk.Suggestion{Command: sdk.Command{Executable: executable, Args: args}, Title: title, Description: title, Reason: reason, Kind: kind, Risk: risk, Priority: priority, Source: "pip", Placeholders: placeholders}
}
