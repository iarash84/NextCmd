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

<div dir="rtl">

# ساخت با Make

فایل Makefile فرمان‌های کوتاه موردنیاز برای توسعه روزمره را فراهم می‌کند. این builder فقط Go toolchain نصب‌شده روی سیستم را فراخوانی می‌کند و dependency دانلود نمی‌کند.

Make، cache مربوط به Go را در `target/.go-cache` نگهداری می‌کند و `make clean` این cache را همراه تمام فایل‌های اجرایی تولیدشده حذف می‌کند.

## پیش‌نیازها

- Go نسخه 1.24 یا جدیدتر
- GNU Make

در Windows می‌توان از GNU Make ارائه‌شده توسط MSYS2، MinGW، Chocolatey، Scoop یا محیط سازگار دیگری استفاده کرد. نصب ابزارها را در همان terminal بررسی کنید:

```text
go version
make --version
```

## فرمان‌ها

برای نمایش فهرست کامل فرمان‌ها و توضیح هرکدام در terminal اجرا کنید:

```text
make help
```

```text
make build       # ساخت نسخه سیستم میزبان در target/
make clean       # حذف target/ و فایل اجرایی ریشه
make test        # اجرای unit testها
make run         # ساخت و اجرای نسخه سیستم میزبان
make build-root  # ساخت نسخه میزبان و کپی فایل اجرایی در ریشه
make build-all   # ساخت هر شش خروجی سیستم‌عامل و معماری پشتیبانی‌شده
```

فرمان‌های توسعه‌ای دیگر:

```text
make format
make vet
make test-race
```

فرمان `make test-race` به CGO و یک C compiler سازگار نیاز دارد.

## خروجی‌ها

`make build` فقط فایل سیستم‌عامل و معماری فعلی را تولید می‌کند:

```text
target/NextCmd.exe    # Windows
target/NextCmd        # Linux و macOS
```

`make build-root` همان فایل میزبان را با نام `NextCmd.exe` یا `NextCmd` در ریشه repository کپی می‌کند و cross-build انجام نمی‌دهد.

`make build-all` شش فایل Windows، Linux و macOS را برای amd64 و arm64 داخل `target/` تولید می‌کند.

</div>
