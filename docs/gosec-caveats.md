# Gosec Version-Specific Caveats

This document describes version-specific behaviors in gosec that affect remediation strategies.

## G120: Unbounded Request Body

### gosec 2.11+ Behavior Changes

Starting with gosec 2.11 (included in golangci-lint 2.11+), the G120 rule has stricter detection:

#### 1. Helper Functions Not Recognized

gosec only recognizes **inline** `http.MaxBytesReader` calls. Helper functions that wrap this call are not detected.

**Does NOT work with gosec 2.11+:**

```go
// Helper function - gosec doesn't trace the call
httputilmore.LimitRequestBody(w, r, httputilmore.DefaultMaxBodySize)
if err := r.ParseForm(); err != nil { ... }  // G120 flagged
```

**Works with gosec 2.11+:**

```go
// Inline call - gosec recognizes this pattern
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
if err := r.ParseForm(); err != nil { ... }  // OK
```

#### 2. FormValue() Flagged After ParseForm()

gosec 2.11+ flags `r.FormValue()` calls even when `ParseForm()` was already called with a limited body. This is because `FormValue()` internally calls `ParseForm()` if not already parsed, and gosec doesn't track that the form is already parsed.

**Does NOT work with gosec 2.11+:**

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
if err := r.ParseForm(); err != nil { ... }
value := r.FormValue("key")  // G120 flagged!
```

**Works with gosec 2.11+:**

```go
r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
if err := r.ParseForm(); err != nil { ... }
value := r.Form.Get("key")  // OK - directly accesses parsed form
```

### Complete Example

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 1. Limit body size INLINE (gosec recognizes this)
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

    // 2. Parse the form
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    // 3. Use r.Form.Get() instead of r.FormValue()
    username := r.Form.Get("username")
    password := r.Form.Get("password")

    // ... handle request
}
```

### Why This Matters

The stricter detection in gosec 2.11+ means:

1. **Helper functions** like `httputilmore.LimitRequestBody()` provide good abstractions but won't satisfy the linter
2. **FormValue()** is a common pattern but triggers false positives after proper limiting
3. **Local/CI version mismatch** can cause CI failures that don't reproduce locally

### Recommendations

1. **Keep golangci-lint versions in sync** between local development and CI
2. **Use inline `http.MaxBytesReader`** rather than helper functions for G120
3. **Use `r.Form.Get()`** instead of `r.FormValue()` after calling `ParseForm()`
4. **Document the pattern** in code comments referencing G120

### Version Reference

| golangci-lint | gosec | G120 Behavior |
|---------------|-------|---------------|
| 2.10.x | 2.21.x | Recognizes inline MaxBytesReader only |
| 2.11.x | 2.24.x | Same + flags FormValue() after ParseForm |

## Keeping Versions in Sync

To avoid local/CI lint mismatches:

```bash
# Check local version
golangci-lint --version

# Install specific version to match CI
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.1
```

Or pin the version in CI to match local:

```yaml
# .github/workflows/lint.yaml
- uses: golangci/golangci-lint-action@v8
  with:
    version: v2.11.1
```
