package docker

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

type commandSpec struct {
	args     []string
	title    string
	risk     sdk.Risk
	priority int
}

var commands = []commandSpec{
	{[]string{"version"}, "Show Docker version", sdk.Safe, 42},
	{[]string{"info"}, "Show Docker system information", sdk.Safe, 44},
	{[]string{"build", "-t", "<image>", "."}, "Build an image", sdk.Mutating, 88},
	{[]string{"run", "--rm", "<image>"}, "Run and remove a container", sdk.Mutating, 82},
	{[]string{"ps"}, "List running containers", sdk.Safe, 84},
	{[]string{"ps", "-a"}, "List all containers", sdk.Safe, 72},
	{[]string{"images"}, "List local images", sdk.Safe, 76},
	{[]string{"logs", "<container>"}, "Show container logs", sdk.Safe, 68},
	{[]string{"exec", "-it", "<container>", "<command>"}, "Run a command in a container", sdk.Mutating, 62},
	{[]string{"stop", "<container>"}, "Stop a container", sdk.Mutating, 54},
	{[]string{"rm", "<container>"}, "Remove a container", sdk.Destructive, 34},
	{[]string{"rmi", "<image>"}, "Remove an image", sdk.Destructive, 30},
	{[]string{"pull", "<image>"}, "Pull an image", sdk.Mutating, 58},
	{[]string{"push", "<image>"}, "Push an image", sdk.Dangerous, 28},
	{[]string{"compose", "config"}, "Validate the Compose configuration", sdk.Safe, 86},
	{[]string{"compose", "up", "-d"}, "Start Compose services", sdk.Mutating, 90},
	{[]string{"compose", "ps"}, "List Compose services", sdk.Safe, 82},
	{[]string{"compose", "logs"}, "Show Compose logs", sdk.Safe, 70},
	{[]string{"compose", "build"}, "Build Compose services", sdk.Mutating, 78},
	{[]string{"compose", "down"}, "Stop and remove Compose resources", sdk.Destructive, 42},
	{[]string{"system", "df"}, "Show Docker disk usage", sdk.Safe, 50},
	{[]string{"system", "prune"}, "Remove unused Docker data", sdk.Destructive, 18},
}

func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
	trimmed := strings.TrimSpace(input.Input)
	if !matches(trimmed, "docker") {
		return nil, nil
	}
	state, _ := input.Project.(State)
	out := make([]sdk.Suggestion, 0, len(commands)+len(state.Services)*5)
	for _, spec := range commands {
		priority, reason := spec.priority, "Matches the current Docker input"
		if len(spec.args) > 0 && spec.args[0] == "compose" && state.ComposeFile == "" {
			priority -= 20
			reason = "No Compose file was detected"
		}
		out = append(out, suggest(spec.args, spec.title, sdk.Completion, spec.risk, priority, reason))
	}
	return append(out, dynamic(input.Input, state)...), nil
}

func (*Plugin) Help() []sdk.CommandHelp {
	out := make([]sdk.CommandHelp, 0, len(commands))
	for _, spec := range commands {
		out = append(out, sdk.CommandHelp{Command: sdk.Command{Executable: "docker", Args: append([]string(nil), spec.args...)}, Description: spec.title, Risk: spec.risk})
	}
	return out
}

func dynamic(input string, state State) []sdk.Suggestion {
	fields := strings.Fields(input)
	if len(fields) < 2 || !strings.EqualFold(fields[0], "docker") {
		return nil
	}
	out := []sdk.Suggestion{}
	if len(fields) >= 2 && strings.EqualFold(fields[1], "compose") {
		verb := "ps"
		if len(fields) >= 3 {
			verb = strings.ToLower(fields[2])
		}
		if verb == "up" || verb == "build" || verb == "logs" || verb == "stop" || verb == "restart" || verb == "rm" {
			for _, service := range state.Services {
				out = append(out, suggest([]string{"compose", verb, service}, verb+" service "+service, sdk.Completion, composeRisk(verb), 98, "Service discovered in "+state.ComposeFile))
			}
		}
	}
	if strings.EqualFold(fields[1], "build") {
		for _, target := range state.Targets {
			out = append(out, suggest([]string{"build", "--target", target, "."}, "Build target "+target, sdk.Completion, sdk.Mutating, 96, "Build stage discovered in a Dockerfile"))
		}
	}
	return out
}

func composeRisk(verb string) sdk.Risk {
	if verb == "rm" {
		return sdk.Destructive
	}
	if verb == "logs" {
		return sdk.Safe
	}
	return sdk.Mutating
}
func matches(input, executable string) bool {
	if input == "" {
		return true
	}
	first := strings.ToLower(strings.Fields(input)[0])
	return strings.HasPrefix(executable, first) || first == executable
}
func suggest(args []string, title string, kind sdk.SuggestionKind, risk sdk.Risk, priority int, reason string) sdk.Suggestion {
	args = append([]string(nil), args...)
	placeholders := []sdk.Placeholder{}
	for i, arg := range args {
		s, e := strings.IndexByte(arg, '<'), strings.IndexByte(arg, '>')
		if s >= 0 && e > s {
			placeholders = append(placeholders, sdk.Placeholder{Name: arg[s+1 : e], ArgIndex: i, Start: s, End: e + 1})
		}
	}
	return sdk.Suggestion{Command: sdk.Command{Executable: "docker", Args: args}, Title: title, Description: title, Reason: reason, Kind: kind, Risk: risk, Priority: priority, Source: "docker", Placeholders: placeholders}
}
