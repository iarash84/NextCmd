package terraform

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nextcmd/sdk"
)

func TestDetectAndDynamicCompletion(t *testing.T) {
	root := t.TempDir()
	source := `resource "aws_instance" "api" {}` + "\n" + `module "network" {` + "\n"
	if err := os.WriteFile(filepath.Join(root, "main.tf"), []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, ".terraform"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".terraform", "environment"), []byte("staging\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	p := New()
	result, err := p.Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: root})
	if err != nil || !result.Detected {
		t.Fatalf("result=%#v err=%v", result, err)
	}
	state := result.Project.(State)
	if len(state.Resources) != 1 || state.Resources[0] != "aws_instance.api" || len(state.Workspaces) != 1 {
		t.Fatalf("state=%#v", state)
	}
	items, _ := p.Complete(context.Background(), sdk.CompletionContext{Input: "terraform workspace ", Project: state})
	assertCommand(t, items, "terraform workspace select staging")
	items, _ = p.Complete(context.Background(), sdk.CompletionContext{Input: "terraform state ", Project: state})
	assertCommand(t, items, "terraform state show aws_instance.api")
}
func TestWorkflowAndRecovery(t *testing.T) {
	p := New()
	next, _ := p.NextActions(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "terraform", Args: []string{"init"}}}})
	if len(next) != 2 {
		t.Fatalf("next=%#v", next)
	}
	recovery, _ := p.Recover(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "terraform"}, Stderr: "Initialization required; run terraform init"}})
	if len(recovery) != 1 || recovery[0].Command.Display() != "terraform init" {
		t.Fatalf("recovery=%#v", recovery)
	}
	if len(p.Help()) < 10 {
		t.Fatal("Terraform help catalog is incomplete")
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
