package pip

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"nextcmd/sdk"
)

type State struct {
	Root                                            string
	RequirementFiles, Packages, VirtualEnvironments []string
	HasPyProject, HasSetup                          bool
}

func (*Plugin) Detect(ctx context.Context, input sdk.ProjectContext) (sdk.DetectionResult, error) {
	root, found, err := findRoot(input.WorkingDirectory)
	if err != nil || !found {
		return sdk.DetectionResult{}, err
	}
	state, err := scan(ctx, root)
	if err != nil {
		return sdk.DetectionResult{}, err
	}
	return sdk.DetectionResult{Detected: true, Project: state, CacheFor: 2 * time.Second}, nil
}
func findRoot(directory string) (string, bool, error) {
	directory, err := filepath.Abs(directory)
	if err != nil {
		return "", false, err
	}
	for current := directory; ; current = filepath.Dir(current) {
		entries, readErr := os.ReadDir(current)
		if readErr != nil {
			return "", false, readErr
		}
		for _, entry := range entries {
			lower := strings.ToLower(entry.Name())
			if !entry.IsDir() && (lower == "pyproject.toml" || lower == "setup.py" || lower == "setup.cfg" || lower == "pipfile" || strings.HasPrefix(lower, "requirements") && strings.HasSuffix(lower, ".txt")) {
				return current, true, nil
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", false, nil
		}
	}
}
func scan(ctx context.Context, root string) (State, error) {
	state := State{Root: root}
	entries, err := os.ReadDir(root)
	if err != nil {
		return State{}, err
	}
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return State{}, ctx.Err()
		default:
		}
		name, lower := entry.Name(), strings.ToLower(entry.Name())
		if entry.IsDir() {
			if lower == ".venv" || lower == "venv" || lower == "env" {
				state.VirtualEnvironments = append(state.VirtualEnvironments, filepath.ToSlash(name))
			}
			continue
		}
		switch lower {
		case "pyproject.toml":
			state.HasPyProject = true
		case "setup.py", "setup.cfg":
			state.HasSetup = true
		}
		if strings.HasPrefix(lower, "requirements") && strings.HasSuffix(lower, ".txt") {
			state.RequirementFiles = append(state.RequirementFiles, filepath.ToSlash(name))
			state.Packages = append(state.Packages, parseRequirements(filepath.Join(root, name))...)
		}
	}
	state.RequirementFiles = uniqueSorted(state.RequirementFiles)
	state.Packages = uniqueSorted(state.Packages)
	sort.Strings(state.VirtualEnvironments)
	return state, nil
}
func parseRequirements(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	out := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(strings.SplitN(scanner.Text(), "#", 2)[0])
		if line == "" || strings.HasPrefix(line, "-") {
			continue
		}
		if name, _, direct := strings.Cut(line, " @ "); direct {
			name = strings.TrimSpace(name)
			if name != "" {
				out = append(out, name)
			}
			continue
		}
		if strings.Contains(line, "://") || strings.HasPrefix(line, ".") || strings.HasPrefix(line, "/") || strings.HasPrefix(line, `\`) {
			continue
		}
		name := line
		for _, separator := range []string{"===", "==", ">=", "<=", "~=", "!=", ">", "<", "[", ";"} {
			if index := strings.Index(name, separator); index >= 0 {
				name = name[:index]
				break
			}
		}
		name = strings.TrimSpace(name)
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}
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
