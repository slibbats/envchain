package redactor_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/redactor"
)

func TestIsSensitive_MatchesKnownPatterns(t *testing.T) {
	r := redactor.New(nil)
	sensitiveKeys := []string{
		"DB_PASSWORD",
		"API_KEY",
		"AUTH_TOKEN",
		"PRIVATE_KEY",
		"AWS_SECRET_ACCESS_KEY",
		"user_credential",
	}
	for _, k := range sensitiveKeys {
		if !r.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_NonSensitive(t *testing.T) {
	r := redactor.New(nil)
	safeKeys := []string{"HOME", "PATH", "USER", "PORT", "LOG_LEVEL"}
	for _, k := range safeKeys {
		if r.IsSensitive(k) {
			t.Errorf("expected %q to be non-sensitive", k)
		}
	}
}

func TestRedact_ReplacesOnlySensitive(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{
		"HOME":        "/home/user",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_KEY":     "abc123",
	}
	out := r.Redact(env)

	if out["HOME"] != "/home/user" {
		t.Errorf("HOME should be unchanged, got %q", out["HOME"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", out["PORT"])
	}
	if out["DB_PASSWORD"] != redactor.Placeholder {
		t.Errorf("DB_PASSWORD should be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != redactor.Placeholder {
		t.Errorf("API_KEY should be redacted, got %q", out["API_KEY"])
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{"API_KEY": "real-value"}
	_ = r.Redact(env)
	if env["API_KEY"] != "real-value" {
		t.Error("original map should not be mutated")
	}
}

func TestRedactValue_SensitiveKey(t *testing.T) {
	r := redactor.New(nil)
	got := r.RedactValue("DB_PASSWORD", "hunter2")
	if got != redactor.Placeholder {
		t.Errorf("expected placeholder, got %q", got)
	}
}

func TestRedactValue_SafeKey(t *testing.T) {
	r := redactor.New(nil)
	got := r.RedactValue("LOG_LEVEL", "debug")
	if got != "debug" {
		t.Errorf("expected %q, got %q", "debug", got)
	}
}

func TestNew_CustomPatterns(t *testing.T) {
	r := redactor.New([]string{"custom"})
	if !r.IsSensitive("MY_CUSTOM_VAR") {
		t.Error("expected MY_CUSTOM_VAR to be sensitive with custom pattern")
	}
	if r.IsSensitive("API_KEY") {
		t.Error("API_KEY should not be sensitive with custom-only patterns")
	}
}
