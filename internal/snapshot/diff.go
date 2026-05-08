package snapshot

import "sort"

// DiffEntry describes a single changed, added, or removed key between two snapshots.
type DiffEntry struct {
	Key    string
	Old    string
	New    string
	Status DiffStatus
}

// DiffStatus categorises the nature of a diff entry.
type DiffStatus string

const (
	Added    DiffStatus = "added"
	Removed  DiffStatus = "removed"
	Modified DiffStatus = "modified"
)

// Diff compares two snapshots and returns the set of differences.
// The result is sorted by key for deterministic output.
func Diff(base, next *Snapshot) []DiffEntry {
	var entries []DiffEntry

	for k, v := range next.Env {
		if old, ok := base.Env[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, New: v, Status: Added})
		} else if old != v {
			entries = append(entries, DiffEntry{Key: k, Old: old, New: v, Status: Modified})
		}
	}

	for k, v := range base.Env {
		if _, ok := next.Env[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, Old: v, Status: Removed})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}
