package app

import (
	"errors"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

var (
	// ErrStudentNotFound indicates that the journal has no student with the specified ID.
	ErrStudentNotFound = errors.New("student not found")
	// ErrStudentAlreadyExists indicates that the journal contains a student with the specified ID.
	ErrStudentAlreadyExists = errors.New("student already exists")
)

// Repository stores the students and rubric of a journal.
// Each read returns an independent copy, and each write stores an independent copy.
type Repository interface {
	// CreateStudent adds a student and reserves the student ID.
	CreateStudent(student *domain.Student) error
	// FindStudent returns the student with the specified ID.
	FindStudent(id string) (*domain.Student, error)
	// ListStudents returns the students in the order of creation.
	ListStudents() ([]*domain.Student, error)
	// SaveStudent updates an existing student.
	SaveStudent(student *domain.Student) error
	// LoadRubric returns the current rubric.
	LoadRubric() (domain.Rubric, error)
	// SaveRubric replaces the rubric with the supplied rubric.
	SaveRubric(rubric domain.Rubric) error
}
