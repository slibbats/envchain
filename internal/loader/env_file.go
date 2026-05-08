package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvFile represents a parsed .env file with key-value pairs.
type EnvFile struct {
	path   string
	values map[string]string
}

// LoadEnvFile reads and parses a .env file from the given path.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are also ignored.
// Each valid line must be in KEY=VALUE format.
func LoadEnvFile(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: open %q: %w", path, err)
	}
	defer f.Close()

	ef := &EnvFile{
		path:   path,
		values: make(map[string]string),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			return nil, fmt.Errorf("loader: %q line %d: missing '=' separator", path, lineNum)
		}

		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("loader: %q line %d: empty key", path, lineNum)
		}

		// Strip optional surrounding quotes from value.
		value = strings.TrimSpace(value)
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		ef.values[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: scan %q: %w", path, err)
	}

	return ef, nil
}

// Values returns a copy of the parsed key-value pairs.
func (ef *EnvFile) Values() map[string]string {
	copy := make(map[string]string, len(ef.values))
	for k, v := range ef.values {
		copy[k] = v
	}
	return copy
}

// Path returns the file path that was loaded.
func (ef *EnvFile) Path() string {
	return ef.path
}
