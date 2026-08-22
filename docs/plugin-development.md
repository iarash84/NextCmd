# Plugin development

Create a normal Go package that imports only `nextcmd/sdk`. Implement metadata first:

```go
type Plugin struct{}

func (Plugin) Info() sdk.PluginInfo {
    return sdk.PluginInfo{ID: "tool", Name: "Tool", Version: "1.0.0"}
}
```

Add only capabilities the tool needs:

```go
func (Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
    return []sdk.Suggestion{{
        Command:  sdk.Command{Executable: "tool", Args: []string{"check"}},
        Title:    "Check the project",
        Reason:   "Matches this project",
        Kind:     sdk.Completion,
        Risk:     sdk.Safe,
        Priority: 70,
    }}, nil
}
```

Do not rank, render, define keys, execute through a shell, or depend on Core packages. For external process inspection, inject a small runner interface so unit tests can use a fake. Respect cancellation through the supplied `context.Context`.

Register the constructor explicitly in `plugins/builtin/plugins.go`. Core needs no change. Test capabilities directly by constructing SDK contexts, for example `sdk.CompletionContext{Input: "tool c", WorkingDirectory: t.TempDir()}`. A complete minimal implementation is in `examples/simple-plugin`.

Suggestions should provide a structured command, concise title and reason, kind, risk, and stable priority. Mark editable arguments such as `<name>` with `Placeholder`; the current UI keeps the accepted command editable and the contract permits richer placeholder navigation later.
