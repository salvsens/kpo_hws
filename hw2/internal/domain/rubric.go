package domain

import (
	"errors"
	"fmt"
	"slices"
)

// maxTotalWeightPercent is the maximum total weight of the criteria.
const maxTotalWeightPercent = 100

var (
	// ErrInvalidCriterion indicates that the criterion has no ID.
	ErrInvalidCriterion = errors.New("invalid criterion")
	// ErrDuplicateCriterionID indicates that the rubric already contains the ID.
	ErrDuplicateCriterionID = errors.New("duplicate criterion ID")
	// ErrTotalWeightExceeded indicates that the total weight would exceed 100 percent.
	ErrTotalWeightExceeded = errors.New("total weight exceeded")
)

// Rubric contains the criteria for a course.
// The zero value is an empty rubric.
type Rubric struct {
	criteria []Criterion
}

// AddCriterion adds a criterion to the end of the list.
// It checks the criterion ID and the total weight before it adds the criterion.
func (r *Rubric) AddCriterion(criterion Criterion) error {
	if criterion.ID() == "" {
		return ErrInvalidCriterion
	}

	var sum int
	for _, c := range r.criteria {
		if c.ID() == criterion.ID() {
			return fmt.Errorf("%w: %s", ErrDuplicateCriterionID, criterion.ID())
		}
		sum += c.Weight().Percent()
	}

	if sum+criterion.Weight().Percent() > maxTotalWeightPercent {
		return ErrTotalWeightExceeded
	}

	r.criteria = append(r.criteria, criterion)

	return nil
}

// Criteria returns a copy of the criteria in the order in which they were added.
func (r *Rubric) Criteria() []Criterion {
	return slices.Clone(r.criteria)
}

// Clone returns a copy of the rubric with an independent list of criteria.
func (r Rubric) Clone() Rubric {
	return Rubric{criteria: r.Criteria()}
}

// FindCriterion returns the criterion with the specified ID.
// It returns a zero criterion and false if the rubric has no criterion with that ID.
func (r *Rubric) FindCriterion(id string) (Criterion, bool) {
	index := slices.IndexFunc(r.criteria, func(criterion Criterion) bool {
		return criterion.ID() == id
	})
	if index < 0 {
		return Criterion{}, false
	}
	return r.criteria[index], true
}
