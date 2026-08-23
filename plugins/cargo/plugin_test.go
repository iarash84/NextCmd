package cargo

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nextcmd/sdk"
)

func TestDetectsWorkspacePackagesAndFeatures(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "Cargo.toml"), "[workspace]\nmembers = [\"app\", \"lib\"]\n")
	mustWrite(t, filepath.Join(root, "Cargo.lock"), "")
	mustWrite(t, filepath.Join(root, "app", "Cargo.toml"), "[package]\nname = \"app\"\nversion = \"0.1.0\"\n\n[features]\nfast = []\n")
	mustWrite(t, filepath.Join(root, "lib", "Cargo.toml"), "[package]\nname = \"shared-lib\"\nversion = \"0.1.0\"\n\n[features]\nserde = []\n")
	mustWrite(t, filepath.Join(root, "target", "ignored", "Cargo.toml"), "[package]\nname = \"ignored\"\n")

	result, err := New().Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: filepath.Join(root, "app")})
	if err != nil {
		t.Fatal(err)
	}
	state, ok := result.Project.(State)
	if !result.Detected || !ok || !state.Workspace || !state.HasLock {
		t.Fatalf("unexpected detection: %#v", result)
	}
	if len(state.Packages) != 2 || state.Packages[0].Name != "app" || state.Packages[1].Name != "shared-lib" {
		t.Fatalf("unexpected packages: %#v", state.Packages)
	}
	if len(state.Features) != 2 || state.Features[0] != "fast" || state.Features[1] != "serde" {
		t.Fatalf("unexpected features: %#v", state.Features)
	}
}

func TestCompletionOutsideProjectAndIncompletePrefix(t *testing.T) {
	items, err := New().Complete(context.Background(), sdk.CompletionContext{Input: "car"})
	if err != nil || !containsCommand(items, "cargo init") || !containsCommand(items, "cargo build") || !containsCommand(items, "cargo test") {
		t.Fatalf("cargo suggestions missing: %#v, %v", items, err)
	}
}

func TestDynamicPackageAndFeatureCompletion(t *testing.T) {
	state := State{Root: "workspace", Workspace: true, Packages: []Package{{Name: "api"}, {Name: "shared"}}, Features: []string{"serde"}}
	packages, _ := New().Complete(context.Background(), sdk.CompletionContext{Input: "cargo build -p", Project: state})
	if !containsCommand(packages, "cargo build -p api") || !containsCommand(packages, "cargo build -p shared") {
		t.Fatalf("workspace package completion missing: %#v", packages)
	}
	prefixed, _ := New().Complete(context.Background(), sdk.CompletionContext{Input: "cargo build -p a", Project: state})
	if !containsCommand(prefixed, "cargo build -p api") {
		t.Fatalf("partially typed package completion missing: %#v", prefixed)
	}
	features, _ := New().Complete(context.Background(), sdk.CompletionContext{Input: "cargo test --features", Project: state})
	if !containsCommand(features, "cargo test --features serde") {
		t.Fatalf("feature completion missing: %#v", features)
	}
}

func TestWorkflowAndRecovery(t *testing.T) {
	plugin := New()
	next, err := plugin.NextActions(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "cargo", Args: []string{"check"}}}})
	if err != nil || !containsCommand(next, "cargo test") || !containsCommand(next, "cargo clippy --all-targets") {
		t.Fatalf("next actions missing: %#v, %v", next, err)
	}
	recovery, err := plugin.Recover(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "cargo", Args: []string{"build"}}, ExitCode: 1, Stderr: "could not find `Cargo.toml`"}})
	if err != nil || !containsCommand(recovery, "cargo init") {
		t.Fatalf("recovery missing: %#v, %v", recovery, err)
	}
}

func TestHelpContainsCoreCargoCommands(t *testing.T) {
	commands := New().Help()
	found := map[string]bool{}
	for _, command := range commands {
		found[command.Command.Display()] = true
	}
	for _, command := range []string{"cargo build", "cargo check", "cargo run", "cargo test", "cargo fmt", "cargo clippy --all-targets"} {
		if !found[command] {
			t.Errorf("help command %q missing", command)
		}
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func containsCommand(items []sdk.Suggestion, command string) bool {
	for _, item := range items {
		if item.Command.Display() == command {
			return true
		}
	}
	return false
}
