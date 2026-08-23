package golang

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

type State struct {
	Root       string
	ModuleFile string
	WorkFile   string
	ModulePath string
	GoVersion  string
	Packages   []string
	GoFiles    []string
	Workspace  bool
	HasTests   bool
	HasVendor  bool
}

func (*Plugin) Detect(ctx context.Context, input sdk.ProjectContext) (sdk.DetectionResult, error) {
	root, moduleFile, workFile, found, err := findRoot(input.WorkingDirectory)
	if err != nil || !found {
		return sdk.DetectionResult{}, err
	}
	state, err := scanProject(ctx, root, moduleFile, workFile)
	if err != nil {
		return sdk.DetectionResult{}, err
	}
	return sdk.DetectionResult{Detected: true, Project: state, CacheFor: 2 * time.Second}, nil
}

func findRoot(directory string) (root, moduleFile, workFile string, found bool, err error) {
	directory, err = filepath.Abs(directory)
	if err != nil {
		return "", "", "", false, err
	}
	nearestModule := ""
	for current := directory; ; current = filepath.Dir(current) {
		work := filepath.Join(current, "go.work")
		if _, statErr := os.Stat(work); statErr == nil {
			return current, nearestModule, work, true, nil
		} else if !errors.Is(statErr, os.ErrNotExist) {
			return "", "", "", false, statErr
		}
		module := filepath.Join(current, "go.mod")
		if _, statErr := os.Stat(module); statErr == nil && nearestModule == "" {
			nearestModule = module
		} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return "", "", "", false, statErr
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	if nearestModule == "" {
		return "", "", "", false, nil
	}
	return filepath.Dir(nearestModule), nearestModule, "", true, nil
}

func scanProject(ctx context.Context, root, moduleFile, workFile string) (State, error) {
	state := State{
		Root:       root,
		ModuleFile: filepath.ToSlash(moduleFile),
		WorkFile:   filepath.ToSlash(workFile),
		Workspace:  workFile != "",
		HasVendor:  directoryExists(filepath.Join(root, "vendor")),
	}
	if moduleFile != "" {
		state.ModulePath, state.GoVersion = parseModule(moduleFile)
	}
	packageSet := map[string]struct{}{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if entry.IsDir() {
			if path != root && ignoredDirectory(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".go") {
			return nil
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		relative = filepath.ToSlash(relative)
		state.GoFiles = append(state.GoFiles, relative)
		state.HasTests = state.HasTests || strings.HasSuffix(strings.ToLower(entry.Name()), "_test.go")
		directory := filepath.ToSlash(filepath.Dir(relative))
		if directory == "." {
			packageSet["."] = struct{}{}
		} else {
			packageSet["./"+directory] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return State{}, err
	}
	for pkg := range packageSet {
		state.Packages = append(state.Packages, pkg)
	}
	sort.Strings(state.Packages)
	sort.Strings(state.GoFiles)
	return state, nil
}

func parseModule(path string) (modulePath, goVersion string) {
	file, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "module":
			modulePath = fields[1]
		case "go":
			goVersion = fields[1]
		}
	}
	return modulePath, goVersion
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func ignoredDirectory(name string) bool {
	switch strings.ToLower(name) {
	case ".git", ".idea", ".vscode", "bin", "node_modules", "obj", "target", "vendor":
		return true
	default:
		return false
	}
}
