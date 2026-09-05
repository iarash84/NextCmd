package kubernetes

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nextcmd/sdk"
)

type fakeRunner map[string]string

func (f fakeRunner) Run(_ context.Context, command sdk.Command) sdk.ExecutionResult {
	value, ok := f[strings.Join(command.Args, " ")]
	if !ok {
		return sdk.ExecutionResult{Command: command, Err: errors.New("not configured"), ExitCode: -1}
	}
	return sdk.ExecutionResult{Command: command, Stdout: value, ExitCode: 0}
}

func TestDetectAndDynamicCompletion(t *testing.T) {
	root := t.TempDir()
	manifest := "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: api\n  namespace: team-a\n"
	if err := os.WriteFile(filepath.Join(root, "deployment.yaml"), []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}
	p := NewWithRunner(fakeRunner{"config get-contexts -o name": "dev\nprod\n"})
	result, err := p.Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: root})
	if err != nil || !result.Detected {
		t.Fatalf("result=%#v err=%v", result, err)
	}
	state := result.Project.(State)
	if len(state.Manifests) != 1 || state.Kinds[0] != "Deployment" || len(state.Contexts) != 2 || state.Namespaces[0] != "team-a" {
		t.Fatalf("state=%#v", state)
	}
	items, _ := p.Complete(context.Background(), sdk.CompletionContext{Input: "kubectl apply ", Project: state})
	assertCommand(t, items, "kubectl apply -f deployment.yaml")
	items, _ = p.Complete(context.Background(), sdk.CompletionContext{Input: "kubectl config ", Project: state})
	assertCommand(t, items, "kubectl config use-context prod")
}
func TestRecoveryAndHelp(t *testing.T) {
	p := New()
	items, _ := p.Recover(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "kubectl"}, Stderr: "connection refused"}})
	if len(items) != 1 || items[0].Kind != sdk.Recovery {
		t.Fatalf("items=%#v", items)
	}
	if len(p.Help()) < 10 {
		t.Fatal("kubectl help catalog is incomplete")
	}
}
func assertCommand(t *testing.T, items []sdk.Suggestion, want string) {
	t.Helper()
	for _, item := range items {
		if item.Command.Display() == want {
			return
		}
	}
	t.Fatalf("command %q missing", want)
}
