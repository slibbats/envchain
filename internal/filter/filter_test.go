package filter_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/filter"
)

func TestMatch_NoRules_AllMatch(t *testing.T) {
	f := filter.New()
	if !f.Match("ANY_KEY") {
		t.Error("expected all keys to match when no rules defined")
	}
}

func TestMatch_Prefix(t *testing.T) {
	f := filter.New().WithPrefix("APP_")
	if !f.Match("APP_SECRET") {
		t.Error("expected APP_SECRET to match prefix APP_")
	}
	if f.Match("DB_HOST") {
		t.Error("expected DB_HOST not to match prefix APP_")
	}
}

func TestMatch_Suffix(t *testing.T) {
	f := filter.New().WithSuffix("_KEY")
	if !f.Match("API_KEY") {
		t.Error("expected API_KEY to match suffix _KEY")
	}
	if f.Match("API_SECRET") {
		t.Error("expected API_SECRET not to match suffix _KEY")
	}
}

func TestMatch_Pattern(t *testing.T) {
	f, err := filter.New().WithPattern(`^DB_.*_HOST$`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match("DB_PRIMARY_HOST") {
		t.Error("expected DB_PRIMARY_HOST to match pattern")
	}
	if f.Match("DB_HOST_PORT") {
		t.Error("expected DB_HOST_PORT not to match pattern")
	}
}

func TestMatch_InvalidPattern(t *testing.T) {
	_, err := filter.New().WithPattern(`[invalid`)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestApply_ReturnsMatchingKeys(t *testing.T) {
	env := map[string]string{
		"APP_TOKEN":  "abc",
		"APP_SECRET": "xyz",
		"DB_HOST":    "localhost",
	}
	f := filter.New().WithPrefix("APP_")
	got := f.Apply(env)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if _, ok := got["DB_HOST"]; ok {
		t.Error("DB_HOST should have been excluded")
	}
}

func TestExclude_ReturnsNonMatchingKeys(t *testing.T) {
	env := map[string]string{
		"APP_TOKEN": "abc",
		"DB_HOST":   "localhost",
		"DB_PORT":   "5432",
	}
	f := filter.New().WithPrefix("DB_")
	got := f.Exclude(env)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if _, ok := got["APP_TOKEN"]; !ok {
		t.Error("APP_TOKEN should be present after exclude")
	}
}

func TestMatch_MultipleRules_AnyMatch(t *testing.T) {
	f := filter.New().WithPrefix("APP_").WithSuffix("_URL")
	if !f.Match("SERVICE_URL") {
		t.Error("expected SERVICE_URL to match via suffix rule")
	}
	if !f.Match("APP_DEBUG") {
		t.Error("expected APP_DEBUG to match via prefix rule")
	}
	if f.Match("DB_HOST") {
		t.Error("expected DB_HOST not to match any rule")
	}
}
