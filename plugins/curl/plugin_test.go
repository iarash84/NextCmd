package curl

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nextcmd/sdk"
)

func TestIncompletePrefixAndCoreCommands(t *testing.T) {
	items, err := New().Complete(context.Background(), sdk.CompletionContext{Input: "cur"})
	if err != nil {
		t.Fatal(err)
	}
	for _, command := range []string{"curl <url>", "curl --head <url>", `curl --request POST --header "Content-Type: application/json" --data <json> <url>`} {
		if !containsCommand(items, command) {
			t.Errorf("command %q missing", command)
		}
	}
	if riskFor(items, "curl --request DELETE <url>") != sdk.Destructive {
		t.Fatal("DELETE request is not destructive")
	}
	if riskFor(items, "curl --insecure <url>") != sdk.Dangerous {
		t.Fatal("insecure TLS request is not dangerous")
	}
}

func TestDetectAndDynamicFileCompletion(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "payload.json"))
	mustWrite(t, filepath.Join(root, "client.pem"))
	mustWrite(t, filepath.Join(root, "request.curl"))
	result, err := New().Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: root})
	if err != nil {
		t.Fatal(err)
	}
	state, ok := result.Project.(State)
	if !result.Detected || !ok || len(state.Files) != 3 {
		t.Fatalf("unexpected state: %#v", result)
	}
	tests := []struct{ input, command string }{
		{"curl --upload-file pay", "curl --upload-file payload.json <url>"},
		{"curl --data @pay", "curl --data-binary @payload.json <url>"},
		{"curl --cacert cli", "curl --cacert client.pem <url>"},
		{"curl -K req", "curl --config request.curl <url>"},
	}
	for _, test := range tests {
		items, completeErr := New().Complete(context.Background(), sdk.CompletionContext{Input: test.input, Project: state})
		if completeErr != nil || !containsCommand(items, test.command) {
			t.Errorf("completion %q missing %q: %#v, %v", test.input, test.command, items, completeErr)
		}
	}
}

func TestWorkflowAndRecovery(t *testing.T) {
	plugin := New()
	next, err := plugin.NextActions(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "curl", Args: []string{"https://example.com"}}}})
	if err != nil || !containsCommand(next, "curl --head https://example.com") {
		t.Fatalf("next actions missing URL: %#v, %v", next, err)
	}
	recovery, err := plugin.Recover(context.Background(), sdk.ExecutionContext{
		Project: State{Certificates: []string{"local-ca.pem"}},
		Result:  sdk.ExecutionResult{Command: sdk.Command{Executable: "curl", Args: []string{"https://example.com"}}, ExitCode: 60, Stderr: "SSL certificate problem"},
	})
	if err != nil || !containsCommand(recovery, "curl --cacert local-ca.pem https://example.com") || riskFor(recovery, "curl --insecure https://example.com") != sdk.Dangerous {
		t.Fatalf("TLS recovery missing: %#v, %v", recovery, err)
	}
}

func TestBestPracticesOnlyApplyToCurlInput(t *testing.T) {
	plugin := New()
	items, _ := plugin.BestPractices(context.Background(), sdk.CommandContext{Input: "git status"})
	if len(items) != 0 {
		t.Fatalf("curl advice leaked into another tool: %#v", items)
	}
	items, _ = plugin.BestPractices(context.Background(), sdk.CommandContext{Input: "curl https://example.com"})
	if !containsCommand(items, "curl --fail-with-body --show-error <url>") {
		t.Fatalf("curl best practice missing: %#v", items)
	}
}

func TestHelpContainsCommonCommands(t *testing.T) {
	found := map[string]bool{}
	for _, item := range New().Help() {
		found[item.Command.Display()] = true
	}
	for _, command := range []string{"curl <url>", "curl --head <url>", "curl --retry 3 --connect-timeout 10 --max-time 60 <url>"} {
		if !found[command] {
			t.Errorf("help command %q missing", command)
		}
	}
}

func mustWrite(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("test"), 0o600); err != nil {
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

func riskFor(items []sdk.Suggestion, command string) sdk.Risk {
	for _, item := range items {
		if item.Command.Display() == command {
			return item.Risk
		}
	}
	return ""
}
