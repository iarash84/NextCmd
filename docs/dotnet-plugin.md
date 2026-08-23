# .NET Plugin

English | [فارسی](#پلاگین-net)

The built-in .NET plugin provides deterministic suggestions for the cross-platform `dotnet` CLI. It uses only the public SDK and the Go standard library; Core contains no .NET-specific type or rule.

## Workspace detection

The detector searches parent directories for `.sln`, `.slnx`, `.csproj`, `.fsproj`, and `.vbproj` files, then scans the workspace while skipping generated or expensive directories such as `bin`, `obj`, `.git`, `node_modules`, and `target`. It records solution paths, project paths, project languages, test projects, and common repository-level configuration files. Detection is cached by Core.

## Suggestions

The first version supports:

- SDK information and health checks;
- creation of solutions, console apps, class libraries, Web APIs, web apps, and xUnit projects;
- restore, build, run, watch, test, clean, publish, pack, and format workflows;
- solution membership and project references;
- NuGet package operations in the compatible verb-first form and the .NET 10 noun-first form;
- local tool restore/list operations;
- common Entity Framework Core migration and database commands.

Project paths are completed dynamically for build, clean, restore, run, publish, and pack. `dotnet test` prioritizes and completes only detected test projects.

Core commands such as `dotnet build`, `dotnet test`, and `dotnet run` remain available even when no workspace is detected. Their context priority is slightly lower and the suggestion explains that a project or solution path may be needed.

## Workflow guidance

Successful restore suggests build without another restore. Successful build suggests test without rebuilding. Test completion suggests formatting verification and publishing. Dependency changes suggest restore and build. EF operations suggest reviewing migrations and running tests.

Best-practice suggestions include `dotnet format --verify-no-changes` and test execution when test projects exist. Recovery handles missing tools or SDK commands, missing restore assets, and directories without a detected project.

## Configuration

The plugin is enabled by default and can be disabled independently:

```json
{"plugins":{"dotnet":false}}
```

---

<div dir="rtl" align="right">

# پلاگین .NET

افزونهٔ داخلی .NET برای ابزار چندسکویی `dotnet` پیشنهادهای ثابت و قابل‌پیش‌بینی ارائه می‌دهد. این افزونه فقط از SDK عمومی و کتابخانهٔ استاندارد Go استفاده می‌کند. هسته هیچ نوع داده یا قانون مخصوص .NET ندارد.

## تشخیص فضای کاری

تشخیص‌دهنده ابتدا در پوشهٔ فعلی و پوشه‌های والد به‌دنبال فایل‌های `.sln`، `.slnx`، `.csproj`، `.fsproj` و `.vbproj` می‌گردد. پس از یافتن فضای کاری، پروژه‌های موجود را بررسی می‌کند و پوشه‌های تولیدی یا سنگین مانند `bin`، `obj`، `.git`، `node_modules` و `target` را نادیده می‌گیرد. مسیر solutionها، زبان و مسیر پروژه‌ها، پروژه‌های تست و فایل‌های تنظیمات مشترک ثبت می‌شوند. هسته این نتیجه را برای مدت کوتاهی نگه می‌دارد تا بررسی با هر کلید تکرار نشود.

## پیشنهادها

نسخه اول موارد زیر را پشتیبانی می‌کند:

- اطلاعات SDK و بررسی وضعیت آن؛
- ساخت solution، برنامهٔ کنسولی، class library، Web API، برنامهٔ وب و پروژهٔ xUnit؛
- دستورهای `restore`، `build`، `run`، `watch`، `test`، `clean`، `publish`، `pack` و `format`؛
- افزودن پروژه به solution و مدیریت reference میان پروژه‌ها؛
- عملیات NuGet با فرم سازگار verb-first و فرم noun-first مربوط به .NET 10؛
- restore و list ابزارهای local؛
- دستورهای رایج migration و پایگاه داده در Entity Framework Core.

مسیر پروژه‌ها برای دستورهای ساخت، پاک‌سازی، بازیابی وابستگی‌ها، اجرا، انتشار و بسته‌بندی به‌صورت پویا تکمیل می‌شود. برای `dotnet test` فقط پروژه‌هایی پیشنهاد می‌شوند که به‌عنوان پروژهٔ تست تشخیص داده شده‌اند.

دستورهای عمومی مانند `dotnet build`، `dotnet test` و `dotnet run` حتی بیرون از فضای کاری شناسایی‌شده نیز نمایش داده می‌شوند. در این حالت امتیاز آن‌ها کمی کمتر است و توضیح پیشنهاد یادآوری می‌کند که شاید لازم باشد مسیر پروژه یا solution مشخص شود.

## پیشنهاد گام بعدی

پس از بازیابی موفق وابستگی‌ها، ساخت پروژه بدون بازیابی دوباره پیشنهاد می‌شود. بعد از ساخت موفق، اجرای تست بدون ساخت مجدد پیشنهاد می‌شود. پس از تست نیز بررسی قالب‌بندی و انتشار در دسترس قرار می‌گیرند. تغییر وابستگی باعث پیشنهاد بازیابی و ساخت می‌شود و پس از عملیات Entity Framework، بررسی migration و اجرای تست پیشنهاد می‌شود.

پیشنهادهای روش بهتر شامل اجرای `dotnet format --verify-no-changes` و اجرای تست در صورت وجود پروژهٔ تست هستند. پیشنهادهای رفع خطا نیز نصب‌نبودن ابزار یا دستور SDK، ساخته‌نشدن فایل‌های بازیابی وابستگی‌ها و نبود پروژه در پوشهٔ فعلی را پوشش می‌دهند.

## تنظیمات

این افزونه به‌صورت پیش‌فرض فعال است و می‌توان آن را مستقل از Git غیرفعال کرد:

<div dir="ltr" align="left">

```json
{"plugins":{"dotnet":false}}
```

</div>

</div>
