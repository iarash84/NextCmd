package dart

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"nextcmd/sdk"
)

type State struct {
	Root      string
	Pubspec   string
	Packages  []string
	DartFiles []string
	TestFiles []string
	IsFlutter bool
	HasLock   bool
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
		if _, err := os.Stat(filepath.Join(current, "pubspec.yaml")); err == nil {
			return current, true, nil
		} else if !os.IsNotExist(err) {
			return "", false, err
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", false, nil
		}
	}
}

func scan(ctx context.Context, root string) (State, error) {
	state := State{Root: root, Pubspec: filepath.Join(root, "pubspec.yaml")}
	if _, err := os.Stat(filepath.Join(root, "pubspec.lock")); err == nil {
		state.HasLock = true
	}
	if data, err := os.ReadFile(state.Pubspec); err == nil {
		state.IsFlutter = strings.Contains(strings.ToLower(string(data)), "flutter:")
	}
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
		if entry.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if ext == ".dart" {
			state.DartFiles = append(state.DartFiles, rel)
			if strings.HasSuffix(strings.ToLower(entry.Name()), "_test.dart") || strings.Contains(rel, "/test/") {
				state.TestFiles = append(state.TestFiles, rel)
			}
		}
		return nil
	})
	if err != nil {
		return State{}, err
	}
	state.DartFiles = unique(state.DartFiles)
	state.TestFiles = unique(state.TestFiles)
	return state, nil
}

func ignoredDirectory(name string) bool {
	switch strings.ToLower(name) {
	case ".git", ".dart_tool", "build", "node_modules", "vendor":
		return true
	default:
		return false
	}
}
func unique(values []string) []string {
	sort.Strings(values)
	out := values[:0]
	for _, v := range values {
		if len(out) == 0 || out[len(out)-1] != v {
			out = append(out, v)
		}
	}
	return out
}
