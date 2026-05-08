package audit_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/audit"
	"github.com/yourorg/envchain/internal/chain"
)

func newChainWithLayer(t *testing.T, name string, kv map[string]string) *chain.Chain {
	t.Helper()
	c := chain.New()
	if err := c.AddLayer(name, kv, 10); err != nil {
		t.Fatalf("AddLayer: %v", err)
	}
	return c
}

func TestWrap_GetFound_LogsResolved(t *testing.T) {
	c := newChainWithLayer(t, "base", map[string]string{"HOST": "localhost"})
	var buf bytes.Buffer
	ac := audit.Wrap(c, audit.New(&buf))
	val, ok := ac.Get("HOST")
	if !ok || val != "localhost" {
		t.Fatalf("unexpected Get result: %q %v", val, ok)
	}
	if !strings.Contains(buf.String(), "RESOLVED") {
		t.Error("expected RESOLVED in audit log")
	}
	if !strings.Contains(buf.String(), "HOST") {
		t.Error("expected key HOST in audit log")
	}
}

func TestWrap_GetMissing_LogsMissing(t *testing.T) {
	c := chain.New()
	var buf bytes.Buffer
	ac := audit.Wrap(c, audit.New(&buf))
	_, ok := ac.Get("NONEXISTENT")
	if ok {
		t.Fatal("expected key to be missing")
	}
	if !strings.Contains(buf.String(), "MISSING") {
		t.Error("expected MISSING in audit log")
	}
}

func TestWrap_NilLogger_UsesDiscard(t *testing.T) {
	c := newChainWithLayer(t, "base", map[string]string{"X": "1"})
	ac := audit.Wrap(c, nil)
	_, ok := ac.Get("X")
	if !ok {
		t.Fatal("expected key to be found")
	}
}

func TestWrap_Resolve_LogsAllKeys(t *testing.T) {
	c := newChainWithLayer(t, "base", map[string]string{"A": "1", "B": "2"})
	var buf bytes.Buffer
	ac := audit.Wrap(c, audit.New(&buf))
	env := ac.Resolve()
	if len(env) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(env))
	}
	log := buf.String()
	if !strings.Contains(log, "RESOLVED") {
		t.Error("expected RESOLVED events in audit log")
	}
}
