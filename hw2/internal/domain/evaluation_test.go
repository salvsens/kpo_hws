package domain_test

import (
	"errors"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

const (
	practiceCriterionID = "C-001"
	projectCriterionID  = "C-002"
)

func TestNewEvaluation(t *testing.T) {
	t.Parallel()

	rubric := newGradeTestRubric(t)
	cases := []struct {
		name        string
		criterionID string
		score       int64
		wantErr     error
	}{
		{
			name:        "zero score",
			criterionID: projectCriterionID,
			score:       0,
		},
		{
			name:        "fractional score",
			criterionID: projectCriterionID,
			score:       825,
		},
		{
			name:        "maximum score",
			criterionID: projectCriterionID,
			score:       1050,
		},
		{
			name:        "negative score",
			criterionID: projectCriterionID,
			score:       -1,
			wantErr:     domain.ErrScoreOutOfRange,
		},
		{
			name:        "above maximum",
			criterionID: projectCriterionID,
			score:       1051,
			wantErr:     domain.ErrScoreOutOfRange,
		},
		{
			name:        "above maximum of another criterion",
			criterionID: practiceCriterionID,
			score:       1050,
			wantErr:     domain.ErrScoreOutOfRange,
		},
		{
			name:        "unknown criterion",
			criterionID: "C-003",
			score:       825,
			wantErr:     domain.ErrCriterionNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			evaluation, err := domain.NewEvaluation(rubric, tc.criterionID, tc.score)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewEvaluation(rubric, %q, %d) error = %v, want %v", tc.criterionID, tc.score, err, tc.wantErr)
			}
			if tc.wantErr != nil {
				if evaluation != (domain.Evaluation{}) {
					t.Errorf("NewEvaluation(rubric, %q, %d) evaluation = %v, want zero value", tc.criterionID, tc.score, evaluation)
				}
				return
			}
			if got := evaluation.CriterionID(); got != tc.criterionID {
				t.Errorf("CriterionID() = %q, want %q", got, tc.criterionID)
			}
			if got := evaluation.ScoreHundredths(); got != tc.score {
				t.Errorf("ScoreHundredths() = %d, want %d", got, tc.score)
			}
		})
	}
}

func newGradeTestRubric(t *testing.T) domain.Rubric {
	t.Helper()

	var rubric domain.Rubric
	practice := newTestCriterion(t, practiceCriterionID, "Практика", 40)
	if err := rubric.AddCriterion(practice); err != nil {
		t.Fatalf("AddCriterion(%q) error = %v, want nil", practice.ID(), err)
	}

	weight, err := domain.NewWeight(60)
	if err != nil {
		t.Fatalf("NewWeight(60) error = %v, want nil", err)
	}
	project, err := domain.NewCriterion(projectCriterionID, "Проект", 1050, weight)
	if err != nil {
		t.Fatalf("NewCriterion(%q, %q, 1050, %v) error = %v, want nil", projectCriterionID, "Проект", weight, err)
	}
	if err := rubric.AddCriterion(project); err != nil {
		t.Fatalf("AddCriterion(%q) error = %v, want nil", project.ID(), err)
	}
	return rubric
}
