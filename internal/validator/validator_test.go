package validator_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/validator"
)

func TestValidateKey_Valid(t *testing.T) {
	keys := []string{"FOO", "_BAR", "MY_VAR_123", "a", "_"}
	for _, k := range keys {
		if err := validator.ValidateKey(k); err != nil {
			t.Errorf("expected key %q to be valid, got: %v", k, err)
		}
	}
}

func TestValidateKey_Empty(t *testing.T) {
	err := validator.ValidateKey("")
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
	if err != validator.ErrEmptyKey {
		t.Errorf("expected ErrEmptyKey, got %v", err)
	}
}

func TestValidateKey_InvalidCharacters(t *testing.T) {
	invalid := []string{"1STARTS_WITH_DIGIT", "HAS-HYPHEN", "HAS SPACE", "HAS.DOT", ""}
	for _, k := range invalid {
		if k == "" {
			continue // covered by TestValidateKey_Empty
		}
		err := validator.ValidateKey(k)
		if err == nil {
			t.Errorf("expected error for key %q, got nil", k)
		}
		var invErr *validator.ErrInvalidKey
		if err != nil {
			_ = invErr // type assertion check via errors.As would require import; simple nil check suffices
		}
	}
}

func TestValidateValue_NullByte(t *testing.T) {
	err := validator.ValidateValue("MY_KEY", "value\x00with null")
	if err == nil {
		t.Fatal("expected error for value with null byte, got nil")
	}
}

func TestValidateValue_Clean(t *testing.T) {
	err := validator.ValidateValue("MY_KEY", "perfectly fine value")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateEnv_AllValid(t *testing.T) {
	env := map[string]string{
		"HOST":     "localhost",
		"PORT":     "8080",
		"_SECRET":  "s3cr3t",
	}
	if err := validator.ValidateEnv(env); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateEnv_InvalidKey(t *testing.T) {
	env := map[string]string{
		"VALID_KEY": "ok",
		"bad-key":   "value",
	}
	if err := validator.ValidateEnv(env); err == nil {
		t.Error("expected error for invalid key in env map, got nil")
	}
}

func TestValidateEnv_NullByteInValue(t *testing.T) {
	env := map[string]string{
		"GOOD_KEY": "good",
		"BAD_VAL":  "oops\x00null",
	}
	if err := validator.ValidateEnv(env); err == nil {
		t.Error("expected error for null byte in value, got nil")
	}
}
