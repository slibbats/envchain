package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/snapshot"
)

func TestNew_CopiesEnv(t *testing.T) {
	original := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.New(original, []string{"base", "prod"})

	original["FOO"] = "mutated"
	if s.Env["FOO"] != "bar" {
		t.Errorf("expected snapshot to be isolated from original map, got %q", s.Env["FOO"])
	}
}

func TestNew_SetsCreatedAt(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.New(map[string]string{}, nil)
	after := time.Now().UTC()

	if s.CreatedAt.Before(before) || s.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v not in expected range [%v, %v]", s.CreatedAt, before, after)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	env := map[string]string{"KEY": "value", "OTHER": "123"}
	layers := []string{"base", "staging"}
	s := snapshot.New(env, layers)

	if err := s.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Env["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", loaded.Env["KEY"])
	}
	if len(loaded.Layers) != 2 || loaded.Layers[1] != "staging" {
		t.Errorf("unexpected layers: %v", loaded.Layers)
	}
}

func TestSave_RestrictedPermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	s := snapshot.New(map[string]string{"SECRET": "abc"}, nil)
	if err := s.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0600)

	if _, err := snapshot.Load(path); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestKeys_ReturnsAllKeys(t *testing.T) {
	s := snapshot.New(map[string]string{"A": "1", "B": "2", "C": "3"}, nil)
	keys := s.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}
