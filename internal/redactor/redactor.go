// Package redactor provides utilities for stripping sensitive values
// from environment maps before they are printed, logged, or exported.
package redactor

import (
	"strings"
)

// DefaultSensitivePatterns holds substrings that, when found in a key
// (case-insensitive), cause the value to be redacted.
var DefaultSensitivePatterns = []string{
	"secret",
	"password",
	"passwd",
	"token",
	"api_key",
	"apikey",
	"private",
	"credential",
	"auth",
}

const Placeholder = "[REDACTED]"

// Redactor filters sensitive values from environment maps.
type Redactor struct {
	patterns []string
}

// New returns a Redactor using the provided patterns.
// Pass nil to use DefaultSensitivePatterns.
func New(patterns []string) *Redactor {
	if patterns == nil {
		patterns = DefaultSensitivePatterns
	}
	return &Redactor{patterns: patterns}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// Redact returns a copy of env where sensitive values are replaced with
// Placeholder. Non-sensitive values are passed through unchanged.
func (r *Redactor) Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.IsSensitive(k) {
			out[k] = Placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactValue returns Placeholder if the key is sensitive, otherwise value.
func (r *Redactor) RedactValue(key, value string) string {
	if r.IsSensitive(key) {
		return Placeholder
	}
	return value
}
