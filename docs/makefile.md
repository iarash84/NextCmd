# Building with Make

English | [فارسی](#ساخت-با-make)

The project Makefile provides the short builder commands used in day-to-day development. It invokes the locally installed Go toolchain and downloads no dependencies.

Make keeps its Go build cache under `target/.go-cache`; `make clean` removes the cache together with all generated executables.

## Prerequisites

- Go 1.24 or newer
- GNU Make

On Windows, use GNU Make from MSYS2, MinGW, Chocolatey, Scoop, or another compatible environment. Confirm both tools are visible in the same terminal:

```text
go version
make --version
```

## Commands

Run `make help` at any time to print the complete command reference directly in the terminal:

```text
make help
```

```text
make build       # Build for the current host into target/
make clean       # Remove target/ and the root executable
make test        # Run unit tests
make run         # Build and run the host executable
make build-root  # Build for the host and copy NextCmd to the repository root
make build-all   # Cross-build all six supported OS/architecture targets
```

Additional development commands are available:

```text
make format
make vet
make test-race
```

`make test-race` requires CGO and a compatible C compiler.

## Outputs

`make build` creates only the executable for the current operating system and architecture:

```text
target/NextCmd.exe    # Windows
target/NextCmd        # Linux and macOS
```

`make build-root` copies that host executable to `NextCmd.exe` or `NextCmd` in the repository root. It does not cross-compile.

`make build-all` creates:

```text
target/NextCmd-windows-amd64.exe
target/NextCmd-windows-arm64.exe
target/NextCmd-linux-amd64
target/NextCmd-linux-arm64
target/NextCmd-darwin-amd64
target/NextCmd-darwin-arm64
```

---

<div dir="rtl" align="right">

# ساخت با Make

فایل Makefile نام‌های کوتاهی برای کارهای رایج توسعه فراهم می‌کند. این فایل فقط ابزارهای Go نصب‌شده روی سیستم را اجرا می‌کند و خودش هیچ وابستگی‌ای دانلود نمی‌کند.

Make حافظهٔ موقت Go را در `target/.go-cache` نگه می‌دارد. دستور `make clean` تمام پوشهٔ `target` و فایل اجرایی کپی‌شده در ریشهٔ پروژه را حذف می‌کند.

## پیش‌نیازها

- Go نسخه 1.24 یا جدیدتر
- GNU Make

در ویندوز می‌توان GNU Make را از طریق MSYS2، MinGW، Chocolatey یا Scoop نصب کرد. برای اطمینان از نصب درست، دستورهای زیر را در همان پایانه‌ای اجرا کنید که قرار است پروژه را بسازید:

<div dir="ltr" align="left">

```text
go version
make --version
```

</div>

## فرمان‌ها

برای مشاهدهٔ همهٔ فرمان‌های پشتیبانی‌شده و توضیح هرکدام، اجرا کنید:

<div dir="ltr" align="left">

```text
make help
```

</div>

<div dir="ltr" align="left">

```text
make build       # ساخت نسخه سیستم میزبان در target/
make clean       # حذف target/ و فایل اجرایی ریشه
make test        # اجرای unit testها
make run         # ساخت و اجرای نسخه سیستم میزبان
make build-root  # ساخت نسخه میزبان و کپی فایل اجرایی در ریشه
make build-all   # ساخت هر شش خروجی سیستم‌عامل و معماری پشتیبانی‌شده
```

</div>

فرمان‌های توسعه‌ای دیگر:

<div dir="ltr" align="left">

```text
make format
make vet
make test-race
```

</div>

فرمان `make test-race` به CGO و یک کامپایلر C سازگار نیاز دارد. اگر این ابزارها نصب نباشند، سایر تست‌ها همچنان قابل اجرا هستند.

## خروجی‌ها

`make build` فقط فایل سیستم‌عامل و معماری فعلی را تولید می‌کند:

<div dir="ltr" align="left">

```text
target/NextCmd.exe    # Windows
target/NextCmd        # Linux و macOS
```

</div>

`make build-root` نسخهٔ مناسب سیستم فعلی را می‌سازد و با نام `NextCmd.exe` یا `NextCmd` در ریشهٔ پروژه کپی می‌کند. این فرمان برای سیستم‌عامل‌های دیگر خروجی نمی‌سازد.

`make build-all` شش فایل مربوط به ویندوز، لینوکس و macOS را برای معماری‌های amd64 و arm64 در پوشهٔ `target/` تولید می‌کند.

</div>
