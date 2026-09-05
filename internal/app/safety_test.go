package app

import (
	"bytes"
	"testing"

	"nextcmd/sdk"
)

func TestAssessCommandSafety(t *testing.T) {
	tests := []struct {
		name    string
		command sdk.Command
		unsafe  bool
	}{
		{"recursive rm", sdk.Command{Executable: "rm", Args: []string{"-rf", "build"}}, true},
		{"separate rm flags", sdk.Command{Executable: "rm", Args: []string{"-r", "-f", "build"}}, true},
		{"plain rm", sdk.Command{Executable: "rm", Args: []string{"file.txt"}}, false},
		{"hard reset", sdk.Command{Executable: "git", Args: []string{"reset", "--hard"}}, true},
		{"soft reset", sdk.Command{Executable: "git", Args: []string{"reset", "--soft"}}, false},
		{"git clean force", sdk.Command{Executable: "git", Args: []string{"clean", "-fd"}}, true},
		{"git clean dry run", sdk.Command{Executable: "git", Args: []string{"clean", "-nd"}}, false},
		{"force push", sdk.Command{Executable: "git", Args: []string{"push", "--force-with-lease"}}, true},
		{"normal push", sdk.Command{Executable: "git", Args: []string{"push"}}, false},
		{"docker system prune", sdk.Command{Executable: "docker", Args: []string{"system", "prune"}}, true},
		{"docker ps", sdk.Command{Executable: "docker", Args: []string{"ps"}}, false},
		{"kubectl delete", sdk.Command{Executable: "kubectl", Args: []string{"delete", "-f", "deployment.yaml"}}, true},
		{"kubectl get", sdk.Command{Executable: "kubectl", Args: []string{"get", "pods"}}, false},
		{"terraform apply", sdk.Command{Executable: "terraform", Args: []string{"apply", "plan.tfplan"}}, true},
		{"terraform destroy", sdk.Command{Executable: "terraform", Args: []string{"destroy"}}, true},
		{"terraform plan", sdk.Command{Executable: "terraform", Args: []string{"plan"}}, false},
		{"PowerShell removal", sdk.Command{Executable: "Remove-Item", Args: []string{"old", "-Recurse", "-Force"}}, true},
		{"Windows removal", sdk.Command{ShellCommand: `rmdir /s /q old`}, true},
		{"dangerous shell segment", sdk.Command{ShellCommand: `echo build && rm -rf build`}, true},
		{"safe shell", sdk.Command{ShellCommand: `echo build && git status`}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := assessCommandSafety(test.command)
			if got.requiresConfirmation != test.unsafe {
				t.Fatalf("assessment=%#v", got)
			}
			if test.unsafe && got.reason == "" {
				t.Fatal("unsafe command has no reason")
			}
		})
	}
}

func TestStripSafetyApproval(t *testing.T) {
	direct, approved := stripSafetyApproval(sdk.Command{Executable: "git", Args: []string{"reset", "--hard", "--yes"}})
	if !approved || direct.Display() != "git reset --hard" {
		t.Fatalf("direct=%#v approved=%v", direct, approved)
	}
	shell, approved := stripSafetyApproval(sdk.Command{ShellCommand: "rm -rf build --yes"})
	if !approved || shell.Display() != "! rm -rf build" {
		t.Fatalf("shell=%#v approved=%v", shell, approved)
	}
	benign, approved := stripSafetyApproval(sdk.Command{Executable: "tool", Args: []string{"--yes", "value"}})
	if approved || benign.Display() != "tool --yes value" {
		t.Fatalf("benign=%#v approved=%v", benign, approved)
	}
}

func TestConfirmUnsafeCommandDefaultsToNo(t *testing.T) {
	safety := commandSafety{requiresConfirmation: true, reason: "test reason"}
	for _, answer := range []string{"\n", "no\n"} {
		var output bytes.Buffer
		confirmed, err := confirmUnsafeCommand(bytes.NewBufferString(answer), &output, sdk.Command{Executable: "rm", Args: []string{"-rf", "old"}}, safety)
		if err != nil || confirmed {
			t.Fatalf("answer=%q confirmed=%v err=%v", answer, confirmed, err)
		}
		if !bytes.Contains(output.Bytes(), []byte("rm -rf old")) || !bytes.Contains(output.Bytes(), []byte("test reason")) {
			t.Fatalf("prompt=%q", output.String())
		}
	}
}

func TestConfirmUnsafeCommandAcceptsExplicitYes(t *testing.T) {
	confirmed, err := confirmUnsafeCommand(bytes.NewBufferString("yes\n"), &bytes.Buffer{}, sdk.Command{Executable: "git", Args: []string{"reset", "--hard"}}, commandSafety{requiresConfirmation: true, reason: "test"})
	if err != nil || !confirmed {
		t.Fatalf("confirmed=%v err=%v", confirmed, err)
	}
}
