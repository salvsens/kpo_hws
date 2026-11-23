package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrCriterionNotFound indicates that the rubric has no criterion with the specified ID.
	ErrCriterionNotFound = errors.New("criterion not found")
	// ErrScoreOutOfRange indicates that the score is below zero or above the maximum.
	ErrScoreOutOfRange = errors.New("score out of range")
)

// ScoreRangeError gives the score and the maximum of the criterion.
type ScoreRangeError struct {
	CriterionName      string
	ScoreHundredths    int64
	MaxScoreHundredths int64
}

// Error returns the error message with the name of the criterion.
func (e ScoreRangeError) Error() string {
	return fmt.Sprintf("%s: %q", ErrScoreOutOfRange, e.CriterionName)
}

// Unwrap returns ErrScoreOutOfRange.
func (e ScoreRangeError) Unwrap() error {
	return ErrScoreOutOfRange
}

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
		return Evaluation{}, ScoreRangeError{
			CriterionName:      criterion.Name(),
			ScoreHundredths:    scoreHundredths,
			MaxScoreHundredths: criterion.MaxScoreHundredths(),
		}
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
