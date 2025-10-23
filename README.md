# beads-webui

Standalone web UI for the beads issue tracker.

## Overview

This is a standalone web interface for [beads](https://github.com/steveyegge/beads), a
dependency-aware issue tracker. It provides a graphical interface for browsing and
visualizing issues, dependencies, and work status.

I've long been enamoured of Fossil-SCM and it's github-in-a-box nature, featuring a first class CLI and a strong web ui with commit timeline and issue tracker all wrapped up in a single executable (plus the db). It strikes me that Beads is excellently poised to do the same thing. This project is an experiment to see what that might entail. Feedback welcome.

The PR which started it: 
https://github.com/steveyegge/beads/pull/77

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

Install directly using `go install`:

```bash
go install github.com/maphew/beads-ui@latest
```

This will install the `beads-ui` binary to your `$GOPATH/bin` (usually `~/go/bin`).

### Building from source

1. Clone this repository:
   ```bash
   git clone https://github.com/maphew/beads-ui.git
   cd beads-ui
   ```

2. Build the web UI:
   ```bash
   go build -o beads-ui .
   ```

### Local development with beads

If you're developing both beads-ui and beads together:

1. Clone both repositories side by side
2. Uncomment the `replace` directive in `go.mod`
3. Run `go mod tidy`

## Usage

Run the web UI with an optional path to a beads database:

```bash
./beads-ui [path/to/.beads/db.sqlite] [port]
```

For example, to use autodiscovery:
```bash
./beads-ui 8080
```

Or specify a path:
```bash
./beads-ui .beads/db.sqlite 8080
```

The web UI will start on `http://127.0.0.1:8080` (or the specified port).

### Autodiscovery

If no database path is provided, the application will automatically search for a beads database in the current directory and standard locations (e.g., `.beads/db.sqlite`).

If no database is found, it will fall back to creating a new empty database.

## Development

To run the web UI in development mode:

```bash
go run main.go /path/to/.beads/db.sqlite
```

To create a test database with sample issues:

```bash
cd cmd
go run create_test_db_main.go /path/to/test.db
```

## Dependencies

The web UI depends on the beads library for database access and issue management. It uses:

- [htmx](https://htmx.org) for dynamic UI updates
- [Graphviz](https://graphviz.org) for dependency graph visualization (server-side)

## License

This project is licensed under the MIT License - see the LICENSE file for details.
