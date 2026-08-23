package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"nextcmd/internal/commandline"
	"nextcmd/internal/completion"
	"nextcmd/internal/execution"
	"nextcmd/internal/history"
	"nextcmd/internal/terminal"
	"nextcmd/sdk"
)

type App struct {
	engine    *completion.Engine
	executor  execution.Executor
	history   *history.Store
	ui        *terminal.UI
	logger    *slog.Logger
	directory string
}

func New(engine *completion.Engine, store *history.Store, ui *terminal.UI, logger *slog.Logger, directory string) *App {
	return &App{engine: engine, history: store, ui: ui, logger: logger, directory: directory}
}
func (a *App) Run(ctx context.Context) error {
	terminal.PrintWelcome(os.Stdout)
	var previous *sdk.ExecutionResult
	for {
		line, err := a.ui.ReadCommand(ctx, a.engine, previous)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if line == "" {
			continue
		}
		if pluginName, help := parseHelpCommand(line); help {
			printHelp(os.Stdout, a.engine, pluginName)
			continue
		}
		if isExitCommand(line) {
			return nil
		}
		if isPrintDirectoryCommand(line) {
			fmt.Println(a.directory)
			continue
		}
		if requested, handled, parseErr := parseListDirectory(line); handled {
			if parseErr != nil {
				fmt.Fprintln(os.Stderr, ":ls:", parseErr)
				continue
			}
			directory := a.directory
			if requested != "" {
				directory, parseErr = ResolveDirectory(a.directory, requested)
				if parseErr != nil {
					fmt.Fprintln(os.Stderr, ":ls:", parseErr)
					continue
				}
			}
			if listErr := printDirectoryListing(os.Stdout, directory); listErr != nil {
				fmt.Fprintln(os.Stderr, ":ls:", listErr)
			}
			continue
		}
		if requested, handled, parseErr := parseChangeDirectory(line); handled {
			if parseErr != nil {
				fmt.Fprintln(os.Stderr, "cd:", parseErr)
				continue
			}
			directory, resolveErr := ResolveDirectory(a.directory, requested)
			if resolveErr != nil {
				fmt.Fprintln(os.Stderr, "cd:", resolveErr)
				continue
			}
			oldDirectory := a.directory
			a.directory = directory
			a.ui.SetDirectory(directory)
			a.engine.Invalidate(oldDirectory)
			a.engine.Invalidate(directory)
			previous = nil
			continue
		}
		command, err := commandline.Parse(line, a.directory)
		if err != nil {
			fmt.Fprintln(os.Stderr, "parse:", err)
			continue
		}
		result := a.executor.Run(ctx, command)
		a.engine.Invalidate(a.directory)
		if result.Stdout != "" {
			fmt.Print(result.Stdout)
		}
		if result.Stderr != "" {
			fmt.Fprint(os.Stderr, result.Stderr)
		}
		if result.Err != nil {
			fmt.Fprintln(os.Stderr, "command failed:", result.Err)
		}
		terminal.PrintExecutionSummary(os.Stdout, result)
		plugin := a.engine.PluginForExecutable(command.Executable)
		if err := a.history.Append(sdk.HistoryEntry{Command: command, WorkingDirectory: a.directory, Timestamp: time.Now(), ExitCode: result.ExitCode, Duration: result.Duration, Plugin: plugin}); err != nil {
			a.logger.Debug("history write failed", "error", err)
		}
		previous = &result
	}
}

func isExitCommand(input string) bool {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "exit", "quit", ":q":
		return true
	default:
		return false
	}
}
