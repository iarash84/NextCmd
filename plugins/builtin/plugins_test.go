package builtin

import "testing"

func TestRegistrationIsExplicitAndOptional(t *testing.T) {
	plugins := All(true, true, true)
	if len(plugins) != 3 || plugins[0].Info().ID != "git" || plugins[1].Info().ID != "dotnet" || plugins[2].Info().ID != "cargo" {
		t.Fatalf("unexpected built-ins: %#v", plugins)
	}
	if len(All(false, false, false)) != 0 {
		t.Fatal("disabled plugins registered")
	}
	if plugins := All(false, true, false); len(plugins) != 1 || plugins[0].Info().ID != "dotnet" {
		t.Fatal("dotnet registration is not independent")
	}
	if plugins := All(false, false, true); len(plugins) != 1 || plugins[0].Info().ID != "cargo" {
		t.Fatal("cargo registration is not independent")
	}
}
