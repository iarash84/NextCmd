package builtin

import "testing"

func TestRegistrationIsExplicitAndOptional(t *testing.T) {
	if len(All(true)) != 1 || All(true)[0].Info().ID != "git" {
		t.Fatal("git not registered")
	}
	if len(All(false)) != 0 {
		t.Fatal("disabled plugin registered")
	}
}
