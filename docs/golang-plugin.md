# Go Plugin

English | [فارسی](#افزونهٔ-go)

The built-in Go plugin provides deterministic Go toolchain suggestions through the public NextCmd SDK. All Go-specific state and workflow rules live under `plugins/golang`; Core only sees standard plugin capabilities and structured suggestions.

## Project detection

The detector searches the active directory and its parents for `go.mod` and `go.work`. The nearest module is used unless an ancestor workspace is found. A detected workspace becomes the scan root while the nested module metadata remains available.

The cached state contains:

- project root and module/workspace file paths;
- module path and declared Go version;
- package directories and local `.go` files;
- whether test files and a `vendor` directory exist;
- whether the project is a Go workspace.

Scanning is read-only and honors context cancellation. Generated or unrelated directories such as `.git`, editor metadata, `bin`, `obj`, `node_modules`, `target`, and `vendor` are skipped. The plugin does not run `go list` for every keystroke and does not access the network.

## Commands

The initial catalog includes:

- `go version` and common `go env` inspection;
- module initialization, tidy, download, verify, graph, and vendor;
- workspace initialization, use, and synchronization;
- build, run, test, race detection, and coverage;
- vet, format, generate, list, and documentation;
- dependency retrieval and command installation;
- build-cache cleanup with a destructive risk label.

Type an incomplete executable prefix such as `g` or `go` to receive Go suggestions. Git and Go may both match `g`; Core merges and ranks both sets deterministically.

## Dynamic completion

For `build`, `test`, `vet`, `fmt`, `generate`, and `list`, the plugin suggests package paths discovered from the current project, such as `.` and `./internal/auth`. If multiple packages exist, `./...` is suggested with high relevance. For `go run`, local non-test `.go` files are suggested. Final filtering and ordering remain the responsibility of Core ranking.

## Workflow guidance

Successful builds suggest tests and running the current main package. Tests suggest vet and build; formatting suggests vet and tests; vet suggests tests; and generation suggests formatting and tests. Changes to module dependencies suggest `go mod verify` and a full test run.

Best-practice suggestions cover formatting, vet, tests, and keeping module files synchronized. Recovery covers a missing module, a required `go mod tidy`, an unresolved module dependency, and build constraints that exclude all Go files. Recovery is intentionally small and does not attempt to parse every compiler diagnostic.

## Configuration and help

The plugin ID is `go`. It is enabled by default and can be disabled without changing Core:

```json
{"plugins":{"go":false}}
```

Inside NextCmd, use the following command to print the static catalog:

```text
:? go
```

---

<div dir="rtl" align="right">

# افزونهٔ Go

افزونهٔ داخلی Go با استفاده از SDK عمومی NextCmd، فرمان‌های ابزار Go را به‌شکلی ثابت و قابل‌پیش‌بینی پیشنهاد می‌دهد. تمام اطلاعات و قواعد مخصوص Go در مسیر `plugins/golang` قرار دارند. هسته فقط قابلیت‌های عمومی افزونه و پیشنهادهای ساختاریافته را می‌شناسد.

## تشخیص پروژه

افزونه در مسیر کاری فعلی و پوشه‌های والد به‌دنبال `go.mod` و `go.work` می‌گردد. اگر workspace پیدا نشود، نزدیک‌ترین module انتخاب می‌شود. اگر یکی از پوشه‌های والد فایل `go.work` داشته باشد، همان پوشه ریشهٔ بررسی خواهد بود و اطلاعات module داخلی نیز نگهداری می‌شود.

اطلاعات cacheشده شامل این موارد است:

- مسیر ریشه و فایل‌های module یا workspace؛
- نام module و نسخهٔ Go نوشته‌شده در `go.mod`؛
- مسیر packageها و فایل‌های محلی `.go`؛
- وجود فایل تست و پوشهٔ `vendor`؛
- workspace بودن پروژه.

بررسی فایل‌ها فقط خواندنی است و لغو عملیات از طریق context را رعایت می‌کند. پوشه‌های نامرتبط یا تولیدشده مانند `.git`، تنظیمات ویرایشگر، `bin`، `obj`، `node_modules`، `target` و `vendor` بررسی نمی‌شوند. افزونه برای هر کلید `go list` را اجرا نمی‌کند و به شبکه دسترسی ندارد.

## فرمان‌ها

فهرست اولیه این گروه‌ها را پوشش می‌دهد:

- مشاهدهٔ نسخه و تنظیمات محیط Go؛
- ساخت module و فرمان‌های tidy، download، verify، graph و vendor؛
- ساخت و مدیریت workspace؛
- build، run، test، race detection و coverage؛
- vet، format، generate، list و مستندات؛
- دریافت dependency و نصب ابزارهای Go؛
- پاک‌کردن build cache با برچسب خطر destructive.

برای دیدن پیشنهادهای Go لازم نیست نام ابزار را کامل بنویسید؛ ورودی `g` یا `go` کافی است. ممکن است `g` هم‌زمان با Git و Go مطابقت داشته باشد. در این حالت هسته پیشنهادهای هر دو افزونه را ترکیب می‌کند و با الگوریتم ثابت مرتب می‌سازد.

## تکمیل پویای آرگومان‌ها

برای فرمان‌های `build`، `test`، `vet`، `fmt`، `generate` و `list`، مسیر packageهای واقعی پروژه مانند `.` و `./internal/auth` پیشنهاد داده می‌شوند. اگر چند package وجود داشته باشد، `./...` اهمیت بیشتری می‌گیرد. برای `go run` نیز فایل‌های محلی Go، به‌جز فایل‌های `_test.go`، پیشنهاد داده می‌شوند. فیلترکردن و ترتیب نهایی پیشنهادها همچنان بر عهدهٔ ranking عمومی هسته است.

## پیشنهادهای گردش کار

پس از build موفق، اجرای تست‌ها و main package پیشنهاد می‌شود. بعد از test، فرمان‌های vet و build؛ بعد از format، فرمان‌های vet و test؛ و بعد از generate، قالب‌بندی و تست پیشنهاد می‌شوند. تغییر dependencyهای module نیز پیشنهاد بررسی moduleها و اجرای همهٔ تست‌ها را ایجاد می‌کند.

پیشنهادهای best practice شامل قالب‌بندی، vet، تست و هماهنگ نگه‌داشتن فایل‌های module هستند. recovery نیز نبود module، نیاز به `go mod tidy`، dependency پیدانشده و ناسازگاری build constraintها را پوشش می‌دهد. این بخش عمداً ساده است و تمام خطاهای compiler را تفسیر نمی‌کند.

## تنظیمات و راهنما

شناسهٔ افزونه `go` است. افزونه به‌طور پیش‌فرض فعال است و بدون تغییر Core می‌توان آن را غیرفعال کرد:

<div dir="ltr" align="left">

```json
{"plugins":{"go":false}}
```

</div>

برای مشاهدهٔ فهرست ثابت فرمان‌ها، داخل NextCmd اجرا کنید:

<div dir="ltr" align="left">

```text
:? go
```

</div>

</div>
