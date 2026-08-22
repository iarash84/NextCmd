# Contributing

Keep changes explicit, cross-platform, deterministic, and easy to test. The SDK must remain UI-independent and standard-library-only. Tool-specific knowledge belongs in its plugin. Avoid global registration and shell command strings.

Before submitting a change, run `gofmt -w .`, `go vet ./...`, `go test ./...`, and `go test -race ./...`. Add unit tests with injected fakes; integration tests may create local repositories with `t.TempDir()` but must not require GitHub or network access. Document new SDK contracts and user-visible controls.
