// Package redactor provides a Redactor type that scans environment variable
// keys for sensitive patterns and replaces their values with a safe
// placeholder before the data leaves the process boundary.
//
// # Usage
//
//	r := redactor.New(nil) // uses DefaultSensitivePatterns
//	safe := r.Redact(resolvedEnv)
//
// # Custom patterns
//
// Supply your own slice of lowercase substrings to New:
//
//	r := redactor.New([]string{"secret", "token", "vault"})
//
// # Integration
//
// Redactor is intentionally decoupled from the masker package so it can be
// used independently in CLI output, audit logs, and snapshot diffs without
// pulling in masker's dependencies.
package redactor
