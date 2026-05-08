package snapshot_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/snapshot"
)

func TestDiff_Added(t *testing.T) {
	base := snapshot.New(map[string]string{"A": "1"}, nil)
	next := snapshot.New(map[string]string{"A": "1", "B": "2"}, nil)

	entries := snapshot.Diff(base, next)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != snapshot.Added || entries[0].Key != "B" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestDiff_Removed(t *testing.T) {
	base := snapshot.New(map[string]string{"A": "1", "B": "2"}, nil)
	next := snapshot.New(map[string]string{"A": "1"}, nil)

	entries := snapshot.Diff(base, next)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != snapshot.Removed || entries[0].Key != "B" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestDiff_Modified(t *testing.T) {
	base := snapshot.New(map[string]string{"KEY": "old"}, nil)
	next := snapshot.New(map[string]string{"KEY": "new"}, nil)

	entries := snapshot.Diff(base, next)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Status != snapshot.Modified || e.Old != "old" || e.New != "new" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	base := snapshot.New(env, nil)
	next := snapshot.New(env, nil)

	if entries := snapshot.Diff(base, next); len(entries) != 0 {
		t.Errorf("expected no diff, got %d entries", len(entries))
	}
}

func TestDiff_SortedByKey(t *testing.T) {
	base := snapshot.New(map[string]string{}, nil)
	next := snapshot.New(map[string]string{"Z": "1", "A": "2", "M": "3"}, nil)

	entries := snapshot.Diff(base, next)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "A" || entries[1].Key != "M" || entries[2].Key != "Z" {
		t.Errorf("entries not sorted: %v %v %v", entries[0].Key, entries[1].Key, entries[2].Key)
	}
}
