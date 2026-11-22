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
)

// Student identifies a student in the gradebook by ID.
type Student struct {
	id   string
	name string
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
