# Plugin Development Guide

English | [فارسی](#راهنمای-توسعه-plugin)

This guide shows how to add a compile-time plugin to NextCmd from package creation through registration, testing, configuration, documentation, and release verification. A plugin author should need the public `sdk` package, not Core implementation details.

## Architecture contract

The dependency direction is:

```text
your plugin --> sdk <-- Core
                    <-- Terminal UI
```

A plugin:

- must implement only `sdk.Plugin`, which consists of `Info()`;
- may independently implement any capability it needs;
- returns structured data and never renders terminal UI;
- does not rank final suggestions;
- does not define key bindings;
- does not invoke commands through a shell;
- must not import `nextcmd/internal/...`;
- should use only the standard library unless a dependency has a clear justification.

Core discovers optional capabilities with type assertions. A plugin is never forced to implement capabilities it does not need.

## Available capabilities

| Capability | Purpose |
|---|---|
| `sdk.CompletionProvider` | Suggestions from the current editor input and project context |
| `sdk.ProjectDetector` | Workspace detection and cached project state |
| `sdk.NextActionProvider` | Suggestions after a successful execution |
| `sdk.BestPracticeProvider` | Optional recommendations that influence ranking |
| `sdk.RecoveryProvider` | Suggestions after a failed execution |
| `sdk.HelpProvider` | Static command catalog used by `:? <plugin>` |

Implement only the interfaces that provide real value for the tool.

## Step 1: Create the package

Create a package under `plugins/<id>`. A practical layout is:

```text
plugins/acme/
├── plugin.go
├── completion.go
├── context.go
├── workflow.go
└── plugin_test.go
```

Only `plugin.go` is required. The other files are separation-of-concern conventions, not framework requirements.

## Step 2: Implement plugin metadata

```go
package acme

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin {
    return &Plugin{}
}

func (*Plugin) Info() sdk.PluginInfo {
    return sdk.PluginInfo{
        ID:          "acme",
        Name:        "Acme CLI",
        Version:     "1.0.0",
        Description: "Commands for Acme projects",
    }
}
```

The ID must be short, unique, stable, and suitable for `:? acme`. Avoid mutable global state and hidden `init()` registration.

## Step 3: Return structured suggestions

Commands are never stored as shell strings. Keep executable and arguments separate:

```go
func suggestion(args []string, title string, risk sdk.Risk, priority int) sdk.Suggestion {
    return sdk.Suggestion{
        Command: sdk.Command{
            Executable: "acme",
            Args:       append([]string(nil), args...),
        },
        Title:       title,
        Description: title,
        Reason:      "Matches the current Acme workflow",
        Kind:        sdk.Completion,
        Risk:        risk,
        Priority:    priority,
        Source:      "acme",
    }
}
```

Use `sdk.Safe`, `sdk.Mutating`, `sdk.Destructive`, or `sdk.Dangerous` accurately. Priority is plugin metadata; Core still owns final deterministic ranking.

## Step 4: Add completion

```go
func (*Plugin) Complete(_ context.Context, input sdk.CompletionContext) ([]sdk.Suggestion, error) {
    text := strings.ToLower(strings.TrimSpace(input.Input))
    if text != "" && !strings.HasPrefix("acme", text) && !strings.HasPrefix(text, "acme") {
        return nil, nil
    }

    return []sdk.Suggestion{
        suggestion([]string{"check"}, "Check the project", sdk.Safe, 80),
        suggestion([]string{"build"}, "Build the project", sdk.Mutating, 70),
        suggestion([]string{"deploy", "<environment>"}, "Deploy the project", sdk.Dangerous, 30),
    }, nil
}
```

`sdk.CompletionContext` contains the current input, working directory, optional detected project state, and history. Respect `context.Context` cancellation when completion performs I/O.

### Placeholders

If an argument contains editable data, describe it explicitly:

```go
sdk.Placeholder{
    Name:     "environment",
    ArgIndex: 1,
    Start:    0,
    End:      len("<environment>"),
}
```

The command remains editable after acceptance. The contract allows richer placeholder navigation in future UIs.

## Step 5: Add the help catalog

Implement `sdk.HelpProvider` so users can run `:? acme`:

```go
func (*Plugin) Help() []sdk.CommandHelp {
    return []sdk.CommandHelp{
        {
            Command:     sdk.Command{Executable: "acme", Args: []string{"check"}},
            Description: "Check the project",
            Risk:        sdk.Safe,
        },
        {
            Command:     sdk.Command{Executable: "acme", Args: []string{"build"}},
            Description: "Build the project",
            Risk:        sdk.Mutating,
        },
    }
}
```

The catalog contains static templates. Actual files, projects, branches, or other dynamic values belong in completion.

## Step 6: Detect a project when needed

Use a plugin-owned state type. Core stores it as `any` and returns it only to the same plugin:

```go
type State struct {
    Root       string
    ConfigFile string
}

func (*Plugin) Detect(_ context.Context, input sdk.ProjectContext) (sdk.DetectionResult, error) {
    config := filepath.Join(input.WorkingDirectory, "acme.yaml")
    if _, err := os.Stat(config); errors.Is(err, os.ErrNotExist) {
        return sdk.DetectionResult{}, nil
    } else if err != nil {
        return sdk.DetectionResult{}, fmt.Errorf("inspect Acme project: %w", err)
    }

    return sdk.DetectionResult{
        Detected: true,
        Project:  State{Root: input.WorkingDirectory, ConfigFile: config},
        CacheFor: 2 * time.Second,
    }, nil
}
```

Keep detection cheap, skip generated directories, return deterministic ordering, and choose a short reasonable cache duration. Never scan expensive trees on every keystroke.

Read the state in completion safely:

```go
state, detected := input.Project.(State)
if detected {
    // Add context-relevant suggestions.
}
```

Core commands should usually remain available even outside a detected workspace, with lower relevance and a clear reason.

## Step 7: Add workflow capabilities

Next action after successful execution:

```go
func (*Plugin) NextActions(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
    if input.Result.Command.Executable != "acme" || len(input.Result.Command.Args) == 0 {
        return nil, nil
    }
    if input.Result.Command.Args[0] == "build" {
        item := suggestion([]string{"test"}, "Test the successful build", sdk.Mutating, 85)
        item.Kind = sdk.NextAction
        return []sdk.Suggestion{item}, nil
    }
    return nil, nil
}
```

Best practice:

```go
func (*Plugin) BestPractices(_ context.Context, input sdk.CommandContext) ([]sdk.Suggestion, error) {
    item := suggestion([]string{"check"}, "Check before publishing", sdk.Safe, 60)
    item.Kind = sdk.BestPractice
    item.Reason = "Recommended before publishing artifacts"
    return []sdk.Suggestion{item}, nil
}
```

Recovery after failure:

```go
func (*Plugin) Recover(_ context.Context, input sdk.ExecutionContext) ([]sdk.Suggestion, error) {
    if input.Result.Command.Executable != "acme" {
        return nil, nil
    }
    if strings.Contains(strings.ToLower(input.Result.Stderr), "not initialized") {
        item := suggestion([]string{"init"}, "Initialize the project", sdk.Mutating, 95)
        item.Kind = sdk.Recovery
        return []sdk.Suggestion{item}, nil
    }
    return nil, nil
}
```

Providers must return errors instead of panicking. Core logs provider failures in debug mode and keeps the UI alive.

## Step 8: Isolate external processes

If detection or completion must call the external tool, inject a runner:

```go
type Runner interface {
    Run(context.Context, string, string, ...string) (string, error)
}

type Plugin struct {
    runner Runner
}

func New() *Plugin {
    return &Plugin{runner: commandRunner{}}
}

func NewWithRunner(runner Runner) *Plugin {
    return &Plugin{runner: runner}
}
```

Here the second string is the working directory. This step replaces the empty `Plugin` and simple `New` definitions from step 2; it does not add a second type with the same name. The production runner must use `exec.CommandContext(ctx, executable, args...)`, assign `cmd.Dir`, and never use `sh -c`, `cmd /c`, or a shell command string. Tests supply a fake runner and require no network account.

## Step 9: Register the plugin explicitly

Registration is the only required composition change. Import the package in `plugins/builtin/plugins.go`, add an enable flag, and append the constructor explicitly:

```go
import "nextcmd/plugins/acme"

func All(gitEnabled, dotnetEnabled, acmeEnabled bool) []sdk.Plugin {
    plugins := []sdk.Plugin{}
    if gitEnabled {
        plugins = append(plugins, git.New())
    }
    if dotnetEnabled {
        plugins = append(plugins, dotnet.New())
    }
    if acmeEnabled {
        plugins = append(plugins, acme.New())
    }
    return plugins
}
```

Do not register with `init()` and do not modify Completion Engine, Ranking, Terminal, History, or Execution for a tool-specific plugin.

## Step 10: Add configuration

Add an independent flag to `internal/config.Config` and enable it in `Default()`:

```go
type Config struct {
    // Existing fields...
    AcmeEnabled bool `json:"acmeEnabled"`
}

func Default() Config {
    return Config{
        // Existing defaults...
        AcmeEnabled: true,
    }
}
```

Pass the flag from `cmd/assistant/main.go` to `builtin.All`:

```go
plugins := builtin.All(cfg.GitEnabled, cfg.DotnetEnabled, cfg.AcmeEnabled)
```

Update `plugins/builtin/plugins_test.go` to prove that the new plugin is registered when enabled and that every plugin can be disabled independently.

If the plugin should always be built in and does not need a switch, configuration changes are optional; explicit registration is still required.

## Step 11: Test the plugin

Minimum recommended tests:

- metadata and explicit registration;
- incomplete executable prefixes, such as `ac` for `acme`;
- completion inside and outside a detected project;
- project parsing and ignored generated directories;
- deterministic dynamic argument ordering;
- placeholder metadata;
- risk values for mutating and destructive commands;
- next actions, best practices, and recovery;
- Help catalog contents;
- provider errors and context cancellation.

Construct SDK contexts directly:

```go
func TestCompletion(t *testing.T) {
    plugin := New()
    items, err := plugin.Complete(context.Background(), sdk.CompletionContext{
        Input:            "acme b",
        WorkingDirectory: t.TempDir(),
        Project:          State{Root: "project"},
    })
    if err != nil {
        t.Fatal(err)
    }
    // Assert structured commands, risk, reason, and priority.
}
```

Use `t.TempDir()` for filesystem tests and fake runners for process tests. Tests must not require GitHub, NuGet, cloud credentials, or another network service.

## Step 12: Document and verify

Add `docs/acme-plugin.md` with English first and Persian inside `<div dir="rtl">`. Update README features, configuration, plugin links, and roadmap. If command behavior depends on a current external CLI version, cite its official documentation.

Run all quality gates:

```text
make format
make vet
make test
make build-all
```

If CGO and a C compiler are available:

```text
make test-race
```

Finally, run NextCmd and verify:

```text
ac
:? acme
```

The incomplete prefix should show suggestions and the Help command should print the plugin catalog.

## Files changed checklist

For a typical configurable built-in plugin, the final change contains:

```text
plugins/acme/plugin.go
plugins/acme/completion.go          # if completion is supported
plugins/acme/context.go             # if detection is supported
plugins/acme/workflow.go            # if workflow providers are supported
plugins/acme/plugin_test.go
plugins/builtin/plugins.go
plugins/builtin/plugins_test.go
internal/config/config.go           # only for an enable/disable setting
cmd/assistant/main.go               # only to pass that setting
docs/acme-plugin.md
README.md
```

No tool-specific change should be necessary in `internal/completion`, `internal/ranking`, `internal/terminal`, `internal/execution`, or `internal/history`.

---

<div dir="rtl">

# راهنمای توسعه Plugin

این راهنما تمام مراحل افزودن یک Plugin جدید به NextCmd را از ساخت package تا registration، configuration، تست، مستندسازی و بررسی نهایی توضیح می‌دهد. توسعه‌دهنده Plugin باید فقط به package عمومی `sdk` نیاز داشته باشد و مجبور به مطالعه implementation داخلی Core نباشد.

## قرارداد معماری

جهت dependencyها به‌شکل زیر است:

```text
your plugin --> sdk <-- Core
                    <-- Terminal UI
```

یک Plugin:

- فقط باید interface پایه `sdk.Plugin` و متد `Info()` را پیاده‌سازی کند؛
- می‌تواند هر capability موردنیاز را مستقل از بقیه پیاده‌سازی کند؛
- فقط داده ساختاریافته برمی‌گرداند و UI را render نمی‌کند؛
- ranking نهایی را انجام نمی‌دهد؛
- keybinding تعریف نمی‌کند؛
- command را از طریق shell اجرا نمی‌کند؛
- نباید `nextcmd/internal/...` را import کند؛
- بهتر است فقط از Go Standard Library استفاده کند.

Core با type assertion capabilityهای اختیاری را تشخیص می‌دهد. هیچ Plugin مجبور نیست interface بزرگی را که به آن نیاز ندارد پیاده‌سازی کند.

## capabilityهای موجود

| Capability | کاربرد |
|---|---|
| `sdk.CompletionProvider` | پیشنهاد بر اساس متن editor و context پروژه |
| `sdk.ProjectDetector` | تشخیص workspace و تولید state قابل cache |
| `sdk.NextActionProvider` | پیشنهاد بعد از اجرای موفق |
| `sdk.BestPracticeProvider` | recommendation اختیاری مؤثر بر ranking |
| `sdk.RecoveryProvider` | پیشنهاد بعد از اجرای ناموفق |
| `sdk.HelpProvider` | کاتالوگ ثابت command برای `:? <plugin>` |

فقط capabilityهایی را پیاده‌سازی کنید که برای ابزار ارزش واقعی دارند.

## مرحله ۱: ساخت package

یک package در `plugins/<id>` بسازید. ساختار پیشنهادی:

```text
plugins/acme/
├── plugin.go
├── completion.go
├── context.go
├── workflow.go
└── plugin_test.go
```

فقط `plugin.go` اجباری است. فایل‌های دیگر صرفاً برای جداسازی مسئولیت‌ها پیشنهاد می‌شوند.

## مرحله ۲: تعریف metadata

```go
package acme

import "nextcmd/sdk"

type Plugin struct{}

func New() *Plugin {
    return &Plugin{}
}

func (*Plugin) Info() sdk.PluginInfo {
    return sdk.PluginInfo{
        ID:          "acme",
        Name:        "Acme CLI",
        Version:     "1.0.0",
        Description: "Commands for Acme projects",
    }
}
```

شناسه باید کوتاه، یکتا، پایدار و مناسب استفاده در `:? acme` باشد. از global state قابل تغییر و registration مخفی با `init()` استفاده نکنید.

## مرحله ۳: ساخت Suggestion

Command را هرگز به‌شکل string مربوط به shell ذخیره نکنید. executable و args باید جدا باشند:

```go
sdk.Suggestion{
    Command: sdk.Command{
        Executable: "acme",
        Args:       []string{"build"},
    },
    Title:       "Build the project",
    Description: "Build the project",
    Reason:      "Matches the current Acme workflow",
    Kind:        sdk.Completion,
    Risk:        sdk.Mutating,
    Priority:    70,
    Source:      "acme",
}
```

Risk را از بین `sdk.Safe`، `sdk.Mutating`، `sdk.Destructive` و `sdk.Dangerous` دقیق انتخاب کنید. Priority فقط metadata مربوط به Plugin است و ranking نهایی در Core انجام می‌شود.

برای آرگومان‌های قابل‌ویرایش مانند `<environment>` یک `sdk.Placeholder` با `ArgIndex` و محدوده `Start/End` اضافه کنید.

## مرحله ۴: پیاده‌سازی Completion

متد `Complete` ورودی فعلی، working directory، project state و history را دریافت می‌کند. prefix ناقص executable مانند `ac` را نیز بپذیرید و برای I/O به cancellation در `context.Context` احترام بگذارید.

Commandهای اصلی را معمولاً خارج از workspace هم نمایش دهید، اما priority وابسته به context را پایین‌تر و Reason را واضح کنید. مقادیر پویا مانند project، branch، file یا environment واقعی فقط زمانی اضافه می‌شوند که state معتبر وجود داشته باشد.

## مرحله ۵: افزودن Help

با پیاده‌سازی `sdk.HelpProvider` کاربران می‌توانند `:? acme` را اجرا کنند. کاتالوگ Help شامل templateهای ثابت command، description و risk است. مقادیر پویا در completion عادی باقی می‌مانند.

## مرحله ۶: تشخیص پروژه

state مخصوص Plugin مانند `State` را داخل همان package تعریف کنید. `Detect` باید `sdk.DetectionResult` برگرداند، عملیات سبک انجام دهد، پوشه‌های تولیدی را نادیده بگیرد، خروجی قطعی داشته باشد و `CacheFor` منطقی تعیین کند.

Core مقدار `Project` را با نوع `any` نگهداری و فقط به همان Plugin برمی‌گرداند:

```go
state, detected := input.Project.(State)
```

## مرحله ۷: راهنمای workflow

- `NextActions` فقط بعد از اجرای موفق فراخوانی می‌شود.
- `BestPractices` recommendation اختیاری ارائه می‌دهد و نباید مزاحم UX باشد.
- `Recover` بعد از شکست command فراخوانی می‌شود و فقط سناریوهای شناخته‌شده را پیشنهاد می‌دهد.

هر provider باید error برگرداند و برای flow معمول panic نکند. Core خطا را در debug mode ثبت می‌کند و UI فعال می‌ماند.

## مرحله ۸: جداسازی process خارجی

اگر Plugin باید executable خارجی را اجرا کند، یک Runner interface کوچک تزریق کنید. implementation واقعی باید از `exec.CommandContext(ctx, executable, args...)` استفاده کند و حق استفاده از `sh -c`، `cmd /c` یا shell string را ندارد. در تست‌ها از fake runner استفاده کنید تا network و account خارجی لازم نباشد.

در این حالت `Plugin` خالی و تابع `New` مرحلهٔ ۲ را با نسخهٔ دارای `Runner` جایگزین کنید؛ type دیگری با همان نام نسازید. امضای Runner باید working directory را هم دریافت کند تا اجرای process به directory سراسری یا mutable وابسته نباشد.

## مرحله ۹: registration صریح

package جدید را در `plugins/builtin/plugins.go` import و constructor را صریحاً اضافه کنید:

```go
if acmeEnabled {
    plugins = append(plugins, acme.New())
}
```

این تنها تغییر composition اجباری است. برای Plugin مخصوص یک ابزار نباید `internal/completion`، `internal/ranking`، `internal/terminal`، `internal/execution` یا `internal/history` تغییر کند.

## مرحله ۱۰: configuration

اگر Plugin قابلیت فعال/غیرفعال شدن دارد:

1. فیلد `AcmeEnabled bool` را با JSON key مناسب به `internal/config.Config` اضافه کنید.
2. مقدار پیش‌فرض را در `Default()` تعیین کنید.
3. مقدار را در `cmd/assistant/main.go` به `builtin.All` پاس دهید.
4. تست registration را برای حالت فعال و غیرفعال به‌روزرسانی کنید.

اگر Plugin همیشه built-in است، تغییر configuration اختیاری است؛ registration صریح همچنان اجباری خواهد بود.

## مرحله ۱۱: تست

حداقل موارد زیر را تست کنید:

- metadata و registration صریح؛
- prefix ناقص executable؛
- completion داخل و خارج workspace؛
- parsing پروژه و نادیده‌گرفتن پوشه‌های تولیدی؛
- ترتیب قطعی dynamic argumentها؛
- Placeholderها و Risk؛
- Next Action، Best Practice و Recovery؛
- محتوای Help؛
- error handling و cancellation.

برای filesystem از `t.TempDir()` و برای process از fake runner استفاده کنید. تست نباید به GitHub، NuGet، credential یا سرویس شبکه نیاز داشته باشد.

## مرحله ۱۲: مستندات و verification

یک سند `docs/acme-plugin.md` با بخش انگلیسی در ابتدا و بخش فارسی داخل `<div dir="rtl">` اضافه کنید. README، configuration، لینک Plugin و roadmap را نیز اصلاح کنید.

سپس اجرا کنید:

```text
make format
make vet
make test
make build-all
```

در صورت وجود CGO و C compiler:

```text
make test-race
```

در پایان برنامه را اجرا و `ac` و `:? acme` را آزمایش کنید.

## چک‌لیست فایل‌ها

یک Plugin قابل تنظیم معمولاً این تغییرات را دارد:

```text
plugins/acme/plugin.go
plugins/acme/completion.go
plugins/acme/context.go
plugins/acme/workflow.go
plugins/acme/plugin_test.go
plugins/builtin/plugins.go
plugins/builtin/plugins_test.go
internal/config/config.go
cmd/assistant/main.go
docs/acme-plugin.md
README.md
```

هیچ تغییر مخصوص ابزار نباید در packageهای عمومی Core لازم باشد.

</div>
