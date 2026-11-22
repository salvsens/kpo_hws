package domain_test

import (
	"errors"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

func TestNewStudent(t *testing.T) {
	t.Parallel()

	cases := []struct {
		label   string
		id      string
		name    string
		wantErr error
	}{
		{label: "full name", id: "S-001", name: "Анна Смирнова"},
		{label: "empty ID", name: "Анна Смирнова", wantErr: domain.ErrInvalidStudentID},
		{label: "blank ID", id: "  ", name: "Анна Смирнова", wantErr: domain.ErrInvalidStudentID},
		{label: "empty name", id: "S-001", wantErr: domain.ErrInvalidStudentName},
		{label: "blank name", id: "S-001", name: " \t ", wantErr: domain.ErrInvalidStudentName},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			t.Parallel()

			stud, err := domain.NewStudent(tc.id, tc.name)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewStudent(%q, %q) error = %v, want %v", tc.id, tc.name, err, tc.wantErr)
			}
			if tc.wantErr != nil {
				if stud != nil {
					t.Errorf("NewStudent(%q, %q) student = %v, want nil", tc.id, tc.name, stud)
				}
				return
			}
			if stud == nil {
				t.Fatalf("NewStudent(%q, %q) student = nil, want a student", tc.id, tc.name)
			}
			if got := stud.ID(); got != tc.id {
				t.Errorf("ID() = %q, want %q", got, tc.id)
			}
			if got := stud.Name(); got != tc.name {
				t.Errorf("Name() = %q, want %q", got, tc.name)
			}
		})
	}
}

func TestNewStudentAllowsSameName(t *testing.T) {
	t.Parallel()

	const (
		name     = "Анна Смирнова"
		firstID  = "S-001"
		secondID = "S-002"
	)

	first, err := domain.NewStudent(firstID, name)
	if err != nil || first == nil {
		t.Fatalf("NewStudent(%q, %q) = (%v, %v), want a student and nil error", firstID, name, first, err)
	}
	second, err := domain.NewStudent(secondID, name)
	if err != nil || second == nil {
		t.Fatalf("NewStudent(%q, %q) = (%v, %v), want a student and nil error", secondID, name, second, err)
	}

	if got := first.ID(); got != firstID {
		t.Errorf("first.ID() = %q, want %q", got, firstID)
	}
	if got := second.ID(); got != secondID {
		t.Errorf("second.ID() = %q, want %q", got, secondID)
	}
	if got := second.Name(); got != name {
		t.Errorf("second.Name() = %q, want %q", got, name)
	}
}
