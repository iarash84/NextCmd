package cargo

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
	{[]string{"--version"}, "Show the installed Cargo version", sdk.Safe, 35},
	{[]string{"new", "<name>"}, "Create a binary Rust package", sdk.Mutating, 68},
	{[]string{"new", "<name>", "--lib"}, "Create a Rust library package", sdk.Mutating, 64},
	{[]string{"init"}, "Initialize a Rust package in this directory", sdk.Mutating, 70},
	{[]string{"build"}, "Build the current package or workspace", sdk.Mutating, 90},
	{[]string{"build", "--release"}, "Build optimized release artifacts", sdk.Mutating, 72},
	{[]string{"check"}, "Check code without producing final binaries", sdk.Mutating, 94},
	{[]string{"run"}, "Build and run the current binary", sdk.Mutating, 86},
	{[]string{"test"}, "Run package or workspace tests", sdk.Mutating, 92},
	{[]string{"test", "--doc"}, "Run documentation tests", sdk.Mutating, 66},
	{[]string{"bench"}, "Run benchmark targets", sdk.Mutating, 52},
	{[]string{"fmt"}, "Format Rust source files", sdk.Mutating, 78},
	{[]string{"fmt", "--check"}, "Verify Rust source formatting", sdk.Safe, 82},
	{[]string{"clippy", "--all-targets"}, "Lint all package targets", sdk.Mutating, 84},
	{[]string{"clippy", "--all-targets", "--", "-D", "warnings"}, "Reject all Clippy warnings", sdk.Mutating, 76},
	{[]string{"doc", "--open"}, "Build and open local documentation", sdk.Mutating, 55},
	{[]string{"clean"}, "Remove generated build artifacts", sdk.Destructive, 38},
	{[]string{"fetch"}, "Fetch project dependencies", sdk.Mutating, 58},
	{[]string{"tree"}, "Display the dependency tree", sdk.Safe, 62},
	{[]string{"metadata", "--no-deps", "--format-version", "1"}, "Show workspace metadata", sdk.Safe, 48},
	{[]string{"update"}, "Update dependencies recorded in Cargo.lock", sdk.Mutating, 46},
	{[]string{"add", "<crate>"}, "Add a crate dependency", sdk.Mutating, 62},
	{[]string{"remove", "<crate>"}, "Remove a crate dependency", sdk.Mutating, 50},
	{[]string{"generate-lockfile"}, "Generate Cargo.lock without building", sdk.Mutating, 45},
	{[]string{"package"}, "Assemble a publishable crate archive", sdk.Mutating, 42},
	{[]string{"publish", "--dry-run"}, "Validate a crate publish without uploading", sdk.Mutating, 36},
	{[]string{"publish"}, "Publish a crate to the configured registry", sdk.Dangerous, 18},
}

func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	trimmed := strings.TrimSpace(input.Input)
	if !matchesExecutable(trimmed, "cargo") {
		return nil, nil
	}
	state, _ := input.Project.(State)
	out := make([]sdk.Suggestion, 0, len(commands)+len(state.Packages)+len(state.Features))
	for _, spec := range commands {
		priority := spec.priority
		reason := "Matches the current Cargo input"
		if requiresManifest(spec.args) && state.Root == "" {
			priority -= 18
			reason = "No Cargo.toml was detected; choose a manifest path if needed"
		}
		if state.Workspace && (spec.args[0] == "check" || spec.args[0] == "test" || spec.args[0] == "build") {
			priority += 10
			reason = "A Cargo workspace was detected"
		}
		out = append(out, suggestion(spec.args, spec.title, sdk.Completion, spec.risk, priority, reason))
	}
	return append(out, dynamic(input.Input, state)...), nil
}

func (*Plugin) Help() []sdk.CommandHelp {
	out := make([]sdk.CommandHelp, 0, len(commands))
	for _, spec := range commands {
		out = append(out, sdk.CommandHelp{Command: sdk.Command{Executable: "cargo", Args: append([]string(nil), spec.args...)}, Description: spec.title, Risk: spec.risk})
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
	fields := strings.Fields(strings.ToLower(input))
	if len(fields) < 2 || fields[0] != "cargo" {
		return nil
	}
	verb := fields[1]
	out := []sdk.Suggestion{}
	if supportsPackage(verb) && wantsValue(fields, "-p", "--package") {
		for _, pkg := range state.Packages {
			out = append(out, suggestion([]string{verb, "-p", pkg.Name}, verb+" package "+pkg.Name, sdk.Completion, dynamicRisk(verb), 98, "Package discovered in the current Cargo workspace"))
		}
	}
	if supportsFeatures(verb) && wantsValue(fields, "--features") {
		for _, feature := range state.Features {
			out = append(out, suggestion([]string{verb, "--features", feature}, verb+" with feature "+feature, sdk.Completion, sdk.Mutating, 96, "Feature discovered in Cargo.toml"))
		}
	}
	return out
}

func dynamicRisk(verb string) sdk.Risk {
	if verb == "clean" {
		return sdk.Destructive
	}
	return sdk.Mutating
}

func wantsValue(fields []string, flags ...string) bool {
	if len(fields) == 2 {
		return true
	}
	last := fields[len(fields)-1]
	for _, flag := range flags {
		if last == flag || strings.HasPrefix(flag, last) {
			return true
		}
		if len(fields) > 3 && fields[len(fields)-2] == flag {
			return true
		}
	}
	return false
}

func supportsPackage(verb string) bool {
	switch verb {
	case "build", "check", "run", "test", "bench", "doc", "clean", "clippy":
		return true
	default:
		return false
	}
}

func supportsFeatures(verb string) bool {
	switch verb {
	case "build", "check", "run", "test", "bench", "doc", "clippy":
		return true
	default:
		return false
	}
}

func requiresManifest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	switch args[0] {
	case "--version", "new", "init":
		return false
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
	return sdk.Suggestion{Command: sdk.Command{Executable: "cargo", Args: copied}, Title: title, Description: title, Reason: reason, Kind: kind, Priority: priority, Risk: risk, Source: "cargo", Placeholders: placeholders}
}
