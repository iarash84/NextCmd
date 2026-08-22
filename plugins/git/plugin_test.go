package git

import (
	"context"
	"errors"
	"nextcmd/sdk"
	"strings"
	"testing"
)

type fakeRunner map[string]string

func (f fakeRunner) Run(_ context.Context, _ string, args ...string) (string, error) {
	key := strings.Join(args, " ")
	value, ok := f[key]
	if !ok {
		return "", errors.New("not configured")
	}
	return value, nil
}
func repositoryRunner() fakeRunner {
	return fakeRunner{"rev-parse --is-inside-work-tree": "true\n", "branch --show-current": "main\n", "status --porcelain=v1": " M main.go\nA  staged.go\n?? new.go\n", "for-each-ref --format=%(refname:short) refs/heads": "main\ndevelop\n", "remote": "origin\n", "rev-parse --abbrev-ref @{upstream}": "origin/main\n", "rev-list --count @{upstream}..HEAD": "1\n"}
}
func TestDetectParsesRepositoryState(t *testing.T) {
	p := NewWithRunner(repositoryRunner())
	result, err := p.Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: "x"})
	state := result.Project.(State)
	if err != nil || !state.Modified || !state.Staged || !state.Untracked || !state.Ahead || len(state.Branches) != 2 {
		t.Fatalf("state=%#v err=%v", state, err)
	}
}
func TestCompletionOutsideRepository(t *testing.T) {
	p := NewWithRunner(fakeRunner{})
	got, err := p.Complete(context.Background(), sdk.CompletionContext{Input: "git"})
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range got {
		if s.Command.Args[0] != "init" && s.Command.Args[0] != "clone" {
			t.Fatalf("repository command leaked: %s", s.Command.Display())
		}
	}
}
func TestDynamicBranchCompletion(t *testing.T) {
	p := New()
	got, _ := p.Complete(context.Background(), sdk.CompletionContext{Input: "git switch ", Project: State{InRepository: true, Branches: []string{"main", "feature/auth"}}})
	found := false
	for _, s := range got {
		if s.Command.Display() == "git switch feature/auth" {
			found = true
		}
	}
	if !found {
		t.Fatal("dynamic branch missing")
	}
}

func TestBranchCreationTemplates(t *testing.T) {
	p := New()
	got, _ := p.Complete(context.Background(), sdk.CompletionContext{Input: "git switch -c", Project: State{InRepository: true}})
	want := map[string]bool{
		"git switch -c feature/<name>":    false,
		"git switch -c bugfix/<name>":     false,
		"git switch -c hotfix/<name>":     false,
		"git switch -c release/<version>": false,
	}
	for _, suggestion := range got {
		if _, ok := want[suggestion.Command.Display()]; ok {
			want[suggestion.Command.Display()] = len(suggestion.Placeholders) == 1
		}
	}
	for command, found := range want {
		if !found {
			t.Errorf("branch template missing or has no placeholder: %s", command)
		}
	}
}

func TestDynamicBranchDeletePreservesFlag(t *testing.T) {
	p := New()
	got, _ := p.Complete(context.Background(), sdk.CompletionContext{Input: "git branch -d ", Project: State{InRepository: true, Branches: []string{"feature/done"}}})
	for _, suggestion := range got {
		if suggestion.Command.Display() == "git branch -d feature/done" && suggestion.Risk == sdk.Destructive {
			return
		}
	}
	t.Fatal("dynamic branch deletion suggestion missing")
}

func TestExpandedCommandCatalog(t *testing.T) {
	p := New()
	got, _ := p.Complete(context.Background(), sdk.CompletionContext{Input: "git", Project: State{InRepository: true}})
	want := map[string]bool{"show": false, "tag": false, "cherry-pick": false, "revert": false, "worktree": false, "submodule": false, "grep": false}
	for _, suggestion := range got {
		if len(suggestion.Command.Args) > 0 {
			if _, ok := want[suggestion.Command.Args[0]]; ok {
				want[suggestion.Command.Args[0]] = true
			}
		}
	}
	for command, found := range want {
		if !found {
			t.Errorf("expanded command missing: git %s", command)
		}
	}
}
func TestNextActionAndBestPractice(t *testing.T) {
	p := New()
	state := State{InRepository: true, Staged: true, HasUpstream: true}
	next, _ := p.NextActions(context.Background(), sdk.ExecutionContext{Project: state, Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "git", Args: []string{"add", "."}}}})
	if len(next) < 3 || next[0].Command.Display() != "git diff --cached" {
		t.Fatalf("next=%#v", next)
	}
	practices, _ := p.BestPractices(context.Background(), sdk.CommandContext{Project: state})
	if len(practices) != 1 || practices[0].Kind != sdk.BestPractice {
		t.Fatalf("practices=%#v", practices)
	}
}
