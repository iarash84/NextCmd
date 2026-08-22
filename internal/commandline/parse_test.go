package commandline

import "testing"

func TestParseQuotedArgument(t *testing.T) {
	got, err := Parse(`git commit -m "hello world"`, "work")
	if err != nil || got.Executable != "git" || got.Args[2] != "hello world" {
		t.Fatalf("got=%#v err=%v", got, err)
	}
}
