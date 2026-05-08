package chain_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/chain"
)

func TestAddLayer_DuplicateName(t *testing.T) {
	c := chain.New()
	if err := c.AddLayer("dev", map[string]string{"KEY": "dev_val"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := c.AddLayer("dev", map[string]string{"KEY": "other"}); err == nil {
		t.Fatal("expected error for duplicate layer name, got nil")
	}
}

func TestAddLayer_EmptyName(t *testing.T) {
	c := chain.New()
	if err := c.AddLayer("", map[string]string{}); err == nil {
		t.Fatal("expected error for empty layer name, got nil")
	}
}

func TestResolve_HigherPriorityOverrides(t *testing.T) {
	c := chain.New()
	_ = c.AddLayer("base", map[string]string{"DB_HOST": "localhost", "LOG_LEVEL": "debug"})
	_ = c.AddLayer("prod", map[string]string{"DB_HOST": "prod.db.internal"})

	resolved := c.Resolve()

	if resolved["DB_HOST"] != "prod.db.internal" {
		t.Errorf("expected prod.db.internal, got %q", resolved["DB_HOST"])
	}
	if resolved["LOG_LEVEL"] != "debug" {
		t.Errorf("expected debug, got %q", resolved["LOG_LEVEL"])
	}
}

func TestGet_KeyFound(t *testing.T) {
	c := chain.New()
	_ = c.AddLayer("dev", map[string]string{"SECRET": "s3cr3t"})

	v, ok := c.Get("SECRET")
	if !ok {
		t.Fatal("expected key to be found")
	}
	if v != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %q", v)
	}
}

func TestGet_KeyMissing(t *testing.T) {
	c := chain.New()
	_, ok := c.Get("NONEXISTENT")
	if ok {
		t.Fatal("expected key to be missing")
	}
}

func TestLayerNames_Order(t *testing.T) {
	c := chain.New()
	_ = c.AddLayer("base", map[string]string{})
	_ = c.AddLayer("staging", map[string]string{})
	_ = c.AddLayer("prod", map[string]string{})

	names := c.LayerNames()
	expected := []string{"base", "staging", "prod"}
	for i, name := range expected {
		if names[i] != name {
			t.Errorf("position %d: expected %q, got %q", i, name, names[i])
		}
	}
}
