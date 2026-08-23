# Interactive Help

English | [فارسی](#راهنمای-تعاملی)

NextCmd provides built-in help without opening a browser or invoking an external command.

## General help

Type either form and press Enter:

```text
:?
:؟
```

The output explains keyboard controls, built-in commands, exit commands, and all loaded plugins.

## Plugin command catalogs

Append a plugin ID to print every statically supported command with its description and risk:

```text
:? git
:? dotnet
:؟ git
:؟ dotnet
```

The catalog comes from the plugin through the public `sdk.HelpProvider` capability. Core and Terminal do not contain Git or .NET command lists. Dynamic values such as actual branches, files, remotes, solutions, and projects remain available through normal completion.

## Executable prefixes

Plugin suggestions appear before the executable name is complete. For example, `g` and `gi` show Git suggestions, while `dot`, `dotn`, and `dotnet` show .NET suggestions. Final ranking remains deterministic and uses the full current input.

---

<div dir="rtl">

# راهنمای تعاملی

NextCmd بدون باز کردن browser یا اجرای command خارجی، راهنمای داخلی ارائه می‌دهد.

## راهنمای عمومی

یکی از فرم‌های زیر را تایپ و Enter را فشار دهید:

```text
:?
:؟
```

خروجی، کنترل‌های keyboard، commandهای داخلی، روش‌های خروج و تمام Pluginهای بارگذاری‌شده را توضیح می‌دهد.

## کاتالوگ commandهای Plugin

برای نمایش تمام commandهای ثابت یک Plugin همراه description و risk، شناسه آن را اضافه کنید:

```text
:? git
:? dotnet
:؟ git
:؟ dotnet
```

کاتالوگ از طریق capability عمومی `sdk.HelpProvider` توسط خود Plugin ارائه می‌شود. Core و Terminal هیچ فهرست مخصوص Git یا .NET ندارند. مقادیر dynamic مانند branch، file، remote، solution و project واقعی از طریق completion عادی نمایش داده می‌شوند.

## prefix نام executable

پیشنهادها پیش از کامل شدن نام executable ظاهر می‌شوند. برای نمونه `g` و `gi` پیشنهادهای Git و `dot`، `dotn` و `dotnet` پیشنهادهای .NET را نمایش می‌دهند. ranking نهایی همچنان قطعی و بر اساس تمام متن فعلی است.

</div>
