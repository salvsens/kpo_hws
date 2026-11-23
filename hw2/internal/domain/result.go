package domain

import (
	"errors"
	"math/big"
)

// maxResultHundredths is the upper limit of the result scale in hundredths.
const maxResultHundredths = 1000

var (
	// ErrNoGrades indicates that the student has no grades in the rubric.
	ErrNoGrades = errors.New("no grades")
	// ErrIncompleteRubric indicates that the total weight is below 100 percent.
	ErrIncompleteRubric = errors.New("incomplete rubric")
)

// MissingGradeError identifies a criterion that requires a grade for the final result.
type MissingGradeError struct {
	CriterionID   string
	CriterionName string
}

// ResultStrategy selects the grades for a result.
type ResultStrategy interface {
	selectGrades(student *Student, rubric Rubric) ([]weightedGrade, error)
}

// CurrentResult selects the criteria that have grades.
type CurrentResult struct{}

// FinalResult requires a grade for each criterion and a total weight of 100 percent.
type FinalResult struct{}

// weightedGrade pairs a grade with the criterion that defines its weight.
type weightedGrade struct {
	criterion  Criterion
	evaluation Evaluation
}

// Error returns the error message with the criterion ID.
func (e MissingGradeError) Error() string {
	return ErrGradeNotFound.Error() + ": " + e.CriterionID
}

// Unwrap returns ErrGradeNotFound.
func (e MissingGradeError) Unwrap() error {
	return ErrGradeNotFound
}

// CalculateResult returns the weighted result in hundredths of a point.
// Use a student from NewStudent and a CurrentResult or FinalResult value.
func CalculateResult(student *Student, rubric Rubric, strategy ResultStrategy) (int64, error) {
	grades, err := strategy.selectGrades(student, rubric)
	if err != nil {
		return 0, err
	}

	var sum big.Rat
	var totalWeight int64
	for _, grade := range grades {
		weight := int64(grade.criterion.Weight().Percent())
		contribution := big.NewRat(grade.evaluation.ScoreHundredths(), grade.criterion.MaxScoreHundredths())
		contribution.Mul(contribution, big.NewRat(weight, 1))
		sum.Add(&sum, contribution)
		totalWeight += weight
	}
	sum.Mul(&sum, big.NewRat(maxResultHundredths, totalWeight))

	// Add one half before the division to round up at the midpoint.
	const halfDenominator = 2
	sum.Add(&sum, big.NewRat(1, halfDenominator))
	var rounded big.Int
	return rounded.Quo(sum.Num(), sum.Denom()).Int64(), nil
}

// selectGrades selects all grades that belong to the rubric.
func (CurrentResult) selectGrades(student *Student, rubric Rubric) ([]weightedGrade, error) {
	grades, _, _ := collectResultGrades(student, rubric)
	if len(grades) == 0 {
		return nil, ErrNoGrades
	}
	return grades, nil
}

// selectGrades selects all grades for a complete rubric.
func (FinalResult) selectGrades(student *Student, rubric Rubric) ([]weightedGrade, error) {
	grades, totalWeight, missing := collectResultGrades(student, rubric)
	if totalWeight != maxTotalWeightPercent {
		return nil, ErrIncompleteRubric
	}
	if missing.ID() != "" {
		return nil, MissingGradeError{
			CriterionID:   missing.ID(),
			CriterionName: missing.Name(),
		}
	}
	return grades, nil
}

// collectResultGrades collects grades, total rubric weight, and the first criterion without a grade.
func collectResultGrades(student *Student, rubric Rubric) ([]weightedGrade, int, Criterion) {
	var grades []weightedGrade
	var totalWeight int
	var missing Criterion
	for _, criterion := range rubric.criteria {
		totalWeight += criterion.Weight().Percent()
		evaluation, ok := student.FindGrade(criterion.ID())
		if ok {
			grades = append(grades, weightedGrade{
				criterion:  criterion,
				evaluation: evaluation,
			})
		} else if missing.ID() == "" {
			missing = criterion
		}
	}
	return grades, totalWeight, missing
}
