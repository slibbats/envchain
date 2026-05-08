// Package masker provides utilities for redacting sensitive environment
// variable values so they are never accidentally printed or logged.
package masker

import "strings"

// DefaultSensitiveKeys contains common key substrings that are considered
// sensitive and whose values should be masked.
var DefaultSensitiveKeys = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"CREDENTIAL",
	"AUTH",
}

const masked = "***"

// Masker redacts the values of sensitive environment variables.
type Masker struct {
	sensitiveKeys []string
}

// New returns a Masker that treats any key containing one of the provided
// substrings (case-insensitive) as sensitive. If sensitiveKeys is empty,
// DefaultSensitiveKeys is used.
func New(sensitiveKeys []string) *Masker {
	if len(sensitiveKeys) == 0 {
		sensitiveKeys = DefaultSensitiveKeys
	}
	upper := make([]string, len(sensitiveKeys))
	for i, k := range sensitiveKeys {
		upper[i] = strings.ToUpper(k)
	}
	return &Masker{sensitiveKeys: upper}
}

// IsSensitive reports whether the given key should be treated as sensitive.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, s := range m.sensitiveKeys {
		if strings.Contains(upper, s) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked placeholder if the key is sensitive, otherwise
// it returns the original value unchanged.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return masked
	}
	return value
}

// MaskEnv returns a copy of the provided environment map with sensitive values
// replaced by the masked placeholder.
func (m *Masker) MaskEnv(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = m.MaskValue(k, v)
	}
	return out
}
