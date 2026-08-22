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

## Workflow guidance

Successful restore suggests build without another restore. Successful build suggests test without rebuilding. Test completion suggests formatting verification and publishing. Dependency changes suggest restore and build. EF operations suggest reviewing migrations and running tests.

Best-practice suggestions include `dotnet format --verify-no-changes` and test execution when test projects exist. Recovery handles missing tools or SDK commands, missing restore assets, and directories without a detected project.

## Configuration

The plugin is enabled by default and can be disabled independently:

```json
{"dotnetEnabled": false}
```

---

# پلاگین .NET

پلاگین built-in مربوط به .NET پیشنهادهای قطعی برای CLI چندسکویی `dotnet` ارائه می‌دهد. این Plugin فقط از SDK عمومی و Go Standard Library استفاده می‌کند و Core هیچ نوع یا قانون مخصوص .NET ندارد.

## تشخیص workspace

Detector در parent directoryها به‌دنبال فایل‌های `.sln`، `.slnx`، `.csproj`، `.fsproj` و `.vbproj` می‌گردد و سپس workspace را scan می‌کند. پوشه‌های تولیدی یا سنگین مانند `bin`، `obj`، `.git`، `node_modules` و `target` نادیده گرفته می‌شوند. مسیر solutionها، مسیر و زبان projectها، test projectها و فایل‌های configuration مشترک ثبت می‌شوند و Core نتیجه را cache می‌کند.

## پیشنهادها

نسخه اول موارد زیر را پشتیبانی می‌کند:

- اطلاعات SDK و بررسی وضعیت آن؛
- ساخت solution، console app، class library، Web API، web app و xUnit project؛
- workflowهای restore، build، run، watch، test، clean، publish، pack و format؛
- عضویت project در solution و project referenceها؛
- عملیات NuGet با فرم سازگار verb-first و فرم noun-first مربوط به .NET 10؛
- restore و list ابزارهای local؛
- دستورهای رایج migration و database در Entity Framework Core.

مسیر projectها برای build، clean، restore، run، publish و pack به‌شکل پویا تکمیل می‌شود. `dotnet test` فقط test projectهای تشخیص‌داده‌شده را پیشنهاد می‌دهد.

## راهنمای workflow

بعد از restore موفق، build بدون restore مجدد پیشنهاد می‌شود. بعد از build موفق، test بدون build مجدد پیشنهاد می‌شود. پس از test، بررسی format و publish پیشنهاد می‌شوند. تغییر dependency به restore و build منجر می‌شود و عملیات EF پیشنهاد بررسی migration و اجرای test می‌دهد.

Best Practiceها شامل `dotnet format --verify-no-changes` و اجرای test در صورت وجود test project هستند. Recovery نیز نبود tool یا SDK command، نبود restore assets و دایرکتوری فاقد project را پوشش می‌دهد.

## تنظیمات

Plugin به‌صورت پیش‌فرض فعال است و مستقل از Git غیرفعال می‌شود:

```json
{"dotnetEnabled": false}
```
