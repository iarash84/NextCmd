package dart

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
	{[]string{"--version"}, "Show Dart SDK version", sdk.Safe, 45},
	{[]string{"analyze"}, "Analyze Dart source", sdk.Safe, 82},
	{[]string{"format", "."}, "Format Dart source", sdk.Mutating, 78},
	{[]string{"test"}, "Run Dart tests", sdk.Mutating, 86},
	{[]string{"run"}, "Run the Dart application", sdk.Mutating, 80},
	{[]string{"compile", "exe", "<file>"}, "Compile a Dart executable", sdk.Mutating, 62},
	{[]string{"pub", "get"}, "Resolve Dart dependencies", sdk.Mutating, 88},
	{[]string{"pub", "upgrade"}, "Upgrade Dart dependencies", sdk.Mutating, 58},
	{[]string{"pub", "outdated"}, "Check outdated Dart dependencies", sdk.Safe, 68},
	{[]string{"pub", "deps"}, "Show dependency tree", sdk.Safe, 52},
	{[]string{"pub", "cache", "repair"}, "Repair the pub cache", sdk.Mutating, 42},
	{[]string{"create", "<name>"}, "Create a Dart package", sdk.Mutating, 60},
}

func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	trimmed := strings.TrimSpace(input.Input)
	if trimmed != "" && !strings.HasPrefix(strings.ToLower(strings.Fields(trimmed)[0]), "dart") && !strings.HasPrefix(strings.ToLower(strings.Fields(trimmed)[0]), "flutter") {
		return nil, nil
	}
	state, _ := input.Project.(State)
	executable := "dart"
	if len(strings.Fields(trimmed)) > 0 && strings.EqualFold(strings.Fields(trimmed)[0], "flutter") {
		executable = "flutter"
	}
	out := []sdk.Suggestion{}
	for _, spec := range commands {
		out = append(out, suggest(executable, spec.args, spec.title, spec.risk, spec.priority, "Matches the current Dart project context"))
	}
	if state.IsFlutter {
		out = append(out, suggest("flutter", []string{"run"}, "Run the Flutter application", sdk.Mutating, 92, "A Flutter project was detected"), suggest("flutter", []string{"pub", "get"}, "Resolve Flutter dependencies", sdk.Mutating, 90, "A Flutter project was detected"))
	}
	return out, nil
}

func (*Plugin) Help() []sdk.CommandHelp {
	out := []sdk.CommandHelp{}
	for _, spec := range commands {
		out = append(out, sdk.CommandHelp{Command: sdk.Command{Executable: "dart", Args: spec.args}, Description: spec.title, Risk: spec.risk})
	}
	return out
}

func suggest(executable string, args []string, title string, risk sdk.Risk, priority int, reason string) sdk.Suggestion {
	args = append([]string(nil), args...)
	placeholders := []sdk.Placeholder{}
	for i, arg := range args {
		if s, e := strings.IndexByte(arg, '<'), strings.IndexByte(arg, '>'); s >= 0 && e > s {
			placeholders = append(placeholders, sdk.Placeholder{Name: arg[s+1 : e], ArgIndex: i, Start: s, End: e + 1})
		}
	}
	return sdk.Suggestion{Command: sdk.Command{Executable: executable, Args: args}, Title: title, Description: title, Reason: reason, Kind: sdk.Completion, Risk: risk, Priority: priority, Source: "dart", Placeholders: placeholders}
}
