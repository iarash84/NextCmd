// Package terraform provides context-aware Terraform CLI suggestions.
package terraform

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"nextcmd/sdk"
)

type Plugin struct{}
type State struct {
	Root                                  string
	Files, Resources, Modules, Workspaces []string
	Initialized                           bool
}

func New() *Plugin { return &Plugin{} }
func (*Plugin) Info() sdk.PluginInfo {
	return sdk.PluginInfo{ID: "terraform", Name: "Terraform", Version: "1.0.0", Description: "Context-aware Terraform commands"}
}

func (*Plugin) Detect(ctx context.Context, input sdk.ProjectContext) (sdk.DetectionResult, error) {
	root, files, resources, modules, err := findProject(ctx, input.WorkingDirectory)
	if err != nil || len(files) == 0 {
		return sdk.DetectionResult{}, err
	}
	state := State{Root: root, Files: files, Resources: resources, Modules: modules, Initialized: directoryExists(filepath.Join(root, ".terraform"))}
	if data, readErr := os.ReadFile(filepath.Join(root, ".terraform", "environment")); readErr == nil {
		if workspace := strings.TrimSpace(string(data)); workspace != "" {
			state.Workspaces = []string{workspace}
		}
	}
	return sdk.DetectionResult{Detected: true, Project: state, CacheFor: 3 * time.Second}, nil
}

func findProject(ctx context.Context, directory string) (string, []string, []string, []string, error) {
	root, err := filepath.Abs(directory)
	if err != nil {
		return "", nil, nil, nil, err
	}
	for current := root; ; current = filepath.Dir(current) {
		entries, readErr := os.ReadDir(current)
		if readErr != nil {
			return "", nil, nil, nil, readErr
		}
		var files, resources, modules []string
		for _, entry := range entries {
			select {
			case <-ctx.Done():
				return "", nil, nil, nil, ctx.Err()
			default:
			}
			if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".tf") {
				continue
			}
			files = append(files, filepath.ToSlash(entry.Name()))
			r, m := parseTerraformFile(filepath.Join(current, entry.Name()))
			resources = append(resources, r...)
			modules = append(modules, m...)
		}
		if len(files) > 0 {
			sort.Strings(files)
			return current, files, unique(resources), unique(modules), nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return "", nil, nil, nil, nil
}
func parseTerraformFile(path string) ([]string, []string) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil
	}
	defer file.Close()
	var resources, modules []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(fields) >= 3 && fields[0] == "resource" {
			kind := strings.Trim(fields[1], `"`)
			name := strings.Trim(fields[2], `"{`)
			resources = append(resources, kind+"."+name)
		} else if len(fields) >= 2 && fields[0] == "module" {
			modules = append(modules, strings.Trim(fields[1], `"{`))
		}
	}
	return resources, modules
}
func directoryExists(path string) bool { info, err := os.Stat(path); return err == nil && info.IsDir() }

type commandSpec struct {
	args     []string
	title    string
	risk     sdk.Risk
	priority int
}

var commands = []commandSpec{
	{[]string{"version"}, "Show Terraform version", sdk.Safe, 48}, {[]string{"init"}, "Initialize the working directory", sdk.Mutating, 90},
	{[]string{"fmt", "-recursive"}, "Format Terraform files", sdk.Mutating, 82}, {[]string{"fmt", "-check", "-recursive"}, "Check Terraform formatting", sdk.Safe, 88},
	{[]string{"validate"}, "Validate the configuration", sdk.Safe, 94}, {[]string{"plan"}, "Preview infrastructure changes", sdk.Safe, 96},
	{[]string{"plan", "-out", "<plan-file>"}, "Save an execution plan", sdk.Mutating, 78}, {[]string{"show"}, "Show state or a plan", sdk.Safe, 65},
	{[]string{"apply"}, "Apply infrastructure changes", sdk.Dangerous, 52}, {[]string{"apply", "<plan-file>"}, "Apply a saved plan", sdk.Dangerous, 62},
	{[]string{"destroy"}, "Destroy managed infrastructure", sdk.Dangerous, 20}, {[]string{"output"}, "Show output values", sdk.Safe, 72},
	{[]string{"state", "list"}, "List resources in state", sdk.Safe, 75}, {[]string{"state", "show", "<resource>"}, "Show a state resource", sdk.Safe, 68},
	{[]string{"workspace", "list"}, "List workspaces", sdk.Safe, 62}, {[]string{"workspace", "select", "<workspace>"}, "Select a workspace", sdk.Mutating, 58},
}

func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	if !matches(input.Input, "terraform") {
		return nil, nil
	}
	state, _ := input.Project.(State)
	out := make([]sdk.Suggestion, 0, len(commands))
	for _, c := range commands {
		priority, reason := c.priority, "Matches the current Terraform input"
		if needsInit(c.args) && state.Root != "" && !state.Initialized {
			priority -= 15
			reason = "Run terraform init first"
		}
		out = append(out, suggest(c.args, c.title, sdk.Completion, c.risk, priority, reason))
	}
	fields := strings.Fields(input.Input)
	if len(fields) >= 2 && strings.EqualFold(fields[1], "workspace") {
		for _, v := range state.Workspaces {
			out = append(out, suggest([]string{"workspace", "select", v}, "Select workspace "+v, sdk.Completion, sdk.Mutating, 98, "Workspace reported by Terraform"))
		}
	}
	if len(fields) >= 2 && strings.EqualFold(fields[1], "state") {
		for _, v := range state.Resources {
			out = append(out, suggest([]string{"state", "show", v}, "Show state for "+v, sdk.Completion, sdk.Safe, 96, "Resource discovered in Terraform configuration"))
		}
	}
	return out, nil
}
func (*Plugin) Help() []sdk.CommandHelp {
	out := make([]sdk.CommandHelp, 0, len(commands))
	for _, c := range commands {
		out = append(out, sdk.CommandHelp{Command: sdk.Command{Executable: "terraform", Args: append([]string(nil), c.args...)}, Description: c.title, Risk: c.risk})
	}
	return out
}
func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "terraform") || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	switch input.Result.Command.Args[0] {
	case "init":
		return []sdk.Suggestion{suggest([]string{"validate"}, "Validate initialized configuration", sdk.NextAction, sdk.Safe, 94, "Terraform initialization completed"), suggest([]string{"plan"}, "Preview infrastructure changes", sdk.NextAction, sdk.Safe, 90, "Terraform initialization completed")}, nil
	case "plan":
		return []sdk.Suggestion{suggest([]string{"apply", "<plan-file>"}, "Apply the reviewed plan", sdk.NextAction, sdk.Dangerous, 55, "Review the plan before applying it")}, nil
	case "fmt":
		return []sdk.Suggestion{suggest([]string{"validate"}, "Validate formatted configuration", sdk.NextAction, sdk.Safe, 92, "Terraform files were formatted")}, nil
	}
	return nil, nil
}
func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if state.Root == "" {
		return nil, nil
	}
	return []sdk.Suggestion{suggest([]string{"fmt", "-check", "-recursive"}, "Check formatting", sdk.BestPractice, sdk.Safe, 76, "Terraform files were detected"), suggest([]string{"validate"}, "Validate configuration", sdk.BestPractice, sdk.Safe, 80, "Terraform files were detected"), suggest([]string{"plan"}, "Review a plan before apply", sdk.BestPractice, sdk.Safe, 84, "Preview infrastructure changes")}, nil
}
func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "terraform") {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	if strings.Contains(message, "terraform init") || strings.Contains(message, "initialization required") {
		return []sdk.Suggestion{suggest([]string{"init"}, "Initialize Terraform", sdk.Recovery, sdk.Mutating, 98, "Terraform reported that initialization is required")}, nil
	}
	if strings.Contains(message, "configuration is not valid") {
		return []sdk.Suggestion{suggest([]string{"validate"}, "Validate the configuration", sdk.Recovery, sdk.Safe, 94, "Terraform reported an invalid configuration")}, nil
	}
	return nil, nil
}
func matches(input, executable string) bool {
	fields := strings.Fields(strings.TrimSpace(input))
	return len(fields) == 0 || strings.HasPrefix(executable, strings.ToLower(fields[0]))
}
func needsInit(args []string) bool {
	return len(args) > 0 && args[0] != "version" && args[0] != "init" && args[0] != "fmt"
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
	return sdk.Suggestion{Command: sdk.Command{Executable: "terraform", Args: copied}, Title: title, Description: title, Reason: reason, Kind: kind, Risk: risk, Priority: priority, Source: "terraform", Placeholders: placeholders}
}
func unique(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range values {
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
