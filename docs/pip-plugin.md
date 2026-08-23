# pip / pip3 Plugin

English | [فارسی](#افزونهٔ-pip--pip3)

One plugin supports both `pip` and `pip3`. It detects Python projects from `pyproject.toml`, setup files, `Pipfile`, and `requirements*.txt`. Local requirement files, declared package names, and common virtual-environment directories are cached without invoking Python or accessing a package index.

Requirement files complete `pip install -r <file>` and `pip3 install -r <file>`. Declared packages complete `show` and `uninstall`. The catalog covers install, upgrade, editable installs, uninstall, list/outdated, show, check, freeze, download, wheel, cache, configuration diagnostics, and version lookup.

Environment changes suggest `pip check`. Recovery covers incompatible package versions, permission failures, TLS configuration problems, and externally managed Python environments. Suggestions never bypass certificate validation automatically.

The plugin ID is `pip`; disabling it disables both executable aliases. Use `{"plugins":{"pip":false}}` or show both catalogs with `:? pip`.

---

<div dir="rtl" align="right">

# افزونهٔ pip / pip3

یک افزونه هر دو فرمان `pip` و `pip3` را پشتیبانی می‌کند. پروژه‌های Python با فایل‌های `pyproject.toml`، فایل‌های setup، `Pipfile` و `requirements*.txt` تشخیص داده می‌شوند. نام فایل‌های requirements، packageهای ثبت‌شده و پوشه‌های رایج virtual environment بدون اجرای Python یا دسترسی به package index در cache قرار می‌گیرند.

فایل‌های requirements برای `pip install -r <file>` و `pip3 install -r <file>` پیشنهاد داده می‌شوند. packageهای ثبت‌شده نیز فرمان‌های `show` و `uninstall` را کامل می‌کنند. فهرست ثابت شامل install، upgrade، editable install، uninstall، list/outdated، show، check، freeze، download، wheel، cache و بررسی تنظیمات است.

پس از تغییر محیط، `pip check` پیشنهاد می‌شود. recovery نسخهٔ ناسازگار package، خطای دسترسی، مشکل تنظیمات TLS و محیط Python مدیریت‌شده توسط سیستم‌عامل را پوشش می‌دهد. افزونه هیچ‌گاه غیرفعال‌کردن بررسی certificate را خودکار پیشنهاد نمی‌دهد.

شناسهٔ افزونه `pip` است و غیرفعال‌کردن آن هر دو فرمان را غیرفعال می‌کند:

<div dir="ltr" align="left">

```json
{"plugins":{"pip":false}}
```

```text
:? pip
```

</div>

</div>
