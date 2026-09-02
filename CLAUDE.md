# CLAUDE.md

Guidance for AI coding agents working in this repository.

## Project

`fl` is an experimental Go CLI for interacting with [Flintlock](https://github.com/liquidmetal-dev/flintlock). See
`README.md` for build/usage details.

## Layout

- `main.go` — entrypoint
- `internal/cmd/` — Cobra CLI commands (`root.go`, `microvm/`, `version/`)
- `pkg/app/` — business logic wrapping the Flintlock client
- `pkg/logging/` — zap-based logging setup

## Build

```sh
make build
```

Tool versions are pinned in `mise.toml`; run `mise install` before building.

## Commit and PR conventions

- **Conventional Commits are required.** Every commit message must use a `type: subject` prefix such as `feat:`,
  `fix:`, `chore:`, `refactor:`, `docs:`, or `test:`, matching the existing git history. `.goreleaser.yaml` relies on
  these prefixes to filter the release changelog.
- **Do not add `Co-Authored-By` or any other AI/agent attribution footer** to commit messages or pull request
  descriptions. Commits and PRs must not mention the tool or agent that produced them.
