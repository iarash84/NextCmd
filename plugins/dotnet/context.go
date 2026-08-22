package dotnet

import (
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

type Project struct {
	Path, Name, Language string
	Test                 bool
}

type State struct {
	Root      string
	Projects  []Project
	Solutions []string
	HasConfig bool
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
	return sdk.DetectionResult{Detected: len(state.Projects) > 0 || len(state.Solutions) > 0, Project: state, CacheFor: 2 * time.Second}, nil
}

func findRoot(directory string) (string, bool, error) {
	directory, err := filepath.Abs(directory)
	if err != nil {
		return "", false, err
	}
	projectFallback := ""
	for current := directory; ; current = filepath.Dir(current) {
		entries, readErr := os.ReadDir(current)
		if readErr != nil {
			return "", false, readErr
		}
		hasProject := false
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".sln" || ext == ".slnx" {
				return current, true, nil
			}
			if isProjectExtension(ext) {
				hasProject = true
			}
		}
		if hasProject && projectFallback == "" {
			projectFallback = current
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return projectFallback, projectFallback != "", nil
}

func scan(ctx context.Context, root string) (State, error) {
	state := State{Root: root}
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
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		switch {
		case isProjectExtension(ext):
			state.Projects = append(state.Projects, inspectProject(path, relative, ext))
		case ext == ".sln" || ext == ".slnx":
			state.Solutions = append(state.Solutions, relative)
		case strings.EqualFold(entry.Name(), "global.json"), strings.EqualFold(entry.Name(), "Directory.Build.props"), strings.EqualFold(entry.Name(), "Directory.Packages.props"):
			state.HasConfig = true
		}
		return nil
	})
	if err != nil && !errors.Is(err, context.Canceled) {
		return State{}, err
	}
	sort.Slice(state.Projects, func(i, j int) bool { return state.Projects[i].Path < state.Projects[j].Path })
	sort.Strings(state.Solutions)
	return state, err
}

func inspectProject(path, relative, extension string) Project {
	project := Project{Path: relative, Name: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))}
	switch extension {
	case ".fsproj":
		project.Language = "F#"
	case ".vbproj":
		project.Language = "Visual Basic"
	default:
		project.Language = "C#"
	}
	data, err := os.ReadFile(path)
	if err == nil {
		content := strings.ToLower(string(data))
		project.Test = strings.Contains(content, "<istestproject>true</istestproject>") || strings.Contains(content, "microsoft.net.test.sdk")
	}
	return project
}

func isProjectExtension(extension string) bool {
	return extension == ".csproj" || extension == ".fsproj" || extension == ".vbproj"
}
func ignoredDirectory(name string) bool {
	switch strings.ToLower(name) {
	case ".git", ".vs", ".idea", "bin", "obj", "node_modules", "target", "packages":
		return true
	default:
		return false
	}
}
