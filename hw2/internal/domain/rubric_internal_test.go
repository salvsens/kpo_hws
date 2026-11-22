package domain

import (
	"slices"
	"testing"
)

func TestRubricCloneKeepsIndependentCriteria(t *testing.T) {
	t.Parallel()

	first := newRubricCloneCriterion(t, "C-001")
	originalOnly := newRubricCloneCriterion(t, "C-002")
	cloneOnly := newRubricCloneCriterion(t, "C-003")
	criteria := make([]Criterion, 1, 2)
	criteria[0] = first
	rubric := Rubric{criteria: criteria}

	clone := rubric.Clone()
	if err := rubric.AddCriterion(originalOnly); err != nil {
		t.Fatalf("original AddCriterion(%q) error = %v, want nil", originalOnly.ID(), err)
	}
	if err := clone.AddCriterion(cloneOnly); err != nil {
		t.Fatalf("clone AddCriterion(%q) error = %v, want nil", cloneOnly.ID(), err)
	}

	wantOriginal := []Criterion{first, originalOnly}
	if got := rubric.Criteria(); !slices.Equal(got, wantOriginal) {
		t.Errorf("original Criteria() = %v, want %v", got, wantOriginal)
	}
	wantClone := []Criterion{first, cloneOnly}
	if got := clone.Criteria(); !slices.Equal(got, wantClone) {
		t.Errorf("clone Criteria() = %v, want %v", got, wantClone)
	}
}

func TestEmptyRubricCloneAcceptsFirstCriterion(t *testing.T) {
	t.Parallel()

	var rubric Rubric
	criterion := newRubricCloneCriterion(t, "C-001")

	clone := rubric.Clone()
	if err := clone.AddCriterion(criterion); err != nil {
		t.Fatalf("clone AddCriterion(%q) error = %v, want nil", criterion.ID(), err)
	}

	if got := rubric.Criteria(); len(got) != 0 {
		t.Errorf("original Criteria() = %v, want an empty list", got)
	}
	if got := clone.Criteria(); !slices.Equal(got, []Criterion{criterion}) {
		t.Errorf("clone Criteria() = %v, want [%v]", got, criterion)
	}
}

func newRubricCloneCriterion(t *testing.T, id string) Criterion {
	t.Helper()

	weight, err := NewWeight(30)
	if err != nil {
		t.Fatalf("NewWeight(30) error = %v, want nil", err)
	}
	criterion, err := NewCriterion(id, id, 1000, weight)
	if err != nil {
		t.Fatalf("NewCriterion(%q, %q, 1000, %v) error = %v, want nil", id, id, weight, err)
	}
	return criterion
}
