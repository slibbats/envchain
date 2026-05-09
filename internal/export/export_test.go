package export_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/export"
)

func TestWrite_DotenvFormat(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "production",
		"PORT":    "8080",
	}
	var sb strings.Builder
	if err := export.Write(&sb, env, export.FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "APP_ENV=production\n") {
		t.Errorf("missing APP_ENV line, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=8080\n") {
		t.Errorf("missing PORT line, got:\n%s", out)
	}
}

func TestWrite_ExportFormat(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	var sb strings.Builder
	if err := export.Write(&sb, env, export.FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "export DB_HOST=localhost\n"
	if sb.String() != want {
		t.Errorf("got %q, want %q", sb.String(), want)
	}
}

func TestWrite_DeterministicOrder(t *testing.T) {
	env := map[string]string{"Z": "last", "A": "first", "M": "middle"}
	var sb strings.Builder
	_ = export.Write(&sb, env, export.FormatDotenv)
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A=") {
		t.Errorf("first line should be A=, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z=") {
		t.Errorf("last line should be Z=, got %q", lines[2])
	}
}

func TestWrite_QuotesValuesWithSpaces(t *testing.T) {
	env := map[string]string{"GREETING": "hello world"}
	var sb strings.Builder
	_ = export.Write(&sb, env, export.FormatDotenv)
	want := "GREETING='hello world'\n"
	if sb.String() != want {
		t.Errorf("got %q, want %q", sb.String(), want)
	}
}

func TestWrite_EmptyValue(t *testing.T) {
	env := map[string]string{"EMPTY": ""}
	var sb strings.Builder
	_ = export.Write(&sb, env, export.FormatDotenv)
	want := "EMPTY=''\n"
	if sb.String() != want {
		t.Errorf("got %q, want %q", sb.String(), want)
	}
}

func TestWrite_EscapesSingleQuoteInValue(t *testing.T) {
	env := map[string]string{"MSG": "it's alive"}
	var sb strings.Builder
	_ = export.Write(&sb, env, export.FormatDotenv)
	out := sb.String()
	// Should contain the escaped form.
	if !strings.Contains(out, `'\''`) {
		t.Errorf("expected escaped single-quote in output, got: %q", out)
	}
}
