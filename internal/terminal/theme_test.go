package terminal

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"nextcmd/sdk"
)

func TestSuggestionBadge(t *testing.T) {
	tests := []struct {
		kind sdk.SuggestionKind
		want string
	}{
		{sdk.Completion, "COMP"},
		{sdk.NextAction, "NEXT"},
		{sdk.BestPractice, "TIP"},
		{sdk.Recovery, "FIX"},
	}
	for _, test := range tests {
		if got := suggestionBadge(test.kind); got != test.want {
			t.Errorf("suggestionBadge(%q) = %q, want %q", test.kind, got, test.want)
		}
	}
}

func TestPrintExecutionSummaryWithoutTerminalColors(t *testing.T) {
	var output bytes.Buffer
	PrintExecutionSummary(&output, sdk.ExecutionResult{ExitCode: 1, Duration: 25 * time.Millisecond})
	text := output.String()
	if !strings.Contains(text, "failed") || !strings.Contains(text, "exit 1") {
		t.Fatalf("unexpected summary: %q", text)
	}
	if strings.Contains(text, "\x1b[") {
		t.Fatalf("non-terminal output contains ANSI colors: %q", text)
	}
}

func TestPrintExecutionSummaryShowsCancellation(t *testing.T) {
	var output bytes.Buffer
	PrintExecutionSummary(&output, sdk.ExecutionResult{ExitCode: -1, Canceled: true})
	if text := output.String(); !strings.Contains(text, "canceled") {
		t.Fatalf("unexpected summary: %q", text)
	}
}
