package domain_test

import (
	"errors"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

func TestNewWeight(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		percent int
		wantErr error
	}{
		{name: "minimum", percent: 1},
		{name: "within range", percent: 40},
		{name: "maximum", percent: 100},
		{name: "negative", percent: -1, wantErr: domain.ErrInvalidWeight},
		{name: "zero", percent: 0, wantErr: domain.ErrInvalidWeight},
		{name: "above maximum", percent: 101, wantErr: domain.ErrInvalidWeight},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			weight, err := domain.NewWeight(tc.percent)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("NewWeight(%d) error = %v, want %v", tc.percent, err, tc.wantErr)
			}
			if tc.wantErr != nil {
				if weight != (domain.Weight{}) {
					t.Errorf("NewWeight(%d) weight = %v, want zero value", tc.percent, weight)
				}
				return
			}
			if got := weight.Percent(); got != tc.percent {
				t.Errorf("Percent() = %d, want %d", got, tc.percent)
			}
		})
	}
}

func TestWeightEquality(t *testing.T) {
	t.Parallel()

	const percent = 40
	first, err := domain.NewWeight(percent)
	if err != nil {
		t.Fatalf("NewWeight(%d) error = %v, want nil", percent, err)
	}
	second, err := domain.NewWeight(percent)
	if err != nil {
		t.Fatalf("NewWeight(%d) error = %v, want nil", percent, err)
	}

	if first != second {
		t.Errorf("Weights for %d percent differ: %v and %v", percent, first, second)
	}
}
