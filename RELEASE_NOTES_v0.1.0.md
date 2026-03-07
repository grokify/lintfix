# Release Notes v0.1.0

Initial release of lintfix - a structured lint remediation database for Go projects.

## Highlights

- Structured JSON database mapping lint rules to remediation strategies
- Nolint comment generators with documented reasons
- Pre-written reason strings for common scenarios
- References to helper packages (mogo) for code-based fixes

## Features

### Remediation Database

Query remediation guidance for lint rules:

```go
db := lintfix.MustLoadRemediations()
fix := db.GetGosec("G120")
// fix.Remediation.Package = "github.com/grokify/mogo/net/http/httputilmore"
// fix.Remediation.Function = "LimitRequestBody"
```

Three remediation types:

- `code` - Fix by adding/changing code (links to helper packages)
- `nolint` - Fix by adding nolint annotation with documentation
- `refactor` - Fix requires broader code changes

### Nolint Generators

Generate properly formatted nolint comments:

```go
comment := gosec.NolintG117(gosec.CommonReasons.OAuthTokenResponse)
// "//nolint:gosec // G117: OAuth token response per RFC 6749"
```

### Supported Rules

| Linter | Rules |
|--------|-------|
| gosec | G101, G117, G118, G120, G401, G501, G601, G704 |
| staticcheck | SA1019, SA4006 |
| errcheck | unchecked |

## Installation

```bash
go get github.com/grokify/lintfix
```
