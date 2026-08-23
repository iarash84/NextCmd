# Contributing

Keep changes explicit, cross-platform, deterministic, and easy to test. The SDK must remain UI-independent and standard-library-only. Tool-specific knowledge belongs in its plugin. Avoid global registration and shell command strings.

Before submitting a change, run `gofmt -w .`, `go vet ./...`, `go test ./...`, and `go test -race ./...`. Add unit tests with injected fakes; integration tests may create local repositories with `t.TempDir()` but must not require GitHub or network access. Document new SDK contracts and user-visible controls.

Every delivered change should include a concise suggested commit message that describes that change set.

Documentation is English-first and bilingual. Put the Persian section after the English section and wrap the complete Persian block in `<div dir="rtl">` and `</div>`. Keep code, commands, paths, and identifiers in backticks so mixed right-to-left and left-to-right text remains readable.

Generated binaries, Make build artifacts, local verification caches, editor metadata, and fallback local configuration/history files are excluded through the repository `.gitignore`.

---

<div dir="rtl">

# مشارکت در پروژه

تغییرات را صریح، چندسکویی، قطعی و قابل تست نگه دارید. SDK باید مستقل از UI و محدود به Go Standard Library باقی بماند. منطق مخصوص هر ابزار در Plugin همان ابزار قرار می‌گیرد. از registration سراسری، magic و commandهای shell به‌شکل string خودداری کنید.

پیش از ارسال تغییر، دستورهای `gofmt -w .`، `go vet ./...`، `go test ./...` و `go test -race ./...` را اجرا کنید. برای unit test از dependencyهای تزریق‌شده و fake استفاده کنید. Integration test می‌تواند با `t.TempDir()` یک repository محلی ایجاد کند، اما نباید به GitHub یا شبکه نیاز داشته باشد.

قرارداد جدید SDK و رفتار قابل مشاهده برای کاربر را مستند کنید. همراه هر مجموعه تغییر تحویل‌شده، یک پیام commit پیشنهادی کوتاه و مرتبط نیز ارائه دهید.

مستندات باید ابتدا انگلیسی و سپس فارسی باشند. تمام بخش فارسی را بین تگ‌های `<div dir="rtl">` و `</div>` قرار دهید. کدها، commandها، مسیرها و identifierها را داخل backtick بنویسید تا ترکیب متن راست‌به‌چپ و چپ‌به‌راست خوانا باقی بماند.

فایل‌های binary تولیدشده، خروجی‌های build مربوط به Make، cacheهای بررسی محلی، metadata ویرایشگر و فایل‌های fallback مربوط به configuration/history توسط `.gitignore` از commit خارج می‌شوند.

</div>
