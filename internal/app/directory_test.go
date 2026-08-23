package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
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

func TestParseListDirectory(t *testing.T) {
	tests := []struct {
		input, want string
		handled     bool
	}{
		{":ls", "", true},
		{":LS ..", "..", true},
		{`:ls "project with spaces"`, "project with spaces", true},
		{"ls", "", false},
		{"git status", "", false},
	}
	for _, test := range tests {
		got, handled, err := parseListDirectory(test.input)
		if err != nil || got != test.want || handled != test.handled {
			t.Errorf("parseListDirectory(%q) = %q, %v, %v", test.input, got, handled, err)
		}
	}
}

func TestParseListDirectoryRejectsUnclosedQuote(t *testing.T) {
	if _, handled, err := parseListDirectory(`:ls "unfinished`); !handled || err == nil {
		t.Fatalf("parseListDirectory() handled=%v, error=%v", handled, err)
	}
}

func TestPrintDirectoryListingShowsDirectoriesBeforeFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "project"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "notes.txt"), []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if err := printDirectoryListing(&output, root); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, want := range []string{"Directory: " + root, "TYPE", "SIZE", "NAME", "DIR", "project", "FILE", "4 B", "notes.txt"} {
		if !strings.Contains(text, want) {
			t.Errorf("listing does not contain %q:\n%s", want, text)
		}
	}
	if strings.Index(text, "project") > strings.Index(text, "notes.txt") {
		t.Errorf("directory must be listed before file:\n%s", text)
	}
}
