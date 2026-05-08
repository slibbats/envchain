// Package audit provides structured event logging for envchain operations.
//
// It records key lifecycle events — resolution, missing lookups, overrides,
// and masking — to an io.Writer so operators can trace how environment
// variables were composed at runtime without exposing secret values.
//
// Usage:
//
//	logger := audit.New(os.Stderr)
//	logger.Record(audit.EventResolved, "DB_URL", "prod", "resolved from prod layer")
//
// For tests or silent operation use audit.Discard():
//
//	logger := audit.Discard()
package audit
