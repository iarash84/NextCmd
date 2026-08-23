package cargo

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "cargo") || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	switch input.Result.Command.Args[0] {
	case "new", "init":
		return []sdk.Suggestion{suggestion([]string{"check"}, "Check the new package", sdk.NextAction, sdk.Mutating, 92, "A Cargo package was created"), suggestion([]string{"run"}, "Run the new binary", sdk.NextAction, sdk.Mutating, 78, "Try the generated package")}, nil
	case "check":
		return []sdk.Suggestion{suggestion([]string{"clippy", "--all-targets"}, "Lint all targets", sdk.NextAction, sdk.Mutating, 86, "The package passed cargo check"), suggestion([]string{"test"}, "Run the test suite", sdk.NextAction, sdk.Mutating, 90, "The package passed cargo check")}, nil
	case "build":
		return []sdk.Suggestion{suggestion([]string{"test"}, "Run tests after the build", sdk.NextAction, sdk.Mutating, 91, "The package built successfully"), suggestion([]string{"run"}, "Run the built binary", sdk.NextAction, sdk.Mutating, 76, "The package built successfully")}, nil
	case "test":
		return []sdk.Suggestion{suggestion([]string{"build", "--release"}, "Build optimized release artifacts", sdk.NextAction, sdk.Mutating, 78, "Tests completed successfully"), suggestion([]string{"doc"}, "Build project documentation", sdk.NextAction, sdk.Mutating, 60, "Tests completed successfully")}, nil
	case "add", "remove", "update":
		return []sdk.Suggestion{suggestion([]string{"check"}, "Check dependency compatibility", sdk.NextAction, sdk.Mutating, 92, "Project dependencies changed"), suggestion([]string{"test"}, "Test after dependency changes", sdk.NextAction, sdk.Mutating, 84, "Project dependencies changed")}, nil
	case "fmt":
		return []sdk.Suggestion{suggestion([]string{"check"}, "Check the formatted source", sdk.NextAction, sdk.Mutating, 78, "Formatting completed")}, nil
	}
	return nil, nil
}

func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if state.Root == "" {
		return nil, nil
	}
	return []sdk.Suggestion{
		suggestion([]string{"fmt", "--check"}, "Verify Rust formatting", sdk.BestPractice, sdk.Safe, 62, "Keep source formatting deterministic"),
		suggestion([]string{"clippy", "--all-targets", "--", "-D", "warnings"}, "Reject Clippy warnings", sdk.BestPractice, sdk.Mutating, 66, "Catch common Rust mistakes before committing"),
		suggestion([]string{"test"}, "Run the Cargo test suite", sdk.BestPractice, sdk.Mutating, 68, "A Cargo package was detected"),
	}, nil
}

func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "cargo") {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	state, _ := input.Project.(State)
	if state.Root == "" || strings.Contains(message, "could not find `cargo.toml`") || strings.Contains(message, "could not find cargo.toml") {
		return []sdk.Suggestion{suggestion([]string{"init"}, "Initialize a Cargo package", sdk.Recovery, sdk.Mutating, 96, "No Cargo.toml was found")}, nil
	}
	if strings.Contains(message, "package id specification") || strings.Contains(message, "did not match any packages") {
		out := make([]sdk.Suggestion, 0, len(state.Packages))
		for _, pkg := range state.Packages {
			out = append(out, suggestion([]string{"build", "-p", pkg.Name}, "Build package "+pkg.Name, sdk.Recovery, sdk.Mutating, 92, "Use a package discovered in this workspace"))
		}
		return out, nil
	}
	if strings.Contains(message, "does not have the feature") || strings.Contains(message, "none of the selected packages contains") {
		out := make([]sdk.Suggestion, 0, len(state.Features))
		for _, feature := range state.Features {
			out = append(out, suggestion([]string{"build", "--features", feature}, "Build with feature "+feature, sdk.Recovery, sdk.Mutating, 90, "Use a feature declared in Cargo.toml"))
		}
		return out, nil
	}
	return nil, nil
}
