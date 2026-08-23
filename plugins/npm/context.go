package npm

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"nextcmd/sdk"
)

type manifest struct {
	Name            string            `json:"name"`
	Scripts         map[string]string `json:"scripts"`
	Workspaces      json.RawMessage   `json:"workspaces"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}
type State struct {
	Root, Manifest, Name              string
	Scripts, Workspaces, Dependencies []string
	HasLock, HasNodeModules           bool
}

func (*Plugin) Detect(ctx context.Context, input sdk.ProjectContext) (sdk.DetectionResult, error) {
	root, manifestPath, found, err := findRoot(input.WorkingDirectory)
	if err != nil || !found {
		return sdk.DetectionResult{}, err
	}
	state, err := scan(ctx, root, manifestPath)
	if err != nil {
		return sdk.DetectionResult{}, err
	}
	return sdk.DetectionResult{Detected: true, Project: state, CacheFor: 2 * time.Second}, nil
}
func findRoot(directory string) (string, string, bool, error) {
	directory, err := filepath.Abs(directory)
	if err != nil {
		return "", "", false, err
	}
	nearest := ""
	for current := directory; ; current = filepath.Dir(current) {
		candidate := filepath.Join(current, "package.json")
		data, readErr := os.ReadFile(candidate)
		if readErr == nil {
			if nearest == "" {
				nearest = candidate
			}
			var m manifest
			if json.Unmarshal(data, &m) == nil && len(m.Workspaces) > 0 {
				return current, candidate, true, nil
			}
		} else if !errors.Is(readErr, os.ErrNotExist) {
			return "", "", false, readErr
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	if nearest == "" {
		return "", "", false, nil
	}
	return filepath.Dir(nearest), nearest, true, nil
}
func scan(ctx context.Context, root, manifestPath string) (State, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return State{}, err
	}
	var rootManifest manifest
	if err := json.Unmarshal(data, &rootManifest); err != nil {
		return State{}, err
	}
	state := State{Root: root, Manifest: filepath.ToSlash(manifestPath), Name: rootManifest.Name, HasLock: fileExists(filepath.Join(root, "package-lock.json")), HasNodeModules: directoryExists(filepath.Join(root, "node_modules"))}
	for script := range rootManifest.Scripts {
		state.Scripts = append(state.Scripts, script)
	}
	for dependency := range rootManifest.Dependencies {
		state.Dependencies = append(state.Dependencies, dependency)
	}
	for dependency := range rootManifest.DevDependencies {
		state.Dependencies = append(state.Dependencies, dependency)
	}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if entry.IsDir() {
			if path != root && ignored(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if path == manifestPath || !strings.EqualFold(entry.Name(), "package.json") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		var child manifest
		if json.Unmarshal(data, &child) == nil && child.Name != "" {
			state.Workspaces = append(state.Workspaces, child.Name)
		}
		return nil
	})
	if err != nil {
		return State{}, err
	}
	state.Scripts = uniqueSorted(state.Scripts)
	state.Workspaces = uniqueSorted(state.Workspaces)
	state.Dependencies = uniqueSorted(state.Dependencies)
	return state, nil
}
func ignored(name string) bool {
	switch strings.ToLower(name) {
	case ".git", ".idea", ".vscode", "coverage", "dist", "build", "node_modules", "target", "vendor":
		return true
	default:
		return false
	}
}
func fileExists(path string) bool      { info, err := os.Stat(path); return err == nil && !info.IsDir() }
func directoryExists(path string) bool { info, err := os.Stat(path); return err == nil && info.IsDir() }
func uniqueSorted(values []string) []string {
	seen := map[string]struct{}{}
	for _, v := range values {
		if v != "" {
			seen[v] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for v := range seen {
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}
