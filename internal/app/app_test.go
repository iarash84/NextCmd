package app

import "testing"

func TestIsExitCommand(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"exit", true},
		{" QUIT ", true},
		{":q", true},
		{"git status", false},
		{"exit now", false},
	}

	for _, test := range tests {
		if got := isExitCommand(test.input); got != test.want {
			t.Errorf("isExitCommand(%q) = %v, want %v", test.input, got, test.want)
		}
	}
}
