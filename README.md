# beady - a ui for beads

A web UI for the beads issue tracker.

## Overview

Beady is a web interface for [beads](https://github.com/steveyegge/beads), a
dependency-aware issue tracker. It provides a graphical interface for browsing and
visualizing issues, dependencies, and work status.

I've long been enamoured of Fossil-SCM and it's github-in-a-box nature, featuring a first class CLI and a strong web ui with commit timeline and issue tracker all wrapped up in a single executable (plus the db). It strikes me that Beads is excellently poised to do the same thing. This project is an experiment to see what that might entail. Feedback welcome.

The PR which started it: 
https://github.com/steveyegge/beads/pull/77


_--> Also see [mantoni/beads-ui](https://github.com/mantoni/beads-ui from @mantoni. **`bdui`** has a higher development velocity than beady and a bigger feature set. You might like that one better. I'm going to keep poking away at beady anyway as I want to pursue the everything in one file idea._

## Features

- **Issue list** with real-time filtering (search, status, priority)
- **Issue detail** pages with dependencies and activity
- **Dependency graphs** visualized with Graphviz
- **Ready work view** (unblocked issues)
- **Blocked issues view** with blocker details
- **Statistics dashboard** showing open/closed/in-progress counts

## Installation

### Prerequisites

- Go 1.21 or later
- A beads database file

### Quick install from Git

Install the latest release:

```bash
go install github.com/maphew/beads-ui/cmd/beady@latest
```

Or install the latest development version from main branch:

```bash
go install github.com/maphew/beads-ui/cmd/beady@main
```

This will install the `beady` binary to your `$GOPATH/bin` (usually `~/go/bin`).

### Building from source

1. Clone this repository:
```bash
git clone https://github.com/maphew/beady.git
cd beady
```

2. Build the web UI:
```bash
go run build.go
```

### Local development with beads

If you're developing both beads-ui and beads together:

1. Clone both repositories side by side
2. Uncomment the `replace` directive in `go.mod`
3. Run `go mod tidy`

## Usage

Run the web UI with an optional path to a beads database:

```bash
./beady [path/to/.beads/name.db] [port]
```

For example, to use autodiscovery:
```bash
./beady 8080
```

Or specify a path:
```bash
./beady .beads/name.db 8080
```

The web UI will start on `http://127.0.0.1:8080` (or the specified port).

### Autodiscovery

If no database path is provided, the application will automatically search for a beads database in the current directory and standard locations (e.g., `.beads/name.db`).

If no database is found, it will fall back to creating a new empty database.

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

### Releasing

This project follows semantic versioning. To release a new version:

1. Update the version in any relevant files (if needed)
2. Create and push a git tag with the version (e.g., `git tag v1.0.0 && git push origin v1.0.0`)
3. This will make the version available via `go install ...@latest`

## Dependencies

The web UI depends on the beads library for database access and issue management. It uses:

- [htmx](https://htmx.org) for dynamic UI updates
- [Graphviz](https://graphviz.org) for dependency graph visualization (server-side)

## License

This project is licensed under the MIT License - see the LICENSE file for details.
