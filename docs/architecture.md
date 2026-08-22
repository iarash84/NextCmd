# Architecture

The dependency direction is deliberately one-way:

```text
plugins ──> sdk <── internal core <── cmd/assistant
                                  ↑
                         plugins/builtin
```

The `sdk` contains stable data contracts and small optional capability interfaces. It imports only the Go standard library. A plugin must implement `Plugin.Info`; it may independently implement completion, detection, next-action, best-practice, or recovery capabilities. The completion engine uses type assertions, isolates provider errors, caches detection briefly, merges suggestions, and delegates final order to Core ranking.

Commands remain `Executable + Args`; rendering and execution are separate. Execution calls `exec.CommandContext` directly and never invokes a shell. Git-specific process calls and parsing live exclusively in `plugins/git`.

`internal/terminal` owns presentation and keyboard behavior. Its raw-mode boundary has separate Windows, Linux, and macOS files. No UI type appears in the SDK. `internal/history` writes portable JSON Lines and redacts common secret arguments and URL user-info before persistence.

Built-ins are composed explicitly by `plugins/builtin.All`; removing Git from that list leaves Core buildable. There is no `init` registration, reflection, mutable global registry, dynamic library, network access, or third-party dependency.
