package app

import (
	"errors"
	"fmt"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

// ResultMode identifies a rule for calculating a result.
type ResultMode string

const (
	// CurrentMode calculates the result from the available grades.
	CurrentMode ResultMode = "current"
	// FinalMode calculates the result when the rubric and all grades are complete.
	FinalMode ResultMode = "final"
)

// ErrInvalidResultMode indicates that the result mode is unknown.
var ErrInvalidResultMode = errors.New("invalid result mode")

// Result gives the student, the mode, and the result in hundredths of a point.
type Result struct {
	Student         StudentSummary
	Mode            ResultMode
	ScoreHundredths int64
}

// CalculateResult returns the weighted result for a student and mode.
func (g *Gradebook) CalculateResult(studentID string, mode ResultMode) (Result, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var strategy domain.ResultStrategy
	switch mode {
	case CurrentMode:
		strategy = domain.CurrentResult{}
	case FinalMode:
		strategy = domain.FinalResult{}
	default:
		return Result{}, fmt.Errorf("calculate result for mode %q: %w", mode, ErrInvalidResultMode)
	}

	student, rubric, err := g.loadStudentAndRubric(studentID)
	if err != nil {
		return Result{}, err
	}
	score, err := domain.CalculateResult(student, rubric, strategy)
	if err != nil {
		return Result{}, err
	}
	return Result{
		Student:         summarizeStudent(student),
		Mode:            mode,
		ScoreHundredths: score,
	}, nil
}
