# NextCmd

English | [فارسی](#فارسی)

NextCmd is a fast, deterministic, cross-platform programming command-line assistant written in Go. It suggests editable commands from the current input, project state, and previous execution. The MVP ships with explicitly registered Git and .NET plugins and uses no AI, network service, shell execution, or third-party Go dependency.

## Features

- Interactive editor: Up/Down highlights a suggestion; Tab, Right Arrow, or the first Enter accepts it into the editor; Left Arrow clears the command line; the next Enter executes an accepted command. Use `exit`, `quit`, `:q`, or Ctrl+C to exit.
- Built-in `:?`/`:؟` help, per-plugin command catalogs, and suggestions from incomplete executable prefixes such as `gi` or `dot`.
- Structured, shell-independent commands and captured stdout, stderr, exit code, and duration.
- Capability-based public plugin SDK and explicit compile-time registration.
- Git and .NET workspace detection, cached context, dynamic completion, next actions, best practices, and recovery.
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
{"historyEnabled":true,"maxSuggestions":8,"debug":false,"gitEnabled":true,"dotnetEnabled":true}
```

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

`sdk` is the only package plugin authors need. Core discovers optional capabilities with type assertions, merges results, and owns ranking. Plugins return data and never render UI. Built-ins are listed explicitly in `plugins/builtin`. See [architecture](docs/architecture.md), [plugin development](docs/plugin-development.md), the [Git plugin](docs/git-plugin.md), and the [.NET plugin](docs/dotnet-plugin.md).

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

Future compile-time plugins may support Go, Cargo, Docker, Kubernetes, Terraform, npm, pnpm, and GitHub CLI. Dynamic plugins, AI, telemetry, marketplaces, and cloud sync are intentionally outside this MVP.

## License

MIT. See [LICENSE](LICENSE).

---

<div dir="rtl">

# فارسی

NextCmd یک دستیار خط فرمان برنامه‌نویسی سریع، قطعی و چندسکویی است که با Go نوشته شده است. برنامه با توجه به متن فعلی، وضعیت پروژه و نتیجه دستور قبلی، دستورهای قابل‌ویرایش پیشنهاد می‌دهد. نسخه MVP شامل Git Plugin و Dotnet Plugin با registration صریح است و از AI، سرویس شبکه، اجرای shell یا dependency خارجی Go استفاده نمی‌کند.

## قابلیت‌ها

- ویرایشگر تعاملی: جهت بالا/پایین یک پیشنهاد را مشخص می‌کند؛ Tab، جهت راست یا اولین Enter آن را وارد editor می‌کند؛ جهت چپ کل سطر command را پاک می‌کند و Enter بعدی command پذیرفته‌شده را اجرا می‌کند. خروج با `exit`، `quit`، `:q` یا Ctrl+C انجام می‌شود و UI پیش از خروج پاک‌سازی می‌شود.
- راهنمای داخلی `:?`/`:؟`، کاتالوگ commandهای هر Plugin و نمایش پیشنهاد پیش از کامل شدن executable مانند `gi` یا `dot`.
- commandهای ساختاریافته و مستقل از shell همراه با ثبت stdout، stderr، exit code و مدت اجرا.
- Plugin SDK مبتنی بر capability و registration صریح در زمان build.
- تشخیص Git repository و .NET workspace، context دارای cache، تکمیل پویا، اقدام بعدی، best practice و recovery.
- ranking قطعی بر اساس prefix/fuzzy و history از نوع JSON Lines همراه با حذف اطلاعات حساس.
- پیاده‌سازی فقط با Go Standard Library و مرزهای terminal مخصوص هر پلتفرم.

## Build و اجرا

Go نسخه 1.24 یا جدیدتر لازم است:

<div dir="ltr">

```text
go build -o nextcmd .
./nextcmd
```

</div>

در Windows فایل `nextcmd.exe` را اجرا کنید. برای logهای debug از `nextcmd --debug` استفاده کنید.

### Builder با Make

برای workflow کوتاه و مستقیم پروژه از Makefile استفاده کنید:

<div dir="ltr">

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

همه خروجی‌ها داخل `target/` قرار می‌گیرند. `build-root` فقط خروجی سیستم میزبان را در ریشه کپی می‌کند و `build-all` تمام targetهای پشتیبانی‌شده را cross-build می‌کند. جزئیات در [راهنمای دو‌زبانه Make](docs/makefile.md) قرار دارد.

## معماری

توسعه‌دهنده Plugin فقط به package عمومی `sdk` نیاز دارد. Core قابلیت‌های اختیاری را با type assertion تشخیص می‌دهد، suggestionها را ترکیب می‌کند و ranking نهایی را انجام می‌دهد. Plugin فقط داده برمی‌گرداند و UI را نمی‌شناسد. فهرست built-inها به‌شکل صریح در `plugins/builtin` قرار دارد.

## توسعه و تست

<div dir="ltr">

```text
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
```

</div>

## CI و انتشار نسخه

GitHub Actions فرمت، vet، تست‌ها، race detector و هر شش cross-build را اجرا می‌کند. با push کردن یک tag مانند `v1.0.0`، GitHub Release شامل فایل‌های بسته‌بندی‌شده و checksumهای SHA-256 ساخته می‌شود. جزئیات در [راهنمای دو‌زبانه CI و Release](docs/ci-release.md) قرار دارد.

## مسیر آینده

در آینده می‌توان Pluginهای Go، Cargo، Docker، Kubernetes، Terraform، npm، pnpm و GitHub CLI را در زمان build اضافه کرد. Plugin پویا، AI، telemetry، marketplace و cloud sync عمداً خارج از MVP هستند.

## مجوز

پروژه تحت مجوز MIT منتشر می‌شود. متن رسمی در [LICENSE](LICENSE) قرار دارد.

</div>
