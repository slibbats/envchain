package masker_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/masker"
)

func TestIsSensitive_DefaultKeys(t *testing.T) {
	m := masker.New(nil)

	sensitive := []string{
		"DB_PASSWORD",
		"API_KEY",
		"GITHUB_TOKEN",
		"AWS_SECRET_ACCESS_KEY",
		"PRIVATE_KEY_PATH",
		"AUTH_HEADER",
	}
	for _, key := range sensitive {
		if !m.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_NonSensitive(t *testing.T) {
	m := masker.New(nil)

	plain := []string{
		"APP_ENV",
		"PORT",
		"LOG_LEVEL",
		"DATABASE_HOST",
	}
	for _, key := range plain {
		if m.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestMaskValue_SensitiveReturnsPlaceholder(t *testing.T) {
	m := masker.New(nil)
	got := m.MaskValue("DB_PASSWORD", "supersecret")
	if got != "***" {
		t.Errorf("expected *** got %q", got)
	}
}

func TestMaskValue_PlainPassthrough(t *testing.T) {
	m := masker.New(nil)
	got := m.MaskValue("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected 'production' got %q", got)
	}
}

func TestMaskEnv_MixedKeys(t *testing.T) {
	m := masker.New(nil)
	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "hunter2",
		"PORT":        "8080",
		"API_KEY":     "abc123",
	}

	masked := m.MaskEnv(env)

	if masked["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged")
	}
	if masked["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged")
	}
	if masked["DB_PASSWORD"] != "***" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if masked["API_KEY"] != "***" {
		t.Errorf("API_KEY should be masked")
	}
}

func TestMaskEnv_DoesNotMutateOriginal(t *testing.T) {
	m := masker.New(nil)
	env := map[string]string{"DB_PASSWORD": "hunter2"}
	m.MaskEnv(env)
	if env["DB_PASSWORD"] != "hunter2" {
		t.Error("original map should not be mutated")
	}
}

func TestNew_CustomKeys(t *testing.T) {
	m := masker.New([]string{"INTERNAL"})
	if !m.IsSensitive("INTERNAL_CONFIG") {
		t.Error("expected INTERNAL_CONFIG to be sensitive with custom keys")
	}
	if m.IsSensitive("API_KEY") {
		t.Error("API_KEY should not be sensitive with custom keys only")
	}
}
