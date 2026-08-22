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
)

func Main() {
	debug := flag.Bool("debug", false, "enable debug logging")
	configPath := flag.String("config", config.DefaultPath(), "configuration file")
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
	plugins := builtin.All(cfg.GitEnabled, cfg.DotnetEnabled)
	for _, plugin := range plugins {
		logger.Debug("loaded plugin", "id", plugin.Info().ID, "version", plugin.Info().Version)
	}
	engine := completion.New(plugins, cfg.MaxSuggestions, logger)
	program := app.New(engine, history.New(history.DefaultPath(), cfg.HistoryEnabled), terminal.New(directory), logger, directory)
	if err := program.Run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
