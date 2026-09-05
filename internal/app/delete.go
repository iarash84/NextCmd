package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

type deleteOptions struct {
	requested string
	dryRun    bool
	permanent bool
	approved  bool
}

type deleteCounts struct {
	files, directories int
}

type deleteResult struct {
	target  deleteCandidate
	trashed string
	dryRun  bool
	counts  deleteCounts
}

type undoDelete struct {
	trashed  string
	original string
	kind     deleteKind
}

func parseDeletePath(input string) (deleteOptions, bool, error) {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	command := ""
	switch {
	case lower == ":del":
		return deleteOptions{}, true, fmt.Errorf("usage: :del [--dry-run] [--permanent] [--yes] <path>")
	case lower == ":trash":
		return deleteOptions{}, true, fmt.Errorf("usage: :trash <path>")
	case strings.HasPrefix(lower, ":del "):
		command = ":del"
	case strings.HasPrefix(lower, ":trash "):
		command = ":trash"
	default:
		return deleteOptions{}, false, nil
	}

	requested := strings.TrimSpace(trimmed[len(command):])
	if requested == "" {
		return deleteOptions{}, true, fmt.Errorf("usage: %s <path>", command)
	}
	options := deleteOptions{}
	for {
		switch {
		case strings.HasPrefix(requested, "--dry-run "):
			options.dryRun = true
			requested = strings.TrimSpace(strings.TrimPrefix(requested, "--dry-run"))
		case requested == "--dry-run":
			return deleteOptions{}, true, fmt.Errorf("usage: %s --dry-run <path>", command)
		case strings.HasPrefix(requested, "--permanent "):
			if command == ":trash" {
				return deleteOptions{}, true, fmt.Errorf(":trash does not accept --permanent")
			}
			options.permanent = true
			requested = strings.TrimSpace(strings.TrimPrefix(requested, "--permanent"))
		case requested == "--permanent":
			return deleteOptions{}, true, fmt.Errorf("usage: :del --permanent <path>")
		case strings.HasPrefix(requested, "--yes "):
			if command == ":trash" {
				return deleteOptions{}, true, fmt.Errorf(":trash does not accept --yes")
			}
			options.approved = true
			requested = strings.TrimSpace(strings.TrimPrefix(requested, "--yes"))
		case requested == "--yes":
			return deleteOptions{}, true, fmt.Errorf("usage: :del --permanent --yes <path>")
		default:
			goto parsedFlags
		}
	}

parsedFlags:
	if options.approved && !options.permanent {
		return deleteOptions{}, true, fmt.Errorf("--yes requires --permanent")
	}
	if requested == "" {
		return deleteOptions{}, true, fmt.Errorf("usage: %s <path>", command)
	}
	if requested[0] == '\'' || requested[0] == '"' {
		quote := requested[0]
		if len(requested) < 2 || requested[len(requested)-1] != quote {
			return deleteOptions{}, true, fmt.Errorf("delete path has an unclosed quote")
		}
		requested = requested[1 : len(requested)-1]
	}
	options.requested = requested
	return options, true, nil
}

func deletePath(current string, options deleteOptions, choose func([]deleteCandidate) (deleteKind, error), confirm func(deleteCandidate, deleteOptions, deleteCounts) (bool, error)) (deleteResult, error) {
	candidates, err := findDeleteCandidates(current, options.requested)
	if err != nil {
		return deleteResult{}, err
	}
	if len(candidates) == 0 {
		return deleteResult{}, fmt.Errorf("%q was not found from %s; use :ls to inspect the current directory", options.requested, current)
	}

	target := candidates[0]
	if len(candidates) > 1 {
		if choose == nil {
			return deleteResult{}, fmt.Errorf("file and directory both match; choose one")
		}
		kind, chooseErr := choose(candidates)
		if chooseErr != nil {
			return deleteResult{}, chooseErr
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
			return deleteResult{}, fmt.Errorf("invalid delete choice %q", kind)
		}
	}

	counts, err := inspectDeleteCounts(target)
	if err != nil {
		return deleteResult{}, err
	}
	if options.dryRun {
		return deleteResult{target: target, dryRun: true, counts: counts}, nil
	}
	if options.approved {
		confirm = nil
	} else if confirm == nil {
		confirm = func(deleteCandidate, deleteOptions, deleteCounts) (bool, error) { return true, nil }
	}
	if confirm != nil {
		ok, err := confirm(target, options, counts)
		if err != nil {
			return deleteResult{}, err
		}
		if !ok {
			return deleteResult{}, fmt.Errorf("cancelled")
		}
	}

	if options.permanent {
		if target.kind == deleteDirectory {
			if err := os.RemoveAll(target.path); err != nil {
				return deleteResult{}, fmt.Errorf("delete directory %q: %w", target.path, err)
			}
			return deleteResult{target: target, counts: counts}, nil
		}
		if err := os.Remove(target.path); err != nil {
			return deleteResult{}, fmt.Errorf("delete file %q: %w", target.path, err)
		}
		return deleteResult{target: target, counts: counts}, nil
	}

	trashed, err := moveToTrash(current, target)
	if err != nil {
		return deleteResult{}, err
	}
	return deleteResult{target: target, trashed: trashed, counts: counts}, nil
}

func restoreDeleted(item undoDelete) error {
	if item.trashed == "" || item.original == "" {
		return fmt.Errorf("nothing to undo")
	}
	if _, err := os.Lstat(item.original); err == nil {
		return fmt.Errorf("cannot restore because %q already exists", item.original)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect restore target %q: %w", item.original, err)
	}
	if err := os.MkdirAll(filepath.Dir(item.original), 0o755); err != nil {
		return fmt.Errorf("create restore parent: %w", err)
	}
	if err := os.Rename(item.trashed, item.original); err != nil {
		return fmt.Errorf("restore %q to %q: %w", item.trashed, item.original, err)
	}
	return nil
}

func moveToTrash(current string, target deleteCandidate) (string, error) {
	root := filepath.Join(current, ".nextcmd-trash")
	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", fmt.Errorf("create trash directory %q: %w", root, err)
	}
	name := time.Now().Format("20060102-150405.000000000") + "-" + filepath.Base(target.path)
	destination := filepath.Join(root, name)
	if err := os.Rename(target.path, destination); err != nil {
		return "", fmt.Errorf("move %q to trash: %w", target.path, err)
	}
	return filepath.Clean(destination), nil
}

func inspectDeleteCounts(target deleteCandidate) (deleteCounts, error) {
	if target.kind == deleteFile {
		return deleteCounts{files: 1}, nil
	}
	counts := deleteCounts{directories: 1}
	err := filepath.WalkDir(target.path, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == target.path {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			counts.files++
			return nil
		}
		if entry.IsDir() {
			counts.directories++
			return nil
		}
		counts.files++
		return nil
	})
	if err != nil {
		return deleteCounts{}, fmt.Errorf("inspect directory %q: %w", target.path, err)
	}
	return counts, nil
}

func confirmDelete(reader io.Reader, writer io.Writer, target deleteCandidate, options deleteOptions, counts deleteCounts) (bool, error) {
	action := "Move to trash"
	if options.permanent {
		action = "Permanently delete"
	}
	fmt.Fprintf(writer, "%s %s: %s\n", action, target.kind, target.path)
	if target.kind == deleteDirectory {
		fmt.Fprintf(writer, "Contents: %d files, %d directories\n", counts.files, counts.directories)
	}
	fmt.Fprint(writer, "Continue? [y/N]: ")
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}
		return false, io.EOF
	}
	answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}

func formatDeleteResult(result deleteResult) string {
	if result.dryRun {
		if result.target.kind == deleteDirectory {
			return fmt.Sprintf("Would delete directory: %s (%d files, %d directories)", result.target.path, result.counts.files, result.counts.directories)
		}
		return fmt.Sprintf("Would delete file: %s", result.target.path)
	}
	if result.trashed != "" {
		return fmt.Sprintf("Moved %s to trash: %s\nRun :undo to restore it.", result.target.kind, result.target.path)
	}
	return fmt.Sprintf("Deleted %s: %s", result.target.kind, result.target.path)
}

func trashRecord(result deleteResult) *undoDelete {
	if result.trashed == "" {
		return nil
	}
	return &undoDelete{trashed: result.trashed, original: result.target.path, kind: result.target.kind}
}

func quotePathArgument(path string) string {
	if path == "" {
		return path
	}
	if strings.ContainsAny(path, " \t\"'") {
		return strconv.Quote(path)
	}
	return path
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
