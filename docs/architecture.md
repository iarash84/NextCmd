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

Git and .NET built-ins are composed explicitly by `plugins/builtin.All`; removing either or both from that list leaves Core buildable. There is no `init` registration, reflection, mutable global registry, dynamic library, network access, or third-party dependency.

---

<div dir="rtl">

# معماری

جهت dependencyها عمداً یک‌طرفه است:

```text
plugins --> sdk <-- internal core <-- cmd/assistant
                                  ^
                         plugins/builtin
```

package عمومی `sdk` شامل قراردادهای داده پایدار و interfaceهای کوچک و اختیاری مبتنی بر capability است و فقط از Go Standard Library استفاده می‌کند. هر Plugin تنها موظف به پیاده‌سازی `Plugin.Info` است و می‌تواند قابلیت‌های completion، detection، next action، best practice، recovery یا کاتالوگ Help را مستقل از یکدیگر ارائه دهد.

Completion Engine قابلیت‌ها را با type assertion تشخیص می‌دهد، خطای هر provider را از برنامه جدا نگه می‌دارد، نتیجه detection را برای مدت کوتاه cache می‌کند، suggestionها را ترکیب می‌کند و ترتیب نهایی را به ranking در Core می‌سپارد.

Command همیشه به‌صورت `Executable + Args` باقی می‌ماند و rendering از execution جدا است. اجرا مستقیماً با `exec.CommandContext` انجام می‌شود و shell فراخوانی نمی‌شود. تمام اجرای process و parsing مخصوص Git فقط در `plugins/git` قرار دارد.

`internal/terminal` مسئول نمایش و keyboard behavior است و برای Windows، Linux و macOS مرز raw-mode جداگانه دارد. هیچ نوع UI وارد SDK نشده است. `internal/history` اطلاعات را به‌شکل JSON Lines ذخیره و secretهای رایج و credential موجود در URL را پیش از ذخیره حذف می‌کند.

Git و .NET به‌صورت صریح توسط `plugins/builtin.All` ساخته می‌شوند. با حذف هرکدام یا هردو از این فهرست، Core همچنان build می‌شود. پروژه از registration مخفی با `init`، reflection، registry سراسری mutable، dynamic library، دسترسی شبکه و dependency خارجی استفاده نمی‌کند.

</div>
