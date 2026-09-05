package app

import (
	"fmt"
	"io"
	"strings"
)

func printBuiltinHelp(writer io.Writer, name string) bool {
	key := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(name)), ":")
	help, ok := map[string]string{
		"del":     "Usage: :del [--dry-run] [--permanent] [--yes] <path>\nMoves a file or directory to .nextcmd-trash after confirmation. Use --dry-run to preview the target, --permanent to remove it without undo support, and --yes with --permanent to bypass confirmation. If both a file and directory match, NextCmd asks which one to use.",
		"trash":   "Usage: :trash <path>\nMoves a file or directory to .nextcmd-trash after confirmation. Run :undo to restore the last trashed item in this session.",
		"undo":    "Usage: :undo\nRestores the last file or directory moved to trash in this NextCmd session, unless the original path already exists.",
		"cd":      "Usage: cd <path> or :cd <path>\nChanges NextCmd's active working directory. Relative, absolute, quoted, .., and ~ paths are supported.",
		"ls":      "Usage: :ls [path]\nLists files and directories without changing the active working directory.",
		"mkdir":   "Usage: :mkdir <path>\nCreates a directory and any missing parents.",
		"history": "Usage: :history [count]\nShows recent redacted command history. Count must be between 1 and 1000.",
		"plugins": "Usage: :plugins\nLists enabled plugins.",
		"clear":   "Usage: :clear\nClears the visible terminal screen.",
		"config":  "Usage: :config\nShows effective runtime configuration.",
		"which":   "Usage: :which <command>\nLocates an executable through PATH.",
		"version": "Usage: :version\nShows version and build information.",
		"q":       "Usage: :q\nExits NextCmd.",
		"?":       "Usage: :? [command-or-plugin]\nShows general help, one built-in command, or one plugin command catalog.",
	}[key]
	if !ok {
		return false
	}
	fmt.Fprintln(writer, help)
	return true
}
