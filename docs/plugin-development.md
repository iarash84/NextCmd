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
- returns structured fields such as `Command`, `Title`, `Reason`, and `Risk`; Core passes those fields to the terminal UI, so the plugin never prints, colors, or lays out terminal content itself;
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

Add `docs/acme-plugin.md` with English first and Persian inside `<div dir="rtl" align="right">`. Keep code blocks left-to-right. Update README features, configuration, plugin links, and roadmap. If command behavior depends on a current external CLI version, cite its official documentation.

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

<div dir="rtl" align="right">

# راهنمای توسعهٔ افزونه

این راهنما مراحل افزودن یک افزونهٔ جدید به NextCmd را از ساخت بسته تا ثبت در برنامه، تنظیمات، تست و مستندسازی توضیح می‌دهد. برای ساخت یک افزونهٔ معمولی فقط باید قراردادهای بستهٔ عمومی `sdk` را بشناسید و نیازی به مطالعه یا تغییر جزئیات داخلی هسته ندارید.

## قرارداد معماری

جهت وابستگی‌ها به‌شکل زیر است:

<div dir="ltr" align="left">

```text
your plugin --> sdk <-- Core
                    <-- Terminal UI
```

</div>

هر افزونه باید این قواعد را رعایت کند:

- رابط پایهٔ `sdk.Plugin` و متد `Info()` را پیاده‌سازی کند تا هسته بتواند افزونه را شناسایی کند؛
- فقط قابلیت‌هایی را پیاده‌سازی کند که واقعاً به آن‌ها نیاز دارد؛
- پیشنهاد را به‌شکل اطلاعات ساختاریافته برگرداند. برای مثال، افزونه مقدارهای `Command`، `Title`، `Reason` و `Risk` را می‌سازد و تحویل هسته می‌دهد. افزونه نباید متن را رنگ‌آمیزی کند، سطر پایانه بسازد یا چیزی مستقیماً روی صفحه چاپ کند؛ این کارها فقط بر عهدهٔ رابط پایانه است؛
- امتیاز اولیه و میزان ارتباط پیشنهاد را مشخص کند، اما مرتب‌سازی نهایی را به هسته بسپارد؛
- کلیدهای صفحه‌کلید یا رفتار ویرایشگر را تعریف نکند؛
- دستورها را از طریق پوسته‌هایی مانند `sh` یا `cmd` اجرا نکند؛
- هیچ بسته‌ای از مسیر `nextcmd/internal/...` وارد نکند؛
- تا جای ممکن فقط از کتابخانهٔ استاندارد Go استفاده کند.

هسته با بررسی رابط‌های پیاده‌سازی‌شده تشخیص می‌دهد افزونه چه قابلیت‌هایی دارد. بنابراین افزونه‌ای که فقط تکمیل دستور ارائه می‌دهد، مجبور نیست قابلیت تشخیص پروژه یا بازیابی پس از خطا را نیز پیاده‌سازی کند.

## قابلیت‌های موجود

<table>
<thead>
<tr><th dir="ltr" align="left">Capability</th><th dir="rtl" align="right">کاربرد</th></tr>
</thead>
<tbody>
<tr><td dir="ltr" align="left"><code>sdk.CompletionProvider</code></td><td dir="rtl" align="right">پیشنهاد دستور براساس متن فعلی و وضعیت پروژه</td></tr>
<tr><td dir="ltr" align="left"><code>sdk.ProjectDetector</code></td><td dir="rtl" align="right">تشخیص پروژه و ساخت وضعیت قابل نگهداری در حافظهٔ موقت</td></tr>
<tr><td dir="ltr" align="left"><code>sdk.NextActionProvider</code></td><td dir="rtl" align="right">پیشنهاد بعد از اجرای موفق</td></tr>
<tr><td dir="ltr" align="left"><code>sdk.BestPracticeProvider</code></td><td dir="rtl" align="right">پیشنهاد اختیاری برای انجام کار به روش بهتر</td></tr>
<tr><td dir="ltr" align="left"><code>sdk.RecoveryProvider</code></td><td dir="rtl" align="right">پیشنهاد راه‌حل پس از اجرای ناموفق</td></tr>
<tr><td dir="ltr" align="left"><code>sdk.HelpProvider</code></td><td dir="rtl" align="right">فهرست ثابت دستورها برای <code dir="ltr">:? &lt;plugin&gt;</code></td></tr>
</tbody>
</table>

فقط رابط‌هایی را پیاده‌سازی کنید که برای ابزار شما کاربرد واقعی دارند.

## مرحلهٔ ۱: ساخت بسته

یک بسته در مسیر `plugins/<id>` بسازید. مقدار `<id>` شناسهٔ کوتاه افزونه، برای مثال `acme`، است. ساختار پیشنهادی:

<div dir="ltr" align="left">

```text
plugins/acme/
├── plugin.go
├── completion.go
├── context.go
├── workflow.go
└── plugin_test.go
```

</div>

فقط فایل `plugin.go` اجباری است. فایل‌های دیگر الزام چارچوب نیستند و صرفاً کمک می‌کنند منطق تکمیل، تشخیص پروژه و گردش کار از هم جدا و خوانا بمانند.

## مرحلهٔ ۲: معرفی افزونه

<div dir="ltr" align="left">

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

</div>

شناسه باید کوتاه، یکتا و پایدار باشد، زیرا کاربر آن را در دستوری مانند `:? acme` به کار می‌برد. افزونه را با `init()` به‌صورت مخفی ثبت نکنید و وضعیت قابل‌تغییر سراسری نسازید.

## مرحلهٔ ۳: ساخت پیشنهاد

کل دستور را در یک رشتهٔ مخصوص پوسته ذخیره نکنید. نام فایل اجرایی و آرگومان‌ها باید جدا باشند:

<div dir="ltr" align="left">

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

</div>

مقدار خطر را با دقت انتخاب کنید: `sdk.Safe` برای دستورهای فقط‌خواندنی، `sdk.Mutating` برای تغییرات عادی، `sdk.Destructive` برای حذف یا بازنویسی اطلاعات و `sdk.Dangerous` برای عملیات بسیار پرخطر است. `Priority` فقط اهمیت پیشنهادی افزونه را بیان می‌کند؛ هسته پس از ترکیب همهٔ پیشنهادها ترتیب نهایی را تعیین می‌کند.

برای آرگومان‌های قابل‌ویرایش مانند `<environment>` یک `sdk.Placeholder` اضافه کنید. `ArgIndex` شمارهٔ آرگومان و `Start` و `End` محدودهٔ متن قابل‌جایگزینی را مشخص می‌کنند.

## مرحلهٔ ۴: پیاده‌سازی تکمیل دستور

متد `Complete` متن فعلی ویرایشگر، پوشهٔ کاری، وضعیت تشخیص‌داده‌شدهٔ پروژه و تاریخچه را دریافت می‌کند. نام ناقص ابزار مانند `ac` را هم بپذیرید تا پیشنهادها پیش از کامل‌شدن `acme` ظاهر شوند. اگر این متد فایل می‌خواند یا برنامه‌ای اجرا می‌کند، لغو درخواست از طریق `context.Context` را رعایت کند.

دستورهای عمومی را معمولاً بیرون از پروژهٔ شناسایی‌شده هم نمایش دهید، اما امتیاز آن‌ها را کمتر کنید و در `Reason` توضیح دهید که ممکن است انتخاب مسیر پروژه لازم باشد. مقدارهای پویا، مانند نام پروژه، شاخه، فایل یا محیط واقعی، فقط زمانی پیشنهاد شوند که وضعیت معتبر پروژه در دسترس باشد.

## مرحلهٔ ۵: افزودن راهنما

با پیاده‌سازی `sdk.HelpProvider` کاربر می‌تواند دستور `:? acme` را اجرا کند و فهرست دستورهای افزونه را ببیند. این فهرست فقط الگوهای ثابت دستور، توضیح و میزان خطر را نگه می‌دارد. مقدارهای وابسته به پروژه باید همچنان از طریق تکمیل عادی ارائه شوند.

## مرحلهٔ ۶: تشخیص پروژه

نوعی مانند `State` را در بستهٔ افزونه تعریف کنید تا اطلاعات مخصوص همان ابزار را نگه دارد. متد `Detect` باید یک `sdk.DetectionResult` برگرداند. این متد باید سریع باشد، پوشه‌های تولیدشده و سنگین را نادیده بگیرد، نتایج را همیشه با ترتیب ثابت تولید کند و در `CacheFor` مدت مناسبی برای نگهداری نتیجه تعیین کند.

هسته مقدار `Project` را بدون شناخت نوع داخلی آن نگهداری می‌کند و هنگام فراخوانی بعدی همان افزونه، آن را برمی‌گرداند. افزونه می‌تواند مقدار را با بررسی امن نوع دریافت کند:

<div dir="ltr" align="left">

```go
state, detected := input.Project.(State)
```

</div>

## مرحلهٔ ۷: پیشنهادهای پس از اجرا

- `NextActions` پس از اجرای موفق فراخوانی می‌شود و گام منطقی بعدی را پیشنهاد می‌دهد؛ برای مثال، پس از ساخت موفق پروژه می‌تواند اجرای تست را پیشنهاد دهد.
- `BestPractices` روش بهتر یا بررسی تکمیلی را پیشنهاد می‌دهد. این پیشنهاد باید اختیاری باشد و کاربر را با هشدارهای دائمی آزار ندهد.
- `Recover` پس از شکست دستور فراخوانی می‌شود و برای خطاهای شناخته‌شده راه‌حل مشخص ارائه می‌دهد.

هر فراهم‌کننده باید خطا را برگرداند و برای وضعیت‌های عادی از `panic` استفاده نکند. هسته خطا را در حالت اشکال‌زدایی ثبت می‌کند و رابط کاربری همچنان فعال می‌ماند.

## مرحلهٔ ۸: جداسازی اجرای برنامهٔ خارجی

اگر افزونه باید ابزار خارجی را اجرا کند، یک رابط کوچک به نام `Runner` تعریف و آن را از طریق سازنده به افزونه تزریق کنید. پیاده‌سازی واقعی باید از `exec.CommandContext(ctx, executable, args...)` استفاده کند و پوشهٔ کاری را در `cmd.Dir` قرار دهد. استفاده از `sh -c`، `cmd /c` یا یک رشتهٔ کامل پوسته مجاز نیست. در تست‌ها یک Runner ساختگی قرار دهید تا تست به ابزار نصب‌شده، حساب کاربری یا شبکه وابسته نباشد.

در این حالت، نوع خالی `Plugin` و تابع `New` مرحلهٔ دوم را با نسخه‌ای که `Runner` دارد جایگزین کنید؛ نوع دیگری با همان نام نسازید. رابط Runner باید پوشهٔ کاری را نیز دریافت کند تا اجرای برنامه به یک متغیر سراسری و قابل‌تغییر وابسته نشود.

## مرحلهٔ ۹: ثبت صریح افزونه

بستهٔ جدید را در `plugins/builtin/plugins.go` وارد کنید و سازندهٔ آن را به‌صورت صریح به فهرست افزونه‌ها اضافه کنید:

<div dir="ltr" align="left">

```go
if acmeEnabled {
    plugins = append(plugins, acme.New())
}
```

</div>

این تنها تغییر اجباری در نقطهٔ اتصال برنامه است. افزودن یک ابزار جدید نباید نیازمند تغییر بسته‌های عمومی `internal/completion`، `internal/ranking`، `internal/terminal`، `internal/execution` یا `internal/history` باشد.

## مرحلهٔ ۱۰: افزودن تنظیمات

اگر می‌خواهید کاربر بتواند افزونه را فعال یا غیرفعال کند:

1. فیلد `AcmeEnabled bool` را با کلید JSON مناسب به `internal/config.Config` اضافه کنید.
2. مقدار پیش‌فرض را در `Default()` تعیین کنید.
3. مقدار را در `cmd/assistant/main.go` به `builtin.All` پاس دهید.
4. تست ثبت افزونه را برای هر دو حالت فعال و غیرفعال به‌روزرسانی کنید.

اگر افزونه همیشه داخلی و فعال است، تغییر تنظیمات اختیاری است؛ بااین‌حال ثبت صریح آن همچنان اجباری خواهد بود.

## مرحلهٔ ۱۱: تست

حداقل موارد زیر را تست کنید:

- اطلاعات معرفی افزونه و ثبت صریح آن؛
- نام ناقص فایل اجرایی؛
- تکمیل دستور درون و بیرون پروژه؛
- خواندن وضعیت پروژه و نادیده‌گرفتن پوشه‌های تولیدی؛
- ترتیب ثابت آرگومان‌های پویا؛
- محل‌های قابل‌ویرایش و میزان خطر؛
- Next Action، Best Practice و Recovery؛
- محتوای راهنما؛
- مدیریت خطا و لغو درخواست.

برای تست فایل‌ها از `t.TempDir()` و برای اجرای برنامه از Runner ساختگی استفاده کنید. تست نباید به GitHub، NuGet، اطلاعات ورود یا سرویس شبکه نیاز داشته باشد.

## مرحلهٔ ۱۲: مستندسازی و بررسی نهایی

یک سند مانند `docs/acme-plugin.md` بسازید. بخش انگلیسی باید در ابتدا و بخش فارسی پس از آن قرار گیرد. README، نمونهٔ تنظیمات، پیوند راهنمای افزونه و مسیر آیندهٔ پروژه را نیز در صورت نیاز به‌روزرسانی کنید.

سپس اجرا کنید:

<div dir="ltr" align="left">

```text
make format
make vet
make test
make build-all
```

</div>

در صورت وجود CGO و C compiler:

<div dir="ltr" align="left">

```text
make test-race
```

</div>

در پایان برنامه را اجرا و `ac` و `:? acme` را آزمایش کنید.

## چک‌لیست فایل‌ها

برای یک افزونهٔ داخلی که امکان فعال یا غیرفعال‌شدن دارد، معمولاً فایل‌های زیر تغییر می‌کنند:

<div dir="ltr" align="left">

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

</div>

نباید هیچ منطق مخصوص ابزار جدیدی به بسته‌های عمومی هسته اضافه شود.

</div>
