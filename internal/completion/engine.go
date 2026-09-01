package completion

import (
	"context"
	"log/slog"
	"nextcmd/internal/ranking"
	"nextcmd/sdk"
	"strings"
	"sync"
	"time"
)

type cachedProject struct {
	value   map[string]any
	expires time.Time
}
type Engine struct {
	plugins []sdk.Plugin
	max     int
	logger  *slog.Logger
	mu      sync.Mutex
	cache   map[string]cachedProject
}

func New(plugins []sdk.Plugin, max int, logger *slog.Logger) *Engine {
	return &Engine{plugins: append([]sdk.Plugin(nil), plugins...), max: max, logger: logger, cache: map[string]cachedProject{}}
}
func (e *Engine) Plugins() []sdk.Plugin { return append([]sdk.Plugin(nil), e.plugins...) }
func (e *Engine) PluginForExecutable(executable string) string {
	for _, plugin := range e.plugins {
		provider, ok := plugin.(sdk.HelpProvider)
		if !ok {
			continue
		}
		for _, command := range provider.Help() {
			if strings.EqualFold(command.Command.Executable, executable) {
				return plugin.Info().ID
			}
		}
	}
	return ""
}
func (e *Engine) Help(pluginName string) (sdk.PluginInfo, []sdk.CommandHelp, bool) {
	for _, plugin := range e.plugins {
		info := plugin.Info()
		if !strings.EqualFold(pluginName, info.ID) && !strings.EqualFold(pluginName, info.Name) {
			continue
		}
		provider, ok := plugin.(sdk.HelpProvider)
		if !ok {
			return info, nil, true
		}
		return info, provider.Help(), true
	}
	return sdk.PluginInfo{}, nil, false
}
func (e *Engine) Invalidate(directory string) { e.mu.Lock(); delete(e.cache, directory); e.mu.Unlock() }
func (e *Engine) project(ctx context.Context, directory string) map[string]any {
	e.mu.Lock()
	cached, ok := e.cache[directory]
	e.mu.Unlock()
	if ok && time.Now().Before(cached.expires) {
		return cached.value
	}
	values := map[string]any{}
	ttl := 750 * time.Millisecond
	for _, plugin := range e.plugins {
		if detector, ok := plugin.(sdk.ProjectDetector); ok {
			result, err := detector.Detect(ctx, sdk.ProjectContext{WorkingDirectory: directory})
			if err != nil {
				e.logger.Debug("project detection failed", "plugin", plugin.Info().ID, "error", err)
				continue
			}
			if result.Detected {
				values[plugin.Info().ID] = result.Project
			}
			if result.CacheFor > 0 && result.CacheFor < ttl {
				ttl = result.CacheFor
			}
		}
	}
	e.mu.Lock()
	e.cache[directory] = cachedProject{values, time.Now().Add(ttl)}
	e.mu.Unlock()
	return values
}
func (e *Engine) Complete(ctx context.Context, input, directory string, previous *sdk.ExecutionResult) []sdk.Suggestion {
	if isBuiltinCommandInput(input) {
		return nil
	}
	if strings.HasPrefix(strings.TrimSpace(input), ":") {
		return ranking.Rank(input, builtinSuggestions(), 0)
	}
	projects := e.project(ctx, directory)
	all := []sdk.Suggestion{}
	for _, plugin := range e.plugins {
		project := projects[plugin.Info().ID]
		if provider, ok := plugin.(sdk.CompletionProvider); ok {
			items, err := provider.Complete(ctx, sdk.CompletionContext{Input: input, WorkingDirectory: directory, Project: project})
			e.add(&all, plugin, items, err)
		}
		if provider, ok := plugin.(sdk.BestPracticeProvider); ok {
			items, err := provider.BestPractices(ctx, sdk.CommandContext{Input: input, WorkingDirectory: directory, Project: project, Previous: previous})
			e.add(&all, plugin, items, err)
		}
		if previous != nil {
			ec := sdk.ExecutionContext{WorkingDirectory: directory, Project: project, Result: *previous}
			if previous.ExitCode == 0 {
				if provider, ok := plugin.(sdk.NextActionProvider); ok {
					items, err := provider.NextActions(ctx, ec)
					e.add(&all, plugin, items, err)
				}
			} else if provider, ok := plugin.(sdk.RecoveryProvider); ok {
				items, err := provider.Recover(ctx, ec)
				e.add(&all, plugin, items, err)
			}
		}
	}
	return ranking.Rank(input, all, e.max)
}

func isBuiltinCommandInput(input string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(input))
	return trimmed == "cd" || trimmed == ":cd" || strings.HasPrefix(trimmed, "cd ") || strings.HasPrefix(trimmed, ":cd ") || trimmed == "pwd" || trimmed == ":pwd" || trimmed == ":ls" || strings.HasPrefix(trimmed, ":ls ") || trimmed == ":mkdir" || strings.HasPrefix(trimmed, ":mkdir ") || trimmed == ":del" || strings.HasPrefix(trimmed, ":del ") || isUtilityCommandInput(trimmed)
}

func isUtilityCommandInput(trimmed string) bool {
	for _, name := range []string{":history", ":plugins", ":clear", ":config", ":which", ":version"} {
		if trimmed == name || strings.HasPrefix(trimmed, name+" ") {
			return true
		}
	}
	return false
}
func (e *Engine) add(all *[]sdk.Suggestion, plugin sdk.Plugin, items []sdk.Suggestion, err error) {
	if err != nil {
		e.logger.Debug("plugin provider failed", "plugin", plugin.Info().ID, "error", err)
		return
	}
	for i := range items {
		if items[i].Source == "" {
			items[i].Source = plugin.Info().ID
		}
	}
	*all = append(*all, items...)
}
