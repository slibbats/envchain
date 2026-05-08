package audit

import (
	"fmt"

	"github.com/yourorg/envchain/internal/chain"
)

// AuditedChain wraps a *chain.Chain and records an audit event on every Get.
type AuditedChain struct {
	c      *chain.Chain
	logger *Logger
}

// Wrap returns an AuditedChain that delegates to c and logs via logger.
func Wrap(c *chain.Chain, logger *Logger) *AuditedChain {
	if logger == nil {
		logger = Discard()
	}
	return &AuditedChain{c: c, logger: logger}
}

// Get retrieves key from the underlying chain and records the outcome.
func (a *AuditedChain) Get(key string) (string, bool) {
	val, ok := a.c.Get(key)
	if ok {
		a.logger.Record(EventResolved, key, "", fmt.Sprintf("value length=%d", len(val)))
	} else {
		a.logger.Record(EventMissing, key, "", "key not found in any layer")
	}
	return val, ok
}

// Resolve returns all resolved key-value pairs and logs an event per key.
func (a *AuditedChain) Resolve() map[string]string {
	env := a.c.Resolve()
	for k, v := range env {
		a.logger.Record(EventResolved, k, "", fmt.Sprintf("value length=%d", len(v)))
	}
	return env
}
