package dotnet

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nextcmd/sdk"
)

func TestDetectsSolutionProjectsAndTests(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "App.sln"), "")
	mustWrite(t, filepath.Join(root, "src", "App", "App.csproj"), `<Project Sdk="Microsoft.NET.Sdk"><PropertyGroup><TargetFramework>net10.0</TargetFramework></PropertyGroup></Project>`)
	mustWrite(t, filepath.Join(root, "tests", "App.Tests", "App.Tests.csproj"), `<Project Sdk="Microsoft.NET.Sdk"><ItemGroup><PackageReference Include="Microsoft.NET.Test.Sdk" /></ItemGroup></Project>`)
	mustWrite(t, filepath.Join(root, "src", "Library", "Library.fsproj"), `<Project Sdk="Microsoft.NET.Sdk" />`)
	if err := os.MkdirAll(filepath.Join(root, "src", "App", "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(root, "src", "App", "bin", "Ignored.csproj"), "")

	result, err := New().Detect(context.Background(), sdk.ProjectContext{WorkingDirectory: filepath.Join(root, "src", "App")})
	if err != nil {
		t.Fatal(err)
	}
	state, ok := result.Project.(State)
	if !result.Detected || !ok || len(state.Solutions) != 1 || len(state.Projects) != 3 {
		t.Fatalf("unexpected state: %#v", result)
	}
	if !state.Projects[2].Test && !state.Projects[1].Test && !state.Projects[0].Test {
		t.Fatal("test project was not detected")
	}
}

func TestCompletionOutsideWorkspaceOffersCreationOnly(t *testing.T) {
	got, err := New().Complete(context.Background(), sdk.CompletionContext{Input: "dotnet"})
	if err != nil {
		t.Fatal(err)
	}
	creation, build := false, false
	for _, item := range got {
		if item.Command.Display() == "dotnet new console -n <name>" {
			creation = true
		}
		if len(item.Command.Args) > 0 && item.Command.Args[0] == "build" {
			build = true
		}
	}
	if !creation || build {
		t.Fatalf("unexpected outside-workspace suggestions: %#v", got)
	}
}

func TestDynamicProjectAndTestCompletion(t *testing.T) {
	state := State{Projects: []Project{{Path: "src/App/App.csproj", Name: "App"}, {Path: "tests/App.Tests/App.Tests.csproj", Name: "App.Tests", Test: true}}}
	plugin := New()
	build, _ := plugin.Complete(context.Background(), sdk.CompletionContext{Input: "dotnet build ", Project: state})
	if !containsCommand(build, "dotnet build src/App/App.csproj") {
		t.Fatal("dynamic build project missing")
	}
	tests, _ := plugin.Complete(context.Background(), sdk.CompletionContext{Input: "dotnet test ", Project: state})
	if !containsCommand(tests, "dotnet test tests/App.Tests/App.Tests.csproj") {
		t.Fatal("dynamic test project missing")
	}
	if containsCommand(tests, "dotnet test src/App/App.csproj") {
		t.Fatal("non-test project offered for dotnet test")
	}
}

func TestNextActionAfterBuild(t *testing.T) {
	items, err := New().NextActions(context.Background(), sdk.ExecutionContext{Project: State{Projects: []Project{{Test: true}}}, Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "dotnet", Args: []string{"build"}}}})
	if err != nil || !containsCommand(items, "dotnet test --no-build") {
		t.Fatalf("items=%#v err=%v", items, err)
	}
}

func TestRecoveryWithoutWorkspace(t *testing.T) {
	items, err := New().Recover(context.Background(), sdk.ExecutionContext{Result: sdk.ExecutionResult{Command: sdk.Command{Executable: "dotnet", Args: []string{"build"}}, ExitCode: 1}})
	if err != nil || !containsCommand(items, "dotnet new sln -n <name>") {
		t.Fatalf("items=%#v err=%v", items, err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
func containsCommand(items []sdk.Suggestion, command string) bool {
	for _, item := range items {
		if item.Command.Display() == command {
			return true
		}
	}
	return false
}
