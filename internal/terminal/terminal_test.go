package terminal

import (
	"bytes"
	"testing"

	"nextcmd/sdk"
)

func TestAcceptSelected(t *testing.T) {
	suggestions := []sdk.Suggestion{{
		Command: sdk.Command{Executable: "git", Args: []string{"status"}},
	}}

	accepted, ok := acceptSelected("git sta", suggestions, 0)
	if !ok || accepted != "git status" {
		t.Fatalf("acceptSelected() = %q, %v", accepted, ok)
	}

	accepted, ok = acceptSelected("git status", suggestions, 0)
	if ok || accepted != "git status" {
		t.Fatalf("an already accepted command must be ready to execute: %q, %v", accepted, ok)
	}
}

func TestReadKeyRecognizesRightArrow(t *testing.T) {
	tests := []struct {
		name     string
		sequence []byte
	}{
		{"ANSI", []byte{27, '[', 'C'}},
		{"Windows", []byte{224, 77}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event, err := readKey(bytes.NewReader(test.sequence))
			if err != nil || event.kind != KeyRight {
				t.Fatalf("readKey() = %#v, %v", event, err)
			}
		})
	}
}

func TestReadKeyIgnoresUnsupportedArrowSequence(t *testing.T) {
	event, err := readKey(bytes.NewReader([]byte{27, '[', 'H'}))
	if err != nil || event.kind != KeyIgnored {
		t.Fatalf("readKey() = %#v, %v", event, err)
	}
}

func TestReadKeyRecognizesLeftArrow(t *testing.T) {
	tests := []struct {
		name     string
		sequence []byte
	}{
		{"ANSI", []byte{27, '[', 'D'}},
		{"Windows", []byte{224, 75}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			event, err := readKey(bytes.NewReader(test.sequence))
			if err != nil || event.kind != KeyLeft {
				t.Fatalf("readKey() = %#v, %v", event, err)
			}
		})
	}
}

func TestAcceptSelectedRejectsInvalidIndex(t *testing.T) {
	line, ok := acceptSelected("git", nil, 0)
	if ok || line != "git" {
		t.Fatalf("invalid selection changed the line: %q, %v", line, ok)
	}
}
