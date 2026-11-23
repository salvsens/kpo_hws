package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunStartsConsoleWithAnEmptyJournal(t *testing.T) {
	t.Parallel()

	var first bytes.Buffer
	if err := run(strings.NewReader("add-student Анна Смирнова\nstudents\nexit\n"), &first); err != nil {
		t.Fatalf("first run() error = %v, want nil", err)
	}
	if !strings.Contains(first.String(), "S-001  Анна Смирнова\n") {
		t.Errorf("first run() output = %q, want student S-001 in list", first.String())
	}

	var second bytes.Buffer
	if err := run(strings.NewReader("students"), &second); err != nil {
		t.Fatalf("second run() error = %v, want nil", err)
	}
	if !strings.Contains(second.String(), "Студентов пока нет.") {
		t.Errorf("second run() output = %q, want an empty journal", second.String())
	}
}
