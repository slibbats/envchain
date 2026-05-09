// Package validator provides utilities for validating environment variable
// keys and values before they are injected into the chain. It enforces naming
// conventions and detects potentially unsafe or malformed entries.
package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// validKeyPattern matches POSIX-compliant environment variable names:
// must start with a letter or underscore, followed by letters, digits, or underscores.
var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ErrEmptyKey is returned when a key is an empty string.
var ErrEmptyKey = errors.New("validator: key must not be empty")

// ErrInvalidKey is returned when a key contains illegal characters.
type ErrInvalidKey struct {
	Key string
}

func (e *ErrInvalidKey) Error() string {
	return fmt.Sprintf("validator: invalid key %q: must match [A-Za-z_][A-Za-z0-9_]*", e.Key)
}

// ErrNullByte is returned when a value contains a null byte.
type ErrNullByte struct {
	Key string
}

func (e *ErrNullByte) Error() string {
	return fmt.Sprintf("validator: value for key %q contains a null byte", e.Key)
}

// ValidateKey checks that k is a non-empty, POSIX-compliant identifier.
func ValidateKey(k string) error {
	if k == "" {
		return ErrEmptyKey
	}
	if !validKeyPattern.MatchString(k) {
		return &ErrInvalidKey{Key: k}
	}
	return nil
}

// ValidateValue checks that v does not contain a null byte.
func ValidateValue(key, value string) error {
	if strings.ContainsRune(value, '\x00') {
		return &ErrNullByte{Key: key}
	}
	return nil
}

// ValidateEnv validates an entire map of environment variables, returning the
// first error encountered. Iteration order is non-deterministic; callers that
// need deterministic error ordering should sort keys before calling.
func ValidateEnv(env map[string]string) error {
	for k, v := range env {
		if err := ValidateKey(k); err != nil {
			return err
		}
		if err := ValidateValue(k, v); err != nil {
			return err
		}
	}
	return nil
}
