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

## Terminal theme

The cyan prompt and arrow identify the editor and selected row. Each suggestion includes a compact kind badge (`COMP`, `NEXT`, `TIP`, or `FIX`), a color-coded risk, and its source plugin. After execution, a green or red status line reports the exit code and duration. Set `NO_COLOR` to any non-empty value for plain output; colors are also omitted when output is redirected.

---

<div dir="rtl">

# راهنمای تعاملی

NextCmd بدون باز کردن browser یا اجرای command خارجی، راهنمای داخلی ارائه می‌دهد.

## راهنمای عمومی

یکی از فرم‌های زیر را تایپ و Enter را فشار دهید:

<div dir="ltr">

```text
:?
:؟
```

</div>

خروجی، کنترل‌های keyboard، commandهای داخلی، روش‌های خروج و تمام Pluginهای بارگذاری‌شده را توضیح می‌دهد.

## کاتالوگ commandهای Plugin

برای نمایش تمام commandهای ثابت یک Plugin همراه description و risk، شناسه آن را اضافه کنید:

<div dir="ltr">

```text
:? git
:? dotnet
:؟ git
:؟ dotnet
```

</div>

کاتالوگ از طریق capability عمومی `sdk.HelpProvider` توسط خود Plugin ارائه می‌شود. Core و Terminal هیچ فهرست مخصوص Git یا .NET ندارند. مقادیر dynamic مانند branch، file، remote، solution و project واقعی از طریق completion عادی نمایش داده می‌شوند.

## prefix نام executable

پیشنهادها پیش از کامل شدن نام executable ظاهر می‌شوند. برای نمونه `g` و `gi` پیشنهادهای Git و `dot`، `dotn` و `dotnet` پیشنهادهای .NET را نمایش می‌دهند. ranking نهایی همچنان قطعی و بر اساس تمام متن فعلی است.

## ظاهر Terminal

prompt و فلش فیروزه‌ای، editor و سطر انتخاب‌شده را مشخص می‌کنند. هر suggestion یک badge کوتاه برای نوع (`COMP`، `NEXT`، `TIP` یا `FIX`)، ریسک رنگی و شناسهٔ Plugin منبع دارد. پس از اجرا نیز یک سطر سبز یا قرمز exit code و مدت اجرا را نمایش می‌دهد. برای خروجی ساده، متغیر استاندارد `NO_COLOR` را روی یک مقدار غیرخالی قرار دهید؛ هنگام redirect شدن خروجی نیز رنگ‌ها خودکار حذف می‌شوند.

</div>
