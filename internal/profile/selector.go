package profile

import (
	"fmt"
	"os"
)

const envProfileKey = "ENVCHAIN_PROFILE"

// Selector resolves which profile to activate.
type Selector struct {
	config    *Config
	envLookup func(string) (string, bool)
}

// NewSelector creates a Selector backed by the given Config.
func NewSelector(cfg *Config) *Selector {
	return &Selector{config: cfg, envLookup: os.LookupEnv}
}

// Select returns the Profile for the given name.
// If name is empty, it falls back to the ENVCHAIN_PROFILE environment variable.
// Returns an error if no profile can be determined or the name is not found.
func (s *Selector) Select(name string) (*Profile, error) {
	if name == "" {
		var ok bool
		name, ok = s.envLookup(envProfileKey)
		if !ok || name == "" {
			return nil, fmt.Errorf("profile: no profile specified and %s is not set", envProfileKey)
		}
	}
	return s.config.Get(name)
}

// ActiveName returns the name that would be selected without resolving the full profile.
func (s *Selector) ActiveName(name string) string {
	if name != "" {
		return name
	}
	if v, ok := s.envLookup(envProfileKey); ok {
		return v
	}
	return ""
}
