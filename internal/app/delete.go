package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type deleteKind string

const (
	deleteFile      deleteKind = "file"
	deleteDirectory deleteKind = "directory"
)

type deleteCandidate struct {
	path string
	kind deleteKind
}

func parseDeletePath(input string) (requested string, handled bool, err error) {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	if lower == ":del" {
		return "", true, fmt.Errorf("usage: :del <path>")
	}
	if !strings.HasPrefix(lower, ":del ") {
		return "", false, nil
	}
	requested = strings.TrimSpace(trimmed[len(":del"):])
	if requested == "" {
		return "", true, fmt.Errorf("usage: :del <path>")
	}
	if requested[0] == '\'' || requested[0] == '"' {
		quote := requested[0]
		if len(requested) < 2 || requested[len(requested)-1] != quote {
			return "", true, fmt.Errorf("delete path has an unclosed quote")
		}
		requested = requested[1 : len(requested)-1]
	}
	return requested, true, nil
}

func deletePath(current, requested string, choose func([]deleteCandidate) (deleteKind, error)) (deleteCandidate, error) {
	candidates, err := findDeleteCandidates(current, requested)
	if err != nil {
		return deleteCandidate{}, err
	}
	if len(candidates) == 0 {
		return deleteCandidate{}, fmt.Errorf("path %q does not exist", requested)
	}

	target := candidates[0]
	if len(candidates) > 1 {
		kind, chooseErr := choose(candidates)
		if chooseErr != nil {
			return deleteCandidate{}, chooseErr
		}
		found := false
		for _, candidate := range candidates {
			if candidate.kind == kind {
				target = candidate
				found = true
				break
			}
		}
		if !found {
			return deleteCandidate{}, fmt.Errorf("invalid delete choice %q", kind)
		}
	}

	if target.kind == deleteDirectory {
		if err := os.RemoveAll(target.path); err != nil {
			return deleteCandidate{}, fmt.Errorf("delete directory %q: %w", target.path, err)
		}
		return target, nil
	}
	if err := os.Remove(target.path); err != nil {
		return deleteCandidate{}, fmt.Errorf("delete file %q: %w", target.path, err)
	}
	return target, nil
}

func findDeleteCandidates(current, requested string) ([]deleteCandidate, error) {
	resolved, err := resolveDeletePath(current, requested)
	if err != nil {
		return nil, err
	}
	if candidate, ok, statErr := inspectDeleteCandidate(resolved); statErr != nil {
		return nil, statErr
	} else if ok {
		return []deleteCandidate{candidate}, nil
	}

	parent := filepath.Dir(resolved)
	base := filepath.Base(resolved)
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil, fmt.Errorf("read directory %q: %w", parent, err)
	}
	candidates := []deleteCandidate{}
	kinds := map[deleteKind]bool{}
	for _, entry := range entries {
		if !strings.EqualFold(entry.Name(), base) {
			continue
		}
		candidate, ok, statErr := inspectDeleteCandidate(filepath.Join(parent, entry.Name()))
		if statErr != nil {
			return nil, statErr
		}
		if ok {
			candidates = append(candidates, candidate)
			kinds[candidate.kind] = true
		}
	}
	if len(candidates) > 1 && kinds[deleteFile] && kinds[deleteDirectory] {
		return candidates, nil
	}
	return nil, nil
}

func resolveDeletePath(current, requested string) (string, error) {
	requested = strings.TrimSpace(requested)
	if requested == "" {
		return "", fmt.Errorf("usage: :del <path>")
	}
	if requested == "~" || strings.HasPrefix(requested, "~/") || strings.HasPrefix(requested, `~\`) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve home directory: %w", err)
		}
		if requested == "~" {
			requested = home
		} else {
			requested = filepath.Join(home, requested[2:])
		}
	}
	if !filepath.IsAbs(requested) {
		requested = filepath.Join(current, requested)
	}
	resolved, err := filepath.Abs(requested)
	if err != nil {
		return "", fmt.Errorf("resolve delete path: %w", err)
	}
	return filepath.Clean(resolved), nil
}

func inspectDeleteCandidate(path string) (deleteCandidate, bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return deleteCandidate{}, false, nil
		}
		return deleteCandidate{}, false, fmt.Errorf("inspect %q: %w", path, err)
	}
	kind := deleteFile
	if info.IsDir() {
		kind = deleteDirectory
	}
	return deleteCandidate{path: filepath.Clean(path), kind: kind}, true, nil
}

func promptDeleteKind(reader io.Reader, writer io.Writer, candidates []deleteCandidate) (deleteKind, error) {
	fmt.Fprintln(writer, "Both a file and a directory match this path.")
	for _, candidate := range candidates {
		fmt.Fprintf(writer, "  %s: %s\n", candidate.kind, candidate.path)
	}
	fmt.Fprint(writer, "Delete file or directory? [f/d]: ")

	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
	case "f", "file":
		return deleteFile, nil
	case "d", "dir", "directory":
		return deleteDirectory, nil
	default:
		return "", fmt.Errorf("expected f/file or d/directory")
	}
}
