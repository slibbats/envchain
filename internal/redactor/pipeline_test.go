package redactor_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/redactor"
)

func TestEnvSlice_ContainsAllKeys(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{
		"HOME":     "/home/user",
		"API_KEY":  "secret123",
		"LOG_LEVEL": "info",
	}
	slice := r.EnvSlice(env)
	if len(slice) != len(env) {
		t.Fatalf("expected %d entries, got %d", len(env), len(slice))
	}
	for _, entry := range slice {
		if !strings.Contains(entry, "=") {
			t.Errorf("entry %q missing '='", entry)
		}
	}
}

func TestEnvSlice_RedactsSensitive(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{"API_KEY": "real-secret"}
	slice := r.EnvSlice(env)
	if len(slice) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if !strings.Contains(slice[0], redactor.Placeholder) {
		t.Errorf("expected placeholder in %q", slice[0])
	}
	if strings.Contains(slice[0], "real-secret") {
		t.Errorf("real secret must not appear in output")
	}
}

func TestFilterKeys_ExcludesSensitive(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{
		"HOME":        "/home/user",
		"DB_PASSWORD": "hunter2",
		"PORT":        "5432",
	}
	out := r.FilterKeys(env)
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been filtered out")
	}
	if out["HOME"] != "/home/user" {
		t.Errorf("HOME should be present, got %q", out["HOME"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("PORT should be present, got %q", out["PORT"])
	}
}

func TestSummary_CorrectCounts(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{
		"HOME":     "/home/user",
		"API_KEY":  "abc",
		"DB_TOKEN": "xyz",
	}
	summary := r.Summary(env)
	if !strings.HasPrefix(summary, "2/3") {
		t.Errorf("unexpected summary: %q", summary)
	}
}

func TestSummary_NoneRedacted(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{"HOME": "/home/user", "PORT": "8080"}
	summary := r.Summary(env)
	if !strings.HasPrefix(summary, "0/2") {
		t.Errorf("unexpected summary: %q", summary)
	}
}
