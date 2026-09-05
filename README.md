# NextCmd

English | [فارسی](#فارسی)

<p align="center">
  <img src="assets/nextcmd.png" alt="NextCmd icon" width="220">
</p>

NextCmd is a fast, deterministic, cross-platform programming command-line assistant written in Go. It suggests editable commands from the current input, project state, and previous execution. It ships with explicitly registered Git, .NET, Cargo, Curl, Go, Docker, npm, pip, Kubernetes, and Terraform plugins and uses no AI, network service, or third-party Go dependency.

## Features

- Interactive editor: Up/Down highlights a suggestion; Tab, Right Arrow, or the first Enter accepts it into the editor; placeholders such as `<container>` become editable fields, typing replaces the active field, and Tab moves to the next unresolved field. Left/Right moves the caret; Ctrl+P/Ctrl+N navigate older/newer command history; Ctrl+R performs interactive fuzzy history search; Ctrl+A/Ctrl+E jump to the start or end; Escape or Ctrl+U clears the command line. History navigation and search preserve the current draft. Use `exit`, `quit`, `:q`, Ctrl+C, or Ctrl+D to exit.
- Color-aware terminal theme with highlighted selection, suggestion and risk badges, plugin source, and execution status. Colors stay out of redirected output and can be disabled with `NO_COLOR`; see the [complete badge and risk reference](docs/help.md#suggestion-kinds).
- Built-in command palette when `:` is typed, `:?`/`:؟` help, per-plugin command catalogs, and suggestions from incomplete executable prefixes such as `gi` or `dot`.
- Built-in workspace utilities for listing files, completing paths, trashing/restoring files, viewing redacted history, inspecting plugins and configuration, locating executables, clearing the screen, and checking version information.
- Filesystem path completion for ordinary command arguments, including nested and quoted paths, without invoking the command or a shell.
- Direct, structured process execution by default, plus explicit cross-platform shell execution with a leading `!`. Stdout and stderr stream live while remaining captured for contextual follow-ups; Ctrl+C cancels the running process and returns to NextCmd.
- Confirmation guard for recognized high-risk commands such as recursive forced removal, destructive Git operations, force-push, and Docker prune. Confirmation defaults to no; append `--yes` only for deliberate non-interactive execution.
- Capability-based public plugin SDK and explicit compile-time registration.
- Injectable `sdk.Runner` and `sdk.StreamingRunner` contracts for deterministic tests without launching real tools.
- Git, .NET, Cargo, Curl, Go, Docker, npm, pip, Kubernetes, and Terraform context detection, cached local state, dynamic completion, next actions, best practices, and recovery.
- Deterministic prefix/fuzzy ranking and JSON-lines history with secret redaction for structured arguments, shell options, environment assignments, authentication headers, and URL credentials.
- Standard-library-only implementation with platform-specific terminal boundaries.

## Build and run

Go 1.24 or newer is required.

```text
go build -o nextcmd .
./nextcmd
```

On Windows run `nextcmd.exe`. Debug logging is enabled with `nextcmd --debug`. Configuration defaults to the OS user configuration directory at `nextcmd/config.json`:

```json
{"historyEnabled":true,"maxSuggestions":8,"debug":false,"plugins":{"git":true,"dotnet":true,"cargo":true,"curl":true,"go":true,"docker":true,"npm":true,"pip":true,"kubernetes":true,"terraform":true}}
```

Every registered plugin is enabled unless its ID is set to `false` in `plugins`. Legacy `gitEnabled`, `dotnetEnabled`, and `cargoEnabled` keys are still read for backward compatibility.

### Working directory

NextCmd shows the active working directory above every prompt. Start in another directory with `nextcmd --directory <path>`, or change it without restarting:

```text
pwd
:ls
:ls ..
:del old.txt
:del --dry-run old-build
:trash old.txt
:undo
cd ..
cd "C:\Users\Admin\source\repos\My Project"
:cd /path/to/project
```

`:ls` lists the files and directories in the active working directory. Pass a relative, absolute, or quoted path to inspect another directory without changing the active one. Directories are shown first, followed by files, with type and size columns.

`:del <path>` moves a file or directory from the active working directory to `.nextcmd-trash` after confirmation. It resolves relative, absolute, quoted, and `~` paths, detects whether the target is a file or directory, and asks which one to remove if both a matching file and directory are found. Use `:del --dry-run <path>` to preview, `:del --permanent <path>` to delete without undo support, `:del --permanent --yes <path>` for deliberate non-interactive deletion, and `:undo` to restore the last trashed item in this session.

`cd` and `:cd` update completion, project detection, command execution, and history together. Running `cd` without a path selects the user home directory. NextCmd keeps this state internally and does not change the parent shell directory.

Prefix a command with `!` when shell syntax or a shell built-in is required. NextCmd uses `cmd.exe` on Windows and `/bin/sh` on Linux and macOS. Because shell commands can expand variables, redirect files, and chain processes, use `!` only with text you trust:

```text
! dir
! echo hello > output.txt
! printf '%s\n' hello | grep hello
```

Command output is displayed as it arrives. While a command is running, press Ctrl+C to cancel that process and return to the NextCmd prompt. At the prompt, Ctrl+C or Ctrl+D still exits NextCmd.

Before writing history, NextCmd redacts common credential-bearing options, sensitive environment assignments, authentication headers, and URL user information in both direct and shell commands. Redaction is a best-effort safeguard; avoid placing secrets directly on command lines when a tool supports a safer input mechanism.

NextCmd asks for confirmation before recognized high-risk commands such as `rm -rf`, `git reset --hard`, `git clean -f`, force-push, Docker prune, and equivalent recursive forced deletion on Windows. The default answer is no. For automation, append `--yes`; NextCmd removes this approval flag before launching the command. The guard is an additional safety net, not a substitute for reviewing commands.

Useful built-in commands:

```text
:history 10
:plugins
:clear
:del old.txt
:trash old.txt
:undo
:config
:which git
:version
```

See the [interactive help guide](docs/help.md#utility-commands) for the complete behavior and output of each command.

Cross-build examples:

```text
GOOS=linux GOARCH=arm64 go build ./...
GOOS=darwin GOARCH=amd64 go build ./...
```

PowerShell uses `$env:GOOS="linux"; $env:GOARCH="arm64"; go build ./...`.

### Make builder

The Makefile provides the short project workflow:

```text
make help
make build
make test
make run
make clean
make build-root
make build-all
```

All artifacts are written to `target/`. `build-root` copies only the current host executable to the repository root, while `build-all` cross-builds every supported target. See the bilingual [Make guide](docs/makefile.md).

## Architecture

`sdk` is the only package plugin authors need. Core discovers optional capabilities with type assertions, merges results, and owns ranking. Plugins return structured commands, descriptions, reasons, and risk metadata; only the terminal package decides how those fields look on screen. Built-ins are listed explicitly in `plugins/builtin`. See [architecture](docs/architecture.md), [plugin development](docs/plugin-development.md), and the plugin guides for [Git](docs/git-plugin.md), [.NET](docs/dotnet-plugin.md), [Cargo](docs/cargo-plugin.md), [Curl](docs/curl-plugin.md), [Go](docs/golang-plugin.md), [Docker](docs/docker-plugin.md), [npm](docs/npm-plugin.md), [pip](docs/pip-plugin.md), [Kubernetes](docs/kubernetes-plugin.md), and [Terraform](docs/terraform-plugin.md).

Code that executes commands can depend on `sdk.Runner`; interactive hosts can use `sdk.StreamingRunner`. Production constructors select the system executor, while tests and embedding applications can inject deterministic implementations through `app.NewWithRunner` and `git.NewWithRunner`.

For keyboard controls and plugin command catalogs, see the bilingual [interactive help guide](docs/help.md) or type `:?` inside NextCmd.

## Development and testing

```text
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
```

## CI and releases

GitHub Actions runs formatting, vet, tests, race detection, and all six cross-builds. Pushing a semantic version tag such as `v1.0.0` creates a GitHub Release with packaged binaries and SHA-256 checksums. See the bilingual [CI and release guide](docs/ci-release.md).

## Roadmap

Future compile-time plugins may support pnpm and GitHub CLI. Dynamic plugins, AI, telemetry, marketplaces, and cloud sync are intentionally outside this MVP.

## License

MIT. See [LICENSE](LICENSE).

---

<div dir="rtl" align="right">

# فارسی

<p align="center">
  <img src="assets/nextcmd.png" alt="آیکون NextCmd" width="220">
</p>

NextCmd یک دستیار خط فرمان سریع و چندسکویی است که با زبان Go نوشته شده است. برنامه هنگام تایپ، متن فعلی، وضعیت پروژه و نتیجهٔ آخرین دستور را بررسی می‌کند و چند دستور قابل‌ویرایش پیشنهاد می‌دهد. افزونه‌های Git، .NET، Cargo، Curl، Go، Docker، npm، pip، Kubernetes و Terraform به‌صورت صریح در زمان ساخت برنامه ثبت می‌شوند. برنامه از هوش مصنوعی، سرویس شبکه یا کتابخانهٔ جانبی Go استفاده نمی‌کند.

## قابلیت‌ها

- ویرایشگر تعاملی: کلیدهای بالا و پایین میان پیشنهادها جابه‌جا می‌شوند. کلید Tab، جهت راست یا اولین Enter پیشنهاد را وارد ویرایشگر می‌کند. جای‌نگهدارهایی مانند `<container>` به فیلد قابل‌ویرایش تبدیل می‌شوند؛ تایپ، فیلد فعال را جایگزین می‌کند و Tab به فیلد حل‌نشدهٔ بعدی می‌رود. جهت چپ و راست نشانگر را داخل خط فرمان جابه‌جا می‌کنند. Ctrl+P دستور قدیمی‌تر و Ctrl+N دستور جدیدتر را از تاریخچه وارد ویرایشگر می‌کند. Ctrl+R جست‌وجوی fuzzy تعاملی را در تاریخچه آغاز می‌کند. پیمایش و جست‌وجو متن اولیهٔ کاربر را حفظ می‌کنند. Ctrl+A و Ctrl+E نشانگر را به ابتدا یا انتهای خط می‌برند. Escape یا Ctrl+U خط فرمان فعلی را پاک می‌کند. برای خروج می‌توان از `exit`، `quit`، `:q`، Ctrl+C یا Ctrl+D استفاده کرد.
- ظاهر رنگی پایانه: پیشنهاد انتخاب‌شده، نوع پیشنهاد، میزان خطر، افزونهٔ پیشنهاددهنده و نتیجهٔ اجرای دستور با رنگ‌های متفاوت نمایش داده می‌شوند. با تنظیم متغیر `NO_COLOR` می‌توان رنگ‌ها را غیرفعال کرد. معنی کامل برچسب‌ها و سطح‌های خطر در [راهنمای تعاملی](docs/help.md#نوع-پیشنهاد) آمده است.
- نمایش فهرست فرمان‌های داخلی با تایپ `:`، راهنمای داخلی با دستورهای `:?` و `:؟` و امکان مشاهدهٔ همهٔ دستورهای شناخته‌شدهٔ هر افزونه.
- ابزارهای داخلی برای نمایش فایل‌ها، تکمیل مسیرها، انتقال فایل‌ها به trash و بازگردانی آن‌ها، مشاهدهٔ تاریخچهٔ پاک‌سازی‌شده، بررسی افزونه‌ها و تنظیمات، یافتن فایل اجرایی، پاک‌کردن صفحه و دیدن اطلاعات نسخه.
- تکمیل مسیر فایل و پوشه برای آرگومان دستورهای عادی، شامل مسیرهای تو‌در‌تو و نقل‌قول‌شده، بدون اجرای دستور یا shell.
- نمایش پیشنهاد پیش از کامل‌شدن نام ابزار؛ برای مثال، با نوشتن `gi` پیشنهادهای Git و با نوشتن `dot` پیشنهادهای .NET ظاهر می‌شوند.
- اجرای مستقیم دستورهای عادی با نگهداری نام برنامه و آرگومان‌های جداگانه و امکان اجرای صریح shell با پیشوند `!`. خروجی عادی و خطا هم‌زمان با اجرا نمایش داده و برای پیشنهادهای بعدی ثبت می‌شوند. هنگام اجرای دستور، Ctrl+C فقط همان process را لغو می‌کند و NextCmd باز می‌ماند.
- درخواست تأیید برای فرمان‌های پرخطر شناخته‌شده مانند حذف بازگشتی اجباری، عملیات مخرب Git، force-push و Docker prune. پاسخ پیش‌فرض منفی است و `--yes` فقط برای اجرای غیرتعاملی آگاهانه استفاده می‌شود.
- قراردادهای قابل‌تزریق `sdk.Runner` و `sdk.StreamingRunner` برای تست قطعی بدون اجرای ابزارهای واقعی.
- تشخیص وضعیت Git، .NET، Cargo، Go، Docker، npm، pip، Kubernetes، Terraform و فایل‌های محلی موردنیاز Curl، نگهداری کوتاه‌مدت context برای افزایش سرعت، تکمیل مقادیر پویا و ارائهٔ پیشنهاد بعد از موفقیت یا شکست دستور.
- مرتب‌سازی ثابت و قابل‌پیش‌بینی پیشنهادها و ذخیرهٔ تاریخچه در قالب JSON Lines پس از پاک‌سازی آرگومان‌ها، optionهای shell، مقداردهی متغیرهای محیطی، headerهای احراز هویت و اطلاعات ورود URL.
- پیاده‌سازی فقط با کتابخانهٔ استاندارد Go و کد جداگانه برای رفتارهای وابسته به هر سیستم‌عامل.

## ساخت و اجرا

Go نسخه 1.24 یا جدیدتر لازم است:

<div dir="ltr" align="left">

```text
go build -o nextcmd .
./nextcmd
```

</div>

در ویندوز فایل `nextcmd.exe` را اجرا کنید. برای مشاهدهٔ گزارش‌های فنی، برنامه را با `nextcmd --debug` اجرا کنید.

تنظیمات افزونه‌ها عمومی است و برای افزونهٔ جدید به فیلد تازه نیاز ندارد. همهٔ افزونه‌ها فعال‌اند، مگر اینکه شناسهٔ آن‌ها در نقشهٔ `plugins` برابر `false` باشد:

<div dir="ltr" align="left">

```json
{"historyEnabled":true,"maxSuggestions":8,"debug":false,"plugins":{"git":true,"dotnet":true,"cargo":true,"curl":true,"go":true,"docker":true,"npm":true,"pip":true,"kubernetes":true,"terraform":true}}
```

</div>

کلیدهای قدیمی `gitEnabled`، `dotnetEnabled` و `cargoEnabled` برای سازگاری با نسخه‌های پیشین همچنان خوانده می‌شوند.

### مسیر کاری

مسیر کاری فعلی بالای هر prompt نمایش داده می‌شود. برای شروع برنامه در مسیر دیگری از `nextcmd --directory <path>` استفاده کنید. داخل برنامه نیز می‌توانید مسیر را بدون جابه‌جاکردن فایل اجرایی یا راه‌اندازی دوباره تغییر دهید:

<div dir="ltr" align="left">

```text
pwd
:ls
:ls ..
:del old.txt
:del --dry-run old-build
:trash old.txt
:undo
cd ..
cd "C:\Users\Admin\source\repos\My Project"
:cd /path/to/project
```

</div>

دستور `:ls` فایل‌ها و پوشه‌های مسیر کاری فعلی را نمایش می‌دهد. برای دیدن محتوای مسیری دیگر، همان مسیر را پس از دستور بنویسید؛ برای مثال `:ls ..` یا `:ls "My Project"`. این کار مسیر کاری فعلی را تغییر نمی‌دهد. در خروجی، ابتدا پوشه‌ها و سپس فایل‌ها همراه با نوع و اندازه نمایش داده می‌شوند.

دستور `:del <path>` فایل یا پوشه را پس از تأیید از مسیر کاری فعال به `.nextcmd-trash` منتقل می‌کند. مسیرهای نسبی، کامل، نقل‌قول‌شده و `~` پشتیبانی می‌شوند. برنامه تشخیص می‌دهد هدف فایل است یا پوشه و اگر هم فایل و هم پوشهٔ مطابق پیدا شود، از کاربر می‌پرسد کدام مورد حذف شود. برای پیش‌نمایش از `:del --dry-run <path>`، برای حذف دائمی بدون امکان undo از `:del --permanent <path>`، برای حذف غیرتعاملی آگاهانه از `:del --permanent --yes <path>` و برای بازگردانی آخرین مورد منتقل‌شده به trash از `:undo` استفاده کنید.

دستورهای `cd` و `:cd` مسیر مورد استفاده برای پیشنهادها، تشخیص پروژه، اجرای دستور و تاریخچه را هم‌زمان تغییر می‌دهند. اجرای `cd` بدون مسیر، پوشهٔ خانگی کاربر را انتخاب می‌کند. این تغییر فقط داخل NextCmd است و مسیر shell والد را تغییر نمی‌دهد.

برای اجرای قابلیت‌های داخلی shell یا استفاده از pipe و redirect، ابتدای دستور `!` قرار دهید. NextCmd در Windows از `cmd.exe` و در Linux و macOS از `/bin/sh` استفاده می‌کند. چون shell می‌تواند متغیرها را گسترش دهد، فایل‌ها را بازنویسی کند و چند process را پشت‌سرهم اجرا کند، فقط متن مورداعتماد را پس از `!` اجرا کنید:

<div dir="ltr" align="left">

```text
! dir
! echo hello > output.txt
! printf '%s\n' hello | grep hello
```

</div>

خروجی دستور هم‌زمان با تولید نمایش داده می‌شود. هنگام اجرای یک دستور، Ctrl+C همان process را لغو می‌کند و برنامه به prompt بعدی برمی‌گردد. در خود prompt، کلیدهای Ctrl+C و Ctrl+D همچنان از NextCmd خارج می‌شوند.

NextCmd پیش از ذخیرهٔ تاریخچه، optionهای رایج حاوی credential، مقداردهی متغیرهای حساس، headerهای احراز هویت و اطلاعات کاربری URL را در فرمان‌های مستقیم و shell پاک‌سازی می‌کند. این قابلیت یک محافظ best-effort است؛ هرگاه ابزار روش امن‌تری برای دریافت secret دارد، از قراردادن مستقیم آن در خط فرمان خودداری کنید.

NextCmd پیش از فرمان‌های پرخطر شناخته‌شده مانند `rm -rf`، `git reset --hard`، `git clean -f`، force-push، Docker prune و حذف بازگشتی اجباری معادل در Windows تأیید می‌گیرد. پاسخ پیش‌فرض منفی است. برای automation می‌توان `--yes` را در انتهای فرمان قرار داد؛ NextCmd این flag تأیید را پیش از اجرای ابزار حذف می‌کند. این محافظ جای بررسی خود فرمان را نمی‌گیرد.

چند دستور داخلی کاربردی:

<div dir="ltr" align="left">

```text
:history 10
:plugins
:clear
:del old.txt
:trash old.txt
:undo
:config
:which git
:version
```

</div>

جزئیات کامل رفتار و خروجی هر دستور در [راهنمای تعاملی](docs/help.md#دستورهای-کاربردی-داخلی) آمده است.

### ساخت با Make

برای اجرای سادهٔ کارهای رایج پروژه از Makefile استفاده کنید:

<div dir="ltr" align="left">

```text
make help
make build
make test
make run
make clean
make build-root
make build-all
```

</div>

همهٔ خروجی‌ها در پوشهٔ `target/` قرار می‌گیرند. دستور `build-root` نسخهٔ مناسب سیستم فعلی را در ریشهٔ پروژه کپی می‌کند. دستور `build-all` نیز نسخهٔ مربوط به همهٔ سیستم‌عامل‌ها و معماری‌های پشتیبانی‌شده را می‌سازد. جزئیات در [راهنمای Make](docs/makefile.md) آمده است.

## معماری

توسعه‌دهندهٔ افزونه فقط به بستهٔ عمومی `sdk` نیاز دارد. هسته تشخیص می‌دهد هر افزونه چه قابلیت‌هایی دارد، پیشنهادهای همهٔ افزونه‌ها را با هم ترکیب می‌کند و ترتیب نهایی نمایش را تعیین می‌کند. افزونه فقط اطلاعاتی مانند دستور، عنوان، دلیل پیشنهاد و میزان خطر را برمی‌گرداند؛ نحوهٔ نمایش این اطلاعات بر عهدهٔ رابط پایانه است. افزونه‌های داخلی به‌صورت صریح در `plugins/builtin` ثبت می‌شوند.

کدی که فرمان اجرا می‌کند می‌تواند به `sdk.Runner` وابسته باشد و میزبان تعاملی از `sdk.StreamingRunner` استفاده کند. سازنده‌های عادی اجراکنندهٔ واقعی سیستم را انتخاب می‌کنند، اما تست‌ها و برنامه‌های میزبان می‌توانند با `app.NewWithRunner` و `git.NewWithRunner` پیاده‌سازی قطعی خود را تزریق کنند.

برای جزئیات بیشتر، [معماری](docs/architecture.md)، [راهنمای توسعهٔ افزونه](docs/plugin-development.md) و راهنماهای [Git](docs/git-plugin.md)، [.NET](docs/dotnet-plugin.md)، [Cargo](docs/cargo-plugin.md)، [Curl](docs/curl-plugin.md)، [Go](docs/golang-plugin.md)، [Docker](docs/docker-plugin.md)، [npm](docs/npm-plugin.md)، [pip](docs/pip-plugin.md)، [Kubernetes](docs/kubernetes-plugin.md) و [Terraform](docs/terraform-plugin.md) را بخوانید.

## توسعه و تست

<div dir="ltr" align="left">

```text
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
```

</div>

## بررسی خودکار و انتشار نسخه

GitHub Actions قالب‌بندی کد، بررسی `go vet`، تست‌ها، تشخیص رقابت داده و ساخت هر شش خروجی را انجام می‌دهد. با ارسال یک برچسب نسخه مانند `v1.0.0`، صفحهٔ انتشار GitHub همراه فایل‌های فشرده و checksumهای SHA-256 ساخته می‌شود. جزئیات در [راهنمای بررسی و انتشار](docs/ci-release.md) آمده است.

## مسیر آینده

در آینده می‌توان افزونه‌های pnpm و GitHub CLI را در زمان ساخت به برنامه اضافه کرد. افزونهٔ پویا، هوش مصنوعی، ارسال داده‌های آماری، بازار افزونه و همگام‌سازی ابری عمداً در محدودهٔ نسخهٔ اولیه نیستند.

## مجوز

پروژه تحت مجوز MIT منتشر می‌شود. متن رسمی در [LICENSE](LICENSE) قرار دارد.

</div>
