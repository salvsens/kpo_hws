package domain_test

import (
	"errors"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

func TestNewCriterion(t *testing.T) {
	t.Parallel()

	const (
		id   = "C-001"
		name = "Итоговый проект"
	)
	weight, err := domain.NewWeight(40)
	if err != nil {
		t.Fatalf("NewWeight(40) error = %v, want nil", err)
	}

	for _, maxScore := range []int64{1, 1050} {
		criterion, err := domain.NewCriterion(id, name, maxScore, weight)
		if err != nil {
			t.Fatalf("NewCriterion(%q, %q, %d, %v) error = %v, want nil", id, name, maxScore, weight, err)
		}
		if got := criterion.ID(); got != id {
			t.Errorf("ID() = %q, want %q", got, id)
		}
		if got := criterion.Name(); got != name {
			t.Errorf("Name() = %q, want %q", got, name)
		}
		if got := criterion.MaxScoreHundredths(); got != maxScore {
			t.Errorf("MaxScoreHundredths() = %d, want %d", got, maxScore)
		}
		if got := criterion.Weight(); got != weight {
			t.Errorf("Weight() = %v, want %v", got, weight)
		}
	}
}

func TestNewCriterionRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	const (
		id       = "C-001"
		name     = "Итоговый проект"
		maxScore = 1050
	)
	weight, err := domain.NewWeight(40)
	if err != nil {
		t.Fatalf("NewWeight(40) error = %v, want nil", err)
	}

	cases := []struct {
		label    string
		id       string
		name     string
		maxScore int64
		weight   domain.Weight
		wantErr  error
	}{
		{
			label:    "empty ID",
			id:       "",
			name:     name,
			maxScore: maxScore,
			weight:   weight,
			wantErr:  domain.ErrInvalidCriterionID,
		},
		{
			label:    "blank ID",
			id:       "  ",
			name:     name,
			maxScore: maxScore,
			weight:   weight,
			wantErr:  domain.ErrInvalidCriterionID,
		},
		{
			label:    "empty name",
			id:       id,
			name:     "",
			maxScore: maxScore,
			weight:   weight,
			wantErr:  domain.ErrInvalidCriterionName,
		},
		{
			label:    "blank name",
			id:       id,
			name:     " \t ",
			maxScore: maxScore,
			weight:   weight,
			wantErr:  domain.ErrInvalidCriterionName,
		},
		{
			label:    "zero maximum",
			id:       id,
			name:     name,
			maxScore: 0,
			weight:   weight,
			wantErr:  domain.ErrInvalidMaxScore,
		},
		{
			label:    "negative maximum",
			id:       id,
			name:     name,
			maxScore: -1,
			weight:   weight,
			wantErr:  domain.ErrInvalidMaxScore,
		},
		{
			label:    "zero weight",
			id:       id,
			name:     name,
			maxScore: maxScore,
			weight:   domain.Weight{},
			wantErr:  domain.ErrInvalidWeight,
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			t.Parallel()

			criterion, err := domain.NewCriterion(tc.id, tc.name, tc.maxScore, tc.weight)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("NewCriterion(%q, %q, %d, %v) error = %v, want %v",
					tc.id, tc.name, tc.maxScore, tc.weight, err, tc.wantErr)
			}
			if criterion != (domain.Criterion{}) {
				t.Errorf("NewCriterion(%q, %q, %d, %v) criterion = %v, want zero value",
					tc.id, tc.name, tc.maxScore, tc.weight, criterion)
			}
		})
	}
}
