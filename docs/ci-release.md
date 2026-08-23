# Continuous Integration and Releases

English | [فارسی](#یکپارچهسازی-و-انتشار)

NextCmd uses two GitHub Actions workflows under `.github/workflows`.

## CI

`ci.yml` runs for pushes and pull requests targeting `main`, and can also be started manually. It:

- runs tests and `go vet` natively on Ubuntu, Windows, and macOS;
- verifies `gofmt` and runs the race detector on Ubuntu;
- invokes `make build-all` for all six supported OS/architecture targets;
- uploads the six binaries as a workflow artifact retained for seven days;
- cancels an older CI run when a newer commit is pushed to the same ref.

The workflow has read-only repository permissions.

## Releases

`release.yml` publishes a release when a semantic version tag beginning with `v` is pushed, for example:

```text
git tag -a v1.0.0 -m "NextCmd v1.0.0"
git push origin v1.0.0
```

The workflow validates that the tag exists and follows semantic-version syntax, runs tests, builds all targets, creates ZIP files for Windows, `tar.gz` files for Linux and macOS, generates SHA-256 checksums, and publishes the assets with generated release notes.

A release can also be re-triggered from the Actions page with `workflow_dispatch`, but the supplied tag must already exist. The release job receives only the `contents: write` permission required to create the GitHub Release. No repository secret is required; it uses the scoped `GITHUB_TOKEN` supplied by GitHub Actions.

---

<div dir="rtl" align="right">

# یکپارچه‌سازی و انتشار

NextCmd در مسیر `.github/workflows` دو گردش کار خودکار دارد: یکی برای بررسی تغییرات و دیگری برای ساخت نسخهٔ قابل انتشار.

## CI

فایل `ci.yml` با هر push یا pull request مربوط به شاخهٔ اصلی `main` اجرا می‌شود و از صفحهٔ Actions نیز می‌توان آن را دستی اجرا کرد. این گردش کار:

- تست‌ها و `go vet` را مستقیماً روی Ubuntu، ویندوز و macOS اجرا می‌کند؛
- قالب‌بندی `gofmt` و تشخیص رقابت داده را روی Ubuntu بررسی می‌کند؛
- با `make build-all` هر شش خروجی پشتیبانی‌شده را می‌سازد؛
- فایل‌های اجرایی ساخته‌شده را به‌مدت هفت روز در بخش artifact نگه می‌دارد؛
- اگر commit جدیدی به همان شاخه برسد، اجرای قدیمی و بی‌استفاده را لغو می‌کند.

این گردش کار فقط اجازهٔ خواندن محتوای مخزن را دارد.

## انتشار نسخه

هنگامی که یک برچسب نسخه با حرف `v` به مخزن ارسال شود، فایل `release.yml` نسخهٔ قابل دانلود می‌سازد:

<div dir="ltr" align="left">

```text
git tag -a v1.0.0 -m "NextCmd v1.0.0"
git push origin v1.0.0
```

</div>

گردش کار ابتدا بررسی می‌کند برچسب از قالب نسخه‌گذاری معنایی، مانند `v1.0.0`، پیروی کند. سپس تست‌ها و ساخت همهٔ سیستم‌عامل‌ها را اجرا می‌کند. خروجی ویندوز در فایل ZIP و خروجی لینوکس و macOS در فایل `tar.gz` قرار می‌گیرد. در پایان checksum از نوع SHA-256 و یادداشت انتشار به‌صورت خودکار ساخته و همراه فایل‌ها منتشر می‌شوند.

می‌توان انتشار را از صفحهٔ Actions دوباره اجرا کرد، اما برچسب واردشده باید از قبل وجود داشته باشد. مرحلهٔ انتشار فقط مجوز `contents: write` دارد و به secret جداگانه نیاز ندارد. احراز هویت با `GITHUB_TOKEN` محدودشدهٔ GitHub Actions انجام می‌شود.

</div>
