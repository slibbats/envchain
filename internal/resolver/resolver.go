package resolver

import (
	"fmt"
	"os"
	"strings"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/loader"
)

// Config holds the configuration for building a resolved chain.
type Config struct {
	// Layers defines ordered env files from lowest to highest priority.
	// Each entry is a named layer mapped to a file path.
	Layers []LayerConfig
	// InjectOS controls whether OS environment variables form the base layer.
	InjectOS bool
}

// LayerConfig represents a single named layer and its source file.
type LayerConfig struct {
	Name string
	FilePath string
}

// Build constructs a Chain from the given Config, loading each layer in order.
func Build(cfg Config) (*chain.Chain, error) {
	c := chain.New()

	if cfg.InjectOS {
		osVars := osEnvironment()
		if err := c.AddLayer("os", osVars); err != nil {
			return nil, fmt.Errorf("resolver: adding os layer: %w", err)
		}
	}

	for _, lc := range cfg.Layers {
		if lc.FilePath == "" {
			continue
		}
		vars, err := loader.LoadEnvFile(lc.FilePath)
		if err != nil {
			if os.IsNotExist(err) {
				// Missing optional layer files are skipped.
				continue
			}
			return nil, fmt.Errorf("resolver: loading layer %q from %q: %w", lc.Name, lc.FilePath, err)
		}
		if err := c.AddLayer(lc.Name, vars); err != nil {
			return nil, fmt.Errorf("resolver: adding layer %q: %w", lc.Name, err)
		}
	}

	return c, nil
}

// osEnvironment reads the current process environment into a map.
func osEnvironment() map[string]string {
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return env
}
