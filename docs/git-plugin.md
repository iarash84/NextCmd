# Git plugin

The Git plugin recognizes common commands for status, staging, commits, history, branches, stashes, remotes, tags, integration, worktrees, submodules, configuration, repository creation, and cloning. The expanded catalog includes show, reflog, blame, grep, annotated tags, cherry-pick, revert, safe reset/clean variants, worktree, and submodule commands.

Branch creation is a first-class workflow. Typing `git switch -c` suggests editable conventions such as `feature/<name>`, `bugfix/<name>`, `hotfix/<name>`, `release/<version>`, `chore/<name>`, `docs/<name>`, and `refactor/<name>`. These conventions are optional and remain editable before execution.

Detection runs Git directly and parses repository membership, current branch, porcelain status, local branches, remotes, upstream, and commits ahead. Core caches the returned state for one second. Dynamic completion covers local branches for switch/checkout/delete/merge/rebase, changed files for add/restore/blame, and remotes for push/pull/fetch. Branch deletion preserves its flag and is marked destructive.

Priority metadata reflects actual state: modified files raise diff/add relevance, staged files raise cached-diff/commit relevance, and ahead commits raise push relevance. Core still owns final ranking. Successful add/commit/pull/status commands produce contextual next actions. Reviewing staged changes is a first-class best practice. Failed repository commands can suggest init, while pathspec/revision failures suggest real branches.

Successful branch switches can recommend publishing a new branch with upstream tracking. Stash, merge, rebase, cherry-pick, revert, and fetch operations also provide contextual next actions.

Limitations: porcelain rename paths receive only minimal parsing; commit hash, tag-name, and remote-branch completion are not yet dynamic; the short cache can briefly show stale state; recovery intentionally recognizes only common failure text.

---

<div dir="rtl" align="right">

# افزونهٔ Git

افزونهٔ Git دستورهای رایج مربوط به مشاهدهٔ وضعیت، افزودن فایل، ثبت تغییر، تاریخچه، شاخه، stash، remote، tag، ادغام، worktree، submodule، تنظیمات، ساخت مخزن و clone را می‌شناسد. دستورهای `show`، `reflog`، `blame`، `grep`، `cherry-pick`، `revert` و حالت‌های امن `reset` و `clean` نیز پشتیبانی می‌شوند.

برای ساخت شاخه، با تایپ `git switch -c` چند الگوی قابل‌ویرایش مانند `feature/<name>`، `bugfix/<name>`، `hotfix/<name>` و `release/<version>` پیشنهاد می‌شود. این نام‌ها اجباری نیستند؛ پیشنهاد ابتدا وارد ویرایشگر می‌شود و کاربر می‌تواند پیش از اجرا آن را تغییر دهد.

افزونه با اجرای مستقیم Git تشخیص می‌دهد پوشهٔ فعلی داخل مخزن است یا نه. سپس شاخهٔ فعلی، فایل‌های تغییرکرده، شاخه‌های محلی، remoteها، upstream و تعداد commitهای ارسال‌نشده را از خروجی Git استخراج می‌کند. هسته این وضعیت را یک ثانیه نگه می‌دارد تا Git با هر کلید دوباره اجرا نشود. نام شاخه‌ها، فایل‌های مناسب و remoteها براساس دستور فعلی به‌صورت پویا پیشنهاد می‌شوند.

افزونه با توجه به وضعیت واقعی مخزن، اهمیت اولیهٔ هر پیشنهاد را مشخص می‌کند. وجود فایل تغییرکرده اهمیت `diff` و `add`، وجود فایل آماده‌شده برای ثبت اهمیت `diff --cached` و `commit` و وجود commit ارسال‌نشده اهمیت `push` را بیشتر می‌کند. ترتیب نهایی همیشه در هسته تعیین می‌شود.

پس از موفقیت دستورهایی مانند `add`، `commit`، `pull` و `status`، افزونه با توجه به وضعیت جدید مخزن گام بعدی را پیشنهاد می‌دهد. برای نمونه، پیش از ثبت تغییرات، بررسی فایل‌های آماده‌شده با `git diff --cached` پیشنهاد می‌شود. اگر دستور بیرون از مخزن اجرا شود، `git init` می‌تواند به‌عنوان راه‌حل نمایش داده شود. خطای نام شاخه نیز باعث پیشنهاد شاخه‌های واقعی می‌شود.

پس از ساخت یا تعویض شاخه، افزونه می‌تواند انتشار شاخه و تنظیم upstream را پیشنهاد دهد. برای `stash`، `merge`، `rebase`، `cherry-pick`، `revert` و `fetch` نیز گام بعدی مناسب با وضعیت مخزن ارائه می‌شود.

محدودیت‌های فعلی: مسیر فایل تغییرنام‌یافته به‌شکل ساده خوانده می‌شود؛ تکمیل پویای شناسهٔ commit، نام tag و شاخهٔ remote هنوز وجود ندارد؛ نگهداری کوتاه‌مدت وضعیت ممکن است برای لحظه‌ای اطلاعات قدیمی نشان دهد؛ پیشنهاد رفع خطا نیز عمداً فقط خطاهای رایج را پوشش می‌دهد.

</div>
