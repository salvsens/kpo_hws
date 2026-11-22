// Package memory stores a journal in memory.
package memory

import (
	"fmt"
	"sync"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

var _ app.Repository = (*Repository)(nil)

// Repository holds the students and rubric for one journal.
// The zero value is an empty repository.
type Repository struct {
	mu       sync.RWMutex
	students map[string]*domain.Student
	order    []string
	rubric   domain.Rubric
}

// NewRepository creates an empty repository.
func NewRepository() *Repository {
	return &Repository{}
}

// CreateStudent adds a copy of the student to the repository.
func (r *Repository) CreateStudent(student *domain.Student) error {
	if student == nil || student.ID() == "" {
		return domain.ErrInvalidStudentID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	id := student.ID()
	if _, exists := r.students[id]; exists {
		return fmt.Errorf("create student %q: %w", id, app.ErrStudentAlreadyExists)
	}
	if r.students == nil {
		r.students = make(map[string]*domain.Student)
	}
	r.students[id] = student.Clone()
	r.order = append(r.order, id)
	return nil
}

// FindStudent returns a copy of the student with the specified ID.
func (r *Repository) FindStudent(id string) (*domain.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	student, exists := r.students[id]
	if !exists {
		return nil, fmt.Errorf("find student %q: %w", id, app.ErrStudentNotFound)
	}
	return student.Clone(), nil
}

// ListStudents returns copies of the students in the order of creation.
func (r *Repository) ListStudents() ([]*domain.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	students := make([]*domain.Student, 0, len(r.order))
	for _, id := range r.order {
		students = append(students, r.students[id].Clone())
	}
	return students, nil
}

// SaveStudent replaces an existing student with a copy of the supplied student.
func (r *Repository) SaveStudent(student *domain.Student) error {
	if student == nil || student.ID() == "" {
		return domain.ErrInvalidStudentID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	id := student.ID()
	if _, exists := r.students[id]; !exists {
		return fmt.Errorf("save student %q: %w", id, app.ErrStudentNotFound)
	}
	r.students[id] = student.Clone()
	return nil
}

// LoadRubric returns a copy of the current rubric.
func (r *Repository) LoadRubric() (domain.Rubric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.rubric.Clone(), nil
}

// SaveRubric stores a copy of the supplied rubric.
func (r *Repository) SaveRubric(rubric domain.Rubric) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rubric = rubric.Clone()
	return nil
}
