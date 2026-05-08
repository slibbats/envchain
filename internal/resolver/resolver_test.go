package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain/envchain/internal/resolver"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestBuild_SingleLayer(t *testing.T) {
	p := writeTempEnv(t, "APP_ENV=dev\nDEBUG=true\n")
	cfg := resolver.Config{
		Layers: []resolver.LayerConfig{{Name: "dev", FilePath: p}},
	}
	c, err := resolver.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, _ := c.Get("APP_ENV"); got != "dev" {
		t.Errorf("APP_ENV: got %q, want %q", got, "dev")
	}
}

func TestBuild_HigherPriorityOverrides(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=dev\nSECRET=base-secret\n")
	prod := writeTempEnv(t, "APP_ENV=prod\n")
	cfg := resolver.Config{
		Layers: []resolver.LayerConfig{
			{Name: "base", FilePath: base},
			{Name: "prod", FilePath: prod},
		},
	}
	c, err := resolver.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, _ := c.Get("APP_ENV"); got != "prod" {
		t.Errorf("APP_ENV: got %q, want %q", got, "prod")
	}
	if got, _ := c.Get("SECRET"); got != "base-secret" {
		t.Errorf("SECRET: got %q, want %q", got, "base-secret")
	}
}

func TestBuild_MissingFileSkipped(t *testing.T) {
	cfg := resolver.Config{
		Layers: []resolver.LayerConfig{
			{Name: "ghost", FilePath: "/nonexistent/.env.ghost"},
		},
	}
	_, err := resolver.Build(cfg)
	if err != nil {
		t.Fatalf("expected missing file to be skipped, got error: %v", err)
	}
}

func TestBuild_InjectOS(t *testing.T) {
	t.Setenv("ENVCHAIN_TEST_VAR", "from-os")
	cfg := resolver.Config{
		InjectOS: true,
		Layers:   []resolver.LayerConfig{},
	}
	c, err := resolver.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, ok := c.Get("ENVCHAIN_TEST_VAR"); !ok || got != "from-os" {
		t.Errorf("ENVCHAIN_TEST_VAR: got %q ok=%v, want %q", got, ok, "from-os")
	}
}

func TestBuild_EmptyFilePath_Skipped(t *testing.T) {
	cfg := resolver.Config{
		Layers: []resolver.LayerConfig{{Name: "empty", FilePath: ""}},
	}
	_, err := resolver.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
