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
:? cargo
:؟ git
:؟ dotnet
:؟ cargo
```

The catalog comes from the plugin through the public `sdk.HelpProvider` capability. Core and Terminal do not contain Git or .NET command lists. Dynamic values such as actual branches, files, remotes, solutions, and projects remain available through normal completion.

## Executable prefixes

Plugin suggestions appear before the executable name is complete. For example, `g` and `gi` show Git suggestions, while `dot`, `dotn`, and `dotnet` show .NET suggestions. Final ranking remains deterministic and uses the full current input.

## Terminal theme

The cyan prompt and arrow identify the editor and selected row. Each suggestion includes a compact kind badge (`COMP`, `NEXT`, `TIP`, or `FIX`), a color-coded risk, and its source plugin. After execution, a green or red status line reports the exit code and duration. Set `NO_COLOR` to any non-empty value for plain output; colors are also omitted when output is redirected.

---

<div dir="rtl" align="right">

# راهنمای تعاملی

NextCmd راهنمای داخلی دارد؛ بنابراین برای دیدن کلیدها و دستورهای پشتیبانی‌شده لازم نیست مرورگر باز کنید یا برنامهٔ دیگری اجرا کنید.

## راهنمای عمومی

یکی از دو دستور زیر را تایپ کنید و Enter را فشار دهید:

<div dir="ltr" align="left">

```text
:?
:؟
```

</div>

برنامه کاربرد کلیدهای صفحه‌کلید، دستورهای داخلی، روش‌های خروج و نام افزونه‌های فعال را نمایش می‌دهد.

## فهرست دستورهای یک افزونه

برای دیدن همهٔ دستورهای ثابت یک افزونه، شناسهٔ آن را پس از `:?` بنویسید. نتیجه، متن دستور، توضیح کوتاه و میزان خطر آن را نمایش می‌دهد:

<div dir="ltr" align="left">

```text
:? git
:? dotnet
:? cargo
:؟ git
:؟ dotnet
:؟ cargo
```

</div>

خود افزونه این فهرست را از طریق رابط عمومی `sdk.HelpProvider` فراهم می‌کند. هسته و رابط پایانه هیچ فهرست ثابت و مخصوص Git یا .NET ندارند. مقادیر وابسته به پروژه، مانند نام شاخه، فایل، remote، solution یا project، در پیشنهادهای عادی و براساس وضعیت واقعی پروژه نمایش داده می‌شوند.

## پیشنهاد با نام ناقص ابزار

لازم نیست نام ابزار را کامل تایپ کنید. برای نمونه، `g` و `gi` پیشنهادهای Git را نشان می‌دهند؛ `dot`، `dotn` و `dotnet` نیز پیشنهادهای .NET را نمایش می‌دهند. ترتیب پیشنهادها همیشه با الگوریتمی ثابت و براساس تمام متن فعلی تعیین می‌شود.

## ظاهر پایانه

نشانه و فلش فیروزه‌ای، محل نوشتن دستور و پیشنهاد انتخاب‌شده را مشخص می‌کنند. کنار هر پیشنهاد یک برچسب کوتاه دیده می‌شود: `COMP` برای تکمیل، `NEXT` برای گام بعدی، `TIP` برای روش پیشنهادی و `FIX` برای رفع خطا. میزان خطر و شناسهٔ افزونهٔ پیشنهاددهنده نیز نمایش داده می‌شود. پس از اجرای دستور، یک سطر سبز یا قرمز کد خروج و مدت اجرا را نشان می‌دهد. برای غیرفعال‌کردن رنگ‌ها، متغیر استاندارد `NO_COLOR` را روی یک مقدار غیرخالی قرار دهید. هنگام انتقال خروجی به فایل یا برنامه‌ای دیگر نیز رنگ‌ها خودکار حذف می‌شوند.

</div>
