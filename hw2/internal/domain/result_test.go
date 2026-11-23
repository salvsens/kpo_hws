package domain_test

import (
	"errors"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

func TestCalculateResult(t *testing.T) {
	t.Parallel()

	rubric := newResultTestRubric(t)
	cases := []struct {
		name          string
		strategy      domain.ResultStrategy
		grades        map[string]int64
		want          int64
		wantErr       error
		wantMissingID string
	}{
		{
			name:     "current with both grades",
			strategy: domain.CurrentResult{},
			grades: map[string]int64{
				practiceCriterionID: 800,
				projectCriterionID:  1500,
			},
			want: 770,
		},
		{
			name:     "current with one missing grade",
			strategy: domain.CurrentResult{},
			grades:   map[string]int64{practiceCriterionID: 800},
			want:     800,
		},
		{
			name:     "current with an explicit zero",
			strategy: domain.CurrentResult{},
			grades: map[string]int64{
				practiceCriterionID: 800,
				projectCriterionID:  0,
			},
			want: 320,
		},
		{
			name:     "current with a single zero",
			strategy: domain.CurrentResult{},
			grades:   map[string]int64{practiceCriterionID: 0},
			want:     0,
		},
		{
			name:     "current without grades",
			strategy: domain.CurrentResult{},
			wantErr:  domain.ErrNoGrades,
		},
		{
			name:     "final with both grades",
			strategy: domain.FinalResult{},
			grades: map[string]int64{
				practiceCriterionID: 800,
				projectCriterionID:  1500,
			},
			want: 770,
		},
		{
			name:     "final with maximum grades",
			strategy: domain.FinalResult{},
			grades: map[string]int64{
				practiceCriterionID: 1000,
				projectCriterionID:  2000,
			},
			want: 1000,
		},
		{
			name:          "final with one missing grade",
			strategy:      domain.FinalResult{},
			grades:        map[string]int64{practiceCriterionID: 800},
			wantErr:       domain.ErrGradeNotFound,
			wantMissingID: projectCriterionID,
		},
		{
			name:          "final without grades",
			strategy:      domain.FinalResult{},
			wantErr:       domain.ErrGradeNotFound,
			wantMissingID: practiceCriterionID,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			student := newGradeTestStudent(t)
			for id, score := range tc.grades {
				setResultTestGrade(t, student, rubric, id, score)
			}

			got, err := domain.CalculateResult(student, rubric, tc.strategy)

			if got != tc.want || !errors.Is(err, tc.wantErr) {
				t.Fatalf("CalculateResult() = (%d, %v), want (%d, %v)", got, err, tc.want, tc.wantErr)
			}
			if tc.wantMissingID != "" {
				missing, ok := errors.AsType[domain.MissingGradeError](err)
				if !ok || missing.CriterionID != tc.wantMissingID {
					t.Errorf("CalculateResult() error = %v, want MissingGradeError for %q", err, tc.wantMissingID)
				}
			}
		})
	}
}

func TestCalculateResultWithIncompleteRubric(t *testing.T) {
	t.Parallel()

	var rubric domain.Rubric
	student := newGradeTestStudent(t)
	checkResult(t, student, rubric, domain.CurrentResult{}, 0, domain.ErrNoGrades)
	checkResult(t, student, rubric, domain.FinalResult{}, 0, domain.ErrIncompleteRubric)

	addResultTestCriterion(t, &rubric, practiceCriterionID, 1000, 40)
	setResultTestGrade(t, student, rubric, practiceCriterionID, 800)

	checkResult(t, student, rubric, domain.CurrentResult{}, 800, nil)
	checkResult(t, student, rubric, domain.FinalResult{}, 0, domain.ErrIncompleteRubric)
}

func TestCalculateResultRounding(t *testing.T) {
	t.Parallel()

	var rubric domain.Rubric
	addResultTestCriterion(t, &rubric, practiceCriterionID, 10_000, 100)
	cases := []struct {
		name  string
		score int64
		want  int64
	}{
		{
			name:  "below half",
			score: 624,
			want:  62,
		},
		{
			name:  "at half",
			score: 625,
			want:  63,
		},
		{
			name:  "above half",
			score: 626,
			want:  63,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			student := newGradeTestStudent(t)
			setResultTestGrade(t, student, rubric, practiceCriterionID, tc.score)

			checkResult(t, student, rubric, domain.CurrentResult{}, tc.want, nil)
			checkResult(t, student, rubric, domain.FinalResult{}, tc.want, nil)
		})
	}
}

func TestCalculateResultRoundsOnlyTheTotal(t *testing.T) {
	t.Parallel()

	var rubric domain.Rubric
	addResultTestCriterion(t, &rubric, practiceCriterionID, 1600, 50)
	addResultTestCriterion(t, &rubric, projectCriterionID, 1600, 50)
	student := newGradeTestStudent(t)
	setResultTestGrade(t, student, rubric, practiceCriterionID, 100)
	setResultTestGrade(t, student, rubric, projectCriterionID, 100)

	checkResult(t, student, rubric, domain.FinalResult{}, 63, nil)
}

func TestCalculateResultAfterGradeChanges(t *testing.T) {
	t.Parallel()

	rubric := newResultTestRubric(t)
	student := newGradeTestStudent(t)
	setResultTestGrade(t, student, rubric, practiceCriterionID, 800)
	setResultTestGrade(t, student, rubric, projectCriterionID, 1500)
	checkResult(t, student, rubric, domain.FinalResult{}, 770, nil)

	if err := student.EditGrade(rubric, projectCriterionID, 2000); err != nil {
		t.Fatalf("EditGrade(rubric, %q, 2000) error = %v, want nil", projectCriterionID, err)
	}
	checkResult(t, student, rubric, domain.FinalResult{}, 920, nil)

	if err := student.CancelGrade(projectCriterionID); err != nil {
		t.Fatalf("CancelGrade(%q) error = %v, want nil", projectCriterionID, err)
	}
	checkResult(t, student, rubric, domain.CurrentResult{}, 800, nil)
	checkResult(t, student, rubric, domain.FinalResult{}, 0, domain.ErrGradeNotFound)

	setResultTestGrade(t, student, rubric, projectCriterionID, 2000)
	checkResult(t, student, rubric, domain.FinalResult{}, 920, nil)
}

func newResultTestRubric(t *testing.T) domain.Rubric {
	t.Helper()

	var rubric domain.Rubric
	addResultTestCriterion(t, &rubric, practiceCriterionID, 1000, 40)
	addResultTestCriterion(t, &rubric, projectCriterionID, 2000, 60)
	return rubric
}

func addResultTestCriterion(t *testing.T, rubric *domain.Rubric, id string, maxScore int64, percent int) {
	t.Helper()

	weight, err := domain.NewWeight(percent)
	if err != nil {
		t.Fatalf("NewWeight(%d) error = %v, want nil", percent, err)
	}
	criterion, err := domain.NewCriterion(id, id, maxScore, weight)
	if err != nil {
		t.Fatalf("NewCriterion(%q, %q, %d, %v) error = %v, want nil", id, id, maxScore, weight, err)
	}
	if err := rubric.AddCriterion(criterion); err != nil {
		t.Fatalf("AddCriterion(%q) error = %v, want nil", id, err)
	}
}

func setResultTestGrade(t *testing.T, student *domain.Student, rubric domain.Rubric, id string, score int64) {
	t.Helper()

	if err := student.SetGrade(rubric, id, score); err != nil {
		t.Fatalf("SetGrade(rubric, %q, %d) error = %v, want nil", id, score, err)
	}
}

func checkResult(t *testing.T, student *domain.Student, rubric domain.Rubric, strategy domain.ResultStrategy, want int64, wantErr error) {
	t.Helper()

	got, err := domain.CalculateResult(student, rubric, strategy)

	if got != want || !errors.Is(err, wantErr) {
		t.Errorf("CalculateResult(..., %T) = (%d, %v), want (%d, %v)", strategy, got, err, want, wantErr)
	}
}
