package domain

import "errors"

var (
	// ErrCriterionNotFound indicates that the rubric has no criterion with the specified ID.
	ErrCriterionNotFound = errors.New("criterion not found")
	// ErrScoreOutOfRange indicates that the score is below zero or above the maximum.
	ErrScoreOutOfRange = errors.New("score out of range")
)

// Evaluation holds a score for one criterion.
type Evaluation struct {
	criterionID     string
	scoreHundredths int64
}

// NewEvaluation creates an evaluation for a criterion in the rubric.
// It accepts scores from zero to the maximum of the criterion, inclusive.
func NewEvaluation(rubric Rubric, criterionID string, scoreHundredths int64) (Evaluation, error) {
	criterion, ok := rubric.FindCriterion(criterionID)
	if !ok {
		return Evaluation{}, ErrCriterionNotFound
	}
	if scoreHundredths < 0 || scoreHundredths > criterion.MaxScoreHundredths() {
		return Evaluation{}, ErrScoreOutOfRange
	}

	return Evaluation{
		criterionID:     criterionID,
		scoreHundredths: scoreHundredths,
	}, nil
}

// CriterionID returns the ID of the criterion.
func (e Evaluation) CriterionID() string {
	return e.criterionID
}

// ScoreHundredths returns the score in hundredths of a point.
func (e Evaluation) ScoreHundredths() int64 {
	return e.scoreHundredths
}
