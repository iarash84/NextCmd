# Architecture

The dependency direction is deliberately one-way:

```text
plugins ──> sdk <── internal core <── cmd/assistant
                                  ↑
                         plugins/builtin
```

The `sdk` contains stable data contracts and small optional capability interfaces. It imports only the Go standard library. A plugin must implement `Plugin.Info`; it may independently implement completion, detection, next-action, best-practice, recovery, or help-catalog capabilities. The completion engine uses type assertions, isolates provider errors, caches detection briefly, merges suggestions, and delegates final order to Core ranking.

Commands remain `Executable + Args`; rendering and execution are separate. Execution calls `exec.CommandContext` directly and never invokes a shell. Git-specific process calls and parsing live exclusively in `plugins/git`.

`internal/terminal` owns presentation and keyboard behavior. Its raw-mode boundary has separate Windows, Linux, and macOS files. No UI type appears in the SDK. `internal/history` writes portable JSON Lines and redacts common secret arguments and URL user-info before persistence.

Git, .NET, and Cargo built-ins are composed explicitly by `plugins/builtin.All`; removing any or all of them from that list leaves Core buildable. There is no `init` registration, reflection, mutable global registry, dynamic library, network access, or third-party dependency.

---

<div dir="rtl" align="right">

# معماری

وابستگی میان بخش‌های پروژه عمداً یک‌طرفه است:

<div dir="ltr" align="left">

```text
plugins --> sdk <-- internal core <-- cmd/assistant
                                  ^
                         plugins/builtin
```

</div>

بستهٔ عمومی `sdk` مدل‌های داده و رابط‌های کوچک و پایدار را تعریف می‌کند و فقط به کتابخانهٔ استاندارد Go وابسته است. هر افزونه فقط باید متد `Plugin.Info` را برای معرفی نام، شناسه و نسخهٔ خود پیاده‌سازی کند. قابلیت‌های دیگر، مانند تکمیل دستور، تشخیص پروژه، پیشنهاد گام بعدی، پیشنهاد روش بهتر، بازیابی پس از خطا و راهنمای دستورها، کاملاً اختیاری و مستقل از یکدیگرند.

موتور تکمیل با بررسی رابط‌های پیاده‌سازی‌شده می‌فهمد هر افزونه چه قابلیت‌هایی دارد. خطای یک افزونه فقط در حالت اشکال‌زدایی ثبت می‌شود و باعث بسته‌شدن برنامه نمی‌شود. نتیجهٔ تشخیص پروژه برای مدت کوتاهی نگهداری می‌شود تا با هر کلید فشرده‌شده دوباره محاسبه نشود. در پایان، هسته پیشنهادها را ترکیب و با الگوریتمی ثابت مرتب می‌کند.

هر دستور به‌صورت دو بخش جدا نگهداری می‌شود: نام فایل اجرایی در `Executable` و آرگومان‌ها در `Args`. متن قابل‌نمایش دستور جداگانه ساخته می‌شود و همان متن برای اجرا به پوسته فرستاده نمی‌شود. برنامه دستور را مستقیماً با `exec.CommandContext` اجرا می‌کند. تمام منطق اجرای Git و تفسیر خروجی آن فقط در `plugins/git` قرار دارد.

بستهٔ `internal/terminal` مسئول نمایش، رنگ‌ها و واکنش به صفحه‌کلید است. کد ورود به حالت تعاملی پایانه برای ویندوز، لینوکس و macOS جدا شده است. SDK هیچ وابستگی‌ای به رابط کاربری ندارد. بستهٔ `internal/history` تاریخچه را در قالب JSON Lines ذخیره می‌کند و پیش از ذخیره، رمزها، tokenها و اطلاعات ورود موجود در URL را می‌پوشاند.

افزونه‌های Git، .NET و Cargo به‌صورت صریح در تابع `plugins/builtin.All` ساخته می‌شوند. اگر هرکدام یا همهٔ آن‌ها از این فهرست حذف شوند، هسته همچنان بدون تغییر ساخته می‌شود. پروژه از ثبت مخفی با `init`، reflection، فهرست سراسری قابل‌تغییر، کتابخانهٔ پویا، دسترسی شبکه یا کتابخانهٔ جانبی استفاده نمی‌کند.

</div>
