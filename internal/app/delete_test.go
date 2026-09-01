package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDeletePath(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		handled bool
		wantErr bool
	}{
		{input: ":del", wantErr: true, handled: true},
		{input: ":del ", wantErr: true, handled: true},
		{input: ":del old.txt", want: "old.txt", handled: true},
		{input: ":del \"old dir\"", want: "old dir", handled: true},
		{input: ":del 'old dir", wantErr: true, handled: true},
		{input: ":delete old.txt", handled: false},
		{input: "del old.txt", handled: false},
	}
	for _, test := range tests {
		got, handled, err := parseDeletePath(test.input)
		if handled != test.handled {
			t.Errorf("parseDeletePath(%q) handled = %v", test.input, handled)
		}
		if test.wantErr && err == nil {
			t.Errorf("parseDeletePath(%q) expected error", test.input)
		}
		if !test.wantErr && err == nil && got != test.want {
			t.Errorf("parseDeletePath(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestDeletePathRemovesFile(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "old.txt")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	deleted, err := deletePath(base, "old.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	if deleted.kind != deleteFile {
		t.Fatalf("deleted kind = %q, want %q", deleted.kind, deleteFile)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("file still exists or unexpected error: %v", err)
	}
}

func TestDeletePathRemovesDirectory(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "old", "nested.txt")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	deleted, err := deletePath(base, "old", nil)
	if err != nil {
		t.Fatal(err)
	}
	if deleted.kind != deleteDirectory {
		t.Fatalf("deleted kind = %q, want %q", deleted.kind, deleteDirectory)
	}
	if _, err := os.Stat(filepath.Join(base, "old")); !os.IsNotExist(err) {
		t.Fatalf("directory still exists or unexpected error: %v", err)
	}
}

func TestPromptDeleteKind(t *testing.T) {
	var output bytes.Buffer
	got, err := promptDeleteKind(bytes.NewBufferString("d\n"), &output, []deleteCandidate{
		{path: "same", kind: deleteFile},
		{path: "same", kind: deleteDirectory},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != deleteDirectory {
		t.Fatalf("promptDeleteKind() = %q, want %q", got, deleteDirectory)
	}
	if !bytes.Contains(output.Bytes(), []byte("file")) || !bytes.Contains(output.Bytes(), []byte("directory")) {
		t.Fatalf("prompt did not describe both choices: %q", output.String())
	}
}
