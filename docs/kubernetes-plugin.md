# Kubernetes Plugin

English | [فارسی](#افزونهٔ-kubernetes)

The built-in Kubernetes plugin provides deterministic `kubectl` suggestions. It detects YAML manifests containing both `apiVersion` and `kind` in the active directory or its parents. Manifest names, kinds, and namespaces are read locally; configured contexts are queried through the injectable Runner and cached briefly. Detection never contacts a cluster.

The catalog covers cluster inspection, contexts, pods, deployments, logs, manifest diff/apply/delete, and rollout status/restart. Dynamic completion suggests detected manifests, contexts, and namespaces. Applying a manifest suggests inspecting pods and rollout status. Best practice recommends `kubectl diff` before apply, while recovery for connection failures suggests inspecting configured contexts.

Risk metadata distinguishes read-only inspection, context changes, manifest mutation, and destructive deletion. Review the selected context and namespace before mutating a cluster.

The plugin ID is `kubernetes`. Disable it with `{"plugins":{"kubernetes":false}}` and view its catalog with `:? kubernetes`.

---

<div dir="rtl" align="right">

# افزونهٔ Kubernetes

افزونهٔ داخلی Kubernetes پیشنهادهای ثابت و قابل‌پیش‌بینی `kubectl` را ارائه می‌دهد. فایل‌های YAML دارای `apiVersion` و `kind` در مسیر فعال یا والدهای آن شناسایی می‌شوند. نام، نوع و namespace مربوط به manifestها محلی خوانده می‌شوند و contextها از طریق Runner قابل‌تزریق دریافت و برای مدت کوتاهی cache می‌شوند. تشخیص پروژه هیچ اتصالی به cluster برقرار نمی‌کند.

فهرست فرمان‌ها شامل بررسی cluster، مدیریت context، مشاهدهٔ pod و deployment، log، عملیات diff/apply/delete روی manifest و rollout است. تکمیل پویا manifest، context و namespace واقعی را پیشنهاد می‌دهد. پس از apply نیز بررسی podها و rollout پیشنهاد می‌شود و در خطای اتصال، مشاهدهٔ contextها ارائه می‌گردد.

سطح خطر میان مشاهدهٔ فقط‌خواندنی، تغییر context، تغییر manifest و حذف مخرب تفاوت می‌گذارد. پیش از تغییر cluster، context و namespace را بررسی کنید.

شناسهٔ افزونه `kubernetes` است. با `{"plugins":{"kubernetes":false}}` غیرفعال و با `:? kubernetes` راهنمای آن نمایش داده می‌شود.

</div>
