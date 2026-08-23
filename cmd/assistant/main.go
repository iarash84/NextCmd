package assistant

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"nextcmd/internal/app"
	"nextcmd/internal/completion"
	"nextcmd/internal/config"
	"nextcmd/internal/history"
	"nextcmd/internal/terminal"
	"nextcmd/plugins/builtin"
	"nextcmd/sdk"
)

func Main() {
	debug := flag.Bool("debug", false, "enable debug logging")
	configPath := flag.String("config", config.DefaultPath(), "configuration file")
	workingDirectory := flag.String("directory", "", "initial working directory")
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "config:", err)
		os.Exit(2)
	}
	if *debug {
		cfg.Debug = true
	}
	level := slog.LevelWarn
	if cfg.Debug {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	directory, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if *workingDirectory != "" {
		directory, err = app.ResolveDirectory(directory, *workingDirectory)
		if err != nil {
			fmt.Fprintln(os.Stderr, "directory:", err)
			os.Exit(2)
		}
	}
	plugins := []sdk.Plugin{}
	for _, plugin := range builtin.All() {
		if cfg.PluginEnabled(plugin.Info().ID) {
			plugins = append(plugins, plugin)
		}
	}
	for _, plugin := range plugins {
		logger.Debug("loaded plugin", "id", plugin.Info().ID, "version", plugin.Info().Version)
	}
	engine := completion.New(plugins, cfg.MaxSuggestions, logger)
	settings := app.RuntimeSettings{
		ConfigPath:      *configPath,
		HistoryEnabled:  cfg.HistoryEnabled,
		MaxSuggestions:  cfg.MaxSuggestions,
		Debug:           cfg.Debug,
		PluginOverrides: cfg.Plugins,
	}
	program := app.New(engine, history.New(history.DefaultPath(), cfg.HistoryEnabled), terminal.New(directory), logger, directory, settings)
	if err := program.Run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
