package history

import (
	"nextcmd/sdk"
	"testing"
)

func TestRedact(t *testing.T) {
	got := Redact(sdk.Command{Executable: "tool", Args: []string{"--token", "abc", "--password=xyz", "https://user:pass@example.com/repo"}})
	for _, arg := range got.Args {
		if arg == "abc" || arg == "--password=xyz" || arg == "https://user:pass@example.com/repo" {
			t.Fatalf("secret leaked: %q", arg)
		}
	}
}
func TestStoreRoundTrip(t *testing.T) {
	store := New(t.TempDir()+"/history.jsonl", true)
	if err := store.Append(sdk.HistoryEntry{Command: sdk.Command{Executable: "git", Args: []string{"--secret", "value"}}}); err != nil {
		t.Fatal(err)
	}
	entries, err := store.Load(10)
	if err != nil || len(entries) != 1 || entries[0].Command.Args[1] != "<redacted>" {
		t.Fatalf("entries=%#v err=%v", entries, err)
	}
}
