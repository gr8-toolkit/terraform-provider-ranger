# Go Conventions

## Language and Tooling

- Go version is declared in `go.mod` — always use `go-version-file: go.mod`
  when setting up Go in CI or local scripts, never hardcode a version.
- All code lives under the `ranger/` package. `main.go` is the thin entry point
  only; no business logic belongs there.
- Tool dependencies (e.g. `tfplugindocs`) are pinned in `tools/tools.go` using
  a `//go:build tools` constraint. Add new dev tools there — never install them
  ad hoc outside `go.mod`.

## Formatting and Linting

- All Go code must be formatted with `gofmt` / `goimports` before committing.
- Pre-commit hooks enforce trailing whitespace, end-of-file newline, and line
  endings — run `pre-commit run --all-files` locally to catch these early.
- There is no `golangci-lint` config yet; avoid introducing one without
  discussing it first.

## Error Handling

- Use `diag.FromErr(err)` to convert errors into Terraform diagnostics inside
  resource CRUD functions.
- Return `diag.Diagnostics` from all CRUD functions — never `error`.
- Check HTTP response status with `resp.IsError()` (resty) before inspecting
  the body. Log or return meaningful messages, not raw status codes alone.
- Do not swallow errors silently. If a 404 is an acceptable state (e.g. resource
  already deleted), document the intent explicitly with a comment.

## Testing

- Tests use `net/http/httptest` to mock the Ranger API — no real Ranger instance
  is required for unit tests.
- Use `github.com/stretchr/testify/assert` for assertions.
- Test file naming follows Go conventions: `<file>_test.go` in the same package.
- Run tests with `make test` (which runs `go test ./... -v`).
- Acceptance tests against a live Ranger instance are not wired up; if you add
  them, guard them with `TF_ACC=1` and document the required environment vars.

## Naming

- Exported functions that return a `*schema.Resource` are named
  `resource<ResourceType>()` (e.g. `resourcePolicy()`).
- The Terraform resource type name follows `ranger_<noun>` (e.g. `ranger_policy`).
- Unexported helpers use camelCase; keep them close to the function that uses them.

## Packages and Imports

- Keep all provider logic in the `ranger` package.
- Group imports: stdlib → external → internal (goimports enforces this).
- Do not use dot-imports (`import . "pkg"`) or blank-identifier imports outside
  `tools/tools.go`.

## Struct and Type Conventions

- `Client` wraps `*resty.Client`. Add new API methods as methods on `Client`,
  not as standalone functions.
- `parsePolicyResponse` and similar pure helpers should be package-level
  functions, not methods, when they don't need receiver state.
- Strip all server-managed fields (`id`, `guid`, `createdBy`, `updatedBy`,
  `createTime`, `updateTime`, `version`, `resourceSignature`) from stored state
  to prevent perpetual drift.
