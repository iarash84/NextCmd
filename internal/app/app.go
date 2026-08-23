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
	fmt.Println("NextCmd — arrows: select, Tab/Right: accept, Enter: run, :? help, exit/quit/:q: exit")
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
		plugin := ""
		if command.Executable == "git" {
			plugin = "git"
		}
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
