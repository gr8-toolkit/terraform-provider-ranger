# Pre-commit and Prek

## Overview

This project uses [pre-commit](https://pre-commit.com/) to enforce code quality
locally and in CI. Hooks are configured in `.pre-commit-config.yaml`. CI runs
them via [prek](https://github.com/j178/prek) (`j178/prek-action@v1`).

## Installing pre-commit

```bash
pip install pre-commit
pre-commit install        # installs git hook into .git/hooks/pre-commit
pre-commit install --hook-type commit-msg  # optional: commit message hooks
```

After `pre-commit install`, hooks run automatically on every `git commit`.

## Running Hooks Manually

Run all hooks against all files (useful before raising a PR):

```bash
pre-commit run --all-files
```

Run a single hook by ID:

```bash
pre-commit run trailing-whitespace --all-files
pre-commit run prettier --all-files
pre-commit run markdownlint --all-files
```

Run only against staged files (default — what happens on `git commit`):

```bash
pre-commit run
```

## Configured Hooks

### `pre-commit/pre-commit-hooks` (v6.0.0)

| Hook ID | What it checks |
| --- | --- |
| `check-merge-conflict` | Blocks committed merge conflict markers |
| `check-added-large-files` | Blocks accidentally committed large files |
| `detect-private-key` | Blocks committed private keys |
| `check-case-conflict` | Catches case-insensitive filename conflicts |
| `mixed-line-ending` | Enforces consistent line endings |
| `trailing-whitespace` | Strips trailing whitespace (markdown linebreak-safe) |
| `end-of-file-fixer` | Ensures every file ends with a newline |

### `mirrors-prettier` (v4.0.0-alpha.8)

- Hook: `prettier` at version `3.1.1`
- Scope: `.yml` and `.yaml` files only
- Formats YAML with opinionated Prettier rules

### `igorshubovych/markdownlint-cli` (v0.45.0)

- Hook: `markdownlint --fix`
- Auto-fixes Markdown files per `.markdownlint.yaml` rules:
  - `MD013`: max line length 120 (tables excluded)
  - `MD033`: allows `<a>`, `<p>`, `<img>` HTML tags

## Prek in CI

CI runs pre-commit via the `prek` job in `.github/workflows/ci.yml`:

```yaml
- uses: j178/prek-action@v1
```

This action installs and runs `pre-commit run --all-files` in CI without
requiring Python or pre-commit to be installed manually. It caches hook
environments for fast re-runs.

The `prek` job runs on every PR and every push to `main`. A failing prek job
blocks merge. Fix locally with `pre-commit run --all-files` before pushing.

## Updating Hook Versions

To update all hooks to their latest tagged versions:

```bash
pre-commit autoupdate
```

Commit the resulting `.pre-commit-config.yaml` change and verify CI still
passes. Test the updated hooks with `pre-commit run --all-files` first.

## Adding a New Hook

1. Find the hook repo and ID on [pre-commit.com/hooks](https://pre-commit.com/hooks.html).
2. Add the repo block (or a new hook under an existing block) in
   `.pre-commit-config.yaml`.
3. Run `pre-commit run --all-files` to validate.
4. If the hook modifies files, commit those changes before pushing.

## Skipping Hooks (Emergency Only)

```bash
git commit --no-verify -m "fix: emergency hotfix"
```

Use only in genuine emergencies. The `prek` CI job will still catch violations.

## Markdownlint Rules

Configuration lives in `.markdownlint.yaml`. Current overrides:

```yaml
MD013:
  line_length: 120   # default is 80
  tables: false      # don't enforce line length inside tables
MD033:
  allowed_elements: [a, p, img]
```

All Markdown written in this repo (including steering files) must conform to
these rules. Run `pre-commit run markdownlint --all-files` to verify.
