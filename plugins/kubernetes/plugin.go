// Package kubernetes provides context-aware kubectl suggestions.
package kubernetes

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"nextcmd/internal/execution"
	"nextcmd/sdk"
)

type Plugin struct{ runner sdk.Runner }

type State struct {
	Root       string
	Manifests  []string
	Kinds      []string
	Contexts   []string
	Namespaces []string
}

func New() *Plugin { return NewWithRunner(execution.Executor{}) }
func NewWithRunner(runner sdk.Runner) *Plugin {
	if runner == nil {
		runner = execution.Executor{}
	}
	return &Plugin{runner: runner}
}
func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "kubernetes", Name: "Kubernetes", Version: "1.0.0", Description: "Context-aware kubectl commands"}
}

func (p *Plugin) Detect(ctx context.Context, input sdk.ProjectContext) (sdk.DetectionResult, error) {
	root, manifests, kinds, namespaces, err := findManifests(ctx, input.WorkingDirectory)
	if err != nil || len(manifests) == 0 {
		return sdk.DetectionResult{}, err
	}
	state := State{Root: root, Manifests: manifests, Kinds: kinds, Namespaces: namespaces}
	state.Contexts = p.runLines(ctx, input.WorkingDirectory, "config", "get-contexts", "-o", "name")
	return sdk.DetectionResult{Detected: true, Project: state, CacheFor: 3 * time.Second}, nil
}

func (p *Plugin) runLines(ctx context.Context, directory string, args ...string) []string {
	result := p.runner.Run(ctx, sdk.Command{Executable: "kubectl", Args: args, WorkingDirectory: directory})
	if result.Err != nil {
		return nil
	}
	return uniqueLines(result.Stdout)
}

func findManifests(ctx context.Context, directory string) (string, []string, []string, []string, error) {
	root, err := filepath.Abs(directory)
	if err != nil {
		return "", nil, nil, nil, err
	}
	for current := root; ; current = filepath.Dir(current) {
		found, kinds, namespaces, scanErr := scanManifestDirectory(ctx, current)
		if scanErr != nil {
			return "", nil, nil, nil, scanErr
		}
		if len(found) > 0 {
			return current, found, kinds, namespaces, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return "", nil, nil, nil, nil
}

func scanManifestDirectory(ctx context.Context, root string) ([]string, []string, []string, error) {
	var manifests, kinds, namespaces []string
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, nil, nil, err
	}
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil, nil, nil, ctx.Err()
		default:
		}
		if entry.IsDir() || (filepath.Ext(entry.Name()) != ".yaml" && filepath.Ext(entry.Name()) != ".yml") {
			continue
		}
		path := filepath.Join(root, entry.Name())
		file, openErr := os.Open(path)
		if openErr != nil {
			continue
		}
		scanner, api, kind := bufio.NewScanner(file), false, ""
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			api = api || strings.HasPrefix(line, "apiVersion:")
			if strings.HasPrefix(line, "kind:") {
				kind = strings.TrimSpace(strings.TrimPrefix(line, "kind:"))
			}
			if strings.HasPrefix(line, "namespace:") {
				namespaces = append(namespaces, strings.TrimSpace(strings.TrimPrefix(line, "namespace:")))
			}
		}
		_ = file.Close()
		if api && kind != "" {
			manifests = append(manifests, filepath.ToSlash(entry.Name()))
			kinds = append(kinds, kind)
		}
	}
	sort.Strings(manifests)
	return manifests, unique(kinds), unique(namespaces), nil
}

type commandSpec struct {
	args     []string
	title    string
	risk     sdk.Risk
	priority int
}

var commands = []commandSpec{
	{[]string{"version", "--client"}, "Show kubectl client version", sdk.Safe, 50},
	{[]string{"cluster-info"}, "Show cluster information", sdk.Safe, 55},
	{[]string{"config", "get-contexts"}, "List configured contexts", sdk.Safe, 65},
	{[]string{"config", "current-context"}, "Show the current context", sdk.Safe, 68},
	{[]string{"config", "use-context", "<context>"}, "Switch Kubernetes context", sdk.Mutating, 72},
	{[]string{"get", "pods"}, "List pods", sdk.Safe, 88},
	{[]string{"get", "deployments"}, "List deployments", sdk.Safe, 82},
	{[]string{"describe", "<resource>"}, "Describe a resource", sdk.Safe, 72},
	{[]string{"logs", "<pod>"}, "Show pod logs", sdk.Safe, 70},
	{[]string{"apply", "-f", "<manifest>"}, "Apply a manifest", sdk.Mutating, 90},
	{[]string{"diff", "-f", "<manifest>"}, "Preview manifest changes", sdk.Safe, 94},
	{[]string{"delete", "-f", "<manifest>"}, "Delete manifest resources", sdk.Destructive, 38},
	{[]string{"rollout", "status", "deployment/<name>"}, "Watch deployment rollout", sdk.Safe, 62},
	{[]string{"rollout", "restart", "deployment/<name>"}, "Restart a deployment", sdk.Mutating, 48},
}

func (p *Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	if !matches(input.Input, "kubectl") {
		return nil, nil
	}
	state, _ := input.Project.(State)
	out := make([]sdk.Suggestion, 0, len(commands)+len(state.Manifests))
	for _, spec := range commands {
		out = append(out, suggest(spec.args, spec.title, sdk.Completion, spec.risk, spec.priority, "Matches the current kubectl input"))
	}
	fields := strings.Fields(input.Input)
	if len(fields) >= 2 {
		verb := strings.ToLower(fields[1])
		if verb == "apply" || verb == "diff" || verb == "delete" {
			for _, file := range state.Manifests {
				out = append(out, suggest([]string{verb, "-f", file}, verb+" manifest "+file, sdk.Completion, manifestRisk(verb), 99, "Kubernetes manifest discovered in the current directory"))
			}
		}
		if verb == "config" {
			for _, value := range state.Contexts {
				out = append(out, suggest([]string{"config", "use-context", value}, "Use context "+value, sdk.Completion, sdk.Mutating, 98, "Context reported by kubectl"))
			}
		}
	}
	for _, namespace := range state.Namespaces {
		out = append(out, suggest([]string{"get", "pods", "-n", namespace}, "List pods in "+namespace, sdk.Completion, sdk.Safe, 92, "Namespace reported by kubectl"))
	}
	return out, nil
}

func (*Plugin) Help() []sdk.CommandHelp {
	out := make([]sdk.CommandHelp, 0, len(commands))
	for _, c := range commands {
		out = append(out, sdk.CommandHelp{Command: sdk.Command{Executable: "kubectl", Args: append([]string(nil), c.args...)}, Description: c.title, Risk: c.risk})
	}
	return out
}
func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "kubectl") || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	if input.Result.Command.Args[0] == "apply" {
		return []sdk.Suggestion{suggest([]string{"get", "pods"}, "Inspect pods after apply", sdk.NextAction, sdk.Safe, 92, "The manifest was applied"), suggest([]string{"rollout", "status", "deployment/<name>"}, "Watch the deployment rollout", sdk.NextAction, sdk.Safe, 88, "The manifest was applied")}, nil
	}
	return nil, nil
}
func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if state.Root == "" {
		return nil, nil
	}
	return []sdk.Suggestion{suggest([]string{"diff", "-f", "<manifest>"}, "Preview changes before apply", sdk.BestPractice, sdk.Safe, 80, "Kubernetes manifests were detected")}, nil
}
func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "kubectl") {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	if strings.Contains(message, "current-context") || strings.Contains(message, "connection refused") {
		return []sdk.Suggestion{suggest([]string{"config", "get-contexts"}, "Inspect configured contexts", sdk.Recovery, sdk.Safe, 96, "kubectl could not connect to the selected cluster")}, nil
	}
	return nil, nil
}

func matches(input, executable string) bool {
	fields := strings.Fields(strings.TrimSpace(input))
	return len(fields) == 0 || strings.HasPrefix(executable, strings.ToLower(fields[0]))
}
func manifestRisk(verb string) sdk.Risk {
	if verb == "delete" {
		return sdk.Destructive
	}
	if verb == "diff" {
		return sdk.Safe
	}
	return sdk.Mutating
}
func suggest(args []string, title string, kind sdk.SuggestionKind, risk sdk.Risk, priority int, reason string) sdk.Suggestion {
	copied := append([]string(nil), args...)
	var placeholders []sdk.Placeholder
	for i, arg := range copied {
		s, e := strings.IndexByte(arg, '<'), strings.IndexByte(arg, '>')
		if s >= 0 && e > s {
			placeholders = append(placeholders, sdk.Placeholder{Name: arg[s+1 : e], ArgIndex: i, Start: s, End: e + 1})
		}
	}
	return sdk.Suggestion{Command: sdk.Command{Executable: "kubectl", Args: copied}, Title: title, Description: title, Reason: reason, Kind: kind, Risk: risk, Priority: priority, Source: "kubernetes", Placeholders: placeholders}
}
func uniqueLines(value string) []string { return unique(strings.Fields(value)) }
func unique(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range values {
		v = strings.Trim(strings.TrimSpace(v), `"'`)
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
