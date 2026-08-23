# npm Plugin

English | [فارسی](#افزونهٔ-npm)

The npm plugin finds the nearest `package.json`, preferring an ancestor workspace root. It reads package scripts, dependency names, local workspace package names, `package-lock.json`, and `node_modules` state with the Go standard library. `node_modules` and generated directories are skipped during scanning.

Dynamic completion supports `npm run <script>`, declared dependencies for `npm uninstall`, and local packages for `npm --workspace <name> run <script>`. The catalog covers install/CI, scripts, tests, dependency maintenance, audit, cache operations, execution, and package publishing. Publish and destructive cache operations receive elevated risk labels.

Dependency changes suggest tests and an audit. Best practices recommend tests, audit, and reproducible `npm ci` installs when a lockfile exists. Recovery handles a missing `package.json`, an unknown script, and a missing or incompatible lockfile.

The plugin is enabled by default. Disable it with `{"plugins":{"npm":false}}` and view its catalog with `:? npm`.

---

<div dir="rtl" align="right">

# افزونهٔ npm

افزونهٔ npm نزدیک‌ترین `package.json` را پیدا می‌کند و اگر workspace در پوشهٔ والد وجود داشته باشد، ریشهٔ workspace را ترجیح می‌دهد. scriptها، dependencyها، نام packageهای محلی، وجود `package-lock.json` و وضعیت `node_modules` فقط با کتابخانهٔ استاندارد Go خوانده می‌شوند. هنگام بررسی نیز `node_modules` و پوشه‌های تولیدشده نادیده گرفته می‌شوند.

تکمیل پویا برای `npm run <script>`، dependencyهای ثبت‌شده در `npm uninstall` و packageهای محلی در `npm --workspace <name> run <script>` کار می‌کند. فهرست ثابت نصب، CI، scriptها، تست، نگهداری dependency، audit، cache، اجرای ابزار و انتشار package را پوشش می‌دهد. انتشار و پاک‌سازی‌های حذف‌کننده سطح خطر بیشتری دارند.

پس از تغییر dependencyها، اجرای تست و audit پیشنهاد می‌شود. recovery نبود `package.json`، script ناشناخته و lockfile ناموجود یا ناسازگار را پوشش می‌دهد.

<div dir="ltr" align="left">

```json
{"plugins":{"npm":false}}
```

```text
:? npm
```

</div>

</div>
