# Repository Guidelines

## Project Structure & Module Organization
Keep executables under `cmd/` with a `main.go` entry point, internal packages in `internal/` (private to this module), and reusable libraries in `pkg/` (public APIs). Supporting data, prompts, or fixture assets belong in `assets/` with subfolders that mirror the consuming package. Tests live alongside code as `*_test.go` files and use the `testdata/` directory for fixtures. Temporary experiments go in `examples/` or a `_scratch/` directory and may be pruned at review time.

## Build, Test, and Development Commands
Initialize the module with `go mod init` if starting fresh, then run `go mod download` to fetch dependencies declared in `go.mod`. Use `go mod tidy` to prune unused dependencies and update `go.sum`. Run `golangci-lint run` or `make lint` for static analysis, `go test ./...` for the full test suite, and `go run main.go` or `go run ./cmd/...` to execute the default entry point. When iterating quickly, use `go test ./internal/<package>` to scope to a specific package, or `go build -o bin/app ./cmd/...` to compile a binary.

## Coding Style & Naming Conventions
Follow Effective Go guidelines with tab indentation, Go's static typing, and doc comments that describe package purpose and exported functions. Use lowercase package names (e.g., `research`, `tools`) and PascalCase for exported types and functions (`ResearchAgent`, `ParseQuery`). Shared utilities should avoid init-time side effects. Keep public functions under 40 lines, favoring helper functions over deeply nested logic. Run `gofmt -s -w .`, `goimports -w .`, and `golangci-lint run` before opening a pull request.

## Testing Guidelines
Write tests using Go's testing package that mirror the behavioral seams of each component. Use `<feature>_test.go` filenames and descriptive test functions (`TestResearchAgentHandlesRateLimits`). Provide test helpers and table-driven tests for common scenarios; place fixture data in `testdata/` directories. Maintain >=90% statement coverage (`go test -cover`) and add regression tests for every bug fix. Include integration tests with build tags (`//go:build integration`) or separate `*_integration_test.go` files whenever a component relies on multiple services.

## Commit & Pull Request Guidelines
Follow Conventional Commits (`feat:`, `fix:`, `chore:`) and keep commit bodies focused on the why and the rollout impact. Squash noisy WIP commits. Pull requests need a summary, testing notes, linked issue IDs, and screenshots or logs if behavior changes. Tag reviewers responsible for the touched components, and ensure CI (lint, tests, type checks) is green before requesting review.

## Security & Configuration Tips
Never commit secrets; load them from `.env` or environment variables using packages like `godotenv` and document required keys in `docs/configuration.md`. Validate outbound API calls and sanitize user input to avoid injection attacks. Rotate API keys quarterly and audit service capabilities whenever new integrations are introduced.

## Issue Tracking with bd (beads)

**IMPORTANT**: This project uses **bd (beads)** for ALL issue tracking. Do NOT use markdown TODOs, task lists, or other tracking methods.

### Why bd?

- Dependency-aware: Track blockers and relationships between issues
- Git-friendly: Auto-syncs to JSONL for version control
- Agent-optimized: JSON output, ready work detection, discovered-from links
- Prevents duplicate tracking systems and confusion

### Quick Start

**Check for ready work:**
```bash
bd ready --json
```

**Create new issues:**
```bash
bd create "Issue title" -t bug|feature|task -p 0-4 --json
bd create "Issue title" -p 1 --deps discovered-from:bd-123 --json
```

**Claim and update:**
```bash
bd update bd-42 --status in_progress --json
bd update bd-42 --priority 1 --json
```

**Complete work:**
```bash
bd close bd-42 --reason "Completed" --json
```

### Issue Types

- `bug` - Something broken
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature with subtasks
- `chore` - Maintenance (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (default, nice-to-have)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Workflow for AI Agents

1. **Check ready work**: `bd ready` shows unblocked issues
2. **Claim your task**: `bd update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - `bd create "Found bug" -p 1 --deps discovered-from:<parent-id>`
5. **Complete**: `bd close <id> --reason "Done"`

### Auto-Sync

bd automatically syncs with git:
- Exports to `.beads/issues.jsonl` after changes (5s debounce)
- Imports from JSONL when newer (e.g., after `git pull`)
- No manual export/import needed!

### MCP Server (Recommended)

If using Claude or MCP-compatible clients, install the beads MCP server:

```bash
pip install beads-mcp
```

Add to MCP config (e.g., `~/.config/claude/config.json`):
```json
{
  "beads": {
    "command": "beads-mcp",
    "args": []
  }
}
```

Then use `mcp__beads__*` functions instead of CLI commands.

### Important Rules

- ✅ Use bd for ALL task tracking
- ✅ Always use `--json` flag for programmatic use
- ✅ Link discovered work with `discovered-from` dependencies
- ✅ Check `bd ready` before asking "what should I work on?"
- ❌ Do NOT create markdown TODO lists
- ❌ Do NOT use external issue trackers
- ❌ Do NOT duplicate tracking systems

For more details, see README.md and QUICKSTART.md.
