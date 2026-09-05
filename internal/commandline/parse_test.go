package commandline

import "testing"

func TestParseQuotedArgument(t *testing.T) {
	got, err := Parse(`git commit -m "hello world"`, "work")
	if err != nil || got.Executable != "git" || got.Args[2] != "hello world" {
		t.Fatalf("got=%#v err=%v", got, err)
	}
}

func TestParseShellCommandPreservesSyntax(t *testing.T) {
	got, err := Parse(` ! echo "hello world" | findstr hello > output.txt `, "work")
	if err != nil || got.ShellCommand != `echo "hello world" | findstr hello > output.txt` || got.WorkingDirectory != "work" {
		t.Fatalf("got=%#v err=%v", got, err)
	}
	if got.Display() != `! echo "hello world" | findstr hello > output.txt` {
		t.Fatalf("display=%q", got.Display())
	}
}

func TestParseRejectsEmptyShellCommand(t *testing.T) {
	if _, err := Parse("!   ", "work"); err == nil {
		t.Fatal("expected an empty shell command error")
	}
}
