package chain

import (
	"errors"
	"fmt"
)

// Layer represents a named environment layer (e.g., "dev", "staging", "prod").
type Layer struct {
	Name   string
	Values map[string]string
}

// Chain holds an ordered list of layers, resolved from lowest to highest priority.
type Chain struct {
	layers []*Layer
}

// New creates an empty Chain.
func New() *Chain {
	return &Chain{}
}

// AddLayer appends a layer to the chain. Layers added later take higher priority.
func (c *Chain) AddLayer(name string, values map[string]string) error {
	if name == "" {
		return errors.New("layer name must not be empty")
	}
	for _, l := range c.layers {
		if l.Name == name {
			return fmt.Errorf("layer %q already exists in chain", name)
		}
	}
	c.layers = append(c.layers, &Layer{Name: name, Values: values})
	return nil
}

// Resolve returns the final merged environment map, with later layers
// overriding earlier ones. No secret values are ever logged here.
func (c *Chain) Resolve() map[string]string {
	result := make(map[string]string)
	for _, layer := range c.layers {
		for k, v := range layer.Values {
			result[k] = v
		}
	}
	return result
}

// Get returns the resolved value for a key and whether it was found.
func (c *Chain) Get(key string) (string, bool) {
	resolved := c.Resolve()
	v, ok := resolved[key]
	return v, ok
}

// LayerNames returns the names of all layers in priority order (lowest first).
func (c *Chain) LayerNames() []string {
	names := make([]string, len(c.layers))
	for i, l := range c.layers {
		names[i] = l.Name
	}
	return names
}
