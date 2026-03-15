# Copilot Instructions — jsonreference

## Project Overview

Go implementation of [JSON Reference](https://datatracker.ietf.org/doc/html/draft-pbryan-zyp-json-ref-03)
(RFC 3986-based URI references with JSON Pointer fragments). JSON References are used
extensively in OpenAPI and JSON Schema specifications to express `$ref` links between
documents and within a single document.

Single module: `github.com/go-openapi/jsonreference`.

### Package layout (single package)

| File | Contents |
|------|----------|
| `reference.go` | `Ref` type, constructors (`New`, `MustCreateRef`), accessors, classification, and `Inherits` for parent-child resolution |
| `internal/normalize_url.go` | URL normalization (lowercase scheme/host, remove default ports, deduplicate slashes) — replaces the deprecated `purell` library |

### Key API

- `Ref` — the core type representing a parsed JSON Reference
- `New(string) (Ref, error)` — parse a JSON Reference string
- `MustCreateRef(string) Ref` — parse or panic
- `(*Ref).GetURL() *url.URL` — the underlying URL
- `(*Ref).GetPointer() *jsonpointer.Pointer` — the JSON Pointer fragment
- `(*Ref).Inherits(child Ref) (*Ref, error)` — resolve a child reference against this parent

### Dependencies

- `github.com/go-openapi/jsonpointer` — JSON Pointer (RFC 6901) implementation
- `github.com/go-openapi/testify/v2` — test-only assertions (zero-dep testify fork)

## Building & testing

```sh
go test ./...
```

## Conventions

Coding conventions are found beneath `.github/copilot`

### Summary

- All `.go` files must have SPDX license headers (Apache-2.0).
- Commits require DCO sign-off (`git commit -s`).
- Linting: `golangci-lint run` — config in `.golangci.yml` (posture: `default: all` with explicit disables).
- Every `//nolint` directive **must** have an inline comment explaining why.
- Tests: `go test ./...`. CI runs on `{ubuntu, macos, windows} x {stable, oldstable}` with `-race`.
- Test framework: `github.com/go-openapi/testify/v2` (not `stretchr/testify`; `testifylint` does not work).

See `.github/copilot/` (symlinked to `.claude/rules/`) for detailed rules on Go conventions, linting, testing, and contributions.
