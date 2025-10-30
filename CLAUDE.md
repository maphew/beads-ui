# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Beady is a web UI for the [beads](https://github.com/steveyegge/beads) issue tracker. The goal is to create a "fossil-like" all-in-one experience with a CLI and web interface bundled in a single executable that embeds all assets at build time.

Key technologies:
- **Backend**: Go 1.24+ with embedded filesystem (embed.FS)
- **Frontend**: HTMX for dynamic updates, Pico CSS for styling
- **Database**: SQLite via the beads library (github.com/steveyegge/beads)
- **Visualization**: Server-side Graphviz for dependency graphs

## Build and Development Commands

### Building

```bash
# Build with branch-aware naming (outputs to bin/beady or bin/beady-<branch>)
go run build.go

# Direct build (outputs to current directory)
go build -o beady ./cmd/beady

# Install globally
go install github.com/maphew/beads-ui/cmd/beady@latest
```

### Running

```bash
# Autodiscover database in current directory
beady 8080

# Specify database path
beady .beads/name.db 8080

# Development mode with live-reload (must run from repo root)
beady --dev

# Run from source without building
go run cmd/beady/main.go /path/to/.beads/name.db
```

### Testing and Development

```bash
# Create test database with sample issues
cd cmd
go run create_test_db_main.go /path/to/test.db

# Run with local beads development
# 1. Uncomment the replace directive in go.mod
# 2. Run: go mod tidy
```

## Architecture

### Single-File Distribution

All web assets (HTML templates, CSS, JS) are embedded into the binary at build time via `assets/beady/embed.go`:
- Templates: `assets/beady/templates/*.html`
- Static files: `assets/beady/static/*.{css,js}`

The `embed.FS` variable holds the embedded filesystem that's loaded at runtime.

### Development vs Production Mode

**Development mode** (`--dev` flag):
- Must run from repository root
- Uses `os.DirFS("assets/beady")` instead of embedded filesystem
- Live-reload via WebSocket watches `assets/beady/templates/` and `assets/beady/static/`
- Auto-opens browser
- Server auto-shuts down 5 seconds after last client disconnects

**Production mode**:
- Uses embedded filesystem from binary
- No file watching or WebSocket endpoint

### Template System

Templates are pre-parsed at startup into a single `template.Template` named "all" with shared FuncMap:
- `lower`, `upper`, `title`, `string` template functions available
- All templates accessed via `tmplAll.ExecuteTemplate(w, "filename.html", data)`
- Templates must be in `templates/` directory with `.html` extension

### HTTP Routes

- `/` - Main issue list with search/filter
- `/ready` - Unblocked issues (ready work view)
- `/blocked` - Blocked issues with blocker details
- `/issue/{id}` - Issue detail page with dependencies/events
- `/graph/{id}` - Graphviz dependency visualization
- `/api/issues` - JSON API (also handles HTMX partial requests)
- `/api/issue/{id}` - Single issue JSON
- `/api/stats` - Statistics JSON
- `/static/*` - Static assets (CSS/JS)
- `/ws` - WebSocket for live-reload (dev mode only)

### Data Enrichment Pattern

Issue lists are enriched with labels and dependency counts via `enrichIssuesWithLabels()`:
- Returns `[]*IssueWithLabels` wrapping `*beads.Issue`
- Adds `Labels []string`, `DepsCount`, `BlockersCount` fields
- Used consistently across all list views (index, ready, blocked)

### Theme System

Supports light/dark/auto modes with browser localStorage persistence:
- Theme selection in header dropdown
- Auto mode follows system preference via CSS media queries
- Preference stored client-side in `localStorage.getItem('theme')`

### Build Naming Convention

`build.go` creates branch-aware binaries:
- Main branch → `bin/beady`
- Feature branches → `bin/beady-<branch-name>`
- Sanitizes branch names (alphanumeric, dash, underscore only)
- Adds `.exe` extension on Windows

## Database Integration

Uses the beads library (`github.com/steveyegge/beads`) for all storage operations:
- Database autodiscovery via `beads.FindDatabasePath()`
- All handlers use `context.Context` from request
- Primary storage interface is `beads.Storage` (global `store` variable)
- Issue status values are lowercase in database: "open", "in_progress", "closed"

## File Organization

```
cmd/
  beady/main.go          # Main server entry point
  create_test_db_main.go # Test database generator
assets/beady/
  embed.go               # Embed directive for templates/static
  templates/*.html       # HTML templates
  static/*.{css,js}      # CSS and JavaScript
build.go                 # Custom build tool with branch naming
```

## Important Implementation Details

1. **File watching requires repo root**: Development mode checks for `assets/beady` directory existence and fails if not found.

2. **Template reloading**: Only `.html` files in `templates/` trigger re-parsing. Static files trigger browser reload but no server-side action.

3. **Issue filtering**: URL parameters `?search=`, `?status=`, `?priority=` are handled in `handleIndex` and `handleAPIIssues`.

4. **HTMX integration**: Check `r.Header.Get("HX-Request")` to return partial HTML (`issues_tbody.html`) instead of JSON.

5. **Graphviz dependency**: Server-side graph generation requires Graphviz installed (used by beads library, not directly by beady).

6. **Context propagation**: Always use `r.Context()` for database operations to support cancellation/timeout.

7. **Static file fallback**: `handleStatic` tries `static/` first, then falls back to `templates/` for backwards compatibility.

8. **Content-Type headers**: Explicitly set for CSS (`text/css`) and JS (`application/javascript`) files to avoid browser issues.
