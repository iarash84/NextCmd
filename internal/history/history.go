package history

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"nextcmd/sdk"
)

type Store struct {
	path    string
	enabled bool
	mu      sync.Mutex
}

func New(path string, enabled bool) *Store { return &Store{path: path, enabled: enabled} }
func (s *Store) Enabled() bool             { return s.enabled }
func (s *Store) Path() string              { return s.path }
func DefaultPath() string {
	directory, err := os.UserConfigDir()
	if err != nil {
		return ".nextcmd-history.jsonl"
	}
	return filepath.Join(directory, "nextcmd", "history.jsonl")
}
func (s *Store) Append(entry sdk.HistoryEntry) error {
	if !s.enabled {
		return nil
	}
	entry.Command = Redact(entry.Command)
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	file, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(entry)
}
func (s *Store) Load(limit int) ([]sdk.HistoryEntry, error) {
	if !s.enabled {
		return nil, nil
	}
	file, err := os.Open(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var entries []sdk.HistoryEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry sdk.HistoryEntry
		if json.Unmarshal(scanner.Bytes(), &entry) == nil {
			entries = append(entries, entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if limit > 0 && len(entries) > limit {
		entries = entries[len(entries)-limit:]
	}
	return entries, nil
}
