# Building NextCmd

English | [فارسی](#ساخت-nextcmd)

NextCmd can be built directly with Go or through the provided CMake facade. CMake does not replace the Go toolchain and does not download dependencies; it invokes the locally installed `go` executable.

## Prerequisites

- Go 1.24 or newer
- CMake 3.20 or newer
- A CMake generator available on the host, such as Ninja, Visual Studio, Unix Makefiles, or NMake

Verify the tools before configuring the project:

```text
go version
cmake --version
```

## Configure and build

Run these commands from the repository root:

```text
cmake -S . -B build
cmake --build build
```

`-S .` selects the repository root as the source directory. `-B build` creates an out-of-source build directory that is excluded by `.gitignore`.

The executable is generated at:

```text
build/bin/NextCmd.exe    # Windows
build/bin/NextCmd        # Linux and macOS
```

Run it on Windows with:

```powershell
.\build\bin\NextCmd.exe
```

Run it on Linux or macOS with:

```sh
./build/bin/NextCmd
```

## Build targets

The default `nextcmd` target builds the executable. Additional development targets are available:

```text
cmake --build build --target nextcmd
cmake --build build --target nextcmd-root
cmake --build build --target format
cmake --build build --target vet
cmake --build build --target test
cmake --build build --target test-race
cmake --build build --target build-all-platforms
```

`test-race` requires CGO and a working C compiler. `build-all-platforms` produces Windows, Linux, and macOS binaries for amd64 and arm64 under `build/bin`.

## Build the executable in the repository root

For a quick local test, build the dedicated `nextcmd-root` target:

```text
cmake --build build --target nextcmd-root
```

On Windows this creates `NextCmd.exe` directly in the repository root, where it can be run with:

```powershell
.\NextCmd.exe
```

The root executable is ignored by Git and removed by CMake's `clean` target. Regular builds remain out-of-source, which is the recommended default.

CMake keeps its Go build cache under `build/go-cache` by default, isolating generated cache files from the repository and the user's global Go cache. Advanced users can override it while configuring with `-DNEXTCMD_GO_CACHE_DIR=<path>`.

To make root output the default for a separate build tree, enable the option while configuring:

```text
cmake -S . -B build-root -DNEXTCMD_OUTPUT_IN_SOURCE_ROOT=ON
cmake --build build-root
```

Individual cross-build targets are also available, for example:

```text
cmake --build build --target build-windows-arm64
cmake --build build --target build-linux-amd64
cmake --build build --target build-darwin-arm64
```

## Select a generator

CMake normally selects a default generator. Ninja can be selected explicitly when it is installed:

```text
cmake -S . -B build -G Ninja
cmake --build build
```

On Windows, a Visual Studio generator can also be selected explicitly:

```text
cmake -S . -B build -G "Visual Studio 17 2022"
cmake --build build --config Release
```

The Go build itself is not affected by CMake's Debug or Release configuration because compiler behavior remains controlled by the Go toolchain.

## Clean rebuild

The CMake build is out-of-source. Remove only the generated `build` directory and configure again:

```powershell
Remove-Item -Recurse -Force .\build
cmake -S . -B build
cmake --build build
```

On Linux or macOS:

```sh
rm -rf ./build
cmake -S . -B build
cmake --build build
```

## Troubleshooting

- `cmake: command not found`: install CMake, then reopen the terminal or IDE so `PATH` is refreshed.
- `Could NOT find GO_EXECUTABLE`: install Go and ensure `go version` works in the same terminal.
- Generator errors: install Ninja, Make, or the requested Visual Studio CMake components, or let CMake select its default generator.
- Race detector errors: enable CGO and install a compatible C compiler, or use the regular `test` target.

---

# ساخت NextCmd

NextCmd را می‌توان مستقیماً با Go یا از طریق رابط CMake موجود در پروژه build کرد. CMake جایگزین Go toolchain نیست، dependency دانلود نمی‌کند و فقط executable محلی `go` را فراخوانی می‌کند.

## پیش‌نیازها

- Go نسخه 1.24 یا جدیدتر
- CMake نسخه 3.20 یا جدیدتر
- یک generator قابل استفاده مانند Ninja، Visual Studio، Unix Makefiles یا NMake

ابتدا نصب ابزارها را بررسی کنید:

```text
go version
cmake --version
```

## پیکربندی و build

دستورهای زیر را در ریشه repository اجرا کنید:

```text
cmake -S . -B build
cmake --build build
```

گزینه `-S .` ریشه پروژه را به‌عنوان source و `-B build` پوشه جداگانه `build` را برای فایل‌های تولیدشده انتخاب می‌کند. این پوشه توسط `.gitignore` نادیده گرفته می‌شود.

فایل اجرایی در مسیر زیر ساخته می‌شود:

```text
build/bin/NextCmd.exe    # Windows
build/bin/NextCmd        # Linux و macOS
```

اجرا در Windows:

```powershell
.\build\bin\NextCmd.exe
```

اجرا در Linux یا macOS:

```sh
./build/bin/NextCmd
```

## targetهای موجود

```text
cmake --build build --target nextcmd
cmake --build build --target nextcmd-root
cmake --build build --target format
cmake --build build --target vet
cmake --build build --target test
cmake --build build --target test-race
cmake --build build --target build-all-platforms
```

target مربوط به `test-race` به CGO و یک C compiler نیاز دارد. target مربوط به `build-all-platforms` خروجی Windows، Linux و macOS را برای معماری‌های amd64 و arm64 در `build/bin` تولید می‌کند.

## ساخت فایل اجرایی در ریشه پروژه

برای تست سریع برنامه، target مستقل زیر را build کنید:

```text
cmake --build build --target nextcmd-root
```

در Windows فایل `NextCmd.exe` مستقیماً در ریشه repository ساخته می‌شود و می‌توان آن را اجرا کرد:

```powershell
.\NextCmd.exe
```

این فایل توسط Git نادیده گرفته می‌شود و target مربوط به `clean` در CMake نیز آن را حذف می‌کند. build معمولی همچنان خارج از source tree انجام می‌شود که روش پیشنهادی است.

CMake به‌صورت پیش‌فرض cache مربوط به Go را در `build/go-cache` نگهداری می‌کند تا فایل‌های تولیدشده از repository و cache سراسری کاربر جدا باشند. کاربران پیشرفته می‌توانند هنگام configure مسیر دیگری را با `-DNEXTCMD_GO_CACHE_DIR=<path>` تعیین کنند.

برای اینکه خروجی ریشه در یک build tree جداگانه به حالت پیش‌فرض تبدیل شود، option زیر را هنگام configure فعال کنید:

```text
cmake -S . -B build-root -DNEXTCMD_OUTPUT_IN_SOURCE_ROOT=ON
cmake --build build-root
```

برای ساخت یک target خاص نیز می‌توان از دستورهایی مانند موارد زیر استفاده کرد:

```text
cmake --build build --target build-windows-arm64
cmake --build build --target build-linux-amd64
cmake --build build --target build-darwin-arm64
```

## انتخاب generator

در صورت نصب Ninja می‌توان آن را صریحاً انتخاب کرد:

```text
cmake -S . -B build -G Ninja
cmake --build build
```

در Windows امکان استفاده از Visual Studio نیز وجود دارد:

```text
cmake -S . -B build -G "Visual Studio 17 2022"
cmake --build build --config Release
```

تنظیم Debug یا Release در CMake رفتار compiler زبان Go را تغییر نمی‌دهد، زیرا build همچنان توسط Go toolchain کنترل می‌شود.

## build تمیز

در Windows:

```powershell
Remove-Item -Recurse -Force .\build
cmake -S . -B build
cmake --build build
```

در Linux یا macOS:

```sh
rm -rf ./build
cmake -S . -B build
cmake --build build
```

## رفع خطاهای رایج

- خطای `cmake: command not found`: ابتدا CMake را نصب و terminal یا IDE را دوباره باز کنید.
- خطای `Could NOT find GO_EXECUTABLE`: مطمئن شوید دستور `go version` در همان terminal اجرا می‌شود.
- خطای generator: Ninja، Make یا componentهای CMake مربوط به Visual Studio را نصب کنید یا generator را مشخص نکنید تا CMake مقدار پیش‌فرض را انتخاب کند.
- خطای race detector: CGO و یک C compiler سازگار را فعال کنید یا از target معمولی `test` استفاده کنید.
