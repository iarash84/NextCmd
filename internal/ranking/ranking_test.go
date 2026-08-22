package ranking

import (
	"nextcmd/sdk"
	"testing"
)

func TestRankIsDeterministicAndPrefersPriority(t *testing.T) {
	input := []sdk.Suggestion{{Command: sdk.Command{Executable: "git", Args: []string{"stash"}}, Priority: 2}, {Command: sdk.Command{Executable: "git", Args: []string{"status"}}, Priority: 8}}
	got := Rank("git sta", input, 10)
	if len(got) != 2 || got[0].Command.Args[0] != "status" {
		t.Fatalf("unexpected ranking: %#v", got)
	}
}
