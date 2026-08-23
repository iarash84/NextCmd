# Docker Plugin

English | [فارسی](#افزونهٔ-docker)

The Docker plugin provides deterministic Docker CLI and Compose suggestions. It detects the nearest `Dockerfile`, `compose.yaml`, `compose.yml`, or legacy `docker-compose` file without contacting the Docker daemon.

It parses local Compose service names and multi-stage Dockerfile targets. These values complete commands such as `docker compose logs <service>`, `docker compose up <service>`, and `docker build --target <stage> .`. The static catalog covers build, run, images, containers, logs, exec, pull/push, Compose workflows, disk usage, and cleanup with appropriate risk metadata.

Successful builds suggest inspecting images and running the result. Starting Compose services suggests status and logs. Best practices validate Compose configuration and Docker build checks. Recovery covers an unavailable daemon, missing Compose configuration, and an unknown service.

The plugin is enabled by default. Disable it with `{"plugins":{"docker":false}}` and view its catalog with `:? docker`.

---

<div dir="rtl" align="right">

# افزونهٔ Docker

افزونهٔ Docker فرمان‌های Docker CLI و Compose را به‌شکلی ثابت پیشنهاد می‌دهد. نزدیک‌ترین فایل `Dockerfile`، `compose.yaml`، `compose.yml` یا فایل قدیمی `docker-compose` بدون اتصال به Docker daemon شناسایی می‌شود.

نام serviceهای Compose و stageهای Dockerfile چندمرحله‌ای از فایل‌های محلی خوانده می‌شوند. این مقادیر در فرمان‌هایی مانند `docker compose logs <service>`، `docker compose up <service>` و `docker build --target <stage> .` پیشنهاد داده می‌شوند. فهرست ثابت نیز build، run، imageها، containerها، log، exec، pull/push، گردش کار Compose، فضای دیسک و پاک‌سازی را با سطح خطر مناسب پوشش می‌دهد.

پس از build موفق، بررسی imageها و اجرای خروجی پیشنهاد می‌شود. پس از راه‌اندازی Compose نیز status و logها پیشنهاد داده می‌شوند. recovery در دسترس نبودن daemon، نبود فایل Compose و نام اشتباه service را پوشش می‌دهد.

افزونه به‌طور پیش‌فرض فعال است. برای غیرفعال‌کردن یا دیدن راهنما از نمونه‌های زیر استفاده کنید:

<div dir="ltr" align="left">

```json
{"plugins":{"docker":false}}
```

```text
:? docker
```

</div>

</div>
