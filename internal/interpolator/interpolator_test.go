package interpolator_test

import (
	"os"
	"testing"

	"github.com/yourorg/envchain/internal/interpolator"
)

func TestExpand_SimpleVariable(t *testing.T) {
	env := map[string]string{"HOME": "/home/user"}
	ip := interpolator.New(env, false)
	got := ip.Expand("/root:${HOME}/bin")
	want := "/root:/home/user/bin"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestExpand_DollarSyntax(t *testing.T) {
	env := map[string]string{"USER": "alice"}
	ip := interpolator.New(env, false)
	got := ip.Expand("Hello $USER!")
	if got != "Hello alice!" {
		t.Errorf("unexpected result: %q", got)
	}
}

func TestExpand_UnresolvedNoFallback(t *testing.T) {
	ip := interpolator.New(map[string]string{}, false)
	got := ip.Expand("${MISSING}")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestExpand_FallbackToOS(t *testing.T) {
	t.Setenv("OS_VAR", "from-os")
	ip := interpolator.New(map[string]string{}, true)
	got := ip.Expand("${OS_VAR}")
	if got != "from-os" {
		t.Errorf("expected 'from-os', got %q", got)
	}
}

func TestExpand_EnvTakesPrecedenceOverOS(t *testing.T) {
	os.Setenv("PRIO", "os-value")
	t.Cleanup(func() { os.Unsetenv("PRIO") })
	env := map[string]string{"PRIO": "env-value"}
	ip := interpolator.New(env, true)
	got := ip.Expand("${PRIO}")
	if got != "env-value" {
		t.Errorf("expected 'env-value', got %q", got)
	}
}

func TestExpandAll_ExpandsAllValues(t *testing.T) {
	env := map[string]string{
		"BASE": "/opt",
		"BIN":  "${BASE}/bin",
		"PLAIN": "no-refs",
	}
	ip := interpolator.New(env, false)
	out := ip.ExpandAll(env)
	if out["BIN"] != "/opt/bin" {
		t.Errorf("BIN: got %q, want '/opt/bin'", out["BIN"])
	}
	if out["PLAIN"] != "no-refs" {
		t.Errorf("PLAIN: got %q, want 'no-refs'", out["PLAIN"])
	}
	if out["BASE"] != "/opt" {
		t.Errorf("BASE: got %q, want '/opt'", out["BASE"])
	}
}

func TestExpandAll_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "${B}", "B": "hello"}
	ip := interpolator.New(env, false)
	_ = ip.ExpandAll(env)
	if env["A"] != "${B}" {
		t.Error("input map was mutated")
	}
}

func TestHasReferences_True(t *testing.T) {
	if !interpolator.HasReferences("${FOO}") {
		t.Error("expected true for '${FOO}'")
	}
}

func TestHasReferences_False(t *testing.T) {
	if interpolator.HasReferences("no-dollar-here") {
		t.Error("expected false for plain string")
	}
}
