# Git plugin

The Git plugin recognizes common commands for status, staging, commits, history, branches, stashes, remotes, tags, integration, worktrees, submodules, configuration, repository creation, and cloning. The expanded catalog includes show, reflog, blame, grep, annotated tags, cherry-pick, revert, safe reset/clean variants, worktree, and submodule commands.

Branch creation is a first-class workflow. Typing `git switch -c` suggests editable conventions such as `feature/<name>`, `bugfix/<name>`, `hotfix/<name>`, `release/<version>`, `chore/<name>`, `docs/<name>`, and `refactor/<name>`. These conventions are optional and remain editable before execution.

Detection runs Git directly and parses repository membership, current branch, porcelain status, local branches, remotes, upstream, and commits ahead. Core caches the returned state for one second. Dynamic completion covers local branches for switch/checkout/delete/merge/rebase, changed files for add/restore/blame, and remotes for push/pull/fetch. Branch deletion preserves its flag and is marked destructive.

Priority metadata reflects actual state: modified files raise diff/add relevance, staged files raise cached-diff/commit relevance, and ahead commits raise push relevance. Core still owns final ranking. Successful add/commit/pull/status commands produce contextual next actions. Reviewing staged changes is a first-class best practice. Failed repository commands can suggest init, while pathspec/revision failures suggest real branches.

Successful branch switches can recommend publishing a new branch with upstream tracking. Stash, merge, rebase, cherry-pick, revert, and fetch operations also provide contextual next actions.

Limitations: porcelain rename paths receive only minimal parsing; commit hash, tag-name, and remote-branch completion are not yet dynamic; the short cache can briefly show stale state; recovery intentionally recognizes only common failure text.

---

<div dir="rtl">

# Git Plugin

Git Plugin مجموعه گسترده‌ای از دستورهای status، staging، commit، history، branch، stash، remote، tag، integration، worktree، submodule، configuration، ساخت repository و clone را می‌شناسد. دستورهای show، reflog، blame، grep، annotated tag، cherry-pick، revert، reset/clean امن، worktree و submodule نیز پشتیبانی می‌شوند.

ساخت branch یک workflow مستقل است. با تایپ `git switch -c` الگوهای قابل‌ویرایش `feature/<name>`، `bugfix/<name>`، `hotfix/<name>`، `release/<version>`، `chore/<name>`، `docs/<name>` و `refactor/<name>` پیشنهاد می‌شوند. این الگوها اجباری نیستند و قبل از اجرا قابل تغییر باقی می‌مانند.

تشخیص context با اجرای مستقیم Git انجام می‌شود و عضویت در repository، branch فعلی، porcelain status، branchهای محلی، remoteها، upstream و commitهای جلوتر را parse می‌کند. Core این وضعیت را برای یک ثانیه cache می‌کند. تکمیل پویا برای branchهای محلی در switch/checkout/delete/merge/rebase، فایل‌های تغییریافته در add/restore/blame و remoteها در push/pull/fetch وجود دارد.

metadata مربوط به priority از وضعیت واقعی repository تأثیر می‌گیرد: وجود فایل modified اهمیت diff/add، وجود فایل staged اهمیت cached-diff/commit و commit جلوتر اهمیت push را افزایش می‌دهد. با این حال ranking نهایی همیشه در Core انجام می‌شود.

دستورهای موفق add، commit، pull و status اقدام‌های بعدی وابسته به context تولید می‌کنند. بررسی staged changes یک Best Practice مستقل است. شکست دستورهای وابسته به repository می‌تواند `git init` را پیشنهاد دهد و خطاهای pathspec/revision می‌توانند branchهای واقعی را نمایش دهند.

پس از ساخت یا تعویض branch، Plugin می‌تواند انتشار branch و تنظیم upstream را پیشنهاد دهد. برای stash، merge، rebase، cherry-pick، revert و fetch نیز اقدام بعدی وابسته به context ارائه می‌شود.

محدودیت‌ها: مسیر rename در porcelain به‌شکل ساده parse می‌شود؛ تکمیل commit hash، نام tag و remote branch هنوز پویا نیست؛ cache کوتاه ممکن است برای لحظه‌ای وضعیت قدیمی نشان دهد؛ recovery عمداً فقط خطاهای رایج را می‌شناسد.

</div>
