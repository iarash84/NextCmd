# Architecture

The dependency direction is deliberately one-way:

```text
plugins ──> sdk <── internal core <── cmd/assistant
                                  ↑
                         plugins/builtin
```

The `sdk` contains stable data contracts and small optional capability interfaces. It imports only the Go standard library. A plugin must implement `Plugin.Info`; it may independently implement completion, detection, next-action, best-practice, recovery, or help-catalog capabilities. The completion engine uses type assertions, isolates provider errors, caches detection briefly, merges suggestions, and delegates final order to Core ranking.

Commands remain `Executable + Args`; rendering and execution are separate. Normal execution calls `exec.CommandContext` directly. A command explicitly prefixed with `!` is instead passed to `cmd.exe /D /S /C` on Windows or `/bin/sh -c` on Linux and macOS. Git-specific process calls and parsing live exclusively in `plugins/git`.

`internal/terminal` owns presentation and keyboard behavior. Its raw-mode boundary has separate Windows, Linux, and macOS files. No UI type appears in the SDK. `internal/history` writes portable JSON Lines and redacts common secret arguments and URL user-info before persistence.

Git, .NET, Cargo, Curl, Go, Docker, npm, and pip built-ins are composed explicitly by the parameterless `plugins/builtin.All`. Configuration enables or disables them through a generic map keyed by `Plugin.Info().ID`. Adding another built-in changes only the plugin package and this explicit composition list; it adds no config field, `main` argument, or Core branch. Removing any or all plugins leaves Core buildable. There is no `init` registration, reflection, mutable global registry, dynamic library, network access, or third-party dependency.

The application owns a mutable working-directory value without calling process-wide `os.Chdir`. The same value is passed to completion, detection, execution, and history, while the terminal only displays and updates its local editor context. This keeps directory changes explicit and race-free.

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

هر دستور عادی به‌صورت دو بخش جدا نگهداری می‌شود: نام فایل اجرایی در `Executable` و آرگومان‌ها در `Args`. متن قابل‌نمایش دستور جداگانه ساخته می‌شود و برنامه آن را مستقیماً با `exec.CommandContext` اجرا می‌کند. فقط دستورهایی که صریحاً با `!` شروع شوند در Windows از طریق `cmd.exe` و در Linux و macOS از طریق `/bin/sh` اجرا می‌شوند. تمام منطق اجرای Git و تفسیر خروجی آن فقط در `plugins/git` قرار دارد.

بستهٔ `internal/terminal` مسئول نمایش، رنگ‌ها و واکنش به صفحه‌کلید است. کد ورود به حالت تعاملی پایانه برای ویندوز، لینوکس و macOS جدا شده است. SDK هیچ وابستگی‌ای به رابط کاربری ندارد. بستهٔ `internal/history` تاریخچه را در قالب JSON Lines ذخیره می‌کند و پیش از ذخیره، رمزها، tokenها و اطلاعات ورود موجود در URL را می‌پوشاند.

افزونه‌های Git، .NET، Cargo، Curl، Go، Docker، npm و pip در تابع بدون آرگومان `plugins/builtin.All` به‌صورت صریح ساخته می‌شوند. تنظیمات با یک نقشهٔ عمومی و براساس `Plugin.Info().ID` آن‌ها را فعال یا غیرفعال می‌کند. برای افزودن افزونهٔ داخلی جدید فقط package افزونه و همین فهرست صریح تغییر می‌کنند؛ افزودن فیلد تنظیمات، آرگومان `main` یا شرط مخصوص در هسته لازم نیست. با حذف هر تعداد افزونه، هسته همچنان ساخته می‌شود. پروژه از ثبت مخفی با `init`، reflection، فهرست سراسری قابل‌تغییر، کتابخانهٔ پویا، دسترسی شبکه یا کتابخانهٔ جانبی استفاده نمی‌کند.

برنامه مسیر کاری را در یک مقدار داخلی نگه می‌دارد و `os.Chdir` سراسری را فراخوانی نمی‌کند. همان مسیر به تکمیل دستور، تشخیص پروژه، اجرای process و تاریخچه داده می‌شود. رابط پایانه فقط مسیر را نمایش می‌دهد و context ویرایشگر خود را به‌روزرسانی می‌کند. به این ترتیب تغییر مسیر صریح و بدون race condition باقی می‌ماند.

</div>
