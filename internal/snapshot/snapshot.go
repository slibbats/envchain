// Package snapshot provides functionality to capture and persist the
// resolved environment variable state to disk for auditing and reproducibility.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of resolved environment variables.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Layers    []string          `json:"layers"`
	Env       map[string]string `json:"env"`
}

// New creates a new Snapshot from the given resolved environment map and layer names.
func New(env map[string]string, layers []string) *Snapshot {
	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}
	return &Snapshot{
		CreatedAt: time.Now().UTC(),
		Layers:    layers,
		Env:       copy,
	}
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write %q: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given JSON file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %q: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// Keys returns a sorted list of all environment variable keys in the snapshot.
func (s *Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.Env))
	for k := range s.Env {
		keys = append(keys, k)
	}
	return keys
}
