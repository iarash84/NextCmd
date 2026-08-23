package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCargoPluginDefaultsAndExplicitDisable(t *testing.T) {
	if !Default().CargoEnabled {
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
	if cfg.CargoEnabled {
		t.Fatal("explicit Cargo disable was ignored")
	}
	if !cfg.GitEnabled || !cfg.DotnetEnabled {
		t.Fatal("missing plugin settings did not preserve their defaults")
	}
}
