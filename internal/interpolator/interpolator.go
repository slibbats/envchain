// Package interpolator expands ${VAR} and $VAR references within env values
// using a provided environment map as the source of substitutions.
package interpolator

import (
	"os"
	"strings"
)

// Interpolator expands variable references in environment values.
type Interpolator struct {
	env    map[string]string
	useOS  bool
}

// New returns an Interpolator that resolves references from env.
// If fallbackOS is true, variables not found in env are looked up
// in the process environment.
func New(env map[string]string, fallbackOS bool) *Interpolator {
	return &Interpolator{env: env, useOS: fallbackOS}
}

// Expand replaces all ${VAR} and $VAR occurrences in s with their resolved
// values. Unresolved references are replaced with an empty string unless
// fallbackOS is enabled.
func (i *Interpolator) Expand(s string) string {
	return os.Expand(s, i.resolve)
}

// ExpandAll returns a new map with every value passed through Expand.
// Keys are left unchanged.
func (i *Interpolator) ExpandAll(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = i.Expand(v)
	}
	return out
}

// resolve is the mapping function supplied to os.Expand.
func (i *Interpolator) resolve(key string) string {
	if v, ok := i.env[key]; ok {
		return v
	}
	if i.useOS {
		return os.Getenv(key)
	}
	return ""
}

// HasReferences reports whether s contains any $VAR or ${VAR} references.
func HasReferences(s string) bool {
	return strings.Contains(s, "$")
}
