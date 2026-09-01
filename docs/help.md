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

Type `:` in the command editor to open the built-in command palette. Continue typing to filter it—for example, `:pl` leaves `:plugins`. Use Up/Down to highlight an item and Tab, Right Arrow, or Enter to place it in the editable command line. Commands that require a value are shown with an editable placeholder, such as `:del <path>`, `:which <command>`, and `:cd <path>`. Because the catalog is small and local, all matching built-ins are shown without running plugins or project detection.

Keyboard controls:

| Key | Behavior |
|---|---|
| Up/Down | Highlight suggestions. |
| Tab/Right | Accept the highlighted suggestion. |
| Left/Right | Move the caret inside the command editor. |
| Enter | Accept a suggestion, then execute on the next press. |
| Backspace | Delete the previous character. |
| Escape | Clear the current command line. |
| Ctrl+C/Ctrl+D | Exit and clean the terminal UI. |

## Working directory commands

The active directory is printed above each prompt. Start NextCmd elsewhere with `nextcmd --directory <path>`. Use `pwd` or `:pwd` to print it again. Use `cd <path>` or `:cd <path>` to change the directory used by completion, project detection, execution, and history. Relative paths, absolute paths, quoted paths with spaces, `..`, and `~` are supported. Running `cd` without a path selects the user home directory.

Use `:ls` to list the active directory, or `:ls <path>` to inspect another directory without changing the active one. The output is deterministic: directories appear first, then files, with type and human-readable size columns. This is implemented with Go's cross-platform filesystem APIs and does not invoke Unix `ls` or Windows `dir`.

```text
:ls
:ls ..
:ls "project with spaces"
:del old.txt
:del "old directory"
```

## Utility commands

These commands are handled inside NextCmd. They are not passed to a shell, are not stored in command history, and do not receive plugin suggestions.

| Command | Purpose | Important behavior |
|---|---|---|
| `:history [count]` | Show recent commands executed through NextCmd. | Defaults to 20 entries and accepts 1 through 1000. It displays time, exit code, duration, plugin, working directory, and the redacted structured command. If history is disabled, it says so explicitly. |
| `:plugins` | Show every currently enabled plugin. | Plugins are sorted by ID and displayed with version, name, and description. Disabled plugins are not listed. |
| `:clear` | Clear the terminal screen and return the cursor to the top-left corner. | It only changes the current display and does not remove history or project state. |
| `:del <path>` | Delete a file or directory. | Resolves relative, absolute, quoted, and `~` paths against the active working directory. It detects whether the target is a file or directory, recursively removes directories, and asks which target to remove if both a matching file and directory are found. |
| `:config` | Show the effective runtime configuration. | Displays the configuration and history paths, history status, suggestion limit, debug status, and sorted plugin overrides. It does not print environment variables or secrets. |
| `:which <command>` | Locate the executable that NextCmd would find through the operating-system `PATH`. | Uses Go's cross-platform executable lookup, including Windows `PATHEXT`, and prints an absolute path when it can resolve one. |
| `:version` | Show build and platform information. | Displays the NextCmd version, Go version, OS, architecture, and the source revision when build metadata is available. Official release builds receive their version from the release tag. |

Examples:

```text
:history
:history 5
:plugins
:clear
:del old.txt
:config
:which dotnet
:version
```

## Plugin command catalogs

Append a plugin ID to print every statically supported command with its description and risk:

```text
:? git
:? dotnet
:? cargo
:? curl
:? go
:? docker
:? npm
:? pip
:؟ git
:؟ dotnet
:؟ cargo
:؟ curl
:؟ go
:؟ docker
:؟ npm
:؟ pip
```

The catalog comes from the plugin through the public `sdk.HelpProvider` capability. Core and Terminal do not contain Git or .NET command lists. Dynamic values such as actual branches, files, remotes, solutions, and projects remain available through normal completion.

## Executable prefixes

Plugin suggestions appear before the executable name is complete. For example, `g` and `gi` show Git suggestions, while `dot`, `dotn`, and `dotnet` show .NET suggestions. Final ranking remains deterministic and uses the full current input.

## Terminal theme

The cyan prompt and arrow identify the editor and selected row. Each suggestion includes a compact kind badge (`COMP`, `NEXT`, `TIP`, or `FIX`), a color-coded risk, and its source plugin. After execution, a green or red status line reports the exit code and duration. Set `NO_COLOR` to any non-empty value for plain output; colors are also omitted when output is redirected.

The metadata after a command can be read as follows:

```text
git add .  COMP MUTATING · git
│          │    │          └─ source plugin
│          │    └──────────── risk level
│          └───────────────── suggestion kind
└──────────────────────────── editable command
```

### Suggestion kinds

| Badge | SDK kind | Meaning | When it appears |
|---|---|---|---|
| `COMP` | `Completion` | Completes or expands the text currently being typed. | During normal command editing, including executable, option, and dynamic argument completion. |
| `NEXT` | `NextAction` | Suggests a logical follow-up to the previous successful command. | For example, test after build or commit after staging files. |
| `TIP` | `BestPractice` | Offers an optional review or safer workflow step. | For example, inspect staged changes before committing. It is advice, not a requirement. |
| `FIX` | `Recovery` | Suggests a possible correction after a failed command. | Only for recognized failures, such as a missing project or invalid branch. |

### Risk levels

| Badge | Default color | Meaning | Typical examples |
|---|---|---|---|
| `SAFE` | Green | Primarily reads or inspects information and is not expected to change important state. | Status, logs, headers, diffs, or version information. |
| `MUTATING` | Yellow | Creates or changes files, dependencies, repository state, build output, or a remote resource. | Build, format, add, commit, POST, upload, or package installation. |
| `DESTRUCTIVE` | Red | Can delete, overwrite, or discard data and may be difficult to undo. | Branch deletion, cleanup, reset variants, or an HTTP DELETE request. |
| `DANGEROUS` | Magenta | Has an unusually high security or data-loss risk and requires careful review. | Disabling TLS verification or publishing sensitive/external changes. |

The final `· git`, `· cargo`, `· curl`, or similar suffix identifies the plugin that produced the suggestion. Kind and risk badges are informational metadata supplied by plugins. The badges affect presentation; ranking uses separate priority and relevance metadata. They do not block execution, request confirmation, or replace reviewing the command yourself. With `NO_COLOR`, all labels remain visible as plain text.

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

برای دیدن فهرست فرمان‌های داخلی، در ویرایشگر دستور فقط `:` را تایپ کنید. با ادامهٔ تایپ، فهرست محدود می‌شود؛ برای نمونه، `:pl` فقط پیشنهاد `:plugins` را باقی می‌گذارد. با کلیدهای بالا و پایین یک مورد را انتخاب کنید و با Tab، جهت راست یا Enter آن را وارد خط فرمان قابل‌ویرایش کنید. فرمان‌هایی که به مقدار نیاز دارند با جای‌نگهدار نمایش داده می‌شوند؛ مانند `:which <command>` و `:cd <path>`. این فهرست کوچک کاملاً داخلی است و برای نمایش آن هیچ افزونه یا عملیات تشخیص پروژه اجرا نمی‌شود.

## دستورهای مسیر کاری

مسیر فعلی بالای هر prompt نمایش داده می‌شود. برای شروع در مسیر دیگر از `nextcmd --directory <path>` استفاده کنید. دستور `pwd` یا `:pwd` مسیر را دوباره چاپ می‌کند. با `cd <path>` یا `:cd <path>` می‌توان مسیر مورد استفاده برای پیشنهادها، تشخیص پروژه، اجرای دستورها و تاریخچه را تغییر داد. مسیر نسبی، مسیر کامل، مسیر نقل‌قول‌شده دارای فاصله، `..` و `~` پشتیبانی می‌شوند. اجرای `cd` بدون آرگومان، پوشهٔ خانگی کاربر را انتخاب می‌کند.

دستور `:ls` محتوای مسیر کاری فعلی را نمایش می‌دهد. برای دیدن محتوای مسیری دیگر، بدون تغییر مسیر کاری فعال، از `:ls <path>` استفاده کنید. خروجی همیشه ترتیب مشخصی دارد: ابتدا پوشه‌ها و سپس فایل‌ها نمایش داده می‌شوند و ستون‌های نوع و اندازه نیز در خروجی وجود دارند. این قابلیت با APIهای چندسکویی Go پیاده‌سازی شده است و دستور `ls` در Unix یا `dir` در Windows را اجرا نمی‌کند.

<div dir="ltr" align="left">

```text
:ls
:ls ..
:ls "project with spaces"
:del old.txt
:del "old directory"
```

</div>

## دستورهای کاربردی داخلی

این دستورها درون خود NextCmd پردازش می‌شوند؛ بنابراین به shell فرستاده نمی‌شوند، در تاریخچهٔ اجرای فرمان‌ها ذخیره نمی‌شوند و پیشنهادهای افزونه‌ها نیز با آن‌ها تداخل پیدا نمی‌کنند.

<table>
<thead>
<tr><th>دستور</th><th>کاربرد</th><th>رفتار مهم</th></tr>
</thead>
<tbody>
<tr><td dir="ltr"><code>:history [count]</code></td><td>نمایش دستورهایی که اخیراً از طریق NextCmd اجرا شده‌اند.</td><td>به‌طور پیش‌فرض ۲۰ مورد را نشان می‌دهد و عددی بین ۱ تا ۱۰۰۰ می‌پذیرد. زمان، کد خروج، مدت اجرا، افزونه، مسیر کاری و متن پاک‌سازی‌شدهٔ دستور نمایش داده می‌شوند. اگر تاریخچه غیرفعال باشد، برنامه این وضعیت را واضح اعلام می‌کند.</td></tr>
<tr><td dir="ltr"><code>:plugins</code></td><td>نمایش همهٔ افزونه‌های فعال.</td><td>افزونه‌ها براساس شناسه مرتب می‌شوند و نسخه، نام و توضیح آن‌ها نمایش داده می‌شود. افزونه‌های غیرفعال در این فهرست نیستند.</td></tr>
<tr><td dir="ltr"><code>:clear</code></td><td>پاک‌کردن صفحهٔ پایانه و انتقال نشانگر به ابتدای صفحه.</td><td>فقط محتوای قابل‌مشاهدهٔ صفحه را پاک می‌کند و تاریخچه یا وضعیت پروژه را تغییر نمی‌دهد.</td></tr>
<tr><td dir="ltr"><code>:del &lt;path&gt;</code></td><td>حذف فایل یا پوشه.</td><td>مسیرهای نسبی، کامل، نقل‌قول‌شده و <code>~</code> را نسبت به مسیر کاری فعال حل می‌کند. برنامه تشخیص می‌دهد هدف فایل است یا پوشه، پوشه‌ها را به‌صورت بازگشتی حذف می‌کند و اگر هم فایل و هم پوشهٔ مطابق پیدا شود از کاربر می‌پرسد کدام حذف شود.</td></tr>
<tr><td dir="ltr"><code>:config</code></td><td>نمایش تنظیماتی که برنامه اکنون با آن‌ها اجرا می‌شود.</td><td>مسیر فایل تنظیمات و تاریخچه، وضعیت تاریخچه، تعداد پیشنهادها، وضعیت debug و تنظیمات صریح افزونه‌ها را نشان می‌دهد. متغیرهای محیطی و اطلاعات محرمانه چاپ نمی‌شوند.</td></tr>
<tr><td dir="ltr"><code>:which &lt;command&gt;</code></td><td>یافتن فایل اجرایی یک دستور از طریق <code>PATH</code> سیستم‌عامل.</td><td>از جست‌وجوی چندسکویی Go استفاده می‌کند و قواعد <code>PATHEXT</code> در Windows را نیز رعایت می‌کند. در صورت امکان مسیر کامل فایل چاپ می‌شود.</td></tr>
<tr><td dir="ltr"><code>:version</code></td><td>نمایش اطلاعات نسخه و محیط ساخت.</td><td>نسخهٔ NextCmd، نسخهٔ Go، سیستم‌عامل، معماری و در صورت موجود بودن شناسهٔ کوتاه commit را نشان می‌دهد. نسخهٔ buildهای رسمی از tag انتشار گرفته می‌شود.</td></tr>
</tbody>
</table>

نمونه‌ها:

<div dir="ltr" align="left">

```text
:history
:history 5
:plugins
:clear
:del old.txt
:config
:which dotnet
:version
```

</div>

## فهرست دستورهای یک افزونه

برای دیدن همهٔ دستورهای ثابت یک افزونه، شناسهٔ آن را پس از `:?` بنویسید. نتیجه، متن دستور، توضیح کوتاه و میزان خطر آن را نمایش می‌دهد:

<div dir="ltr" align="left">

```text
:? git
:? dotnet
:? cargo
:? curl
:? go
:? docker
:? npm
:? pip
:؟ git
:؟ dotnet
:؟ cargo
:؟ curl
:؟ go
:؟ docker
:؟ npm
:؟ pip
```

</div>

خود افزونه این فهرست را از طریق رابط عمومی `sdk.HelpProvider` فراهم می‌کند. هسته و رابط پایانه هیچ فهرست ثابت و مخصوص Git یا .NET ندارند. مقادیر وابسته به پروژه، مانند نام شاخه، فایل، remote، solution یا project، در پیشنهادهای عادی و براساس وضعیت واقعی پروژه نمایش داده می‌شوند.

## پیشنهاد با نام ناقص ابزار

لازم نیست نام ابزار را کامل تایپ کنید. برای نمونه، `g` و `gi` پیشنهادهای Git را نشان می‌دهند؛ `dot`، `dotn` و `dotnet` نیز پیشنهادهای .NET را نمایش می‌دهند. ترتیب پیشنهادها همیشه با الگوریتمی ثابت و براساس تمام متن فعلی تعیین می‌شود.

## ظاهر پایانه

نشانه و فلش فیروزه‌ای، محل نوشتن دستور و پیشنهاد انتخاب‌شده را مشخص می‌کنند. کنار هر پیشنهاد یک برچسب کوتاه دیده می‌شود: `COMP` برای تکمیل، `NEXT` برای گام بعدی، `TIP` برای روش پیشنهادی و `FIX` برای رفع خطا. میزان خطر و شناسهٔ افزونهٔ پیشنهاددهنده نیز نمایش داده می‌شود. پس از اجرای دستور، یک سطر سبز یا قرمز کد خروج و مدت اجرا را نشان می‌دهد. برای غیرفعال‌کردن رنگ‌ها، متغیر استاندارد `NO_COLOR` را روی یک مقدار غیرخالی قرار دهید. هنگام انتقال خروجی به فایل یا برنامه‌ای دیگر نیز رنگ‌ها خودکار حذف می‌شوند.

ساختار اطلاعات نمایش‌داده‌شده پس از هر دستور به‌شکل زیر است:

<div dir="ltr" align="left">

```text
git add .  COMP MUTATING · git
│          │    │          └─ source plugin
│          │    └──────────── risk level
│          └───────────────── suggestion kind
└──────────────────────────── editable command
```

</div>

### نوع پیشنهاد

<table>
<thead>
<tr><th dir="ltr" align="left">برچسب</th><th dir="ltr" align="left">نوع در SDK</th><th dir="rtl" align="right">معنی</th><th dir="rtl" align="right">زمان نمایش</th></tr>
</thead>
<tbody>
<tr><td dir="ltr" align="left"><code>COMP</code></td><td dir="ltr" align="left"><code>Completion</code></td><td dir="rtl" align="right">متن در حال تایپ را کامل می‌کند یا شکل کامل‌تر دستور را می‌سازد.</td><td dir="rtl" align="right">هنگام ویرایش عادی دستور، تکمیل نام ابزار، option یا آرگومان پویا.</td></tr>
<tr><td dir="ltr" align="left"><code>NEXT</code></td><td dir="ltr" align="left"><code>NextAction</code></td><td dir="rtl" align="right">گام منطقی بعد از آخرین دستور موفق را پیشنهاد می‌دهد.</td><td dir="rtl" align="right">برای نمونه، اجرای تست پس از build یا commit پس از آماده‌کردن فایل‌ها.</td></tr>
<tr><td dir="ltr" align="left"><code>TIP</code></td><td dir="ltr" align="left"><code>BestPractice</code></td><td dir="rtl" align="right">یک بررسی اختیاری یا روش بهتر برای ادامهٔ کار پیشنهاد می‌دهد.</td><td dir="rtl" align="right">برای نمونه، مشاهدهٔ تغییرات آماده‌شده پیش از commit. این پیشنهاد اجباری نیست.</td></tr>
<tr><td dir="ltr" align="left"><code>FIX</code></td><td dir="ltr" align="left"><code>Recovery</code></td><td dir="rtl" align="right">پس از شکست دستور، یک راه احتمالی برای رفع مشکل ارائه می‌دهد.</td><td dir="rtl" align="right">فقط برای خطاهای شناخته‌شده، مانند نبود پروژه یا اشتباه‌بودن نام branch.</td></tr>
</tbody>
</table>

### سطح خطر

<table>
<thead>
<tr><th dir="ltr" align="left">برچسب</th><th dir="rtl" align="right">رنگ پیش‌فرض</th><th dir="rtl" align="right">معنی</th><th dir="rtl" align="right">نمونه‌های معمول</th></tr>
</thead>
<tbody>
<tr><td dir="ltr" align="left"><code>SAFE</code></td><td dir="rtl" align="right">سبز</td><td dir="rtl" align="right">معمولاً فقط اطلاعات را می‌خواند یا بررسی می‌کند و انتظار نمی‌رود وضعیت مهمی را تغییر دهد.</td><td dir="rtl" align="right">مشاهدهٔ status، log، header، diff یا نسخهٔ ابزار.</td></tr>
<tr><td dir="ltr" align="left"><code>MUTATING</code></td><td dir="rtl" align="right">زرد</td><td dir="rtl" align="right">فایل، dependency، وضعیت repository، خروجی build یا یک منبع remote را ایجاد یا تغییر می‌دهد.</td><td dir="rtl" align="right">build، format، add، commit، درخواست POST، upload یا نصب package.</td></tr>
<tr><td dir="ltr" align="left"><code>DESTRUCTIVE</code></td><td dir="rtl" align="right">قرمز</td><td dir="rtl" align="right">ممکن است اطلاعات را حذف، بازنویسی یا دور بریزد و بازگرداندن نتیجه دشوار باشد.</td><td dir="rtl" align="right">حذف branch، cleanup، بعضی حالت‌های reset یا درخواست HTTP DELETE.</td></tr>
<tr><td dir="ltr" align="left"><code>DANGEROUS</code></td><td dir="rtl" align="right">بنفش</td><td dir="rtl" align="right">خطر امنیتی یا احتمال ازدست‌رفتن اطلاعات در آن بسیار بیشتر است و باید با دقت ویژه بررسی شود.</td><td dir="rtl" align="right">غیرفعال‌کردن بررسی TLS یا انتشار تغییرات حساس و خارجی.</td></tr>
</tbody>
</table>

بخش پایانی مانند `· git`، `· cargo` یا `· curl` شناسهٔ افزونه‌ای است که پیشنهاد را ساخته است. نوع پیشنهاد و سطح خطر فقط اطلاعات کمکی هستند که افزونه ارائه می‌دهد. خود این برچسب‌ها روی ظاهر اثر دارند؛ مرتب‌سازی با metadata جداگانهٔ priority و relevance انجام می‌شود. این اطلاعات اجرای دستور را مسدود نمی‌کنند، تأیید جداگانه نمی‌گیرند و جای بررسی خود دستور توسط کاربر را نمی‌گیرند. با فعال‌کردن `NO_COLOR` رنگ‌ها حذف می‌شوند، ولی متن همهٔ برچسب‌ها همچنان نمایش داده می‌شود.

</div>
