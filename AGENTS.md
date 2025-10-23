# Repository Guidelines

## Architecture & Codebase Structure
This is a standalone Go web application providing a UI for the beads issue tracker. Key components:
- **Main application**: Single `main.go` file with embedded HTML templates (`templates/*.html`) and static assets (`static/*.css`, `static/*.js`)
- **Database**: SQLite via beads library (modernc.org/sqlite)
- **Web framework**: Standard library HTTP with gorilla/websocket for real-time updates
- **UI**: HTMX for dynamic updates, Graphviz for dependency visualization
- **Development**: File watcher with live reload in dev mode (`-d` flag)

## Build, Test, and Development Commands
- **Build**: `go build -o bd-ui .`
- **Run**: `./bd-ui [database-path] [port]` (autodiscovers `.beads/db.sqlite` if no path)
- **Development**: `./bd-ui -d` for live reload mode
- **Test database**: `cd cmd && go run create_test_db_main.go /path/to/test.db`
- **Single test**: `go test -run TestFunctionName` (when tests exist)
- **All tests**: `go test ./...` (no tests currently implemented)

## Coding Style & Naming Conventions
Follow Effective Go with these project specifics:
- Single main package with all HTTP handlers in `main.go`
- Use `context.Context` for all database operations
- Error handling: Return errors from handlers, let HTTP framework handle them
- Templates use Go html/template with embedded FS (`//go:embed`)
- WebSocket connections managed with sync.Mutex for thread safety
- Function organization: HTTP handlers first, then utilities/helpers
- Run `gofmt -s -w .` and `goimports -w .` before commits

## Reference information
`~/OneDrive/dev/llms-txt/` - Instructions for LLMs and AI code editors on how to use various tools and libraries.

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
