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
		dryRun  bool
		perm    bool
	}{
		{input: ":del", wantErr: true, handled: true},
		{input: ":del ", wantErr: true, handled: true},
		{input: ":del old.txt", want: "old.txt", handled: true},
		{input: ":del --dry-run old.txt", want: "old.txt", handled: true, dryRun: true},
		{input: ":del --permanent old.txt", want: "old.txt", handled: true, perm: true},
		{input: ":trash old.txt", want: "old.txt", handled: true},
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
		if !test.wantErr && err == nil && (got.requested != test.want || got.dryRun != test.dryRun || got.permanent != test.perm) {
			t.Errorf("parseDeletePath(%q) = %#v, want path=%q dryRun=%v permanent=%v", test.input, got, test.want, test.dryRun, test.perm)
		}
	}
}

func TestDeletePathMovesFileToTrash(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "old.txt")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	deleted, err := deletePath(base, deleteOptions{requested: "old.txt"}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if deleted.target.kind != deleteFile || deleted.trashed == "" {
		t.Fatalf("delete result = %#v", deleted)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("file still exists or unexpected error: %v", err)
	}
	if _, err := os.Stat(deleted.trashed); err != nil {
		t.Fatalf("trashed file missing: %v", err)
	}
}

func TestDeletePathPermanentlyRemovesDirectory(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "old", "nested.txt")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	deleted, err := deletePath(base, deleteOptions{requested: "old", permanent: true}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if deleted.target.kind != deleteDirectory {
		t.Fatalf("deleted kind = %q, want %q", deleted.target.kind, deleteDirectory)
	}
	if _, err := os.Stat(filepath.Join(base, "old")); !os.IsNotExist(err) {
		t.Fatalf("directory still exists or unexpected error: %v", err)
	}
}

func TestDeletePathDryRunDoesNotRemoveTarget(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "old.txt")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	deleted, err := deletePath(base, deleteOptions{requested: "old.txt", dryRun: true}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !deleted.dryRun || deleted.target.kind != deleteFile {
		t.Fatalf("delete result = %#v", deleted)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("dry-run removed target: %v", err)
	}
}

func TestRestoreDeletedMovesTrashBack(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "old.txt")
	if err := os.WriteFile(target, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	deleted, err := deletePath(base, deleteOptions{requested: "old.txt"}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	record := trashRecord(deleted)
	if record == nil {
		t.Fatal("expected undo record")
	}
	if err := restoreDeleted(*record); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("restore target missing: %v", err)
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
