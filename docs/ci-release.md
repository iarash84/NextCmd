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

<div dir="rtl">

# یکپارچه‌سازی و انتشار

NextCmd دو workflow در مسیر `.github/workflows` دارد.

## CI

فایل `ci.yml` برای push و pull request روی branch اصلی `main` اجرا می‌شود و امکان اجرای دستی نیز دارد. این workflow:

- تست و `go vet` را به‌صورت native روی Ubuntu، Windows و macOS اجرا می‌کند؛
- فرمت `gofmt` و race detector را روی Ubuntu بررسی می‌کند؛
- با `make build-all` هر شش خروجی پشتیبانی‌شده را می‌سازد؛
- binaryها را برای هفت روز به‌عنوان artifact نگهداری می‌کند؛
- با رسیدن commit جدید، اجرای قدیمی همان branch را لغو می‌کند.

دسترسی این workflow به repository فقط خواندنی است.

## Release

فایل `release.yml` هنگام push شدن tag نسخه‌ای که با `v` شروع شود release می‌سازد:

<div dir="ltr">

```text
git tag -a v1.0.0 -m "NextCmd v1.0.0"
git push origin v1.0.0
```

</div>

Workflow معتبر بودن tag و قالب semantic version را بررسی می‌کند، تست‌ها را اجرا می‌کند، تمام targetها را می‌سازد، برای Windows فایل ZIP و برای Linux/macOS فایل `tar.gz` ایجاد می‌کند، checksum از نوع SHA-256 می‌سازد و همه خروجی‌ها را همراه release notes خودکار منتشر می‌کند.

امکان اجرای دوباره از صفحه Actions نیز وجود دارد، اما tag واردشده باید از قبل وجود داشته باشد. job انتشار فقط مجوز `contents: write` دارد و به secret جداگانه نیاز ندارد؛ احراز هویت با `GITHUB_TOKEN` محدودشده GitHub Actions انجام می‌شود.

</div>
