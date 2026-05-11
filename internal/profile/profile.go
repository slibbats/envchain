// Package profile manages named environment profiles (dev, staging, prod)
// and the mapping from profile names to their ordered layer file paths.
package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
)

var validName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Profile holds a named set of env-file layers in priority order (lowest first).
type Profile struct {
	Name   string   `json:"name"`
	Layers []string `json:"layers"`
}

// Config is the top-level profiles configuration.
type Config struct {
	Profiles []Profile `json:"profiles"`
}

// Validate returns an error if the profile name or layer list is invalid.
func (p *Profile) Validate() error {
	if p.Name == "" {
		return errors.New("profile name must not be empty")
	}
	if !validName.MatchString(p.Name) {
		return fmt.Errorf("profile name %q contains invalid characters", p.Name)
	}
	if len(p.Layers) == 0 {
		return fmt.Errorf("profile %q must have at least one layer", p.Name)
	}
	return nil
}

// LoadConfig reads and parses a JSON profiles config file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("profile: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("profile: parse config: %w", err)
	}
	for i := range cfg.Profiles {
		if err := cfg.Profiles[i].Validate(); err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

// Get returns the named profile or an error if not found.
func (c *Config) Get(name string) (*Profile, error) {
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			return &c.Profiles[i], nil
		}
	}
	return nil, fmt.Errorf("profile: %q not found", name)
}

// Names returns the names of all configured profiles.
func (c *Config) Names() []string {
	names := make([]string, len(c.Profiles))
	for i, p := range c.Profiles {
		names[i] = p.Name
	}
	return names
}
