# Terraform Plugin

English | [فارسی](#افزونهٔ-terraform)

The built-in Terraform plugin detects `.tf` files in the active directory or its parents. It extracts resource and module addresses without evaluating configuration. Whether `.terraform` exists and its locally selected workspace are recorded and cached briefly. Detection does not initialize Terraform or access a backend.

The catalog covers init, formatting, validation, planning, showing plans/state, apply, destroy, outputs, state inspection, and workspaces. Dynamic completion suggests detected resource addresses and workspaces. Successful init suggests validation and planning; formatting suggests validation. Recovery recognizes missing initialization and invalid configuration.

`plan`, validation, and state inspection are read-only suggestions. `apply` and `destroy` carry dangerous risk labels because they can materially change remote infrastructure. Always review a saved plan before applying it.

The plugin ID is `terraform`. Disable it with `{"plugins":{"terraform":false}}` and view its catalog with `:? terraform`.

---

<div dir="rtl" align="right">

# افزونهٔ Terraform

افزونهٔ داخلی Terraform فایل‌های `.tf` را در مسیر فعال یا والدهای آن شناسایی می‌کند. آدرس resource و module بدون ارزیابی configuration استخراج می‌شود. وجود پوشهٔ `.terraform` و workspace انتخاب‌شدهٔ محلی ثبت و برای مدت کوتاهی cache می‌شوند. تشخیص پروژه Terraform را initialize نمی‌کند و به backend دسترسی ندارد.

فهرست فرمان‌ها init، قالب‌بندی، validate، plan، نمایش plan و state، apply، destroy، output، بررسی state و workspaceها را پوشش می‌دهد. تکمیل پویا resourceها و workspaceهای واقعی را پیشنهاد می‌کند. پس از init، اعتبارسنجی و plan پیشنهاد می‌شوند و recovery نیز initialization ناقص و configuration نامعتبر را تشخیص می‌دهد.

فرمان‌های plan، validate و مشاهدهٔ state فقط‌خواندنی هستند. apply و destroy به‌دلیل امکان تغییر جدی زیرساخت راه‌دور برچسب خطر dangerous دارند. همیشه پیش از apply، plan ذخیره‌شده را بررسی کنید.

شناسهٔ افزونه `terraform` است. با `{"plugins":{"terraform":false}}` غیرفعال و با `:? terraform` راهنمای آن نمایش داده می‌شود.

</div>
