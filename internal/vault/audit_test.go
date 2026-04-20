package vault

import (
	"errors"
	"testing"
	"time"
)

func TestAuditEntry_String_OK(t *testing.T) {
	e := AuditEntry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Operation: "READ",
		Mount:     "secret",
		Path:      "myapp/config",
	}
	got := e.String()
	want := "[2024-01-15T10:00:00Z] READ secret/myapp/config OK"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestAuditEntry_String_WithError(t *testing.T) {
	e := AuditEntry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Operation: "DELETE",
		Mount:     "secret",
		Path:      "myapp/config",
		Error:     "permission denied",
	}
	got := e.String()
	want := "[2024-01-15T10:00:00Z] DELETE secret/myapp/config ERROR: permission denied"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestAuditLog_Record_NoError(t *testing.T) {
	log := NewAuditLog()
	log.Record("READ", "secret", "myapp/db", nil)

	entries := log.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Operation != "READ" || e.Mount != "secret" || e.Path != "myapp/db" || e.Error != "" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestAuditLog_Record_WithError(t *testing.T) {
	log := NewAuditLog()
	log.Record("WRITE", "kv", "prod/creds", errors.New("forbidden"))

	entries := log.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Error != "forbidden" {
		t.Errorf("expected error 'forbidden', got %q", entries[0].Error)
	}
}

func TestAuditLog_Entries_IsCopy(t *testing.T) {
	log := NewAuditLog()
	log.Record("READ", "secret", "path/one", nil)

	a := log.Entries()
	a[0].Operation = "MUTATED"

	b := log.Entries()
	if b[0].Operation == "MUTATED" {
		t.Error("Entries() should return a copy, not a reference")
	}
}

func TestAuditLog_MultipleEntries(t *testing.T) {
	log := NewAuditLog()
	log.Record("READ", "secret", "a", nil)
	log.Record("WRITE", "secret", "b", nil)
	log.Record("DELETE", "secret", "c", errors.New("not found"))

	if got := len(log.Entries()); got != 3 {
		t.Errorf("expected 3 entries, got %d", got)
	}
}
