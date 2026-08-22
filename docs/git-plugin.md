# Git plugin

The Git MVP recognizes status, add, commit, diff, log, branch, switch, checkout, restore, stash, pull, push, fetch, remote, merge, rebase, init, and clone.

Detection runs Git directly and parses repository membership, current branch, porcelain status, local branches, remotes, upstream, and commits ahead. Core caches the returned state for one second. Dynamic completion covers local branches for switch/checkout/delete/merge/rebase, changed files for add/restore, and remotes for push/pull/fetch.

Priority metadata reflects actual state: modified files raise diff/add relevance, staged files raise cached-diff/commit relevance, and ahead commits raise push relevance. Core still owns final ranking. Successful add/commit/pull/status commands produce contextual next actions. Reviewing staged changes is a first-class best practice. Failed repository commands can suggest init, while pathspec/revision failures suggest real branches.

Limitations: porcelain rename paths receive only minimal parsing; remote branches and detailed option completion are not included; the short cache can briefly show stale state; recovery intentionally recognizes only common failure text.

---

# Git Plugin

نسخه MVP دستورهای status، add، commit، diff، log، branch، switch، checkout، restore، stash، pull، push، fetch، remote، merge، rebase، init و clone را می‌شناسد.

تشخیص context با اجرای مستقیم Git انجام می‌شود و عضویت در repository، branch فعلی، porcelain status، branchهای محلی، remoteها، upstream و commitهای جلوتر را parse می‌کند. Core این وضعیت را برای یک ثانیه cache می‌کند. تکمیل پویا برای branchهای محلی در switch/checkout/delete/merge/rebase، فایل‌های تغییریافته در add/restore و remoteها در push/pull/fetch وجود دارد.

metadata مربوط به priority از وضعیت واقعی repository تأثیر می‌گیرد: وجود فایل modified اهمیت diff/add، وجود فایل staged اهمیت cached-diff/commit و commit جلوتر اهمیت push را افزایش می‌دهد. با این حال ranking نهایی همیشه در Core انجام می‌شود.

دستورهای موفق add، commit، pull و status اقدام‌های بعدی وابسته به context تولید می‌کنند. بررسی staged changes یک Best Practice مستقل است. شکست دستورهای وابسته به repository می‌تواند `git init` را پیشنهاد دهد و خطاهای pathspec/revision می‌توانند branchهای واقعی را نمایش دهند.

محدودیت‌ها: مسیر rename در porcelain به‌شکل ساده parse می‌شود؛ remote branch و تکمیل جامع optionها وجود ندارد؛ cache کوتاه ممکن است برای لحظه‌ای وضعیت قدیمی نشان دهد؛ recovery عمداً فقط خطاهای رایج را می‌شناسد.
