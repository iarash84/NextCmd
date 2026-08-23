package docker

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nextcmd/sdk"
)

func TestDetectsComposeServicesAndDockerTargets(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "compose.yaml"), "services:\n  api:\n    build: .\n  database:\n    image: postgres\nnetworks:\n  default:\n")
	write(t, filepath.Join(root, "Dockerfile"), "FROM golang:1.24 AS builder\nFROM scratch AS runtime\n")
	result, err := New().Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: root})
	if err != nil {
		t.Fatal(err)
	}
	state := result.Project.(State)
	if !result.Detected || !has(state.Services, "api") || !has(state.Services, "database") || !has(state.Targets, "builder") {
		t.Fatalf("unexpected state: %#v", state)
	}
}
func TestDynamicComposeAndBuildCompletion(t *testing.T) {
	state := State{Root: "project", ComposeFile: "compose.yaml", Services: []string{"api"}, Targets: []string{"builder"}}
	items, _ := New().Complete(context.Background(), sdk.CompletionContext{Input: "docker compose logs a", Project: state})
	if !hasCommand(items, "docker compose logs api") {
		t.Fatalf("service completion missing: %#v", items)
	}
	items, _ = New().Complete(context.Background(), sdk.CompletionContext{Input: "docker build --t", Project: state})
	if !hasCommand(items, "docker build --target builder .") {
		t.Fatalf("target completion missing: %#v", items)
	}
}
func TestDockerWorkflowAndRecovery(t *testing.T) {
	plugin := New()
	next, _ := plugin.NextActions(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "docker", Args: []string{"build", "."}}}})
	if !hasCommand(next, "docker images") {
		t.Fatalf("next action missing: %#v", next)
	}
	recovery, _ := plugin.Recover(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "docker", Args: []string{"ps"}}, Stderr: "Cannot connect to the Docker daemon"}})
	if !hasCommand(recovery, "docker info") {
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
