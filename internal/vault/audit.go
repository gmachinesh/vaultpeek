package vault

import (
	"context"
	"fmt"
	"time"
)

// AuditEntry represents a single recorded secret operation.
type AuditEntry struct {
	Timestamp time.Time
	Operation string
	Mount     string
	Path      string
	Error     string
}

// String returns a human-readable representation of the audit entry.
func (e AuditEntry) String() string {
	ts := e.Timestamp.Format(time.RFC3339)
	if e.Error != "" {
		return fmt.Sprintf("[%s] %s %s/%s ERROR: %s", ts, e.Operation, e.Mount, e.Path, e.Error)
	}
	return fmt.Sprintf("[%s] %s %s/%s OK", ts, e.Operation, e.Mount, e.Path)
}

// AuditLog holds an in-memory list of audit entries.
type AuditLog struct {
	entries []AuditEntry
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// Record appends a new entry to the audit log.
func (a *AuditLog) Record(op, mount, path string, err error) {
	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Mount:     mount,
		Path:      path,
	}
	if err != nil {
		entry.Error = err.Error()
	}
	a.entries = append(a.entries, entry)
}

// Entries returns a copy of all recorded audit entries.
func (a *AuditLog) Entries() []AuditEntry {
	out := make([]AuditEntry, len(a.entries))
	copy(out, a.entries)
	return out
}

// FetchAuditLog retrieves recent audit log entries from the Vault audit device.
// It returns an error if the Vault API call fails or the audit device is unavailable.
func FetchAuditLog(ctx context.Context, c *Client) ([]AuditEntry, error) {
	secret, err := c.Logical().ReadWithContext(ctx, "sys/audit")
	if err != nil {
		return nil, fmt.Errorf("reading audit devices: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, nil
	}
	var entries []AuditEntry
	for key := range secret.Data {
		entries = append(entries, AuditEntry{
			Timestamp: time.Now().UTC(),
			Operation: "audit-device",
			Mount:     key,
			Path:      "sys/audit",
		})
	}
	return entries, nil
}
