package vault

import (
	"fmt"
	"sync"
	"time"
)

// AuditEntry records a single operation performed against Vault.
type AuditEntry struct {
	Timestamp time.Time
	Operation string
	Mount     string
	Path      string
	Error     error
}

// String returns a human-readable representation of the audit entry.
func (e AuditEntry) String() string {
	status := "OK"
	if e.Error != nil {
		status = fmt.Sprintf("ERROR: %s", e.Error.Error())
	}
	return fmt.Sprintf("[%s] %s %s/%s -> %s",
		e.Timestamp.UTC().Format(time.RFC3339),
		e.Operation, e.Mount, e.Path, status)
}

// AuditLog is a thread-safe in-memory log of vault operations.
type AuditLog struct {
	mu      sync.Mutex
	entries []AuditEntry
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// Record appends an entry to the audit log.
func (l *AuditLog) Record(op, mount, path string, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, AuditEntry{
		Timestamp: time.Now(),
		Operation: op,
		Mount:     mount,
		Path:      path,
		Error:     err,
	})
}

// Entries returns a snapshot copy of all audit entries.
func (l *AuditLog) Entries() []AuditEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	copy := make([]AuditEntry, len(l.entries))
	for i, e := range l.entries {
		copy[i] = e
	}
	return copy
}

// FetchAuditLog retrieves the current entries from the log and returns them.
func FetchAuditLog(log *AuditLog) []AuditEntry {
	if log == nil {
		return nil
	}
	return log.Entries()
}
