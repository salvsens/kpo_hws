package domain_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

func TestRubricAddCriterion(t *testing.T) {
	t.Parallel()

	first := newTestCriterion(t, "C-002", "Практика", 40)
	overweight := newTestCriterion(t, "C-001", "Проект", 61)
	second := newTestCriterion(t, "C-001", "Проект", 60)
	var rubric domain.Rubric

	if err := rubric.AddCriterion(first); err != nil {
		t.Fatalf("AddCriterion(%q) error = %v, want nil", first.ID(), err)
	}
	if err := rubric.AddCriterion(overweight); !errors.Is(err, domain.ErrTotalWeightExceeded) {
		t.Fatalf("AddCriterion(%q) error = %v, want %v", overweight.ID(), err, domain.ErrTotalWeightExceeded)
	}
	if got := rubric.Criteria(); !slices.Equal(got, []domain.Criterion{first}) {
		t.Fatalf("Criteria() after rejected addition = %v, want [%v]", got, first)
	}
	if err := rubric.AddCriterion(second); err != nil {
		t.Fatalf("AddCriterion(%q) error = %v, want nil", second.ID(), err)
	}

	want := []domain.Criterion{first, second}
	if got := rubric.Criteria(); !slices.Equal(got, want) {
		t.Errorf("Criteria() = %v, want %v", got, want)
	}
}

func TestRubricAddCriterionRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	first := newTestCriterion(t, "C-001", "Практика", 40)
	duplicate := newTestCriterion(t, first.ID(), "Проект", 30)
	cases := []struct {
		name     string
		existing []domain.Criterion
		input    domain.Criterion
		wantErr  error
	}{
		{
			name:    "zero criterion",
			input:   domain.Criterion{},
			wantErr: domain.ErrInvalidCriterion,
		},
		{
			name:     "duplicate ID",
			existing: []domain.Criterion{first},
			input:    duplicate,
			wantErr:  domain.ErrDuplicateCriterionID,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var rubric domain.Rubric
			for _, criterion := range tc.existing {
				if err := rubric.AddCriterion(criterion); err != nil {
					t.Fatalf("AddCriterion(%q) error = %v, want nil", criterion.ID(), err)
				}
			}

			err := rubric.AddCriterion(tc.input)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("AddCriterion(%q) error = %v, want %v", tc.input.ID(), err, tc.wantErr)
			}
			if got := rubric.Criteria(); !slices.Equal(got, tc.existing) {
				t.Errorf("Criteria() = %v, want %v", got, tc.existing)
			}
		})
	}
}

func TestRubricCriteriaReturnsCopy(t *testing.T) {
	t.Parallel()

	original := newTestCriterion(t, "C-001", "Практика", 40)
	replacement := newTestCriterion(t, "C-002", "Проект", 60)
	var rubric domain.Rubric
	if err := rubric.AddCriterion(original); err != nil {
		t.Fatalf("AddCriterion(%q) error = %v, want nil", original.ID(), err)
	}

	criteria := rubric.Criteria()
	if len(criteria) != 1 {
		t.Fatalf("len(Criteria()) = %d, want 1", len(criteria))
	}
	criteria[0] = replacement

	want := []domain.Criterion{original}
	if got := rubric.Criteria(); !slices.Equal(got, want) {
		t.Errorf("Criteria() = %v, want %v", got, want)
	}
}

func TestRubricFindCriterion(t *testing.T) {
	t.Parallel()

	first := newTestCriterion(t, "C-001", "Практика", 40)
	second := newTestCriterion(t, "C-002", "Проект", 60)
	criteria := []domain.Criterion{first, second}
	cases := []struct {
		name     string
		criteria []domain.Criterion
		id       string
		want     domain.Criterion
		wantOK   bool
	}{
		{
			name:     "second criterion",
			criteria: criteria,
			id:       second.ID(),
			want:     second,
			wantOK:   true,
		},
		{
			name:     "unknown ID",
			criteria: criteria,
			id:       "C-003",
			want:     domain.Criterion{},
			wantOK:   false,
		},
		{
			name:   "empty rubric",
			id:     first.ID(),
			want:   domain.Criterion{},
			wantOK: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var rubric domain.Rubric
			for _, criterion := range tc.criteria {
				if err := rubric.AddCriterion(criterion); err != nil {
					t.Fatalf("AddCriterion(%q) error = %v, want nil", criterion.ID(), err)
				}
			}

			got, ok := rubric.FindCriterion(tc.id)

			if got != tc.want || ok != tc.wantOK {
				t.Errorf("FindCriterion(%q) = (%v, %t), want (%v, %t)", tc.id, got, ok, tc.want, tc.wantOK)
			}
		})
	}
}

func newTestCriterion(t *testing.T, id, name string, percent int) domain.Criterion {
	t.Helper()

	weight, err := domain.NewWeight(percent)
	if err != nil {
		t.Fatalf("NewWeight(%d) error = %v, want nil", percent, err)
	}
	const maxScore = 1000
	criterion, err := domain.NewCriterion(id, name, maxScore, weight)
	if err != nil {
		t.Fatalf("NewCriterion(%q, %q, %d, %v) error = %v, want nil", id, name, maxScore, weight, err)
	}
	return criterion
}
