# Contributing

Keep changes explicit, cross-platform, deterministic, and easy to test. The SDK must remain UI-independent and standard-library-only. Tool-specific knowledge belongs in its plugin. Avoid global registration and shell command strings.

Before submitting a change, run `gofmt -w .`, `go vet ./...`, `go test ./...`, and `go test -race ./...`. Add unit tests with injected fakes; integration tests may create local repositories with `t.TempDir()` but must not require GitHub or network access. Document new SDK contracts and user-visible controls.

Every delivered change should include a concise suggested commit message that describes that change set.

Documentation is English-first and bilingual. Put the Persian section after the English section and wrap it in `<div dir="rtl" align="right">` and `</div>`. Wrap fenced code, command output, directory trees, and standalone English text inside that section in `<div dir="ltr" align="left">`. Keep inline code, commands, paths, and identifiers in backticks so mixed-direction text remains readable.

Generated binaries, Make build artifacts, local verification caches, editor metadata, and fallback local configuration/history files are excluded through the repository `.gitignore`.

---

<div dir="rtl" align="right">

# مشارکت در پروژه

تغییرات باید روشن، چندسکویی، قابل‌پیش‌بینی و قابل تست باشند. SDK نباید به رابط کاربری وابسته شود و باید فقط از کتابخانهٔ استاندارد Go استفاده کند. هر منطق مخصوص یک ابزار باید در افزونهٔ همان ابزار قرار گیرد. از ثبت سراسری و مخفی، رفتارهای جادویی و نگهداری دستور پوسته در یک رشته خودداری کنید.

پیش از ارسال تغییر، دستورهای `gofmt -w .`، `go vet ./...`، `go test ./...` و `go test -race ./...` را اجرا کنید. در تست واحد، وابستگی‌ها را از بیرون تزریق و نسخهٔ ساختگی آن‌ها را استفاده کنید. تست یکپارچه می‌تواند با `t.TempDir()` یک مخزن موقت محلی بسازد، اما نباید به GitHub یا شبکه نیاز داشته باشد.

هر قرارداد جدید SDK و هر تغییر قابل مشاهده برای کاربر را مستند کنید. برای هر مجموعه تغییر نیز یک پیام commit کوتاه و مرتبط پیشنهاد دهید.

در مستندات دوزبانه، بخش انگلیسی ابتدا و بخش فارسی پس از آن قرار می‌گیرد. بلوک فارسی باید با `<div dir="rtl" align="right">` راست‌به‌چپ و راست‌چین شود. بلوک‌های کد، خروجی دستور، درخت پوشه‌ها و متن مستقل انگلیسی باید بیرون از جهت فارسی یا در `<div dir="ltr" align="left">` قرار گیرند. دستورها، مسیرها و شناسه‌های کوتاه درون جمله نیز باید داخل backtick باشند تا ترتیب نمایش آن‌ها به هم نریزد.

فایل‌های اجرایی تولیدشده، خروجی‌های Make، حافظه‌های موقت بررسی، تنظیمات ویرایشگر و فایل‌های محلی تنظیمات و تاریخچه در `.gitignore` قرار دارند و نباید commit شوند.

</div>
