# lintfix

[![Go Reference][goref-svg]][goref-url]
[![Build Status][build-svg]][build-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![License][license-svg]][license-url]

A structured lint remediation database for Go projects using golangci-lint.

## Overview

lintfix provides a data layer that maps lint errors to remediation strategies, including:

- **Remediation types**: code fix, nolint annotation, or refactor
- **Helper package references**: links to packages like [mogo](https://github.com/grokify/mogo) for actual fixes
- **Nolint generators**: properly formatted comments with documented reasons
- **Pre-written reasons**: common scenarios like OAuth responses, shutdown handlers

## Installation

```bash
go get github.com/grokify/lintfix
```

## Usage

### Query the Remediation Database

```go
import "github.com/grokify/lintfix"

db := lintfix.MustLoadRemediations()

// Get remediation for a specific rule
fix := db.GetGosec("G120")
fmt.Println(fix.Name)                    // "Unbounded request body"
fmt.Println(fix.Remediation.Type)        // "code"
fmt.Println(fix.Remediation.Package)     // "github.com/grokify/mogo/net/http/httputilmore"
fmt.Println(fix.Remediation.Function)    // "LimitRequestBody"

// Check if there's a helper function available
if fix.HasHelper() {
    fmt.Printf("Use %s.%s()\n", fix.Remediation.Package, fix.Remediation.Function)
}
```

### Generate Nolint Comments

```go
import "github.com/grokify/lintfix/gosec"

// Using pre-written reasons
comment := gosec.NolintG117(gosec.CommonReasons.OAuthTokenResponse)
// Returns: "//nolint:gosec // G117: OAuth token response per RFC 6749"

// Using custom reasons
comment := gosec.NolintG704("URL from validated internal config")
// Returns: "//nolint:gosec // G704: URL from validated internal config"

// Generic nolint generator
comment := gosec.Nolint("G118", "Shutdown handler runs after context cancelled")
```

### Common Reasons

Pre-written reason strings for common scenarios:

```go
// G117 - Secret in JSON response
gosec.CommonReasons.OAuthTokenResponse        // "OAuth token response per RFC 6749"
gosec.CommonReasons.OAuthRegistrationResponse // "OAuth registration response per RFC 7591"

// G118 - context.Background in goroutine
gosec.CommonReasons.ShutdownHandler           // "Shutdown handler runs after request context is cancelled"
gosec.CommonReasons.BackgroundJob             // "Background job outlives request lifecycle"

// G704 - SSRF
gosec.CommonReasons.HttptestServer            // "Test uses httptest server URL"
gosec.CommonReasons.ValidatedAllowlist        // "URL from validated allowlist"

// G101 - Hardcoded credentials (false positives)
gosec.CommonReasons.URLPathNotCredential      // "URL path, not a credential"
gosec.CommonReasons.TestFixture               // "Test fixture with fake credentials"
```

## Supported Linters

| Linter | Rules |
|--------|-------|
| gosec | G101, G117, G118, G120, G401, G501, G601, G704 |
| staticcheck | SA1019, SA4006 |
| errcheck | unchecked |

## Remediation Types

| Type | Description | Example |
|------|-------------|---------|
| `code` | Add or modify code | `LimitRequestBody()` for G120 |
| `nolint` | Add nolint annotation | `//nolint:gosec // G117: reason` |
| `refactor` | Broader code changes | Move hardcoded secrets to env vars |

## Version-Specific Caveats

Some lint rules have version-specific behaviors. See [docs/gosec-caveats.md](docs/gosec-caveats.md) for details.

### G120 (gosec 2.11+)

gosec 2.11+ has stricter G120 detection:

1. **Only inline `http.MaxBytesReader` is recognized** - helper functions are not detected
2. **`r.FormValue()` is flagged** even after `ParseForm()` - use `r.Form.Get()` instead

```go
// Correct pattern for gosec 2.11+
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
if err := r.ParseForm(); err != nil { ... }
value := r.Form.Get("key")  // Not r.FormValue("key")
```

### Keeping Versions in Sync

Keep local golangci-lint version in sync with CI to avoid surprises:

```bash
# Check version
golangci-lint --version

# Install specific version
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.1
```

## Related Packages

- [mogo/net/http/httputilmore](https://github.com/grokify/mogo) - Runtime helpers like `LimitRequestBody()`

## License

MIT License. See [LICENSE](LICENSE) file.

[goref-svg]: https://pkg.go.dev/badge/github.com/grokify/lintfix.svg
[goref-url]: https://pkg.go.dev/github.com/grokify/lintfix
[build-svg]: https://github.com/grokify/lintfix/actions/workflows/ci.yaml/badge.svg
[build-url]: https://github.com/grokify/lintfix/actions/workflows/ci.yaml
[goreport-svg]: https://goreportcard.com/badge/github.com/grokify/lintfix
[goreport-url]: https://goreportcard.com/report/github.com/grokify/lintfix
[license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://github.com/grokify/lintfix/blob/main/LICENSE
