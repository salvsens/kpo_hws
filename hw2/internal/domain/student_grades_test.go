package domain_test

import (
	"errors"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

func TestStudentSetGrade(t *testing.T) {
	t.Parallel()

	rubric := newGradeTestRubric(t)
	student := newGradeTestStudent(t)
	checkMissingGrade(t, student, practiceCriterionID)

	if err := student.SetGrade(rubric, practiceCriterionID, 0); err != nil {
		t.Fatalf("SetGrade(rubric, %q, 0) error = %v, want nil", practiceCriterionID, err)
	}
	checkGrade(t, student, practiceCriterionID, 0)

	err := student.SetGrade(rubric, practiceCriterionID, 800)

	if !errors.Is(err, domain.ErrGradeAlreadySet) {
		t.Errorf("SetGrade(rubric, %q, 800) error = %v, want %v", practiceCriterionID, err, domain.ErrGradeAlreadySet)
	}
	checkGrade(t, student, practiceCriterionID, 0)
}

func TestStudentSetGradeRejectsInvalidEvaluation(t *testing.T) {
	t.Parallel()

	rubric := newGradeTestRubric(t)
	cases := []struct {
		name        string
		criterionID string
		score       int64
		wantErr     error
	}{
		{
			name:        "score above maximum",
			criterionID: practiceCriterionID,
			score:       1200,
			wantErr:     domain.ErrScoreOutOfRange,
		},
		{
			name:        "unknown criterion",
			criterionID: "C-003",
			score:       800,
			wantErr:     domain.ErrCriterionNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			student := newGradeTestStudent(t)

			err := student.SetGrade(rubric, tc.criterionID, tc.score)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("SetGrade(rubric, %q, %d) error = %v, want %v", tc.criterionID, tc.score, err, tc.wantErr)
			}
			checkMissingGrade(t, student, tc.criterionID)
		})
	}
}

func TestStudentSetGradeRejectsZeroStudent(t *testing.T) {
	t.Parallel()

	rubric := newGradeTestRubric(t)
	var student domain.Student

	err := student.SetGrade(rubric, practiceCriterionID, 800)

	if !errors.Is(err, domain.ErrInvalidStudentID) {
		t.Errorf("SetGrade(rubric, %q, 800) error = %v, want %v", practiceCriterionID, err, domain.ErrInvalidStudentID)
	}
	checkMissingGrade(t, &student, practiceCriterionID)
}

func TestStudentEditGrade(t *testing.T) {
	t.Parallel()

	rubric := newGradeTestRubric(t)
	student := newGradeTestStudent(t)
	if err := student.SetGrade(rubric, practiceCriterionID, 800); err != nil {
		t.Fatalf("SetGrade(rubric, %q, 800) error = %v, want nil", practiceCriterionID, err)
	}

	if err := student.EditGrade(rubric, practiceCriterionID, 1200); !errors.Is(err, domain.ErrScoreOutOfRange) {
		t.Errorf("EditGrade(rubric, %q, 1200) error = %v, want %v", practiceCriterionID, err, domain.ErrScoreOutOfRange)
	}
	checkGrade(t, student, practiceCriterionID, 800)

	if err := student.EditGrade(domain.Rubric{}, practiceCriterionID, 900); !errors.Is(err, domain.ErrCriterionNotFound) {
		t.Errorf("EditGrade(empty rubric, %q, 900) error = %v, want %v", practiceCriterionID, err, domain.ErrCriterionNotFound)
	}
	checkGrade(t, student, practiceCriterionID, 800)

	if err := student.EditGrade(rubric, practiceCriterionID, 900); err != nil {
		t.Fatalf("EditGrade(rubric, %q, 900) error = %v, want nil", practiceCriterionID, err)
	}
	checkGrade(t, student, practiceCriterionID, 900)
}

func TestStudentEditGradeRejectsMissingGrade(t *testing.T) {
	t.Parallel()

	rubric := newGradeTestRubric(t)
	student := newGradeTestStudent(t)

	err := student.EditGrade(rubric, practiceCriterionID, 900)

	if !errors.Is(err, domain.ErrGradeNotFound) {
		t.Errorf("EditGrade(rubric, %q, 900) error = %v, want %v", practiceCriterionID, err, domain.ErrGradeNotFound)
	}
	checkMissingGrade(t, student, practiceCriterionID)
}

func TestStudentCancelGrade(t *testing.T) {
	t.Parallel()

	rubric := newGradeTestRubric(t)
	student := newGradeTestStudent(t)
	if err := student.SetGrade(rubric, practiceCriterionID, 800); err != nil {
		t.Fatalf("SetGrade(rubric, %q, 800) error = %v, want nil", practiceCriterionID, err)
	}
	if err := student.SetGrade(rubric, projectCriterionID, 825); err != nil {
		t.Fatalf("SetGrade(rubric, %q, 825) error = %v, want nil", projectCriterionID, err)
	}

	if err := student.CancelGrade(practiceCriterionID); err != nil {
		t.Fatalf("CancelGrade(%q) error = %v, want nil", practiceCriterionID, err)
	}
	checkMissingGrade(t, student, practiceCriterionID)
	checkGrade(t, student, projectCriterionID, 825)

	if err := student.CancelGrade(practiceCriterionID); !errors.Is(err, domain.ErrGradeNotFound) {
		t.Errorf("CancelGrade(%q) error = %v, want %v", practiceCriterionID, err, domain.ErrGradeNotFound)
	}
	if err := student.SetGrade(rubric, practiceCriterionID, 900); err != nil {
		t.Fatalf("SetGrade(rubric, %q, 900) error = %v, want nil", practiceCriterionID, err)
	}
	checkGrade(t, student, practiceCriterionID, 900)
}

func newGradeTestStudent(t *testing.T) *domain.Student {
	t.Helper()

	const id = "S-001"
	const name = "Анна Смирнова"
	student, err := domain.NewStudent(id, name)
	if err != nil {
		t.Fatalf("NewStudent(%q, %q) error = %v, want nil", id, name, err)
	}
	return student
}

func checkGrade(t *testing.T, student *domain.Student, criterionID string, score int64) {
	t.Helper()

	evaluation, ok := student.FindGrade(criterionID)
	if !ok || evaluation.CriterionID() != criterionID || evaluation.ScoreHundredths() != score {
		t.Errorf("FindGrade(%q) = (%v, %t), want an evaluation with score %d and true", criterionID, evaluation, ok, score)
	}
}

func checkMissingGrade(t *testing.T, student *domain.Student, criterionID string) {
	t.Helper()

	if evaluation, ok := student.FindGrade(criterionID); ok || evaluation != (domain.Evaluation{}) {
		t.Errorf("FindGrade(%q) = (%v, %t), want zero evaluation and false", criterionID, evaluation, ok)
	}
}
