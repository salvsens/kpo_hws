package cli

import (
	"fmt"
	"strings"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
)

// setGrade records a new grade from the command arguments.
func (r *REPL) setGrade(args string) (string, error) {
	return applyScore(args, "Оценка выставлена", r.gradebook.SetGrade)
}

// editGrade changes a grade from the command arguments.
func (r *REPL) editGrade(args string) (string, error) {
	return applyScore(args, "Оценка исправлена", r.gradebook.EditGrade)
}

// applyScore reads a score and runs the grade operation.
func applyScore(args, message string, change func(string, string, int64) (app.CriterionSummary, error)) (string, error) {
	fields := strings.Fields(args)
	if len(fields) != 3 {
		return "", errInvalidArguments
	}
	score, err := parseHundredths(fields[2])
	if err != nil {
		return "", err
	}
	criterion, err := change(fields[0], fields[1], score)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s (%s): %s из %s.", message, criterion.Name, formatScore(score),
		formatScore(criterion.MaxScoreHundredths)), nil
}

// cancelGrade removes the grade for the student and criterion.
func (r *REPL) cancelGrade(args string) (string, error) {
	fields := strings.Fields(args)
	if len(fields) != 2 {
		return "", errInvalidArguments
	}
	criterion, err := r.gradebook.CancelGrade(fields[0], fields[1])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Оценка отменена (%s).", criterion.Name), nil
}

// listGrades returns the scores of a student in rubric order.
func (r *REPL) listGrades(args string) (string, error) {
	fields := strings.Fields(args)
	if len(fields) != 1 {
		return "", errInvalidArguments
	}
	grades, err := r.gradebook.ListGrades(fields[0])
	if err != nil {
		return "", err
	}
	if len(grades) == 0 {
		return "Оценок пока нет. Выставьте оценку командой set-grade.", nil
	}
	lines := make([]string, 0, len(grades))
	for _, grade := range grades {
		criterion := grade.Criterion
		lines = append(lines, fmt.Sprintf("%s  %s  %s из %s  вес %d%%",
			criterion.ID, criterion.Name, formatScore(grade.ScoreHundredths),
			formatScore(criterion.MaxScoreHundredths), criterion.WeightPercent))
	}
	return strings.Join(lines, "\n"), nil
}

// showResult returns the current or final result.
func (r *REPL) showResult(args string) (string, error) {
	fields := strings.Fields(args)
	if len(fields) != 2 {
		return "", errInvalidArguments
	}
	result, err := r.gradebook.CalculateResult(fields[0], app.ResultMode(fields[1]))
	if err != nil {
		return "", err
	}
	title := "Текущий итог"
	if result.Mode == app.FinalMode {
		title = "Финальный итог"
	}
	return fmt.Sprintf("%s (%s): %s из 10.", title, result.Student.ID, formatHundredths(result.ScoreHundredths)), nil
}
