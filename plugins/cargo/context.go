package cargo

import (
	"bufio"
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"nextcmd/sdk"
)

type Package struct {
	Name     string
	Manifest string
}

type State struct {
	Root      string
	Manifest  string
	Packages  []Package
	Features  []string
	Workspace bool
	HasLock   bool
}

func (*Plugin) Detect(ctx context.Context, input sdk.ProjectContext) (sdk.DetectionResult, error) {
	root, manifest, found, err := findRoot(input.WorkingDirectory)
	if err != nil || !found {
		return sdk.DetectionResult{}, err
	}
	state, err := scanWorkspace(ctx, root, manifest)
	if err != nil {
		return sdk.DetectionResult{}, err
	}
	return sdk.DetectionResult{Detected: true, Project: state, CacheFor: 2 * time.Second}, nil
}

func findRoot(directory string) (root, manifest string, found bool, err error) {
	directory, err = filepath.Abs(directory)
	if err != nil {
		return "", "", false, err
	}
	nearest := ""
	for current := directory; ; current = filepath.Dir(current) {
		candidate := filepath.Join(current, "Cargo.toml")
		data, readErr := os.ReadFile(candidate)
		if readErr == nil {
			if nearest == "" {
				nearest = candidate
			}
			if hasSection(string(data), "workspace") {
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

func scanWorkspace(ctx context.Context, root, manifest string) (State, error) {
	state := State{Root: root, Manifest: filepath.ToSlash(manifest), HasLock: fileExists(filepath.Join(root, "Cargo.lock"))}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if entry.IsDir() && path != root && ignoredDirectory(entry.Name()) {
			return filepath.SkipDir
		}
		if entry.IsDir() || !strings.EqualFold(entry.Name(), "Cargo.toml") {
			return nil
		}
		name, features, workspace, parseErr := parseManifest(path)
		if parseErr != nil {
			return parseErr
		}
		state.Workspace = state.Workspace || workspace
		state.Features = append(state.Features, features...)
		if name != "" {
			relative, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return relErr
			}
			state.Packages = append(state.Packages, Package{Name: name, Manifest: filepath.ToSlash(relative)})
		}
		return nil
	})
	if err != nil {
		return State{}, err
	}
	sort.Slice(state.Packages, func(i, j int) bool { return state.Packages[i].Name < state.Packages[j].Name })
	state.Features = uniqueSorted(state.Features)
	return state, nil
}

func parseManifest(path string) (name string, features []string, workspace bool, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil, false, err
	}
	defer file.Close()
	section := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(stripComment(scanner.Text()))
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.Trim(line, "[]"))
			workspace = workspace || section == "workspace"
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if section == "package" && key == "name" {
			name = unquote(value)
		}
		if section == "features" && key != "default" {
			features = append(features, unquote(key))
		}
	}
	return name, features, workspace, scanner.Err()
}

func stripComment(line string) string {
	quoted := false
	for i, char := range line {
		if char == '"' {
			quoted = !quoted
		}
		if char == '#' && !quoted {
			return line[:i]
		}
	}
	return line
}

func unquote(value string) string {
	return strings.Trim(strings.TrimSpace(value), "\"'")
}

func hasSection(content, section string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(stripComment(line)) == "["+section+"]" {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ignoredDirectory(name string) bool {
	switch strings.ToLower(name) {
	case ".git", ".idea", ".vscode", "node_modules", "target", "vendor":
		return true
	default:
		return false
	}
}

func uniqueSorted(values []string) []string {
	seen := map[string]struct{}{}
	for _, value := range values {
		if value != "" {
			seen[value] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for value := range seen {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
