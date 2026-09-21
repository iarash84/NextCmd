# Dart / Flutter Plugin

English | [فارسی](#افزونهٔ-dart--flutter)

The built-in Dart / Flutter plugin provides deterministic suggestions for Dart and Flutter development workflows. It detects the nearest directory containing `pubspec.yaml`, reads local project metadata without contacting external services, and exposes Dart-specific behavior through the public NextCmd SDK.

## Project detection

The plugin searches the active directory and its parents for `pubspec.yaml`. The nearest matching directory becomes the project root. A project is recognized as a Flutter project when its `pubspec.yaml` declares a `flutter` dependency or section.

The cached project state contains:

- the project root and `pubspec.yaml` path;
- whether `pubspec.lock` exists;
- whether the project is a Flutter project;
- discovered local `.dart` files;
- discovered test files, including files ending in `_test.dart` and files under a `test` directory.

Scanning is read-only and honors context cancellation. Generated or unrelated directories such as `.git`, `.dart_tool`, `build`, `node_modules`, and `vendor` are skipped. The plugin does not contact pub.dev, run dependency resolution, or invoke the Dart or Flutter SDK while scanning.

## Commands

The static catalog includes:

- Dart SDK version inspection with `dart --version`;
- source analysis with `dart analyze`;
- formatting with `dart format .`;
- test and application execution with `dart test` and `dart run`;
- executable compilation with `dart compile exe <file>`;
- dependency resolution with `dart pub get`;
- dependency upgrades and diagnostics with `dart pub upgrade`, `dart pub outdated`, and `dart pub deps`;
- pub cache repair with `dart pub cache repair`;
- package creation with `dart create <name>`.

For Flutter projects, additional `flutter run` and `flutter pub get` suggestions are provided. The plugin recognizes both `dart` and `flutter` executable prefixes during completion.

## Workflow guidance

After `dart pub get` or `dart pub upgrade`, the plugin suggests analysis and tests. After formatting, it suggests analysis; after analysis, it suggests tests; and after tests, it suggests formatting the source.

Best-practice suggestions for a detected project include formatting Dart source, running static analysis, and running the test suite. Flutter projects receive the equivalent recommendations using the `flutter` executable when appropriate.

Recovery suggestions cover missing `pubspec.yaml` project metadata and dependency solver failures. For dependency resolution errors, the plugin suggests inspecting constraints with `pub outdated` rather than changing versions automatically.

## Configuration and help

The plugin ID is `dart`. It is enabled by default and can be disabled without changing Core:

```json
{"plugins":{"dart":false}}
```

Inside NextCmd, use the following command to print the static catalog:

```text
:? dart
```

---

<div dir="rtl" align="right">

# افزونهٔ Dart / Flutter

افزونهٔ داخلی Dart و Flutter پیشنهادهای ثابت و قابل‌پیش‌بینی برای گردش‌کار توسعهٔ Dart و Flutter ارائه می‌دهد. افزونه نزدیک‌ترین پوشهٔ دارای `pubspec.yaml` را تشخیص می‌دهد، اطلاعات محلی پروژه را بدون اتصال به سرویس‌های خارجی می‌خواند و رفتار مخصوص Dart را از طریق SDK عمومی NextCmd در اختیار هسته قرار می‌دهد.

## تشخیص پروژه

افزونه در مسیر کاری فعلی و پوشه‌های والد آن به‌دنبال `pubspec.yaml` می‌گردد. نزدیک‌ترین پوشهٔ دارای این فایل به‌عنوان ریشهٔ پروژه انتخاب می‌شود. اگر فایل `pubspec.yaml` وابستگی یا بخش `flutter` داشته باشد، پروژه به‌عنوان پروژهٔ Flutter شناسایی می‌شود.

وضعیت cacheشدهٔ پروژه شامل موارد زیر است:

- مسیر ریشهٔ پروژه و فایل `pubspec.yaml`؛
- وجود یا نبود `pubspec.lock`؛
- Flutter بودن پروژه؛
- فایل‌های محلی `.dart` شناسایی‌شده؛
- فایل‌های تست شناسایی‌شده، شامل فایل‌های دارای پسوند `_test.dart` و فایل‌های داخل پوشهٔ `test`.

بررسی فایل‌ها فقط خواندنی است و لغو عملیات از طریق context را رعایت می‌کند. پوشه‌های تولیدشده یا نامرتبط مانند `.git`، `.dart_tool`، `build`، `node_modules` و `vendor` بررسی نمی‌شوند. افزونه هنگام اسکن به pub.dev متصل نمی‌شود، حل dependency را اجرا نمی‌کند و SDK مربوط به Dart یا Flutter را فراخوانی نمی‌کند.

## فرمان‌ها

فهرست ثابت فرمان‌ها این موارد را پوشش می‌دهد:

- مشاهدهٔ نسخهٔ SDK با `dart --version`؛
- تحلیل کد با `dart analyze`؛
- قالب‌بندی با `dart format .`؛
- اجرای تست و برنامه با `dart test` و `dart run`؛
- ساخت فایل اجرایی با `dart compile exe <file>`؛
- دریافت dependencyها با `dart pub get`؛
- ارتقا و بررسی dependencyها با `dart pub upgrade`، `dart pub outdated` و `dart pub deps`؛
- تعمیر cache مربوط به pub با `dart pub cache repair`؛
- ساخت package با `dart create <name>`.

برای پروژه‌های Flutter، پیشنهادهای `flutter run` و `flutter pub get` نیز ارائه می‌شوند. افزونه هنگام تکمیل فرمان، پیشوندهای اجرایی `dart` و `flutter` را تشخیص می‌دهد.

## پیشنهادهای گردش کار

پس از اجرای `dart pub get` یا `dart pub upgrade`، تحلیل و تست پیشنهاد می‌شود. پس از قالب‌بندی، تحلیل؛ پس از تحلیل، تست؛ و پس از تست، قالب‌بندی source پیشنهاد می‌شود.

برای پروژهٔ شناسایی‌شده، پیشنهادهای best practice شامل قالب‌بندی sourceهای Dart، اجرای تحلیل ایستا و اجرای مجموعهٔ تست‌ها است. برای پروژه‌های Flutter، در صورت مناسب بودن، همین پیشنهادها با executable مربوط به `flutter` ارائه می‌شوند.

بخش recovery نبود metadata پروژه یعنی `pubspec.yaml` و شکست حل dependencyها را پوشش می‌دهد. در خطاهای مربوط به حل dependency، افزونه بررسی محدودیت نسخه‌ها با `pub outdated` را پیشنهاد می‌کند و نسخه‌ها را به‌صورت خودکار تغییر نمی‌دهد.

## تنظیمات و راهنما

شناسهٔ افزونه `dart` است. افزونه به‌طور پیش‌فرض فعال است و بدون تغییر Core می‌توان آن را غیرفعال کرد:

<div dir="ltr" align="left">

```json
{"plugins":{"dart":false}}
```

```text
:? dart
```

</div>

</div>
