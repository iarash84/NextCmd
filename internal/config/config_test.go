package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCargoPluginDefaultsAndExplicitDisable(t *testing.T) {
	if !Default().PluginEnabled("cargo") || !Default().PluginEnabled("future-plugin") {
		t.Fatal("Cargo plugin must be enabled by default")
	}
	path := filepath.Join(t.TempDir(), "nextcmd.json")
	if err := os.WriteFile(path, []byte(`{"cargoEnabled":false}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PluginEnabled("cargo") {
		t.Fatal("explicit Cargo disable was ignored")
	}
	if !cfg.PluginEnabled("git") || !cfg.PluginEnabled("dotnet") {
		t.Fatal("missing plugin settings did not preserve their defaults")
	}
}

func TestGenericPluginConfigurationNeedsNoConfigField(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nextcmd.json")
	if err := os.WriteFile(path, []byte(`{"plugins":{"future-plugin":false,"cargo":true}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PluginEnabled("future-plugin") || !cfg.PluginEnabled("CARGO") {
		t.Fatalf("unexpected generic plugin settings: %#v", cfg.Plugins)
	}
}
