package docker

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"nextcmd/sdk"
)

type State struct {
	Root, ComposeFile string
	Dockerfiles       []string
	Services          []string
	Targets           []string
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
		for _, name := range []string{"compose.yaml", "compose.yml", "docker-compose.yaml", "docker-compose.yml", "Dockerfile"} {
			if _, statErr := os.Stat(filepath.Join(current, name)); statErr == nil {
				return current, true, nil
			} else if !errors.Is(statErr, os.ErrNotExist) {
				return "", false, statErr
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
		if entry.IsDir() {
			continue
		}
		name, lower := entry.Name(), strings.ToLower(entry.Name())
		path := filepath.Join(root, name)
		if lower == "dockerfile" || strings.HasPrefix(lower, "dockerfile.") {
			state.Dockerfiles = append(state.Dockerfiles, filepath.ToSlash(name))
			state.Targets = append(state.Targets, parseTargets(path)...)
		}
		if state.ComposeFile == "" && (lower == "compose.yaml" || lower == "compose.yml" || lower == "docker-compose.yaml" || lower == "docker-compose.yml") {
			state.ComposeFile = filepath.ToSlash(name)
			state.Services = parseServices(path)
		}
	}
	sort.Strings(state.Dockerfiles)
	state.Services = uniqueSorted(state.Services)
	state.Targets = uniqueSorted(state.Targets)
	return state, nil
}

func parseTargets(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 4 && strings.EqualFold(fields[0], "FROM") && strings.EqualFold(fields[len(fields)-2], "AS") {
			targets = append(targets, fields[len(fields)-1])
		}
	}
	return targets
}

func parseServices(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	servicesIndent, childIndent := -1, -1
	var services []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " \t\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if servicesIndent < 0 {
			if trimmed == "services:" {
				servicesIndent = indent
			}
			continue
		}
		if indent <= servicesIndent {
			break
		}
		if childIndent < 0 {
			childIndent = indent
		}
		if indent == childIndent && strings.HasSuffix(trimmed, ":") {
			name := strings.Trim(strings.TrimSuffix(trimmed, ":"), "\"'")
			if name != "" {
				services = append(services, name)
			}
		}
	}
	return services
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
