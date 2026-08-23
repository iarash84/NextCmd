package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	HistoryEnabled bool `json:"historyEnabled"`
	MaxSuggestions int  `json:"maxSuggestions"`
	Debug          bool `json:"debug"`
	GitEnabled     bool `json:"gitEnabled"`
	DotnetEnabled  bool `json:"dotnetEnabled"`
	CargoEnabled   bool `json:"cargoEnabled"`
}

func Default() Config {
	return Config{HistoryEnabled: true, MaxSuggestions: 8, GitEnabled: true, DotnetEnabled: true, CargoEnabled: true}
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
	if cfg.MaxSuggestions <= 0 {
		cfg.MaxSuggestions = 8
	}
	return cfg, nil
}
func DefaultPath() string {
	directory, err := os.UserConfigDir()
	if err != nil {
		return "nextcmd.json"
	}
	return filepath.Join(directory, "nextcmd", "config.json")
}
