package app

import (
	"fmt"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

// GradeSummary gives the criterion and the score in hundredths of a point.
type GradeSummary struct {
	Criterion       CriterionSummary
	ScoreHundredths int64
}

// SetGrade records the first grade for a student and criterion.
// It returns the details of the criterion after it saves the grade.
func (g *Gradebook) SetGrade(studentID, criterionID string, scoreHundredths int64) (CriterionSummary, error) {
	return g.changeGrade(studentID, criterionID, func(student *domain.Student, rubric domain.Rubric) error {
		return student.SetGrade(rubric, criterionID, scoreHundredths)
	})
}

// EditGrade changes a grade for a student and criterion.
// It returns the details of the criterion after it saves the grade.
func (g *Gradebook) EditGrade(studentID, criterionID string, scoreHundredths int64) (CriterionSummary, error) {
	return g.changeGrade(studentID, criterionID, func(student *domain.Student, rubric domain.Rubric) error {
		return student.EditGrade(rubric, criterionID, scoreHundredths)
	})
}

// CancelGrade removes a grade for a student and criterion.
// It returns the details of the criterion after it saves the student.
func (g *Gradebook) CancelGrade(studentID, criterionID string) (CriterionSummary, error) {
	return g.changeGrade(studentID, criterionID, func(student *domain.Student, _ domain.Rubric) error {
		return student.CancelGrade(criterionID)
	})
}

// ListGrades returns the grades of the student in rubric order.
func (g *Gradebook) ListGrades(studentID string) ([]GradeSummary, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	student, rubric, err := g.loadStudentAndRubric(studentID)
	if err != nil {
		return nil, err
	}
	var grades []GradeSummary
	for _, criterion := range rubric.Criteria() {
		if evaluation, ok := student.FindGrade(criterion.ID()); ok {
			grades = append(grades, GradeSummary{
				Criterion:       summarizeCriterion(criterion),
				ScoreHundredths: evaluation.ScoreHundredths(),
			})
		}
	}
	return grades, nil
}

// changeGrade saves the student after the grade operation succeeds.
func (g *Gradebook) changeGrade(studentID, criterionID string, change func(*domain.Student, domain.Rubric) error) (CriterionSummary, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	student, rubric, err := g.loadStudentAndRubric(studentID)
	if err != nil {
		return CriterionSummary{}, err
	}
	criterion, ok := rubric.FindCriterion(criterionID)
	if !ok {
		return CriterionSummary{}, fmt.Errorf("find criterion %q: %w", criterionID, domain.ErrCriterionNotFound)
	}
	if err := change(student, rubric); err != nil {
		return CriterionSummary{}, err
	}
	if err := g.repository.SaveStudent(student); err != nil {
		return CriterionSummary{}, err
	}
	return summarizeCriterion(criterion), nil
}
