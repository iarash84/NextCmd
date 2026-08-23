package npm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nextcmd/sdk"
)

func TestDetectsScriptsWorkspacesAndDependencies(t *testing.T) {
	root := t.TempDir()
	write(t, filepath.Join(root, "package.json"), `{"name":"root","workspaces":["packages/*"],"scripts":{"build":"tool build","test":"tool test"},"dependencies":{"react":"1.0.0"}}`)
	write(t, filepath.Join(root, "package-lock.json"), "{}")
	write(t, filepath.Join(root, "packages", "api", "package.json"), `{"name":"@app/api"}`)
	result, err := New().Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: filepath.Join(root, "packages", "api")})
	if err != nil {
		t.Fatal(err)
	}
	state := result.Project.(State)
	if !result.Detected || state.Name != "root" || !state.HasLock || !has(state.Scripts, "build") || !has(state.Workspaces, "@app/api") || !has(state.Dependencies, "react") {
		t.Fatalf("unexpected state: %#v", state)
	}
}
func TestDynamicScriptDependencyAndWorkspaceCompletion(t *testing.T) {
	state := State{Root: "project", Scripts: []string{"build"}, Dependencies: []string{"react"}, Workspaces: []string{"@app/api"}}
	items, _ := New().Complete(context.Background(), sdk.CompletionContext{Input: "npm run b", Project: state})
	if !hasCommand(items, "npm run build") {
		t.Fatalf("script completion missing: %#v", items)
	}
	items, _ = New().Complete(context.Background(), sdk.CompletionContext{Input: "npm uninstall r", Project: state})
	if !hasCommand(items, "npm uninstall react") {
		t.Fatalf("dependency completion missing: %#v", items)
	}
	items, _ = New().Complete(context.Background(), sdk.CompletionContext{Input: "npm -w", Project: state})
	if !hasCommand(items, "npm --workspace @app/api run <script>") {
		t.Fatalf("workspace completion missing: %#v", items)
	}
}
func TestNPMRecoveryUsesDeclaredScripts(t *testing.T) {
	items, _ := New().Recover(context.Background(), sdk.ExecutionContext{Project: State{Root: "project", Scripts: []string{"build"}}, Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "npm", Args: []string{"run", "missing"}}, Stderr: "npm error Missing script"}})
	if !hasCommand(items, "npm run build") {
		t.Fatalf("recovery missing: %#v", items)
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
