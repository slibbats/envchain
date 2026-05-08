package audit_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/audit"
)

func TestNew_WritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Record(audit.EventResolved, "DB_URL", "prod", "value resolved")
	if buf.Len() == 0 {
		t.Fatal("expected output, got none")
	}
}

func TestRecord_ContainsExpectedFields(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Record(audit.EventResolved, "API_KEY", "staging", "ok")
	out := buf.String()
	for _, want := range []string{"RESOLVED", "API_KEY", "staging", "ok"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got: %s", want, out)
		}
	}
}

func TestRecord_MissingEvent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Record(audit.EventMissing, "UNKNOWN_KEY", "", "not found in any layer")
	if !strings.Contains(buf.String(), "MISSING") {
		t.Error("expected MISSING in output")
	}
}

func TestRecord_OverriddenEvent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Record(audit.EventOverridden, "PORT", "dev", "overridden by higher priority layer")
	if !strings.Contains(buf.String(), "OVERRIDDEN") {
		t.Error("expected OVERRIDDEN in output")
	}
}

func TestDiscard_ProducesNoOutput(t *testing.T) {
	l := audit.Discard()
	// Should not panic and should produce no visible side effects.
	l.Record(audit.EventMasked, "SECRET", "prod", "value masked")
}

func TestNew_NilUsesStderr(t *testing.T) {
	// Just ensure New(nil) doesn't panic.
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}
