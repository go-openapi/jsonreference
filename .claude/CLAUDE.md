# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go implementation of [JSON Reference](https://datatracker.ietf.org/doc/html/draft-pbryan-zyp-json-ref-03)
(RFC 3986-based URI references with JSON Pointer fragments). JSON References are used
extensively in OpenAPI and JSON Schema specifications to express `$ref` links between
documents and within a single document.

The `Ref` type parses a reference string into its URL and JSON Pointer components, classifies
it (full URL, path-only, fragment-only, file scheme), and supports inheritance (resolving a
child reference against a parent).

See [docs/MAINTAINERS.md](../docs/MAINTAINERS.md) for CI/CD, release process, and repo structure details.

### Package layout (single package)

| File | Contents |
|------|----------|
| `reference.go` | `Ref` type, constructors (`New`, `MustCreateRef`), accessors (`GetURL`, `GetPointer`, `String`), classification (`IsRoot`, `IsCanonical`), and `Inherits` for parent-child resolution |
| `internal/normalize_url.go` | URL normalization (lowercase scheme/host, remove default ports, deduplicate slashes) — replaces the deprecated `purell` library |

### Key API

- `Ref` — the core type representing a parsed JSON Reference
- `New(string) (Ref, error)` — parse a JSON Reference string
- `MustCreateRef(string) Ref` — parse or panic
- `(*Ref).GetURL() *url.URL` — the underlying URL
- `(*Ref).GetPointer() *jsonpointer.Pointer` — the JSON Pointer fragment
- `(*Ref).Inherits(child Ref) (*Ref, error)` — resolve a child reference against this parent
- `(*Ref).IsRoot() bool` — true if this is a root document reference
- `(*Ref).IsCanonical() bool` — true if the reference starts with `http(s)://` or `file://`

### Dependencies

- `github.com/go-openapi/jsonpointer` — JSON Pointer (RFC 6901) implementation
- `github.com/go-openapi/testify/v2` — test-only assertions (zero-dep testify fork)

### Notable design decisions

- **Replaced `purell` with internal normalization** — the `internal/normalize_url.go` file
  replaces the unmaintained `purell` library. It performs only the safe normalizations that
  were previously used: lowercase scheme/host, remove default HTTP(S) ports, and deduplicate
  path slashes.
- **Classification flags on `Ref`** — rather than re-parsing on each query, the `parse` method
  sets boolean flags (`HasFullURL`, `HasURLPathOnly`, `HasFragmentOnly`, `HasFileScheme`,
  `HasFullFilePath`) once at construction time.
- **JSON Pointer errors are silently ignored** — if the URL fragment is not a valid JSON Pointer,
  the pointer is left as zero-value. This allows the type to represent any URI reference, not
  just those with valid JSON Pointer fragments.
