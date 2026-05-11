// Package filter provides utilities for selecting, excluding, and
// matching environment variable keys by prefix, suffix, or pattern.
package filter

import (
	"regexp"
	"strings"
)

// Filter holds compiled rules for matching environment variable keys.
type Filter struct {
	prefixes []string
	suffixes []string
	patterns []*regexp.Regexp
}

// New returns a Filter with no rules. Use With* methods to add rules.
func New() *Filter {
	return &Filter{}
}

// WithPrefix adds a prefix rule. Keys with this prefix will match.
func (f *Filter) WithPrefix(prefix string) *Filter {
	f.prefixes = append(f.prefixes, prefix)
	return f
}

// WithSuffix adds a suffix rule. Keys with this suffix will match.
func (f *Filter) WithSuffix(suffix string) *Filter {
	f.suffixes = append(f.suffixes, suffix)
	return f
}

// WithPattern adds a regex pattern rule. Keys matching the pattern will match.
// Returns an error if the pattern is invalid.
func (f *Filter) WithPattern(pattern string) (*Filter, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return f, err
	}
	f.patterns = append(f.patterns, re)
	return f, nil
}

// Match reports whether key satisfies at least one rule.
// If no rules are defined, all keys match.
func (f *Filter) Match(key string) bool {
	if len(f.prefixes) == 0 && len(f.suffixes) == 0 && len(f.patterns) == 0 {
		return true
	}
	for _, p := range f.prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	for _, s := range f.suffixes {
		if strings.HasSuffix(key, s) {
			return true
		}
	}
	for _, re := range f.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// Apply returns a new map containing only the entries whose keys match the filter.
func (f *Filter) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if f.Match(k) {
			out[k] = v
		}
	}
	return out
}

// Exclude returns a new map containing only entries whose keys do NOT match.
func (f *Filter) Exclude(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if !f.Match(k) {
			out[k] = v
		}
	}
	return out
}
