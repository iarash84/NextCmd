package golang

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nextcmd/sdk"
)

func TestDetectsModulePackagesFilesAndTests(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "go.mod"), "module example.com/service\n\ngo 1.24\n")
	mustWrite(t, filepath.Join(root, "main.go"), "package main\n")
	mustWrite(t, filepath.Join(root, "internal", "auth", "auth.go"), "package auth\n")
	mustWrite(t, filepath.Join(root, "internal", "auth", "auth_test.go"), "package auth\n")
	mustWrite(t, filepath.Join(root, "vendor", "ignored", "ignored.go"), "package ignored\n")

	result, err := New().Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: filepath.Join(root, "internal", "auth")})
	if err != nil {
		t.Fatal(err)
	}
	state, ok := result.Project.(State)
	if !result.Detected || !ok || state.Root != root || state.ModulePath != "example.com/service" || state.GoVersion != "1.24" || !state.HasTests {
		t.Fatalf("unexpected detection: %#v", result)
	}
	if !contains(state.Packages, ".") || !contains(state.Packages, "./internal/auth") {
		t.Fatalf("packages were not discovered: %#v", state.Packages)
	}
	if contains(state.GoFiles, "vendor/ignored/ignored.go") {
		t.Fatalf("vendor file was scanned: %#v", state.GoFiles)
	}
}

func TestDetectsWorkspaceFromNestedModule(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "go.work"), "go 1.24\n\nuse ./service\n")
	mustWrite(t, filepath.Join(root, "service", "go.mod"), "module example.com/service\n\ngo 1.24\n")
	mustWrite(t, filepath.Join(root, "service", "main.go"), "package main\n")

	result, err := New().Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: filepath.Join(root, "service")})
	if err != nil {
		t.Fatal(err)
	}
	state := result.Project.(State)
	if !state.Workspace || state.Root != root || state.ModulePath != "example.com/service" {
		t.Fatalf("workspace was not detected: %#v", state)
	}
}

func TestCompletionOutsideProjectAndIncompletePrefix(t *testing.T) {
	items, err := New().Complete(context.Background(), sdk.CompletionContext{Input: "g"})
	if err != nil || !containsCommand(items, "go mod init <module>") || !containsCommand(items, "go build ./...") || !containsCommand(items, "go test ./...") {
		t.Fatalf("Go suggestions missing: %#v, %v", items, err)
	}
}

func TestDynamicPackageAndFileCompletion(t *testing.T) {
	state := State{Root: "project", Packages: []string{".", "./cmd/api", "./internal/auth"}, GoFiles: []string{"main.go", "tools.go", "main_test.go"}}
	packages, _ := New().Complete(context.Background(), sdk.CompletionContext{Input: "go test ./c", Project: state})
	if !containsCommand(packages, "go test ./cmd/api") {
		t.Fatalf("package completion missing: %#v", packages)
	}
	packagesWithFlag, _ := New().Complete(context.Background(), sdk.CompletionContext{Input: "go test -race ./c", Project: state})
	if !containsCommand(packagesWithFlag, "go test -race ./cmd/api") {
		t.Fatalf("package completion dropped an existing flag: %#v", packagesWithFlag)
	}
	files, _ := New().Complete(context.Background(), sdk.CompletionContext{Input: "go run t", Project: state})
	if !containsCommand(files, "go run tools.go") || containsCommand(files, "go run main_test.go") {
		t.Fatalf("file completion is incorrect: %#v", files)
	}
}

func TestWorkflowBestPracticesAndRecovery(t *testing.T) {
	plugin := New()
	next, err := plugin.NextActions(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "go", Args: []string{"fmt", "./..."}}}})
	if err != nil || !containsCommand(next, "go vet ./...") || !containsCommand(next, "go test ./...") {
		t.Fatalf("next actions missing: %#v, %v", next, err)
	}
	practices, err := plugin.BestPractices(context.Background(), sdk.CommandContext{Project: State{Root: "project"}})
	if err != nil || !containsCommand(practices, "go fmt ./...") || !containsCommand(practices, "go vet ./...") {
		t.Fatalf("best practices missing: %#v, %v", practices, err)
	}
	recovery, err := plugin.Recover(context.Background(), sdk.ExecutionContext{Project: State{Root: "project"}, Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "go", Args: []string{"test", "./..."}}, ExitCode: 1, Stderr: "go: updates to go.mod needed; to update it: go mod tidy"}})
	if err != nil || !containsCommand(recovery, "go mod tidy") {
		t.Fatalf("recovery missing: %#v, %v", recovery, err)
	}
}

func TestHelpContainsCoreGoCommands(t *testing.T) {
	found := map[string]bool{}
	for _, command := range New().Help() {
		found[command.Command.Display()] = true
	}
	for _, command := range []string{"go build ./...", "go test ./...", "go vet ./...", "go fmt ./...", "go mod tidy", "go run ."} {
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

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func containsCommand(items []sdk.Suggestion, command string) bool {
	for _, item := range items {
		if item.Command.Display() == command {
			return true
		}
	}
	return false
}
