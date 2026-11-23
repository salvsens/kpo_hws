// Package app provides the operations of the gradebook.
package app

import (
	"fmt"
	"sync"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

// StudentSummary gives the ID and name of a student.
type StudentSummary struct {
	ID   string
	Name string
}

// CriterionSummary gives the details of a criterion.
type CriterionSummary struct {
	ID                 string
	Name               string
	MaxScoreHundredths int64
	WeightPercent      int
}

// Gradebook runs the operations of one journal.
type Gradebook struct {
	mu         sync.RWMutex
	repository Repository
	exporter   ResultExporter
}

// NewGradebook creates a gradebook with the supplied repository and exporter.
// Use one gradebook for each repository.
func NewGradebook(repository Repository, exporter ResultExporter) *Gradebook {
	return &Gradebook{
		repository: repository,
		exporter:   exporter,
	}
}

// AddStudent creates a student with the next student ID.
func (g *Gradebook) AddStudent(name string) (StudentSummary, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	students, err := g.repository.ListStudents()
	if err != nil {
		return StudentSummary{}, err
	}
	id := fmt.Sprintf("S-%03d", len(students)+1)
	student, err := domain.NewStudent(id, name)
	if err != nil {
		return StudentSummary{}, err
	}
	if err := g.repository.CreateStudent(student); err != nil {
		return StudentSummary{}, err
	}
	return summarizeStudent(student), nil
}

// ListStudents returns the students in the order of creation.
func (g *Gradebook) ListStudents() ([]StudentSummary, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	students, err := g.repository.ListStudents()
	if err != nil {
		return nil, err
	}
	summaries := make([]StudentSummary, 0, len(students))
	for _, student := range students {
		summaries = append(summaries, summarizeStudent(student))
	}
	return summaries, nil
}

// AddCriterion adds a criterion with the next criterion ID.
// Give the maximum score in hundredths of a point and the weight in percent.
func (g *Gradebook) AddCriterion(name string, maxScoreHundredths int64, weightPercent int) (CriterionSummary, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	rubric, err := g.repository.LoadRubric()
	if err != nil {
		return CriterionSummary{}, err
	}
	weight, err := domain.NewWeight(weightPercent)
	if err != nil {
		return CriterionSummary{}, err
	}
	id := fmt.Sprintf("C-%03d", len(rubric.Criteria())+1)
	criterion, err := domain.NewCriterion(id, name, maxScoreHundredths, weight)
	if err != nil {
		return CriterionSummary{}, err
	}
	if err := rubric.AddCriterion(criterion); err != nil {
		return CriterionSummary{}, err
	}
	if err := g.repository.SaveRubric(rubric); err != nil {
		return CriterionSummary{}, err
	}
	return summarizeCriterion(criterion), nil
}

// ListCriteria returns the criteria in rubric order.
func (g *Gradebook) ListCriteria() ([]CriterionSummary, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	rubric, err := g.repository.LoadRubric()
	if err != nil {
		return nil, err
	}
	criteria := rubric.Criteria()
	summaries := make([]CriterionSummary, 0, len(criteria))
	for _, criterion := range criteria {
		summaries = append(summaries, summarizeCriterion(criterion))
	}
	return summaries, nil
}

// loadStudentAndRubric returns the student and the current rubric.
func (g *Gradebook) loadStudentAndRubric(studentID string) (*domain.Student, domain.Rubric, error) {
	student, err := g.repository.FindStudent(studentID)
	if err != nil {
		return nil, domain.Rubric{}, err
	}
	rubric, err := g.repository.LoadRubric()
	if err != nil {
		return nil, domain.Rubric{}, err
	}
	return student, rubric, nil
}

// summarizeStudent returns the ID and name of the student.
func summarizeStudent(student *domain.Student) StudentSummary {
	return StudentSummary{
		ID:   student.ID(),
		Name: student.Name(),
	}
}

// summarizeCriterion returns the details of the criterion.
func summarizeCriterion(criterion domain.Criterion) CriterionSummary {
	return CriterionSummary{
		ID:                 criterion.ID(),
		Name:               criterion.Name(),
		MaxScoreHundredths: criterion.MaxScoreHundredths(),
		WeightPercent:      criterion.Weight().Percent(),
	}
}
