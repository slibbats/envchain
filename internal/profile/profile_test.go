package profile_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain/internal/profile"
)

func writeTempConfig(t *testing.T, cfg profile.Config) string {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	p := filepath.Join(t.TempDir(), "profiles.json")
	if err := os.WriteFile(p, data, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func TestLoadConfig_ValidProfiles(t *testing.T) {
	cfg := profile.Config{
		Profiles: []profile.Profile{
			{Name: "dev", Layers: []string{".env", ".env.dev"}},
			{Name: "prod", Layers: []string{".env", ".env.prod"}},
		},
	}
	path := writeTempConfig(t, cfg)
	loaded, err := profile.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded.Profiles) != 2 {
		t.Fatalf("want 2 profiles, got %d", len(loaded.Profiles))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := profile.LoadConfig("/nonexistent/profiles.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "profiles.json")
	os.WriteFile(p, []byte("not json"), 0o600)
	_, err := profile.LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestValidate_EmptyName(t *testing.T) {
	p := profile.Profile{Name: "", Layers: []string{".env"}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestValidate_InvalidCharacters(t *testing.T) {
	p := profile.Profile{Name: "my env!", Layers: []string{".env"}}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for invalid characters")
	}
}

func TestValidate_NoLayers(t *testing.T) {
	p := profile.Profile{Name: "dev", Layers: nil}
	if err := p.Validate(); err == nil {
		t.Fatal("expected error for no layers")
	}
}

func TestConfig_Get_Found(t *testing.T) {
	cfg := &profile.Config{
		Profiles: []profile.Profile{
			{Name: "staging", Layers: []string{".env.staging"}},
		},
	}
	p, err := cfg.Get("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "staging" {
		t.Errorf("want staging, got %q", p.Name)
	}
}

func TestConfig_Get_NotFound(t *testing.T) {
	cfg := &profile.Config{}
	_, err := cfg.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestConfig_Names(t *testing.T) {
	cfg := &profile.Config{
		Profiles: []profile.Profile{
			{Name: "dev", Layers: []string{".env"}},
			{Name: "prod", Layers: []string{".env"}},
		},
	}
	names := cfg.Names()
	if len(names) != 2 || names[0] != "dev" || names[1] != "prod" {
		t.Errorf("unexpected names: %v", names)
	}
}
