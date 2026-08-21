# CI/CD Workflows

## CI Workflow (`.github/workflows/ci.yml`)

Triggers on every pull request and every push to `main`.
Permissions: `contents: read` (read-only, minimal blast radius).

### Jobs

#### `prek` (ubuntu-latest)

Runs all pre-commit hooks across the entire repository using prek:

```yaml
- uses: j178/prek-action@v1
```

Catches formatting, whitespace, large files, private keys, and markdownlint
violations before code even reaches review. See `pre-commit.md` for the full
hook inventory.

#### `test` (ubuntu-latest)

1. Checks out the repository.
2. Sets up Go from `go.mod` (`go-version-file: go.mod`) — version is never
   hardcoded.
3. Runs `go build -o ./dist/` to verify the code compiles cleanly.
4. Runs tests via `robherley/go-test-action@v0` which provides structured
   GitHub annotations for failing tests.

No acceptance tests run in CI — only unit tests that mock the Ranger API via
`httptest`.

### Merging Rules

Both `prek` and `test` jobs must be green before merging. PRs must follow
semantic title conventions enforced by `.github/semantic.yml`
(e.g. `feat:`, `fix:`, `chore:`).

## Release Workflow (`.github/workflows/release.yml`)

Triggers on tags matching `v*` (e.g. `v1.2.3`).
Permissions: `contents: write` (needed to publish the GitHub release).

### Steps

1. Checkout with `fetch-depth: 0` — full git history is required by GoReleaser
   for changelog generation.
2. Set up Go from `go.mod`.
3. Import GPG signing key via `crazy-max/ghaction-import-gpg@v6.3.0` using
   secrets `GPG_PRIVATE_KEY` and `PASSPHRASE`. Outputs `fingerprint`.
4. Run GoReleaser via `goreleaser/goreleaser-action@v6.4.0`:
   - Command: `goreleaser release --clean`
   - Env: `GITHUB_TOKEN` (auto-provided) and `GPG_FINGERPRINT`

### Required Secrets

| Secret | Description |
| --- | --- |
| `GPG_PRIVATE_KEY` | GPG private key for signing release artifacts |
| `PASSPHRASE` | Passphrase for the GPG private key |
| `GITHUB_TOKEN` | Auto-provided by GitHub Actions |

### Release Artifacts

GoReleaser (`goreleaser.yml`) produces:

- Binaries for `linux`, `darwin`, `windows`, `freebsd` × `amd64`, `386`,
  `arm`, `arm64` (excluding `darwin/386`)
- CGO disabled (`CGO_ENABLED=0`)
- ZIP archives per platform
- `SHA256SUMS` checksum file
- GPG-signed checksum file
- `terraform-registry-manifest.json` bundled in every archive

Changelog generation is disabled in `goreleaser.yml` — the GitHub release
notes are generated from the tag/PR history instead.

## Creating a Release

1. Ensure `main` is green (both CI jobs passing).
2. Create and push a semver tag:

   ```bash
   git tag v1.2.3
   git push origin v1.2.3
   ```

3. The release workflow triggers automatically.
4. Verify the GitHub release page has all platform archives and a signed
   `SHA256SUMS`.

## Dependabot

`.github/dependabot.yml` configures weekly Go module updates (`gomod`) from
the repository root. Dependabot PRs go through the same CI checks. Merge them
promptly to stay current with security patches.
