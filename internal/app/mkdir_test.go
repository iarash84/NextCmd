package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMakeDirectory(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		handled bool
		wantErr bool
	}{
		{input: ":mkdir", wantErr: true, handled: true},
		{input: ":mkdir ", wantErr: true, handled: true},
		{input: ":mkdir tools", want: "tools", handled: true},
		{input: ":mkdir \"my dir\"", want: "my dir", handled: true},
		{input: ":mkdir 'my dir", wantErr: true, handled: true},
		{input: ":mkdirx tools", handled: false},
		{input: "mkdir tools", handled: false},
	}
	for _, test := range tests {
		got, handled, err := parseMakeDirectory(test.input)
		if handled != test.handled {
			t.Errorf("parseMakeDirectory(%q) handled = %v", test.input, handled)
		}
		if test.wantErr && err == nil {
			t.Errorf("parseMakeDirectory(%q) expected error", test.input)
		}
		if !test.wantErr && err == nil && got != test.want {
			t.Errorf("parseMakeDirectory(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestMakeDirectoryCreatesNestedDirectories(t *testing.T) {
	base := t.TempDir()
	created, err := makeDirectory(base, filepath.Join("a", "b", "c"))
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(created)
	if err != nil || !info.IsDir() {
		t.Fatalf("created directory missing: %q err=%v", created, err)
	}
	if !filepath.IsAbs(created) {
		t.Fatalf("created path must be absolute: %q", created)
	}
	// Relative to base: a/b/c should exist under base.
	if filepath.Base(created) != "c" {
		t.Fatalf("unexpected created path: %q", created)
	}
}

func TestMakeDirectoryRejectsEmpty(t *testing.T) {
	if _, err := makeDirectory(t.TempDir(), "  "); err == nil {
		t.Fatal("empty path must be rejected")
	}
}
