package app

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"nextcmd/internal/buildinfo"
	"nextcmd/internal/completion"
	"nextcmd/internal/history"
)

const defaultHistoryCount = 20

type RuntimeSettings struct {
	ConfigPath      string
	HistoryEnabled  bool
	MaxSuggestions  int
	Debug           bool
	PluginOverrides map[string]bool
}

type internalCommand struct {
	name  string
	arg   string
	count int
}

func parseUtilityCommand(input string) (internalCommand, bool, error) {
	trimmed := strings.TrimSpace(input)
	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return internalCommand{}, false, nil
	}
	name := strings.ToLower(fields[0])
	switch name {
	case ":plugins", ":clear", ":config", ":version":
		if len(fields) != 1 {
			return internalCommand{name: name}, true, fmt.Errorf("%s does not accept arguments", name)
		}
		return internalCommand{name: name}, true, nil
	case ":history":
		if len(fields) > 2 {
			return internalCommand{name: name}, true, fmt.Errorf("usage: :history [count]")
		}
		count := defaultHistoryCount
		if len(fields) == 2 {
			value, err := strconv.Atoi(fields[1])
			if err != nil || value < 1 || value > 1000 {
				return internalCommand{name: name}, true, fmt.Errorf("history count must be between 1 and 1000")
			}
			count = value
		}
		return internalCommand{name: name, count: count}, true, nil
	case ":which":
		if len(fields) != 2 {
			return internalCommand{name: name}, true, fmt.Errorf("usage: :which <command>")
		}
		return internalCommand{name: name, arg: fields[1]}, true, nil
	default:
		return internalCommand{}, false, nil
	}
}

func printHistory(writer io.Writer, store *history.Store, count int) error {
	if !store.Enabled() {
		fmt.Fprintln(writer, "History is disabled in the configuration.")
		return nil
	}
	entries, err := store.Load(count)
	if err != nil {
		return fmt.Errorf("load history: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(writer, "History is empty.")
		return nil
	}
	table := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	fmt.Fprintln(table, "TIME\tEXIT\tDURATION\tPLUGIN\tDIRECTORY\tCOMMAND")
	for _, entry := range entries {
		plugin := entry.Plugin
		if plugin == "" {
			plugin = "-"
		}
		fmt.Fprintf(table, "%s\t%d\t%s\t%s\t%s\t%s\n", entry.Timestamp.Local().Format("2006-01-02 15:04:05"), entry.ExitCode, entry.Duration, plugin, entry.WorkingDirectory, entry.Command.Display())
	}
	if err := table.Flush(); err != nil {
		return fmt.Errorf("write history: %w", err)
	}
	return nil
}

func printPlugins(writer io.Writer, engine *completion.Engine) {
	plugins := engine.Plugins()
	sort.SliceStable(plugins, func(i, j int) bool { return plugins[i].Info().ID < plugins[j].Info().ID })
	if len(plugins) == 0 {
		fmt.Fprintln(writer, "No plugins are enabled.")
		return
	}
	table := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	fmt.Fprintln(table, "ID\tVERSION\tNAME\tDESCRIPTION")
	for _, plugin := range plugins {
		info := plugin.Info()
		fmt.Fprintf(table, "%s\t%s\t%s\t%s\n", info.ID, info.Version, info.Name, info.Description)
	}
	_ = table.Flush()
}

func printConfig(writer io.Writer, settings RuntimeSettings, historyPath string) {
	fmt.Fprintf(writer, "Configuration: %s\n", settings.ConfigPath)
	fmt.Fprintf(writer, "History:      %t\n", settings.HistoryEnabled)
	fmt.Fprintf(writer, "History file: %s\n", historyPath)
	fmt.Fprintf(writer, "Suggestions:  %d\n", settings.MaxSuggestions)
	fmt.Fprintf(writer, "Debug:        %t\n", settings.Debug)
	if len(settings.PluginOverrides) == 0 {
		fmt.Fprintln(writer, "Plugin overrides: none (registered plugins default to enabled)")
		return
	}
	keys := make([]string, 0, len(settings.PluginOverrides))
	for id := range settings.PluginOverrides {
		keys = append(keys, id)
	}
	sort.Strings(keys)
	fmt.Fprintln(writer, "Plugin overrides:")
	for _, id := range keys {
		fmt.Fprintf(writer, "  %-12s %t\n", id, settings.PluginOverrides[id])
	}
}

func findExecutable(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("command %q was not found in PATH", name)
	}
	absolute, err := filepath.Abs(path)
	if err == nil {
		path = absolute
	}
	return filepath.Clean(path), nil
}

func printVersion(writer io.Writer) {
	info := buildinfo.Current()
	fmt.Fprintf(writer, "NextCmd %s\n", info.Version)
	fmt.Fprintf(writer, "Go:   %s\n", info.GoVersion)
	fmt.Fprintf(writer, "Host: %s/%s\n", info.OS, info.Architecture)
	if info.Revision != "" {
		fmt.Fprintf(writer, "Commit: %s\n", info.Revision)
	}
}
