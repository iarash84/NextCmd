package history

import (
	"strings"
	"testing"

	"nextcmd/sdk"
)

func TestRedactStructuredCommand(t *testing.T) {
	original := sdk.Command{
		Executable: "curl",
		Args: []string{
			"--token", "token-value",
			"--client-secret=client-value",
			"-H", "Authorization: Bearer header-value",
			"--user", "alice:user-password",
			"https://user:url-password@example.com/repo",
			"OUTPUT_TOKEN=environment-style-value",
		},
		Environment: map[string]string{"AWS_SECRET_ACCESS_KEY": "environment-value", "MODE": "test"},
	}
	got := Redact(original)
	serialized := got.Display() + " " + got.Environment["AWS_SECRET_ACCESS_KEY"]
	for _, secret := range []string{"token-value", "client-value", "header-value", "user-password", "url-password", "environment-style-value", "environment-value"} {
		if strings.Contains(serialized, secret) {
			t.Errorf("secret %q leaked in %q", secret, serialized)
		}
	}
	if got.Environment["MODE"] != "test" {
		t.Fatalf("non-sensitive environment changed: %#v", got.Environment)
	}
	if original.Args[1] != "token-value" || original.Environment["AWS_SECRET_ACCESS_KEY"] != "environment-value" {
		t.Fatal("Redact mutated the original command")
	}
}

func TestRedactShellCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		secrets []string
		keep    []string
	}{
		{
			name:    "quoted options and pipeline",
			command: `curl -H "Authorization: Bearer header-secret" --token 'token-secret' https://user:url-secret@example.com | jq .`,
			secrets: []string{"header-secret", "token-secret", "url-secret"},
			keep:    []string{"curl", "| jq ."},
		},
		{
			name:    "environment assignments",
			command: `export API_TOKEN="first-secret" && set DATABASE_PASSWORD=second-secret; echo done`,
			secrets: []string{"first-secret", "second-secret"},
			keep:    []string{"export", "&&", "; echo done"},
		},
		{
			name:    "powershell environment assignment",
			command: `$env:CLIENT_SECRET='power-secret'; $token = 'spaced-secret'; Invoke-WebRequest https://example.com`,
			secrets: []string{"power-secret", "spaced-secret"},
			keep:    []string{"Invoke-WebRequest", "https://example.com"},
		},
		{
			name:    "non-sensitive values",
			command: `echo "hello world" && curl -H "Accept: application/json" https://example.com`,
			keep:    []string{`echo "hello world"`, `Accept: application/json`, "https://example.com"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Redact(sdk.Command{ShellCommand: test.command}).ShellCommand
			for _, secret := range test.secrets {
				if strings.Contains(got, secret) {
					t.Errorf("secret %q leaked in %q", secret, got)
				}
			}
			for _, value := range test.keep {
				if !strings.Contains(got, value) {
					t.Errorf("expected %q to remain in %q", value, got)
				}
			}
		})
	}
}

func TestStoreRoundTrip(t *testing.T) {
	store := New(t.TempDir()+"/history.jsonl", true)
	if err := store.Append(sdk.HistoryEntry{Command: sdk.Command{ShellCommand: `curl --api-key persisted-secret https://example.com`}}); err != nil {
		t.Fatal(err)
	}
	entries, err := store.Load(10)
	if err != nil || len(entries) != 1 || strings.Contains(entries[0].Command.ShellCommand, "persisted-secret") || !strings.Contains(entries[0].Command.ShellCommand, redacted) {
		t.Fatalf("entries=%#v err=%v", entries, err)
	}
}
