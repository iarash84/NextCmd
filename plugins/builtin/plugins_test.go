package builtin

import (
	"strings"
	"testing"
)

func TestRegistrationMetadataIsValidAndUnique(t *testing.T) {
	plugins := All()
	if len(plugins) == 0 {
		t.Fatal("no built-in plugins registered")
	}
	seen := map[string]bool{}
	for _, plugin := range plugins {
		info := plugin.Info()
		if info.ID == "" || info.Name == "" || info.Version == "" {
			t.Fatalf("plugin has incomplete metadata: %#v", info)
		}
		id := strings.ToLower(info.ID)
		if seen[id] {
			t.Fatalf("duplicate plugin ID %q", info.ID)
		}
		seen[id] = true
	}
	if !seen["go"] {
		t.Fatal("Go plugin is not registered")
	}
	for _, id := range []string{"docker", "npm", "pip"} {
		if !seen[id] {
			t.Errorf("%s plugin is not registered", id)
		}
	}
}
