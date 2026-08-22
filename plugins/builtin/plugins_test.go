package builtin

import "testing"

func TestRegistrationIsExplicitAndOptional(t *testing.T) {
	plugins := All(true, true)
	if len(plugins) != 2 || plugins[0].Info().ID != "git" || plugins[1].Info().ID != "dotnet" {
		t.Fatalf("unexpected built-ins: %#v", plugins)
	}
	if len(All(false, false)) != 0 {
		t.Fatal("disabled plugins registered")
	}
	if plugins := All(false, true); len(plugins) != 1 || plugins[0].Info().ID != "dotnet" {
		t.Fatal("dotnet registration is not independent")
	}
}
