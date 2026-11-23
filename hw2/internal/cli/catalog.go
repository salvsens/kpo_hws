package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// addStudent adds a student with the supplied name.
func (r *REPL) addStudent(name string) (string, error) {
	student, err := r.gradebook.AddStudent(name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Студент добавлен (%s): %s.", student.ID, student.Name), nil
}

// listStudents returns the names and IDs in the journal.
func (r *REPL) listStudents(_ string) (string, error) {
	students, err := r.gradebook.ListStudents()
	if err != nil {
		return "", err
	}
	if len(students) == 0 {
		return "Студентов пока нет. Добавьте студента командой add-student.", nil
	}
	lines := make([]string, 0, len(students))
	for _, student := range students {
		lines = append(lines, student.ID+"  "+student.Name)
	}
	return strings.Join(lines, "\n"), nil
}

// addCriterion reads the maximum, weight, and name and adds the criterion.
func (r *REPL) addCriterion(args string) (string, error) {
	maximumText, rest := splitWord(args)
	weightText, name := splitWord(rest)
	if name == "" {
		return "", errInvalidArguments
	}
	maximum, err := parseHundredths(maximumText)
	if err != nil {
		return "", fmt.Errorf("максимум: %w", err)
	}
	weight, err := strconv.Atoi(weightText)
	if err != nil {
		return "", fmt.Errorf("вес %q задан неверно. Введите целое число от 1 до 100", weightText)
	}
	criterion, err := r.gradebook.AddCriterion(name, maximum, weight)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Критерий добавлен (%s): %s, максимум %s, вес %d%%.",
		criterion.ID, criterion.Name, formatScore(criterion.MaxScoreHundredths), criterion.WeightPercent), nil
}

// listCriteria returns the criteria in rubric order.
func (r *REPL) listCriteria(_ string) (string, error) {
	criteria, err := r.gradebook.ListCriteria()
	if err != nil {
		return "", err
	}
	if len(criteria) == 0 {
		return "Рубрикатор пуст. Добавьте критерий командой add-criterion.", nil
	}
	lines := make([]string, 0, len(criteria))
	for _, criterion := range criteria {
		lines = append(lines, fmt.Sprintf("%s  %s  максимум %s  вес %d%%",
			criterion.ID, criterion.Name, formatScore(criterion.MaxScoreHundredths), criterion.WeightPercent))
	}
	return strings.Join(lines, "\n"), nil
}
