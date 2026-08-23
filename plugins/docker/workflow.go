package docker

import (
	"context"
	"strings"

	"nextcmd/sdk"
)

func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "docker") || len(input.Result.Command.Args) == 0 {
		return nil, nil
	}
	a := input.Result.Command.Args
	if a[0] == "build" {
		return []sdk.Suggestion{suggest([]string{"images"}, "Inspect local images", sdk.NextAction, sdk.Safe, 84, "An image was built"), suggest([]string{"run", "--rm", "<image>"}, "Run the built image", sdk.NextAction, sdk.Mutating, 88, "An image was built")}, nil
	}
	if a[0] == "compose" && len(a) > 1 && a[1] == "up" {
		return []sdk.Suggestion{suggest([]string{"compose", "ps"}, "Inspect service status", sdk.NextAction, sdk.Safe, 90, "Compose services were started"), suggest([]string{"compose", "logs"}, "Inspect service logs", sdk.NextAction, sdk.Safe, 84, "Compose services were started")}, nil
	}
	return nil, nil
}
func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
	state, _ := input.Project.(State)
	if state.Root == "" {
		return nil, nil
	}
	out := []sdk.Suggestion{}
	if state.ComposeFile != "" {
		out = append(out, suggest([]string{"compose", "config"}, "Validate Compose configuration", sdk.BestPractice, sdk.Safe, 72, "A Compose file was detected"))
	}
	if len(state.Dockerfiles) > 0 {
		out = append(out, suggest([]string{"build", "--check", "."}, "Validate Docker build checks", sdk.BestPractice, sdk.Safe, 62, "A Dockerfile was detected"))
	}
	return out, nil
}
func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
	if !strings.EqualFold(input.Result.Command.Executable, "docker") {
		return nil, nil
	}
	message := strings.ToLower(input.Result.Stdout + input.Result.Stderr)
	switch {
	case strings.Contains(message, "cannot connect to the docker daemon"), strings.Contains(message, "docker daemon is not running"):
		return []sdk.Suggestion{suggest([]string{"info"}, "Check Docker daemon connectivity", sdk.Recovery, sdk.Safe, 96, "The Docker daemon could not be reached")}, nil
	case strings.Contains(message, "no configuration file provided"), strings.Contains(message, "can't find a suitable configuration file"):
		return []sdk.Suggestion{suggest([]string{"compose", "config"}, "Locate and validate Compose configuration", sdk.Recovery, sdk.Safe, 90, "No Compose configuration was found")}, nil
	case strings.Contains(message, "no such service"):
		state, _ := input.Project.(State)
		out := []sdk.Suggestion{}
		for _, service := range state.Services {
			out = append(out, suggest([]string{"compose", "up", "-d", service}, "Start service "+service, sdk.Recovery, sdk.Mutating, 92, "Use a service discovered in the Compose file"))
		}
		return out, nil
	}
	return nil, nil
}
