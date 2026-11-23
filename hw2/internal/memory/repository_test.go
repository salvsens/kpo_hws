package memory_test

import (
	"errors"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
	"github.com/salvsens/kpo_hws/hw2/internal/domain"
	"github.com/salvsens/kpo_hws/hw2/internal/memory"
)

func TestRepositoryCopiesStudents(t *testing.T) {
	t.Parallel()

	repository := memory.NewRepository()
	student := newTestStudent(t, "S-001", "Анна Смирнова")
	criterion := newTestCriterion(t, "C-001", 100)
	var rubric domain.Rubric
	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("AddCriterion(%q) error = %v, want nil", criterion.ID(), err)
	}
	const initialScore = 800
	if err := student.SetGrade(rubric, criterion.ID(), initialScore); err != nil {
		t.Fatalf("SetGrade(%q, %d) error = %v, want nil", criterion.ID(), initialScore, err)
	}
	if err := repository.CreateStudent(student); err != nil {
		t.Fatalf("CreateStudent(%q) error = %v, want nil", student.ID(), err)
	}

	editTestGrade(t, student, rubric, criterion.ID(), 850)
	assertStoredGrade(t, repository, student.ID(), criterion.ID(), initialScore)

	loaded, err := repository.FindStudent(student.ID())
	if err != nil {
		t.Fatalf("FindStudent(%q) error = %v, want nil", student.ID(), err)
	}
	const updatedScore = 900
	editTestGrade(t, loaded, rubric, criterion.ID(), updatedScore)
	assertStoredGrade(t, repository, student.ID(), criterion.ID(), initialScore)

	if err := repository.SaveStudent(loaded); err != nil {
		t.Fatalf("SaveStudent(%q) error = %v, want nil", student.ID(), err)
	}
	editTestGrade(t, loaded, rubric, criterion.ID(), 950)
	assertStoredGrade(t, repository, student.ID(), criterion.ID(), updatedScore)

	listed, err := repository.ListStudents()
	if err != nil || len(listed) != 1 {
		t.Fatalf("ListStudents() = (%v, %v), want one student and nil error", listed, err)
	}
	editTestGrade(t, listed[0], rubric, criterion.ID(), initialScore)
	assertStoredGrade(t, repository, student.ID(), criterion.ID(), updatedScore)
}

func TestRepositoryStudentErrorsAndOrder(t *testing.T) {
	t.Parallel()

	repository := memory.NewRepository()
	first := newTestStudent(t, "S-002", "Анна Смирнова")
	second := newTestStudent(t, "S-001", "Иван Петров")
	duplicate := newTestStudent(t, first.ID(), "Другой студент")

	if got, err := repository.FindStudent(first.ID()); got != nil || !errors.Is(err, app.ErrStudentNotFound) {
		t.Fatalf("FindStudent(%q) = (%v, %v), want (nil, %v)", first.ID(), got, err, app.ErrStudentNotFound)
	}
	if err := repository.SaveStudent(first); !errors.Is(err, app.ErrStudentNotFound) {
		t.Fatalf("SaveStudent(%q) error = %v, want %v", first.ID(), err, app.ErrStudentNotFound)
	}
	for _, student := range []*domain.Student{first, second} {
		if err := repository.CreateStudent(student); err != nil {
			t.Fatalf("CreateStudent(%q) error = %v, want nil", student.ID(), err)
		}
	}
	if err := repository.CreateStudent(duplicate); !errors.Is(err, app.ErrStudentAlreadyExists) {
		t.Errorf("CreateStudent(%q) error = %v, want %v", duplicate.ID(), err, app.ErrStudentAlreadyExists)
	}
	if err := repository.SaveStudent(first); err != nil {
		t.Fatalf("SaveStudent(%q) error = %v, want nil", first.ID(), err)
	}

	students, err := repository.ListStudents()
	if err != nil || len(students) != 2 {
		t.Fatalf("ListStudents() = (%v, %v), want two students and nil error", students, err)
	}
	if students[0].ID() != first.ID() || students[0].Name() != first.Name() || students[1].ID() != second.ID() {
		t.Errorf("ListStudents() = (%v, %v), want (%v, %v)", students[0], students[1], first, second)
	}
}

func TestRepositoryCopiesRubric(t *testing.T) {
	t.Parallel()

	repository := memory.NewRepository()
	rubric, err := repository.LoadRubric()
	if err != nil || len(rubric.Criteria()) != 0 {
		t.Fatalf("LoadRubric() = (%v, %v), want an empty rubric and nil error", rubric, err)
	}
	first := newTestCriterion(t, "C-001", 40)
	second := newTestCriterion(t, "C-002", 60)
	if err := rubric.AddCriterion(first); err != nil {
		t.Fatalf("AddCriterion(%q) error = %v, want nil", first.ID(), err)
	}
	if err := repository.SaveRubric(rubric); err != nil {
		t.Fatalf("SaveRubric() error = %v, want nil", err)
	}

	loaded, err := repository.LoadRubric()
	if err != nil {
		t.Fatalf("LoadRubric() error = %v, want nil", err)
	}
	if err := loaded.AddCriterion(second); err != nil {
		t.Fatalf("AddCriterion(%q) error = %v, want nil", second.ID(), err)
	}

	stored, err := repository.LoadRubric()
	if err != nil {
		t.Fatalf("LoadRubric() error = %v, want nil", err)
	}
	if got := stored.Criteria(); len(got) != 1 || got[0] != first {
		t.Errorf("stored Criteria() = %v, want [%v]", got, first)
	}
}

func newTestStudent(t *testing.T, id, name string) *domain.Student {
	t.Helper()

	student, err := domain.NewStudent(id, name)
	if err != nil {
		t.Fatalf("NewStudent(%q, %q) error = %v, want nil", id, name, err)
	}
	return student
}

func newTestCriterion(t *testing.T, id string, percent int) domain.Criterion {
	t.Helper()

	weight, err := domain.NewWeight(percent)
	if err != nil {
		t.Fatalf("NewWeight(%d) error = %v, want nil", percent, err)
	}
	criterion, err := domain.NewCriterion(id, id, 1000, weight)
	if err != nil {
		t.Fatalf("NewCriterion(%q, %q, 1000, %v) error = %v, want nil", id, id, weight, err)
	}
	return criterion
}

func editTestGrade(t *testing.T, student *domain.Student, rubric domain.Rubric, criterionID string, score int64) {
	t.Helper()

	if err := student.EditGrade(rubric, criterionID, score); err != nil {
		t.Fatalf("EditGrade(%q, %d) error = %v, want nil", criterionID, score, err)
	}
}

func assertStoredGrade(t *testing.T, repository *memory.Repository, studentID, criterionID string, want int64) {
	t.Helper()

	student, err := repository.FindStudent(studentID)
	if err != nil {
		t.Fatalf("FindStudent(%q) error = %v, want nil", studentID, err)
	}
	grade, ok := student.FindGrade(criterionID)
	if !ok || grade.ScoreHundredths() != want {
		t.Errorf("FindGrade(%q) = (%v, %t), want score %d and true", criterionID, grade, ok, want)
	}
}
