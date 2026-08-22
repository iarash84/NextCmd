package git

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nextcmd/sdk"
)

type Runner interface {
	Run(context.Context, string, ...string) (string, error)
}
type State struct {
	InRepository                bool
	Branch                      string
	Modified, Staged, Untracked bool
	Branches, Remotes, Files    []string
	HasUpstream, Ahead          bool
}

func (p *Plugin) Detect(ctx context.Context, input sdk.ProjectContext) (sdk.DetectionResult, error) {
	state, err := p.readState(ctx, input.WorkingDirectory)
	if err != nil {
		return sdk.DetectionResult{}, err
	}
	return sdk.DetectionResult{Detected: state.InRepository, Project: state, CacheFor: time.Second}, nil
}
func (p *Plugin) readState(ctx context.Context, directory string) (State, error) {
	inside, err := p.runner.Run(ctx, directory, "rev-parse", "--is-inside-work-tree")
	if err != nil || strings.TrimSpace(inside) != "true" {
		return State{}, nil
	}
	state := State{InRepository: true}
	state.Branch, _ = p.runner.Run(ctx, directory, "branch", "--show-current")
	state.Branch = strings.TrimSpace(state.Branch)
	status, err := p.runner.Run(ctx, directory, "status", "--porcelain=v1")
	if err != nil {
		return State{}, fmt.Errorf("git status: %w", err)
	}
	for _, line := range strings.Split(strings.TrimRight(status, "\r\n"), "\n") {
		if len(line) < 3 {
			continue
		}
		x, y := line[0], line[1]
		if x == '?' && y == '?' {
			state.Untracked = true
		} else {
			if x != ' ' {
				state.Staged = true
			}
			if y != ' ' {
				state.Modified = true
			}
		}
		state.Files = append(state.Files, strings.TrimSpace(line[3:]))
	}
	branches, _ := p.runner.Run(ctx, directory, "for-each-ref", "--format=%(refname:short)", "refs/heads")
	state.Branches = lines(branches)
	remotes, _ := p.runner.Run(ctx, directory, "remote")
	state.Remotes = lines(remotes)
	if _, err := p.runner.Run(ctx, directory, "rev-parse", "--abbrev-ref", "@{upstream}"); err == nil {
		state.HasUpstream = true
		ahead, _ := p.runner.Run(ctx, directory, "rev-list", "--count", "@{upstream}..HEAD")
		state.Ahead = strings.TrimSpace(ahead) != "" && strings.TrimSpace(ahead) != "0"
	}
	return state, nil
}
func lines(value string) []string {
	fields := strings.Split(strings.TrimSpace(value), "\n")
	out := make([]string, 0, len(fields))
	for _, field := range fields {
		if value := strings.TrimSpace(field); value != "" {
			out = append(out, value)
		}
	}
	return out
}
