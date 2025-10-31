# beady - a ui for beads

A web UI for the beads issue tracker.

## Overview

Beady is a web interface for [beads](https://github.com/steveyegge/beads), the issue
tracker built for LLM agents. Beady provides a graphical interface for browsing and
visualizing bead issues, dependencies, and work status.

I've long been enamoured of Fossil-SCM and it's github-in-a-box nature, featuring a first class CLI and a strong web ui with commit timeline and issue tracker all wrapped up in a single executable (plus the db). It strikes me that Beads is excellently poised to do the same thing. This project is an experiment to see what that might entail. Feedback welcome.

The PR which started it: 
https://github.com/steveyegge/beads/pull/77


_--> Also see [mantoni/beads-ui](https://github.com/mantoni/beads-ui from @mantoni. **`bdui`** has a higher development velocity than beady and a bigger feature set. You might like that one better. I'm going to keep poking away at beady anyway as I want to pursue the everything in one file idea._

## Features

### Read Operations
- **Issue list** with real-time filtering (search, status, priority)
- **Issue detail** pages with dependencies and activity
- **Dependency graphs** visualized with Graphviz
- **Ready work view** (unblocked issues)
- **Blocked issues view** with blocker details
- **Statistics dashboard** showing open/closed/in-progress counts
- **Theme customization** with light/dark/auto modes and persistent preferences
- **Graceful shutdown** via UI button (no need for task manager or kill commands)

### Write Operations (NEW!)
Beady now supports creating and modifying issues through the web UI:

- **Create new issues** with full form (title, type, priority, description, design, acceptance, labels)
- **Update status** via inline dropdown (open, in progress, closed)
- **Change priority** via inline dropdown (P0-P4)
- **Close issues** with optional reason
- **Add comments** with username attribution
- **Edit notes** with collapsible form
- **Manage labels** - add/remove labels inline
- **Manage dependencies** - add/remove blockers and dependencies

All write operations are performed by executing the `bd` CLI, ensuring guaranteed compatibility with the CLI and inheriting all validation logic. For bulk operations, use the `bd` CLI directly.

## Installation

### Prerequisites

- **bd CLI** in PATH (required for write operations) - install from [github.com/steveyegge/beads](https://github.com/steveyegge/beads)
- A beads database file (will be auto-discovered from `.beads/` directory)

### Download Pre-built Binaries (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/maphew/beady/releases):

- **Windows**: `beady_VERSION_Windows_x86_64.zip` (or `i386` for 32-bit)
- **macOS**: `beady_VERSION_Darwin_x86_64.tar.gz` (or `arm64` for Apple Silicon)
- **Linux**: `beady_VERSION_Linux_x86_64.tar.gz` (or `arm64`, `i386`)

Extract the binary and add it to your PATH.

### Install via Go

Install the latest release:

```bash
go install github.com/maphew/beady/cmd/beady@latest
```

Or install a specific version:

```bash
go install github.com/maphew/beady/cmd/beady@v1.0.0
```

This will install the `beady` binary to your `$GOPATH/bin` (usually `~/go/bin`).

### Building from Source

1. Clone this repository:
```bash
git clone https://github.com/maphew/beady.git
cd beady
```

2. Build the web UI:
```bash
go run build.go
# executable will be in ./bin/
```

### Local development with beads

If you're developing both beads-ui and beads together:

1. Clone both repositories side by side
2. Uncomment the `replace` directive in `go.mod`
3. Run `go mod tidy`

## Usage

Run the web UI with an optional path to a beads database:

```bash
beady [path/to/.beads/name.db] [port]
```

For example, to use autodiscovery:
```bash
beady 8080
```

Or specify a path:
```bash
beady .beads/name.db 8080
```

The web UI will start on `http://127.0.0.1:8080` (or the specified port).

### Theme Customization

Beady supports three theme modes for comfortable viewing in different environments:

- **Auto** (default): Automatically follows your system's dark/light preference
- **Light**: Forces light theme regardless of system setting
- **Dark**: Forces dark theme regardless of system setting

You can switch themes using the dropdown in the header. Your preference is automatically saved and will persist across browser sessions.

### Autodiscovery

If no database path is provided, the application will automatically search for a beads database in the current directory and standard locations (e.g., `.beads/name.db`).

If no database is found, it will fall back to creating a new empty database.

### Shutting Down

To stop the Beady server, you have two options:

1. **Via the UI**: Click the "Shutdown" button in the header of any page. You'll be prompted to confirm before the server shuts down gracefully.

2. **Via keyboard**: Press `Ctrl+C` in the terminal where Beady is running.

The UI shutdown button is particularly useful when running Beady in the background or when you don't have easy access to the terminal. It performs a graceful shutdown without needing to use task manager or system kill commands.

## Development

To run the web UI in development mode:

```bash
# from binary with live-reload (e.g. for template work)
beady --dev 

# or run from code without binary
go run cmd/beady/main.go /path/to/.beads/name.db
```

To create a test database with sample issues:

```bash
cd cmd
go run create_test_db_main.go /path/to/test.db
```

### API Endpoints

Beady provides the following HTTP endpoints:

#### Web Pages
- `GET /` - Main issue list with filtering (search, status, priority)
- `GET /ready` - Ready work view (unblocked issues)
- `GET /blocked` - Blocked issues view
- `GET /issue/{id}` - Issue detail page with dependencies and events
- `GET /graph/{id}` - Dependency graph visualization

#### API (JSON)

**Read Endpoints:**
- `GET /api/issues` - List all issues (supports `?search=`, `?status=`, `?priority=` filters)
- `GET /api/issue/{id}` - Get single issue details
- `GET /api/stats` - Get statistics (total, open, in-progress, closed counts)
- `POST /api/shutdown` - Gracefully shutdown the server

**Write Endpoints** (require bd CLI in PATH):
- `POST /api/issues/create` - Create new issue
- `POST /api/issue/status/{id}` - Update issue status
- `POST /api/issue/priority/{id}` - Update issue priority
- `POST /api/issue/close/{id}` - Close issue with reason
- `POST /api/issue/comments/{id}` - Add comment
- `POST /api/issue/notes/{id}` - Update notes
- `POST /api/issue/labels/{id}` - Add labels
- `DELETE /api/issue/labels/{id}/{label}` - Remove label
- `POST /api/issue/dependencies/{id}` - Add dependency
- `DELETE /api/issue/dependencies/{id}/{depSpec}` - Remove dependency

All write endpoints accept JSON request bodies with a `username` field for attribution. See [CLAUDE.md](CLAUDE.md) for detailed API documentation.

#### Static Assets
- `GET /static/*` - CSS, JavaScript, and other static files

#### Development Only
- `GET /ws` - WebSocket endpoint for live-reload (only in `--dev` mode)

### Releasing

This project uses [GoReleaser](https://goreleaser.com/) with GitHub Actions for automated releases. To create a new release:

1. Ensure all changes are committed and pushed to `main` branch
2. Create and push a version tag following semantic versioning:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. GitHub Actions will automatically:
   - Build binaries for all platforms (Linux, Windows, macOS)
   - Generate checksums
   - Create a GitHub release with all artifacts
   - Make the version available via `go install ...@latest`

The release will include:
- Multi-platform binaries (amd64, arm64, 386)
- Archives (`.tar.gz` for Unix, `.zip` for Windows)
- Checksums for verification
- Automated changelog from commits

**Testing a release**: To test the release process without creating an official release, push a tag with a `-test` suffix (e.g., `v0.0.1-test`).

## Dependencies

The web UI depends on the beads library for database access and issue management. It uses:

- [beads](https://github.com/steveyegge/beads) - an issue tracker for LLM agents
- [htmx](https://htmx.org) for dynamic UI updates
- [picocss](https://picocss.com) for styling and widgets
- [Graphviz](https://graphviz.org) for dependency graph visualization (server-side)

## License

This project is licensed under the MIT License - see the LICENSE file for details.
