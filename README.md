# NextCmd

NextCmd is a fast, deterministic, cross-platform programming command-line assistant written in Go. It suggests editable commands from the current input, project state, and previous execution. The MVP ships with an explicitly registered Git plugin and uses no AI, network service, shell execution, or third-party Go dependency.

## Features

- Interactive editor: Up/Down selects, Tab accepts into the editor, Enter executes, Ctrl+C exits.
- Structured, shell-independent commands and captured stdout, stderr, exit code, and duration.
- Capability-based public plugin SDK and explicit compile-time registration.
- Git repository detection, cached context, branch/file/remote completion, next actions, best practices, and basic recovery.
- Deterministic prefix/fuzzy ranking and JSON-lines history with secret redaction.
- Standard-library-only implementation with platform-specific terminal boundaries.

## Build and run

Go 1.24 or newer is required.

```text
go build -o nextcmd .
./nextcmd
```

On Windows run `nextcmd.exe`. Debug logging is enabled with `nextcmd --debug`. Configuration defaults to the OS user configuration directory at `nextcmd/config.json`:

```json
{"historyEnabled":true,"maxSuggestions":8,"debug":false,"gitEnabled":true}
```

Cross-build examples:

```text
GOOS=linux GOARCH=arm64 go build ./...
GOOS=darwin GOARCH=amd64 go build ./...
```

PowerShell uses `$env:GOOS="linux"; $env:GOARCH="arm64"; go build ./...`.

## Architecture

`sdk` is the only package plugin authors need. Core discovers optional capabilities with type assertions, merges results, and owns ranking. Plugins return data and never render UI. Built-ins are listed explicitly in `plugins/builtin`. See [architecture](docs/architecture.md), [plugin development](docs/plugin-development.md), and the [Git plugin](docs/git-plugin.md).

## Development and testing

```text
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
```

## Roadmap

Future compile-time plugins may support Go, .NET, Cargo, Docker, Kubernetes, Terraform, npm, pnpm, and GitHub CLI. Dynamic plugins, AI, telemetry, marketplaces, and cloud sync are intentionally outside this MVP.

## License

MIT. See [LICENSE](LICENSE).
