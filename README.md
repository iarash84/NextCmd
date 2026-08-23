# NextCmd

English | [فارسی](#فارسی)

NextCmd is a fast, deterministic, cross-platform programming command-line assistant written in Go. It suggests editable commands from the current input, project state, and previous execution. It ships with explicitly registered Git, .NET, and Cargo plugins and uses no AI, network service, shell execution, or third-party Go dependency.

## Features

- Interactive editor: Up/Down highlights a suggestion; Tab, Right Arrow, or the first Enter accepts it into the editor; Left Arrow clears the command line; the next Enter executes an accepted command. Use `exit`, `quit`, `:q`, or Ctrl+C to exit.
- Color-aware terminal theme with highlighted selection, suggestion and risk badges, plugin source, and execution status. Colors stay out of redirected output and can be disabled with `NO_COLOR`.
- Built-in `:?`/`:؟` help, per-plugin command catalogs, and suggestions from incomplete executable prefixes such as `gi` or `dot`.
- Structured, shell-independent commands and captured stdout, stderr, exit code, and duration.
- Capability-based public plugin SDK and explicit compile-time registration.
- Git, .NET, and Cargo workspace detection, cached context, dynamic completion, next actions, best practices, and recovery.
- Deterministic prefix/fuzzy ranking and JSON-lines history with secret redaction.
- Standard-library-only implementation with platform-specific terminal boundaries.

## Build and run

Go 1.24 or newer is required.

```text
go build -o nextcmd .
./nextcmd
```

On Windows run `nextcmd.exe`. Debug logging is enabled with `nextcmd --debug`. Configuration defaults to the OS user configuration directory at `nextcmd/config.json`:

```json
{"historyEnabled":true,"maxSuggestions":8,"debug":false,"plugins":{"git":true,"dotnet":true,"cargo":true}}
```

Every registered plugin is enabled unless its ID is set to `false` in `plugins`. Legacy `gitEnabled`, `dotnetEnabled`, and `cargoEnabled` keys are still read for backward compatibility.

### Working directory

NextCmd shows the active working directory above every prompt. Start in another directory with `nextcmd --directory <path>`, or change it without restarting:

```text
pwd
cd ..
cd "C:\Users\Admin\source\repos\My Project"
:cd /path/to/project
```

`cd` and `:cd` update completion, project detection, command execution, and history together. Running `cd` without a path selects the user home directory. NextCmd keeps this state internally and does not change the parent shell directory.

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

`sdk` is the only package plugin authors need. Core discovers optional capabilities with type assertions, merges results, and owns ranking. Plugins return structured commands, descriptions, reasons, and risk metadata; only the terminal package decides how those fields look on screen. Built-ins are listed explicitly in `plugins/builtin`. See [architecture](docs/architecture.md), [plugin development](docs/plugin-development.md), the [Git plugin](docs/git-plugin.md), the [.NET plugin](docs/dotnet-plugin.md), and the [Cargo plugin](docs/cargo-plugin.md).

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

Future compile-time plugins may support Go, Docker, Kubernetes, Terraform, npm, pnpm, and GitHub CLI. Dynamic plugins, AI, telemetry, marketplaces, and cloud sync are intentionally outside this MVP.

## License

MIT. See [LICENSE](LICENSE).

---

<div dir="rtl" align="right">

# فارسی

NextCmd یک دستیار خط فرمان سریع و چندسکویی است که با زبان Go نوشته شده است. برنامه هنگام تایپ، متن فعلی، وضعیت پروژه و نتیجهٔ آخرین دستور را بررسی می‌کند و چند دستور قابل‌ویرایش پیشنهاد می‌دهد. افزونه‌های Git، .NET و Cargo به‌صورت صریح در زمان ساخت برنامه ثبت می‌شوند. برنامه از هوش مصنوعی، سرویس شبکه، اجرای دستور از طریق پوسته و کتابخانهٔ جانبی Go استفاده نمی‌کند.

## قابلیت‌ها

- ویرایشگر تعاملی: کلیدهای بالا و پایین میان پیشنهادها جابه‌جا می‌شوند. کلید Tab، جهت راست یا اولین Enter پیشنهاد را وارد ویرایشگر می‌کند. جهت چپ کل متن را پاک می‌کند و Enter بعدی دستور را اجرا می‌کند. برای خروج می‌توان از `exit`، `quit`، `:q` یا Ctrl+C استفاده کرد.
- ظاهر رنگی پایانه: پیشنهاد انتخاب‌شده، نوع پیشنهاد، میزان خطر، افزونهٔ پیشنهاددهنده و نتیجهٔ اجرای دستور با رنگ‌های متفاوت نمایش داده می‌شوند. با تنظیم متغیر `NO_COLOR` می‌توان رنگ‌ها را غیرفعال کرد.
- راهنمای داخلی با دستورهای `:?` و `:؟` و امکان مشاهدهٔ همهٔ دستورهای شناخته‌شدهٔ هر افزونه.
- نمایش پیشنهاد پیش از کامل‌شدن نام ابزار؛ برای مثال، با نوشتن `gi` پیشنهادهای Git و با نوشتن `dot` پیشنهادهای .NET ظاهر می‌شوند.
- نگهداری دستور به‌صورت نام برنامه و آرگومان‌های جداگانه؛ بنابراین دستورها بدون عبور از پوسته اجرا می‌شوند. خروجی عادی، خروجی خطا، کد خروج و مدت اجرا ثبت می‌شود.
- تشخیص مخزن Git و فضای کاری .NET و Cargo، نگهداری کوتاه‌مدت وضعیت پروژه برای افزایش سرعت، تکمیل مقادیر پویا و ارائهٔ پیشنهاد بعد از موفقیت یا شکست دستور.
- مرتب‌سازی ثابت و قابل‌پیش‌بینی پیشنهادها و ذخیرهٔ تاریخچه در قالب JSON Lines پس از حذف اطلاعات حساس.
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
{"historyEnabled":true,"maxSuggestions":8,"debug":false,"plugins":{"git":true,"dotnet":true,"cargo":true}}
```

</div>

کلیدهای قدیمی `gitEnabled`، `dotnetEnabled` و `cargoEnabled` برای سازگاری با نسخه‌های پیشین همچنان خوانده می‌شوند.

### مسیر کاری

مسیر کاری فعلی بالای هر prompt نمایش داده می‌شود. برای شروع برنامه در مسیر دیگری از `nextcmd --directory <path>` استفاده کنید. داخل برنامه نیز می‌توانید مسیر را بدون جابه‌جاکردن فایل اجرایی یا راه‌اندازی دوباره تغییر دهید:

<div dir="ltr" align="left">

```text
pwd
cd ..
cd "C:\Users\Admin\source\repos\My Project"
:cd /path/to/project
```

</div>

دستورهای `cd` و `:cd` مسیر مورد استفاده برای پیشنهادها، تشخیص پروژه، اجرای دستور و تاریخچه را هم‌زمان تغییر می‌دهند. اجرای `cd` بدون مسیر، پوشهٔ خانگی کاربر را انتخاب می‌کند. این تغییر فقط داخل NextCmd است و مسیر shell والد را تغییر نمی‌دهد.

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

در آینده می‌توان افزونه‌های Go، Docker، Kubernetes، Terraform، npm، pnpm و GitHub CLI را در زمان ساخت به برنامه اضافه کرد. افزونهٔ پویا، هوش مصنوعی، ارسال داده‌های آماری، بازار افزونه و همگام‌سازی ابری عمداً در محدودهٔ نسخهٔ اولیه نیستند.

## مجوز

پروژه تحت مجوز MIT منتشر می‌شود. متن رسمی در [LICENSE](LICENSE) قرار دارد.

</div>
