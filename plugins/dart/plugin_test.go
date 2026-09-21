package dart

import (
	"context"
	"testing"

	"nextcmd/sdk"
)

func TestCompleteAcceptsExecutablePrefixes(t *testing.T) {
	plugin := New()
	for _, input := range []string{"da", "dart", "f", "flutter"} {
		suggestions, err := plugin.Complete(context.Background(), sdk.CompletionContext{Input: input})
		if err != nil {
			t.Fatalf("Complete(%q) returned an error: %v", input, err)
		}
		if len(suggestions) == 0 {
			t.Fatalf("Complete(%q) returned no suggestions", input)
		}
	}
}

func TestCompleteRejectsUnrelatedExecutable(t *testing.T) {
	suggestions, err := New().Complete(context.Background(), sdk.CompletionContext{Input: "docker"})
	if err != nil {
		t.Fatalf("Complete returned an error: %v", err)
	}
	if len(suggestions) != 0 {
		t.Fatalf("Complete(%q) returned %d suggestions; want none", "docker", len(suggestions))
	}
}
