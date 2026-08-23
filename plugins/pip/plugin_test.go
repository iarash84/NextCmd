package pip

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nextcmd/sdk"
)

func TestDetectsRequirementsPackagesAndVirtualEnvironment(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "requirements-dev.txt"), "requests==2.32\npytest>=8\nprivate-lib @ https://user:secret@example.com/package.whl\nhttps://user:secret@example.com/raw.whl\n-e ./local\n")
	if err := os.Mkdir(filepath.Join(root, ".venv"), 0o755); err != nil {
		t.Fatal(err)
	}
	result, err := New().Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: root})
	if err != nil {
		t.Fatal(err)
	}
	state := result.Project.(State)
	if !result.Detected || !has(state.RequirementFiles, "requirements-dev.txt") || !has(state.Packages, "requests") || !has(state.Packages, "pytest") || !has(state.VirtualEnvironments, ".venv") {
		t.Fatalf("unexpected state: %#v", state)
	}
	if !has(state.Packages, "private-lib") || has(state.Packages, "https://user:secret@example.com/raw.whl") {
		t.Fatalf("direct URL parsing leaked or lost package metadata: %#v", state.Packages)
	}
}
func TestPipAndPip3DynamicCompletion(t *testing.T) {
	state := State{Root: "project", RequirementFiles: []string{"requirements.txt"}, Packages: []string{"requests"}}
	items, _ := New().Complete(context.Background(), sdk.CompletionContext{Input: "pip3 install -r", Project: state})
	if !hasCommand(items, "pip3 install -r requirements.txt") {
		t.Fatalf("pip3 requirements completion missing: %#v", items)
	}
	items, _ = New().Complete(context.Background(), sdk.CompletionContext{Input: "pip show req", Project: state})
	if !hasCommand(items, "pip show requests") {
		t.Fatalf("package completion missing: %#v", items)
	}
}
func TestPipWorkflowAndRecovery(t *testing.T) {
	next, _ := New().NextActions(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "pip3", Args: []string{"install", "requests"}}}})
	if !hasCommand(next, "pip3 check") {
		t.Fatalf("next action missing: %#v", next)
	}
	recovery, _ := New().Recover(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "pip", Args: []string{"install", "missing"}}, Stderr: "No matching distribution found"}})
	if !hasCommand(recovery, "pip index versions <package>") {
		t.Fatalf("recovery missing: %#v", recovery)
	}
}
func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
func has(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
func hasCommand(items []sdk.Suggestion, want string) bool {
	for _, item := range items {
		if item.Command.Display() == want {
			return true
		}
	}
	return false
}
