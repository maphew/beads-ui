# beads-webui

Standalone web UI for the beads issue tracker.

## Overview

This is a standalone web interface for [beads](https://github.com/steveyegge/beads), a
dependency-aware issue tracker. It provides a graphical interface for browsing and
visualizing issues, dependencies, and work status.

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

### Building from source

1. Clone this repository:
   ```bash
   git clone https://github.com/maphew/beads-webui.git
   cd beads-webui
   ```

2. Build the web UI:
   ```bash
   go build -o bd-ui .
   ```

## Usage

Run the web UI with a path to a beads database:

```bash
./beads-webui /path/to/.beads/db.sqlite [port]
```

For example:
```bash
./beads-webui .beads/db.sqlite 8080
```

The web UI will start on `http://127.0.0.1:8080` (or the specified port).

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
