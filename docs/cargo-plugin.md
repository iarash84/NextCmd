# Cargo Plugin

English | [فارسی](#افزونهٔ-cargo)

The built-in Cargo plugin provides deterministic Rust workflow suggestions through the public NextCmd SDK. Core contains no Cargo-specific command, context type, or ranking rule.

## Project detection

The detector searches the current directory and its parents for `Cargo.toml`. An ancestor containing `[workspace]` becomes the root; otherwise the nearest manifest is used. It scans local manifests while skipping `.git`, editor metadata, `node_modules`, `target`, and `vendor`.

The cached state contains package names, manifest paths, declared non-default features, workspace status, and the presence of `Cargo.lock`. Manifest parsing is deliberately small and read-only; Cargo is not launched on every keystroke.

## Commands and dynamic completion

The initial catalog covers package creation, init, build, check, run, test, bench, formatting, Clippy, documentation, clean, dependency fetch/tree/update/add/remove, lockfile generation, packaging, and publishing. Destructive or externally visible operations receive higher risk levels and lower default priority.

For workspace-aware verbs, `-p` and `--package` complete package names discovered from local manifests. `--features` completes declared feature names. Incomplete executable prefixes such as `car` also show Cargo suggestions. Core commands remain visible outside a Rust project, with reduced contextual priority.

## Workflow guidance

Successful creation suggests check and run. Check suggests Clippy and tests; build suggests tests and run; tests suggest release build and documentation. Dependency changes suggest check and test. Best-practice suggestions cover formatting verification, strict Clippy warnings, and tests.

Recovery handles a missing `Cargo.toml`, invalid workspace package selections, and unknown features using locally discovered values. It intentionally does not parse every Cargo or rustc diagnostic.

## Configuration

Cargo is enabled by default and can be disabled independently:

```json
{"plugins":{"cargo":false}}
```

Use `:? cargo` inside NextCmd to view its static command catalog.

---

<div dir="rtl" align="right">

# افزونهٔ Cargo

افزونهٔ داخلی Cargo با استفاده از SDK عمومی NextCmd، دستورهای مناسب برای پروژه‌های Rust را به‌شکلی ثابت و قابل‌پیش‌بینی پیشنهاد می‌دهد. هسته هیچ دستور، نوع وضعیت یا قانون مرتب‌سازی مخصوص Cargo ندارد.

## تشخیص پروژه

تشخیص‌دهنده در پوشهٔ فعلی و پوشه‌های والد به‌دنبال `Cargo.toml` می‌گردد. اگر یکی از فایل‌ها بخش `[workspace]` داشته باشد، پوشهٔ آن به‌عنوان ریشه انتخاب می‌شود؛ در غیر این صورت نزدیک‌ترین manifest استفاده خواهد شد. هنگام بررسی فایل‌ها، پوشه‌های `.git`، تنظیمات ویرایشگر، `node_modules`، `target` و `vendor` نادیده گرفته می‌شوند.

وضعیت ذخیره‌شده شامل نام packageها، مسیر manifestها، featureهای تعریف‌شده، workspace بودن پروژه و وجود `Cargo.lock` است. خواندن manifest عمداً ساده و فقط‌خواندنی است و برنامه با هر کلید Cargo را دوباره اجرا نمی‌کند.

## دستورها و تکمیل پویا

فهرست اولیه ساخت package، مقداردهی اولیه، build، check، run، test، benchmark، قالب‌بندی، Clippy، مستندات، پاک‌سازی، مدیریت وابستگی‌ها، ساخت lockfile، بسته‌بندی و انتشار را پوشش می‌دهد. عملیات حذف‌کننده یا قابل‌انتشار با خطر بیشتر و اولویت پیش‌فرض کمتر مشخص می‌شوند.

برای دستورهای مربوط به workspace، گزینه‌های `-p` و `--package` نام packageهای واقعی را پیشنهاد می‌دهند. گزینهٔ `--features` نیز featureهای تعریف‌شده در manifest را تکمیل می‌کند. نام ناقص ابزار، مانند `car`، برای نمایش پیشنهادهای Cargo کافی است. دستورهای عمومی بیرون از پروژهٔ Rust نیز با اهمیت کمتر نمایش داده می‌شوند.

## پیشنهاد گام بعدی

پس از ساخت package، اجرای check و run پیشنهاد می‌شود. پس از check، اجرای Clippy و تست؛ پس از build، تست و run؛ و پس از تست، ساخت نسخهٔ release و مستندات پیشنهاد می‌شوند. تغییر وابستگی‌ها نیز پیشنهاد check و test را ایجاد می‌کند. پیشنهادهای روش بهتر شامل بررسی قالب‌بندی، ردکردن warningهای Clippy و اجرای تست‌ها هستند.

پیشنهادهای رفع خطا نبود `Cargo.toml`، نام اشتباه package در workspace و feature ناشناخته را با اطلاعات محلی پوشش می‌دهند. این نسخه عمداً تمام پیام‌های خطای Cargo یا rustc را تفسیر نمی‌کند.

## تنظیمات

Cargo به‌صورت پیش‌فرض فعال است و می‌توان آن را مستقل از افزونه‌های دیگر غیرفعال کرد:

<div dir="ltr" align="left">

```json
{"plugins":{"cargo":false}}
```

</div>

برای دیدن فهرست ثابت دستورها، داخل NextCmd دستور `:? cargo` را اجرا کنید.

</div>
