// Package audit provides a simple event log for tracking
// environment variable access and resolution events within envchain.
package audit

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// EventKind classifies the type of audit event.
type EventKind string

const (
	EventResolved  EventKind = "RESOLVED"
	EventMissing   EventKind = "MISSING"
	EventOverridden EventKind = "OVERRIDDEN"
	EventMasked    EventKind = "MASKED"
)

// Event represents a single audit log entry.
type Event struct {
	Time      time.Time
	Kind      EventKind
	Key       string
	Layer     string
	Message   string
}

// Logger records audit events to a writer.
type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

// New creates a Logger writing to out. Pass nil to use os.Stderr.
func New(out io.Writer) *Logger {
	if out == nil {
		out = os.Stderr
	}
	return &Logger{out: out}
}

// Record appends an event to the audit log.
func (l *Logger) Record(kind EventKind, key, layer, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := Event{
		Time:    time.Now().UTC(),
		Kind:    kind,
		Key:     key,
		Layer:   layer,
		Message: message,
	}
	fmt.Fprintf(l.out, "%s [%s] key=%q layer=%q %s\n",
		e.Time.Format(time.RFC3339), e.Kind, e.Key, e.Layer, e.Message)
}

// Discard returns a Logger that silently drops all events.
func Discard() *Logger {
	return &Logger{out: io.Discard}
}
