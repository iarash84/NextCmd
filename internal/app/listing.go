package app

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

func parseListDirectory(input string) (requested string, handled bool, err error) {
	trimmed := strings.TrimSpace(input)
	lower := strings.ToLower(trimmed)
	if lower == ":ls" {
		return "", true, nil
	}
	if !strings.HasPrefix(lower, ":ls ") {
		return "", false, nil
	}

	requested = strings.TrimSpace(trimmed[len(":ls"):])
	if requested == "" {
		return "", true, nil
	}
	if requested[0] == '\'' || requested[0] == '"' {
		quote := requested[0]
		if len(requested) < 2 || requested[len(requested)-1] != quote {
			return "", true, fmt.Errorf("directory path has an unclosed quote")
		}
		requested = requested[1 : len(requested)-1]
	}
	return requested, true, nil
}

func printDirectoryListing(writer io.Writer, directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("read directory %q: %w", directory, err)
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	fmt.Fprintf(writer, "Directory: %s\n", directory)
	if len(entries) == 0 {
		fmt.Fprintln(writer, "  (empty directory)")
		return nil
	}

	table := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	fmt.Fprintln(table, "TYPE\tSIZE\tNAME")
	for _, entry := range entries {
		kind := "FILE"
		size := "-"
		if entry.Type()&os.ModeSymlink != 0 {
			kind = "LINK"
		} else if entry.IsDir() {
			kind = "DIR"
		} else {
			info, infoErr := entry.Info()
			if infoErr != nil {
				return fmt.Errorf("inspect %q: %w", entry.Name(), infoErr)
			}
			size = formatFileSize(info.Size())
		}
		fmt.Fprintf(table, "%s\t%s\t%s\n", kind, size, entry.Name())
	}
	if err := table.Flush(); err != nil {
		return fmt.Errorf("write directory listing: %w", err)
	}
	return nil
}

func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	value := float64(size)
	for _, unit := range []string{"KB", "MB", "GB", "TB"} {
		value /= 1024
		if value < 1024 || unit == "TB" {
			return fmt.Sprintf("%.1f %s", value, unit)
		}
	}
	return fmt.Sprintf("%d B", size)
}
