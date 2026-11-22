package app_test

import (
	"errors"
	"slices"
	"sync"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
	"github.com/salvsens/kpo_hws/hw2/internal/domain"
	"github.com/salvsens/kpo_hws/hw2/internal/memory"
)

func TestGradebookScenarios(t *testing.T) {
	t.Parallel()

	gradebook := app.NewGradebook(memory.NewRepository())
	student := addStudent(t, gradebook, "Анна Смирнова")
	second := addStudent(t, gradebook, "Анна Смирнова")
	practice := addCriterion(t, gradebook, "Практика", 1000, 40)
	project := addCriterion(t, gradebook, "Проект", 2000, 60)

	if student.ID != "S-001" || second.ID != "S-002" || practice.ID != "C-001" || project.ID != "C-002" {
		t.Fatalf("student IDs = (%q, %q), criterion IDs = (%q, %q), want (S-001, S-002) and (C-001, C-002)",
			student.ID, second.ID, practice.ID, project.ID)
	}
	students, err := gradebook.ListStudents()
	if err != nil || !slices.Equal(students, []app.StudentSummary{student, second}) {
		t.Fatalf("ListStudents() = (%v, %v), want ([%v %v], nil)", students, err, student, second)
	}
	criteria, err := gradebook.ListCriteria()
	if err != nil || !slices.Equal(criteria, []app.CriterionSummary{practice, project}) {
		t.Fatalf("ListCriteria() = (%v, %v), want ([%v %v], nil)", criteria, err, practice, project)
	}

	setGrade(t, gradebook, student.ID, project, 1500)
	setGrade(t, gradebook, student.ID, practice, 800)
	assertGrades(t, gradebook, student.ID, []app.GradeSummary{
		{
			Criterion:       practice,
			ScoreHundredths: 800,
		},
		{
			Criterion:       project,
			ScoreHundredths: 1500,
		},
	})
	assertGrades(t, gradebook, second.ID, nil)
	assertResult(t, gradebook, student, app.FinalMode, 770)

	if _, err := gradebook.EditGrade(student.ID, practice.ID, 1200); !errors.Is(err, domain.ErrScoreOutOfRange) {
		t.Errorf("EditGrade(%q, %q, 1200) error = %v, want %v", student.ID, practice.ID, err, domain.ErrScoreOutOfRange)
	}
	assertResult(t, gradebook, student, app.FinalMode, 770)

	if got, err := gradebook.EditGrade(student.ID, project.ID, 2000); err != nil || got != project {
		t.Fatalf("EditGrade(%q, %q, 2000) = (%v, %v), want (%v, nil)", student.ID, project.ID, got, err, project)
	}
	assertResult(t, gradebook, student, app.FinalMode, 920)

	if got, err := gradebook.CancelGrade(student.ID, project.ID); err != nil || got != project {
		t.Fatalf("CancelGrade(%q, %q) = (%v, %v), want (%v, nil)", student.ID, project.ID, got, err, project)
	}
	assertResult(t, gradebook, student, app.CurrentMode, 800)
	if _, err := gradebook.CalculateResult(student.ID, app.FinalMode); !errors.Is(err, domain.ErrGradeNotFound) {
		t.Errorf("CalculateResult(%q, final) error = %v, want %v", student.ID, err, domain.ErrGradeNotFound)
	}
	setGrade(t, gradebook, student.ID, project, 2000)
	assertResult(t, gradebook, student, app.FinalMode, 920)
}

func TestGradebookReportsUnknownIDsAndModes(t *testing.T) {
	t.Parallel()

	gradebook := app.NewGradebook(memory.NewRepository())
	student := addStudent(t, gradebook, "Анна Смирнова")
	practice := addCriterion(t, gradebook, "Практика", 1000, 100)
	const unknownStudent = "S-002"
	const unknownCriterion = "C-002"

	if _, err := gradebook.SetGrade(unknownStudent, practice.ID, 800); !errors.Is(err, app.ErrStudentNotFound) {
		t.Errorf("SetGrade(%q, %q, 800) error = %v, want %v", unknownStudent, practice.ID, err, app.ErrStudentNotFound)
	}
	if _, err := gradebook.ListGrades(unknownStudent); !errors.Is(err, app.ErrStudentNotFound) {
		t.Errorf("ListGrades(%q) error = %v, want %v", unknownStudent, err, app.ErrStudentNotFound)
	}
	if _, err := gradebook.CalculateResult(unknownStudent, app.CurrentMode); !errors.Is(err, app.ErrStudentNotFound) {
		t.Errorf("CalculateResult(%q, current) error = %v, want %v", unknownStudent, err, app.ErrStudentNotFound)
	}
	if _, err := gradebook.SetGrade(student.ID, unknownCriterion, 800); !errors.Is(err, domain.ErrCriterionNotFound) {
		t.Errorf("SetGrade(%q, %q, 800) error = %v, want %v", student.ID, unknownCriterion, err, domain.ErrCriterionNotFound)
	}
	if _, err := gradebook.CalculateResult(student.ID, "other"); !errors.Is(err, app.ErrInvalidResultMode) {
		t.Errorf("CalculateResult(%q, other) error = %v, want %v", student.ID, err, app.ErrInvalidResultMode)
	}
	assertGrades(t, gradebook, student.ID, nil)
}

func TestGradebookKeepsIDsAfterRejectedAdds(t *testing.T) {
	t.Parallel()

	gradebook := app.NewGradebook(memory.NewRepository())
	if _, err := gradebook.AddStudent(" "); !errors.Is(err, domain.ErrInvalidStudentName) {
		t.Errorf("AddStudent(blank name) error = %v, want %v", err, domain.ErrInvalidStudentName)
	}
	student := addStudent(t, gradebook, "Анна Смирнова")
	if student.ID != "S-001" {
		t.Errorf("AddStudent() ID = %q, want S-001", student.ID)
	}

	addCriterion(t, gradebook, "Практика", 1000, 40)
	if _, err := gradebook.AddCriterion("Проект", 2000, 61); !errors.Is(err, domain.ErrTotalWeightExceeded) {
		t.Errorf("AddCriterion(Проект, 2000, 61) error = %v, want %v", err, domain.ErrTotalWeightExceeded)
	}
	project := addCriterion(t, gradebook, "Проект", 2000, 60)
	if project.ID != "C-002" {
		t.Errorf("AddCriterion() ID = %q, want C-002", project.ID)
	}
}

func TestGradebookKeepsStoredGradeWhenSaveFails(t *testing.T) {
	t.Parallel()

	repository := &saveStudentRepository{Repository: memory.NewRepository()}
	gradebook := app.NewGradebook(repository)
	student := addStudent(t, gradebook, "Анна Смирнова")
	practice := addCriterion(t, gradebook, "Практика", 1000, 100)
	setGrade(t, gradebook, student.ID, practice, 800)
	repository.saveErr = errors.New("write failed")

	if _, err := gradebook.EditGrade(student.ID, practice.ID, 900); !errors.Is(err, repository.saveErr) {
		t.Errorf("EditGrade(%q, %q, 900) error = %v, want %v", student.ID, practice.ID, err, repository.saveErr)
	}
	assertResult(t, gradebook, student, app.FinalMode, 800)
}

func TestGradebookKeepsConcurrentGradeUpdates(t *testing.T) {
	t.Parallel()

	gradebook := app.NewGradebook(memory.NewRepository())
	student := addStudent(t, gradebook, "Анна Смирнова")
	practice := addCriterion(t, gradebook, "Практика", 1000, 40)
	project := addCriterion(t, gradebook, "Проект", 2000, 60)

	var updates sync.WaitGroup
	for _, criterion := range []app.CriterionSummary{practice, project} {
		updates.Go(func() {
			if _, err := gradebook.SetGrade(student.ID, criterion.ID, criterion.MaxScoreHundredths); err != nil {
				t.Errorf("SetGrade(%q, %q, %d) error = %v, want nil", student.ID, criterion.ID, criterion.MaxScoreHundredths, err)
			}
		})
	}
	updates.Wait()

	assertResult(t, gradebook, student, app.FinalMode, 1000)
}

type saveStudentRepository struct {
	*memory.Repository
	saveErr error
}

func (r *saveStudentRepository) SaveStudent(student *domain.Student) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	return r.Repository.SaveStudent(student)
}

func addStudent(t *testing.T, gradebook *app.Gradebook, name string) app.StudentSummary {
	t.Helper()

	student, err := gradebook.AddStudent(name)
	if err != nil {
		t.Fatalf("AddStudent(%q) error = %v, want nil", name, err)
	}
	return student
}

func addCriterion(t *testing.T, gradebook *app.Gradebook, name string, maximum int64, weight int) app.CriterionSummary {
	t.Helper()

	criterion, err := gradebook.AddCriterion(name, maximum, weight)
	if err != nil {
		t.Fatalf("AddCriterion(%q, %d, %d) error = %v, want nil", name, maximum, weight, err)
	}
	return criterion
}

func setGrade(t *testing.T, gradebook *app.Gradebook, studentID string, criterion app.CriterionSummary, score int64) {
	t.Helper()

	got, err := gradebook.SetGrade(studentID, criterion.ID, score)
	if err != nil || got != criterion {
		t.Fatalf("SetGrade(%q, %q, %d) = (%v, %v), want (%v, nil)", studentID, criterion.ID, score, got, err, criterion)
	}
}

func assertGrades(t *testing.T, gradebook *app.Gradebook, studentID string, want []app.GradeSummary) {
	t.Helper()

	got, err := gradebook.ListGrades(studentID)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("ListGrades(%q) = (%v, %v), want (%v, nil)", studentID, got, err, want)
	}
}

func assertResult(t *testing.T, gradebook *app.Gradebook, student app.StudentSummary, mode app.ResultMode, score int64) {
	t.Helper()

	want := app.Result{
		Student:         student,
		Mode:            mode,
		ScoreHundredths: score,
	}
	got, err := gradebook.CalculateResult(student.ID, mode)
	if err != nil || got != want {
		t.Errorf("CalculateResult(%q, %q) = (%v, %v), want (%v, nil)", student.ID, mode, got, err, want)
	}
}
