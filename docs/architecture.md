# Architecture

The dependency direction is deliberately one-way:

```text
plugins ──> sdk <── internal core <── cmd/assistant
                                  ↑
                         plugins/builtin
```

The `sdk` contains stable data contracts and small optional capability interfaces. It imports only the Go standard library. A plugin must implement `Plugin.Info`; it may independently implement completion, detection, next-action, best-practice, recovery, or help-catalog capabilities. The completion engine uses type assertions, isolates provider errors, caches detection briefly, merges suggestions, and delegates final order to Core ranking.

Commands remain `Executable + Args`; rendering and execution are separate. Normal execution calls `exec.CommandContext` directly. A command explicitly prefixed with `!` is instead passed to `cmd.exe /D /S /C` on Windows or `/bin/sh -c` on Linux and macOS. Git-specific process calls and parsing live exclusively in `plugins/git`.

Execution dependencies use the public `sdk.Runner` contract. Interactive execution depends on its `sdk.StreamingRunner` extension so output can still be mirrored live. Default constructors install the system executor; `app.NewWithRunner` and `git.NewWithRunner` accept deterministic test doubles or host-provided implementations.

The executor streams stdout and stderr to the terminal while simultaneously retaining both streams in `ExecutionResult` for recovery and next-action providers. During execution, the application derives a signal-aware context from its session context; Ctrl+C cancels the child command without canceling the application session.

Before execution, the application assesses the final direct or shell command syntax against a focused set of high-risk forms. Recognized commands require explicit confirmation with a default answer of no. A trailing `--yes` is treated as deliberate non-interactive approval and removed before an external or shell command is launched. Plugin risk metadata remains presentation and ranking input; this Core guard is independent so manually typed commands receive the same protection.

`internal/terminal` owns presentation and keyboard behavior. Its raw-mode boundary has separate Windows, Linux, and macOS files. No UI type appears in the SDK. `internal/history` writes portable JSON Lines and performs non-evaluating, best-effort redaction of structured arguments, shell words, sensitive environment values, authentication headers, and URL user-info before persistence. Redaction returns detached argument slices and environment maps so it never mutates the executed command.

Git, .NET, Cargo, Curl, Go, Docker, npm, and pip built-ins are composed explicitly by the parameterless `plugins/builtin.All`. Configuration enables or disables them through a generic map keyed by `Plugin.Info().ID`. Adding another built-in changes only the plugin package and this explicit composition list; it adds no config field, `main` argument, or Core branch. Removing any or all plugins leaves Core buildable. There is no `init` registration, reflection, mutable global registry, dynamic library, network access, or third-party dependency.

The application owns a mutable working-directory value without calling process-wide `os.Chdir`. The same value is passed to completion, detection, execution, and history, while the terminal only displays and updates its local editor context. This keeps directory changes explicit and race-free.

Core provides low-priority filesystem completion for the current argument of ordinary commands. It reads only the active local directory through Go's cross-platform filesystem APIs and merges those candidates with higher-priority plugin suggestions; it never invokes the command or a shell for generic path discovery.

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

وابستگی‌های اجرای فرمان از قرارداد عمومی `sdk.Runner` استفاده می‌کنند. اجرای تعاملی به رابط توسعه‌یافتهٔ `sdk.StreamingRunner` وابسته است تا خروجی همچنان زنده نمایش داده شود. سازنده‌های پیش‌فرض اجراکنندهٔ سیستم را قرار می‌دهند و `app.NewWithRunner` و `git.NewWithRunner` امکان تزریق fake قطعی یا پیاده‌سازی میزبان را فراهم می‌کنند.

اجراکننده خروجی عادی و خطا را هم‌زمان به terminal می‌فرستد و هر دو را در `ExecutionResult` نیز نگه می‌دارد تا افزونه‌های recovery و پیشنهاد گام بعدی همچنان به نتیجه دسترسی داشته باشند. هنگام اجرا، context حساس به signal از context نشست ساخته می‌شود؛ بنابراین Ctrl+C فرمان فرزند را لغو می‌کند، اما نشست برنامه باز می‌ماند.

پیش از اجرا، برنامه نحو نهایی فرمان مستقیم یا shell را با مجموعه‌ای محدود از الگوهای پرخطر بررسی می‌کند. فرمان شناخته‌شده به تأیید صریح با پاسخ پیش‌فرض منفی نیاز دارد. `--yes` در انتهای فرمان به‌عنوان تأیید آگاهانهٔ غیرتعاملی در نظر گرفته می‌شود و پیش از اجرای فرمان خارجی یا shell حذف می‌گردد. اطلاعات ریسک افزونه همچنان برای نمایش و مرتب‌سازی است؛ این محافظ هسته مستقل عمل می‌کند تا فرمان‌های تایپ‌شدهٔ دستی نیز همان حفاظت را داشته باشند.

بستهٔ `internal/terminal` مسئول نمایش، رنگ‌ها و واکنش به صفحه‌کلید است. کد ورود به حالت تعاملی پایانه برای ویندوز، لینوکس و macOS جدا شده است. SDK هیچ وابستگی‌ای به رابط کاربری ندارد. بستهٔ `internal/history` تاریخچه را در قالب JSON Lines ذخیره می‌کند و بدون ارزیابی shell، آرگومان‌های ساختاریافته، واژه‌های shell، مقادیر حساس محیط، headerهای احراز هویت و اطلاعات ورود URL را به‌صورت best-effort پاک‌سازی می‌کند. این عملیات کپی جداگانه‌ای از آرگومان‌ها و environment می‌سازد و فرمان اجراشده را تغییر نمی‌دهد.

افزونه‌های Git، .NET، Cargo، Curl، Go، Docker، npm و pip در تابع بدون آرگومان `plugins/builtin.All` به‌صورت صریح ساخته می‌شوند. تنظیمات با یک نقشهٔ عمومی و براساس `Plugin.Info().ID` آن‌ها را فعال یا غیرفعال می‌کند. برای افزودن افزونهٔ داخلی جدید فقط package افزونه و همین فهرست صریح تغییر می‌کنند؛ افزودن فیلد تنظیمات، آرگومان `main` یا شرط مخصوص در هسته لازم نیست. با حذف هر تعداد افزونه، هسته همچنان ساخته می‌شود. پروژه از ثبت مخفی با `init`، reflection، فهرست سراسری قابل‌تغییر، کتابخانهٔ پویا، دسترسی شبکه یا کتابخانهٔ جانبی استفاده نمی‌کند.

برنامه مسیر کاری را در یک مقدار داخلی نگه می‌دارد و `os.Chdir` سراسری را فراخوانی نمی‌کند. همان مسیر به تکمیل دستور، تشخیص پروژه، اجرای process و تاریخچه داده می‌شود. رابط پایانه فقط مسیر را نمایش می‌دهد و context ویرایشگر خود را به‌روزرسانی می‌کند. به این ترتیب تغییر مسیر صریح و بدون race condition باقی می‌ماند.

هسته برای آرگومان فعلی دستورهای عادی، تکمیل مسیر filesystem را با اولویت پایین ارائه می‌کند. مسیرها فقط با API چندسکویی Go از پوشهٔ کاری فعال خوانده و با پیشنهادهای دارای اولویت بیشتر افزونه‌ها ادغام می‌شوند؛ برای کشف عمومی مسیر هیچ دستور یا shell اجرا نمی‌شود.

</div>
