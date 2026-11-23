package cli

import (
	"errors"
	"fmt"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

// errorHints gives a reason and an action for each expected error.
var errorHints = []struct {
	err  error
	hint string
}{
	{err: domain.ErrInvalidStudentName, hint: "Имя студента пустое. Введите имя после add-student."},
	{err: domain.ErrInvalidCriterionName, hint: "Название критерия пустое. Введите название после веса."},
	{err: domain.ErrInvalidMaxScore, hint: "Максимум критерия должен быть больше нуля. Введите положительное число."},
	{err: domain.ErrInvalidWeight, hint: "Вес критерия вне диапазона. Введите целое число от 1 до 100."},
	{err: domain.ErrTotalWeightExceeded, hint: "Сумма весов превысит 100%. Укажите меньший вес критерия."},
	{err: app.ErrStudentNotFound, hint: "Студент не найден. Проверьте ID командой students."},
	{err: domain.ErrCriterionNotFound, hint: "Критерий не найден. Проверьте ID командой criteria."},
	{err: domain.ErrGradeAlreadySet, hint: "Оценка уже выставлена. Для изменения используйте edit-grade."},
	{err: domain.ErrGradeNotFound, hint: "Оценка не найдена. Проверьте список командой grades."},
	{err: domain.ErrNoGrades, hint: "Текущий итог недоступен. Выставьте хотя бы одну оценку командой set-grade."},
	{err: domain.ErrIncompleteRubric, hint: "Финальный итог недоступен. Добавьте критерии, чтобы сумма весов стала 100%."},
	{err: app.ErrInvalidResultMode, hint: "Неизвестный режим итога. Укажите current или final."},
}

// explainError returns a message for a command error.
func explainError(err error) string {
	if missing, ok := errors.AsType[domain.MissingGradeError](err); ok {
		return fmt.Sprintf("Финальный итог недоступен. Выставьте оценку по критерию «%s» (%s).",
			missing.CriterionName, missing.CriterionID)
	}
	if score, ok := errors.AsType[domain.ScoreRangeError](err); ok {
		return explainScoreRange(score)
	}
	for _, entry := range errorHints {
		if errors.Is(err, entry.err) {
			return entry.hint
		}
	}
	return err.Error() + "."
}

// explainScoreRange gives the valid score range for a criterion.
func explainScoreRange(err domain.ScoreRangeError) string {
	score := formatScore(err.ScoreHundredths)
	maximum := formatScore(err.MaxScoreHundredths)
	reason := "превышает максимум критерия"
	if err.ScoreHundredths < 0 {
		reason = "меньше нуля для критерия"
	}
	return fmt.Sprintf("Балл %s %s «%s». Введите балл от 0 до %s.", score, reason, err.CriterionName, maximum)
}
