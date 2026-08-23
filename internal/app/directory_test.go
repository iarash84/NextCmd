package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDirectorySupportsRelativePaths(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "project with spaces")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := ResolveDirectory(root, "project with spaces")
	if err != nil || got != child {
		t.Fatalf("ResolveDirectory() = %q, %v", got, err)
	}
}

func TestParseChangeDirectory(t *testing.T) {
	tests := []struct {
		input, want string
		handled     bool
	}{
		{`cd "C:\work tree"`, `C:\work tree`, true},
		{`:cd ../project`, `../project`, true},
		{"cd", "", true},
		{"cargo build", "", false},
	}
	for _, test := range tests {
		got, handled, err := parseChangeDirectory(test.input)
		if err != nil || got != test.want || handled != test.handled {
			t.Errorf("parseChangeDirectory(%q) = %q, %v, %v", test.input, got, handled, err)
		}
	}
}

func TestResolveDirectoryRejectsFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveDirectory(filepath.Dir(path), path); err == nil {
		t.Fatal("file was accepted as a working directory")
	}
}
