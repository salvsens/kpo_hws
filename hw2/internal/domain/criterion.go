package domain

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidCriterionID indicates that the criterion ID is blank.
	ErrInvalidCriterionID = errors.New("criterion id is blank")
	// ErrInvalidCriterionName indicates that the criterion name is blank.
	ErrInvalidCriterionName = errors.New("criterion name is blank")
	// ErrInvalidMaxScore indicates that the maximum score is zero or less.
	ErrInvalidMaxScore = errors.New("maximum score must be positive")
)

// Criterion defines the maximum score and weight for one part of the gradebook.
type Criterion struct {
	id                 string
	name               string
	maxScoreHundredths int64
	weight             Weight
}

// NewCriterion creates a criterion with a maximum score in hundredths of a point.
// It returns an error if the ID or name is blank, the maximum is zero or less, or the weight is invalid.
func NewCriterion(id, name string, maxScoreHundredths int64, weight Weight) (Criterion, error) {
	if strings.TrimSpace(id) == "" {
		return Criterion{}, ErrInvalidCriterionID
	}
	if strings.TrimSpace(name) == "" {
		return Criterion{}, ErrInvalidCriterionName
	}
	if maxScoreHundredths <= 0 {
		return Criterion{}, ErrInvalidMaxScore
	}
	if weight.Percent() == 0 {
		return Criterion{}, ErrInvalidWeight
	}

	return Criterion{
		id:                 id,
		name:               name,
		maxScoreHundredths: maxScoreHundredths,
		weight:             weight,
	}, nil
}

// ID returns the criterion ID.
func (c Criterion) ID() string {
	return c.id
}

// Name returns the criterion name.
func (c Criterion) Name() string {
	return c.name
}

// MaxScoreHundredths returns the maximum score in hundredths of a point.
func (c Criterion) MaxScoreHundredths() int64 {
	return c.maxScoreHundredths
}

// Weight returns the criterion weight.
func (c Criterion) Weight() Weight {
	return c.weight
}
