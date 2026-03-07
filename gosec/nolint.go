// Copyright 2026 John Wang. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gosec

import "fmt"

// Nolint formats a nolint:gosec comment with the given rule and reason.
//
// Example:
//
//	comment := gosec.Nolint("G117", "OAuth token response per RFC 6749")
//	// Returns: "//nolint:gosec // G117: OAuth token response per RFC 6749"
func Nolint(rule, reason string) string {
	return fmt.Sprintf("//nolint:gosec // %s: %s", rule, reason)
}

// NolintG117 returns a nolint comment for G117 (secret in JSON response).
//
// Use this when marshaling structs with intentional secret fields like
// OAuth access_token, client_secret, etc.
//
// Example reasons:
//   - "OAuth token response per RFC 6749"
//   - "OAuth registration response per RFC 7591"
//   - "API key response for authenticated user"
func NolintG117(reason string) string {
	return Nolint("G117", reason)
}

// NolintG118 returns a nolint comment for G118 (context.Background in goroutine).
//
// Use this when a goroutine intentionally uses context.Background because
// the request context is not appropriate (e.g., shutdown handlers).
//
// Example reasons:
//   - "Shutdown handler runs after request context is cancelled"
//   - "Background job outlives request lifecycle"
//   - "Cleanup routine needs independent timeout"
func NolintG118(reason string) string {
	return Nolint("G118", reason)
}

// NolintG704 returns a nolint comment for G704 (SSRF via taint analysis).
//
// Use this when making HTTP requests to URLs from trusted sources.
//
// Example reasons:
//   - "Test uses httptest server URL"
//   - "URL from validated allowlist"
//   - "Internal service URL from config"
//   - "URL constructed from trusted constants"
func NolintG704(reason string) string {
	return Nolint("G704", reason)
}

// NolintG101 returns a nolint comment for G101 (hardcoded credentials).
//
// Use this when a string matches credential patterns but is not actually
// a credential (e.g., URL paths, test fixtures, documentation).
//
// Example reasons:
//   - "URL path, not a credential"
//   - "Test fixture with fake credentials"
//   - "Documentation example"
func NolintG101(reason string) string {
	return Nolint("G101", reason)
}

// CommonReasons provides pre-written reason strings for common scenarios.
//
//nolint:gosec // G101: These are reason strings, not credentials
var CommonReasons = struct {
	// G117 reasons
	OAuthTokenResponse        string
	OAuthRegistrationResponse string

	// G118 reasons
	ShutdownHandler   string
	BackgroundJob     string
	CleanupRoutine    string
	IndependentCancel string

	// G704 reasons
	HttptestServer      string
	ValidatedAllowlist  string
	InternalServiceURL  string
	TrustedConstantsURL string

	// G101 reasons
	URLPathNotCredential string
	TestFixture          string
	DocumentationExample string
}{
	// G117
	OAuthTokenResponse:        "OAuth token response per RFC 6749",
	OAuthRegistrationResponse: "OAuth registration response per RFC 7591",

	// G118
	ShutdownHandler:   "Shutdown handler runs after request context is cancelled",
	BackgroundJob:     "Background job outlives request lifecycle",
	CleanupRoutine:    "Cleanup routine needs independent timeout",
	IndependentCancel: "Requires independent cancellation from request",

	// G704
	HttptestServer:      "Test uses httptest server URL",
	ValidatedAllowlist:  "URL from validated allowlist",
	InternalServiceURL:  "Internal service URL from config",
	TrustedConstantsURL: "URL constructed from trusted constants",

	// G101
	URLPathNotCredential: "URL path, not a credential",
	TestFixture:          "Test fixture with fake credentials",
	DocumentationExample: "Documentation example",
}
