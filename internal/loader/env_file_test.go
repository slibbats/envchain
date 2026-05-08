package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestLoadEnvFile_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	ef, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vals := ef.Values()
	if vals["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", vals["FOO"])
	}
	if vals["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", vals["BAZ"])
	}
}

func TestLoadEnvFile_CommentsAndBlanks(t *testing.T) {
	content := "# this is a comment\n\nKEY=value\n"
	path := writeTempEnv(t, content)
	ef, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vals := ef.Values()
	if len(vals) != 1 {
		t.Errorf("expected 1 key, got %d", len(vals))
	}
	if vals["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", vals["KEY"])
	}
}

func TestLoadEnvFile_QuotedValues(t *testing.T) {
	content := `DB_URL="postgres://localhost/mydb"` + "\nSECRET='top secret'\n"
	path := writeTempEnv(t, content)
	ef, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vals := ef.Values()
	if vals["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("unexpected DB_URL: %q", vals["DB_URL"])
	}
	if vals["SECRET"] != "top secret" {
		t.Errorf("unexpected SECRET: %q", vals["SECRET"])
	}
}

func TestLoadEnvFile_MissingSeparator(t *testing.T) {
	path := writeTempEnv(t, "INVALIDLINE\n")
	_, err := LoadEnvFile(path)
	if err == nil {
		t.Fatal("expected error for missing '=' separator, got nil")
	}
}

func TestLoadEnvFile_FileNotFound(t *testing.T) {
	_, err := LoadEnvFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadEnvFile_PathReturned(t *testing.T) {
	path := writeTempEnv(t, "X=1\n")
	ef, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Path() != path {
		t.Errorf("expected path %q, got %q", path, ef.Path())
	}
}
