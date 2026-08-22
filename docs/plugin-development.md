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

---

# توسعه Plugin

یک package معمولی Go بسازید که فقط `nextcmd/sdk` را import کند. ابتدا metadata را پیاده‌سازی کنید:

```go
type Plugin struct{}

func (Plugin) Info() sdk.PluginInfo {
    return sdk.PluginInfo{ID: "tool", Name: "Tool", Version: "1.0.0"}
}
```

فقط capabilityهای موردنیاز ابزار را اضافه کنید:

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

Plugin نباید ranking، rendering، keybinding یا اجرای shell را انجام دهد و نباید به packageهای داخلی Core وابسته باشد. برای بررسی processهای خارجی، یک Runner interface کوچک تزریق کنید تا unit test بتواند fake بسازد. cancellation را از طریق `context.Context` رعایت کنید.

constructor را به‌صورت صریح در `plugins/builtin/plugins.go` ثبت کنید؛ Core نیازی به تغییر ندارد. capabilityها را مستقیماً با ساخت contextهای SDK تست کنید. نمونه کامل و کوچک در `examples/simple-plugin` قرار دارد.

هر Suggestion بهتر است command ساختاریافته، عنوان و دلیل کوتاه، kind، risk و priority پایدار داشته باشد. آرگومان‌های قابل‌ویرایش مانند `<name>` را با `Placeholder` مشخص کنید. UI فعلی command پذیرفته‌شده را قابل ویرایش نگه می‌دارد و قرارداد SDK امکان navigation پیشرفته‌تر در آینده را فراهم می‌کند.
