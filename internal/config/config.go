package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	HistoryEnabled bool            `json:"historyEnabled"`
	MaxSuggestions int             `json:"maxSuggestions"`
	Debug          bool            `json:"debug"`
	Plugins        map[string]bool `json:"plugins,omitempty"`
}

func Default() Config {
	return Config{HistoryEnabled: true, MaxSuggestions: 8, Plugins: map[string]bool{}}
}

// PluginEnabled defaults unknown plugins to enabled, so adding a built-in does
// not require a new Core configuration field.
func (c Config) PluginEnabled(id string) bool {
	for configuredID, enabled := range c.Plugins {
		if strings.EqualFold(configuredID, id) {
			return enabled
		}
	}
	return true
}
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	if cfg.Plugins == nil {
		cfg.Plugins = map[string]bool{}
	}
	// Read the pre-v0.4 plugin flags as a migration path. New plugins only use
	// the generic plugins map and never add another field here.
	legacy := struct {
		GitEnabled    *bool `json:"gitEnabled"`
		DotnetEnabled *bool `json:"dotnetEnabled"`
		CargoEnabled  *bool `json:"cargoEnabled"`
	}{}
	if err := json.Unmarshal(data, &legacy); err != nil {
		return cfg, err
	}
	applyLegacyPlugin(cfg.Plugins, "git", legacy.GitEnabled)
	applyLegacyPlugin(cfg.Plugins, "dotnet", legacy.DotnetEnabled)
	applyLegacyPlugin(cfg.Plugins, "cargo", legacy.CargoEnabled)
	if cfg.MaxSuggestions <= 0 {
		cfg.MaxSuggestions = 8
	}
	return cfg, nil
}

func applyLegacyPlugin(plugins map[string]bool, id string, legacy *bool) {
	if !pluginConfigured(plugins, id) && legacy != nil {
		plugins[id] = *legacy
	}
}

func pluginConfigured(plugins map[string]bool, id string) bool {
	for configuredID := range plugins {
		if strings.EqualFold(configuredID, id) {
			return true
		}
	}
	return false
}
func DefaultPath() string {
	directory, err := os.UserConfigDir()
	if err != nil {
		return "nextcmd.json"
	}
	return filepath.Join(directory, "nextcmd", "config.json")
}
