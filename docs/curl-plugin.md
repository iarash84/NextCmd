# Curl Plugin

English | [فارسی](#افزونهٔ-curl)

The built-in Curl plugin provides deterministic HTTP transfer, API request, download, upload, and diagnostic suggestions. It uses only the public NextCmd SDK and the Go standard library. Core contains no Curl-specific rule or type.

The command catalog follows the official [curl command-line manual](https://curl.se/docs/manpage.html).

## Completion

Typing `cur` or `curl` exposes common templates for:

- GET, HEAD, POST, PUT, PATCH, and DELETE requests;
- JSON, form, URL-encoded, and file-backed request bodies;
- custom headers, redirects, compression, and response headers;
- output files and remote filenames with partial-file cleanup;
- retries, connection timeout, total timeout, status, and timing;
- verbose diagnostics, local config, netrc, and CA certificates.

DELETE is marked destructive. `--insecure` is marked dangerous and deliberately receives very low priority because it disables TLS identity verification.

## Local file completion

The plugin reads only the immediate working directory and caches the result for two seconds. It dynamically completes:

- regular files for `--upload-file` and file-backed `--data`;
- `.curlrc`, `.curl`, and `.conf` files for `--config`;
- `.pem`, `.crt`, and `.cer` files for `--cacert`.

No request is sent and Curl is not launched during completion. Credentials are not embedded in suggestions; `--netrc-file` points to a local credential file instead.

## Workflow and recovery

After a successful transfer, the plugin can suggest inspecting headers or measuring the status code and duration. Best-practice suggestions use `--fail-with-body`, `--show-error`, bounded timeouts, and conservative retry behavior.

Recovery recognizes common DNS resolution, timeout, TLS certificate, local-file, and HTTP error failures. A local CA certificate is preferred for TLS recovery. `--insecure` remains a low-priority diagnostic fallback and is visibly marked dangerous.

## Configuration and help

Curl is enabled by default. Disable it through the generic plugin map:

```json
{"plugins":{"curl":false}}
```

Use `:? curl` inside NextCmd to print the static command catalog.

---

<div dir="rtl" align="right">

# افزونهٔ Curl

افزونهٔ داخلی Curl برای درخواست‌های HTTP، فراخوانی API، دانلود، بارگذاری فایل و بررسی خطاهای اتصال پیشنهادهای ثابت و قابل‌پیش‌بینی ارائه می‌دهد. این افزونه فقط از SDK عمومی NextCmd و کتابخانهٔ استاندارد Go استفاده می‌کند و هیچ قانون یا نوع مخصوص Curl به هسته اضافه نمی‌شود.

فهرست دستورها براساس [راهنمای رسمی خط فرمان curl](https://curl.se/docs/manpage.html) تهیه شده است.

## تکمیل دستور

با تایپ `cur` یا `curl` الگوهای رایج زیر نمایش داده می‌شوند:

- درخواست‌های GET، HEAD، POST، PUT، PATCH و DELETE؛
- بدنهٔ JSON، فرم، دادهٔ URL-encoded و دادهٔ خوانده‌شده از فایل؛
- header سفارشی، دنبال‌کردن redirect، فشرده‌سازی و نمایش header پاسخ؛
- ذخیرهٔ خروجی در فایل و حذف فایل ناقص پس از خطا؛
- retry، محدودیت زمان اتصال، محدودیت کل زمان، status و مدت انتقال؛
- گزارش verbose، فایل تنظیمات محلی، netrc و گواهی CA.

درخواست DELETE با خطر حذف‌کننده مشخص می‌شود. گزینهٔ `--insecure` نیز خطرناک است و اولویت بسیار کمی دارد، زیرا بررسی هویت TLS را غیرفعال می‌کند.

## تکمیل فایل‌های محلی

افزونه فقط فایل‌های مستقیم مسیر کاری را می‌خواند و نتیجه را دو ثانیه نگه می‌دارد. موارد زیر به‌صورت پویا تکمیل می‌شوند:

- فایل‌های عادی برای `--upload-file` و `--data` مبتنی بر فایل؛
- فایل‌های `.curlrc`، `.curl` و `.conf` برای `--config`؛
- فایل‌های `.pem`، `.crt` و `.cer` برای `--cacert`.

هنگام تکمیل هیچ درخواست شبکه‌ای ارسال و Curl اجرا نمی‌شود. اطلاعات ورود نیز داخل پیشنهاد قرار نمی‌گیرند؛ گزینهٔ `--netrc-file` فقط به فایل محلی اطلاعات ورود اشاره می‌کند.

## پیشنهاد گام بعدی و رفع خطا

پس از انتقال موفق، افزونه می‌تواند مشاهدهٔ headerها یا اندازه‌گیری status و مدت انتقال را پیشنهاد دهد. پیشنهادهای روش بهتر شامل `--fail-with-body`، `--show-error`، محدودیت زمانی و retry محافظه‌کارانه هستند.

پیشنهاد رفع خطا، مشکلات رایج DNS، timeout، گواهی TLS، فایل محلی و پاسخ خطای HTTP را می‌شناسد. برای خطای TLS ابتدا گواهی CA موجود در مسیر کاری پیشنهاد می‌شود. گزینهٔ `--insecure` فقط به‌عنوان راه‌حل تشخیصی کم‌اولویت باقی می‌ماند و با خطر بالا نمایش داده می‌شود.

## تنظیمات و راهنما

Curl به‌صورت پیش‌فرض فعال است. برای غیرفعال‌کردن آن از نقشهٔ عمومی افزونه‌ها استفاده کنید:

<div dir="ltr" align="left">

```json
{"plugins":{"curl":false}}
```

</div>

برای مشاهدهٔ فهرست ثابت دستورها، داخل NextCmd دستور `:? curl` را اجرا کنید.

</div>
