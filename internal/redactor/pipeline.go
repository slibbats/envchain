package redactor

import "fmt"

// EnvSlice converts a key=value map to a sorted slice of strings, redacting
// sensitive keys. The output format matches os/exec Cmd.Env expectations.
func (r *Redactor) EnvSlice(env map[string]string) []string {
	slice := make([]string, 0, len(env))
	for k, v := range env {
		safe := r.RedactValue(k, v)
		slice = append(slice, fmt.Sprintf("%s=%s", k, safe))
	}
	return slice
}

// FilterKeys returns a new map containing only the entries whose keys are NOT
// sensitive. This is useful when constructing child-process environments where
// secrets should be excluded entirely rather than replaced.
func (r *Redactor) FilterKeys(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if !r.IsSensitive(k) {
			out[k] = v
		}
	}
	return out
}

// Summary returns a human-readable string describing how many keys were
// redacted out of the total.
func (r *Redactor) Summary(env map[string]string) string {
	total := len(env)
	redacted := 0
	for k := range env {
		if r.IsSensitive(k) {
			redacted++
		}
	}
	return fmt.Sprintf("%d/%d keys redacted", redacted, total)
}
