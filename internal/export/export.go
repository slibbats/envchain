// Package export provides utilities for serialising a resolved environment
// map into common shell-compatible formats so that it can be sourced or
// consumed by external processes.
package export

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format controls the output syntax produced by Write.
type Format string

const (
	// FormatExport emits POSIX `export KEY=VALUE` lines suitable for eval.
	FormatExport Format = "export"
	// FormatDotenv emits plain `KEY=VALUE` lines compatible with .env loaders.
	FormatDotenv Format = "dotenv"
)

// Write serialises env to w using the requested Format.
// Keys are written in lexicographic order so output is deterministic.
// Values that contain whitespace or special shell characters are
// single-quoted to prevent unintended word-splitting.
func Write(w io.Writer, env map[string]string, format Format) error {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := quote(env[k])
		var line string
		switch format {
		case FormatExport:
			line = fmt.Sprintf("export %s=%s\n", k, v)
		default: // FormatDotenv
			line = fmt.Sprintf("%s=%s\n", k, v)
		}
		if _, err := io.WriteString(w, line); err != nil {
			return fmt.Errorf("export: write key %q: %w", k, err)
		}
	}
	return nil
}

// quote wraps v in single-quotes when it contains characters that a POSIX
// shell would interpret, escaping any literal single-quotes inside.
func quote(v string) string {
	if v == "" {
		return "''"
	}
	needsQuoting := strings.ContainsAny(v, " \t\n$`\\\"'!#&;|<>(){}")
	if !needsQuoting {
		return v
	}
	// Escape embedded single-quotes: ' → '\''.
	escaped := strings.ReplaceAll(v, "'", `'\''`)
	return "'" + escaped + "'"
}
