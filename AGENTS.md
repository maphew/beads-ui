# Repository Guidelines

## Project Structure & Module Organization
Keep user-facing code under `src/`, grouped by agent or service (`src/agents/<agent_name>/` for behaviors, `src/tools/` for shared integrations). Supporting data, prompts, or fixture assets belong in `assets/` with subfolders that mirror the agent consuming them. Tests live in `tests/` and should reflect the module tree (`tests/agents/test_<agent_name>.py`). Temporary experiments go in `sandbox/` and may be pruned at review time.

## Build, Test, and Development Commands
Create a virtual environment (`python -m venv .venv && source .venv/bin/activate`) before installing dependencies. Install runtime requirements with `python -m pip install -r requirements.txt`; add tooling extras to `requirements-dev.txt`. Run `make lint` for static analysis, `make test` for the full pytest suite, and `make run` to execute the default agent entrypoint in `src/main.py`. When iterating quickly, use `pytest tests/agents -k <agent_name>` to scope to an individual agent.

## Coding Style & Naming Conventions
Adopt PEP 8 with four-space indentation, type hints, and descriptive docstrings that state inputs, outputs, and failure modes. Name agents with snake_case modules and PascalCase classes (`ResearchAgent`). Shared utilities should avoid side effects at import time. Keep public functions under 40 lines, favoring helper functions over deeply nested logic. Run `ruff check` and `black src tests` before opening a pull request.

## Testing Guidelines
Write pytest-based unit tests that mirror the behavioral seams of each agent. Use `test_<feature>.py` filenames and descriptive test names (`test_research_agent_handles_rate_limits`). Provide fixtures for external services in `tests/fixtures/`. Maintain >=90% statement coverage and add regression tests for every bug fix. Include integration smoke tests in `tests/integration/` whenever an agent relies on multiple tools.

## Commit & Pull Request Guidelines
Follow Conventional Commits (`feat:`, `fix:`, `chore:`) and keep commit bodies focused on the why and the rollout impact. Squash noisy WIP commits. Pull requests need a summary, testing notes, linked issue IDs, and screenshots or logs if behavior changes. Tag reviewers responsible for the touched agents, and ensure CI (lint, tests, type checks) is green before requesting review.

## Security & Configuration Tips
Never commit secrets; load them from `.env.local` and document required keys in `docs/configuration.md`. Validate outbound tool calls and sanitize model prompts to avoid prompt-injection carryover. Rotate API keys quarterly and audit agent capabilities whenever new tools are introduced.

## Issue Tracking with Beads
We now use Beads (`bd`) for all issue tracking and work planning. The `bd` binary is already available on the system `PATH`, so invoke it directly from the repository root. Stick to a single-workstream setup for now—upstream v0.9.x has known data duplication/loss bugs when multiple repositories share a Beads database. If you detect another `.beads/` directory in the tree or need multi-repo coordination, pause and ask for a human decision before proceeding.

BEFORE ANYTHING ELSE: run 'bd onboard' and follow the instructions"
