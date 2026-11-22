// Package domain defines the rules for the gradebook.
package domain

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidStudentID indicates that the student ID is blank.
	ErrInvalidStudentID = errors.New("invalid student id")
	// ErrInvalidStudentName indicates that the student name is blank.
	ErrInvalidStudentName = errors.New("invalid student name")
	// ErrGradeAlreadySet indicates that the student already has a grade for the criterion.
	ErrGradeAlreadySet = errors.New("grade already set")
	// ErrGradeNotFound indicates that the student has no grade for the criterion.
	ErrGradeNotFound = errors.New("grade not found")
)

// Student holds the ID, name, and grades of a student in the gradebook.
type Student struct {
	id     string
	name   string
	grades map[string]Evaluation
}

// NewStudent creates a student with the specified ID and name.
// It returns an error if the ID or name is blank.
func NewStudent(id, name string) (*Student, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidStudentID
	}
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidStudentName
	}
	return &Student{id: id, name: name}, nil
}

// ID returns the student ID.
func (s Student) ID() string {
	return s.id
}

// Name returns the student name.
func (s Student) Name() string {
	return s.name
}

// SetGrade records the first grade for the criterion after it checks the score.
// It returns ErrGradeAlreadySet if the student has a grade for the criterion.
func (s *Student) SetGrade(rubric Rubric, criterionID string, scoreHundredths int64) error {
	if s.id == "" {
		return ErrInvalidStudentID
	}
	if _, exists := s.grades[criterionID]; exists {
		return ErrGradeAlreadySet
	}

	evaluation, err := NewEvaluation(rubric, criterionID, scoreHundredths)
	if err != nil {
		return err
	}

	if s.grades == nil {
		s.grades = make(map[string]Evaluation)
	}
	s.grades[criterionID] = evaluation
	return nil
}

// FindGrade returns the evaluation for the criterion.
// It returns a zero evaluation and false if the student has no grade for the criterion.
func (s Student) FindGrade(criterionID string) (Evaluation, bool) {
	evaluation, ok := s.grades[criterionID]
	return evaluation, ok
}

// EditGrade changes a grade after it checks the new score.
// It returns ErrGradeNotFound if the student has no grade for the criterion.
func (s *Student) EditGrade(rubric Rubric, criterionID string, scoreHundredths int64) error {
	if _, exists := s.grades[criterionID]; !exists {
		return ErrGradeNotFound
	}

	evaluation, err := NewEvaluation(rubric, criterionID, scoreHundredths)
	if err != nil {
		return err
	}

	s.grades[criterionID] = evaluation
	return nil
}

// CancelGrade removes the grade for the criterion.
// It returns ErrGradeNotFound if the student has no grade for the criterion.
func (s *Student) CancelGrade(criterionID string) error {
	if _, exists := s.grades[criterionID]; !exists {
		return ErrGradeNotFound
	}

	delete(s.grades, criterionID)
	return nil
}
