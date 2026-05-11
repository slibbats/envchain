package profile_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/profile"
)

func makeConfig(names ...string) *profile.Config {
	profiles := make([]profile.Profile, len(names))
	for i, n := range names {
		profiles[i] = profile.Profile{Name: n, Layers: []string{".env"}}
	}
	return &profile.Config{Profiles: profiles}
}

func TestSelect_ExplicitName(t *testing.T) {
	cfg := makeConfig("dev", "prod")
	sel := profile.NewSelector(cfg)
	p, err := sel.Select("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "dev" {
		t.Errorf("want dev, got %q", p.Name)
	}
}

func TestSelect_FallbackToEnv(t *testing.T) {
	cfg := makeConfig("staging")
	sel := profile.NewSelector(cfg)
	// inject env lookup
	sel2 := &profile.Selector{} // use exported test helper below
	_ = sel2

	// Use the exported selector with env override via t.Setenv
	t.Setenv("ENVCHAIN_PROFILE", "staging")
	p, err := sel.Select("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "staging" {
		t.Errorf("want staging, got %q", p.Name)
	}
}

func TestSelect_NoNameNoEnv(t *testing.T) {
	t.Setenv("ENVCHAIN_PROFILE", "")
	cfg := makeConfig("dev")
	sel := profile.NewSelector(cfg)
	_, err := sel.Select("")
	if err == nil {
		t.Fatal("expected error when no profile provided")
	}
}

func TestSelect_UnknownProfile(t *testing.T) {
	cfg := makeConfig("dev")
	sel := profile.NewSelector(cfg)
	_, err := sel.Select("unknown")
	if err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestActiveName_ExplicitTakesPrecedence(t *testing.T) {
	t.Setenv("ENVCHAIN_PROFILE", "prod")
	cfg := makeConfig("dev", "prod")
	sel := profile.NewSelector(cfg)
	if got := sel.ActiveName("dev"); got != "dev" {
		t.Errorf("want dev, got %q", got)
	}
}

func TestActiveName_FallsBackToEnv(t *testing.T) {
	t.Setenv("ENVCHAIN_PROFILE", "prod")
	cfg := makeConfig("prod")
	sel := profile.NewSelector(cfg)
	if got := sel.ActiveName(""); got != "prod" {
		t.Errorf("want prod, got %q", got)
	}
}
