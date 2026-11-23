package cli

import "github.com/salvsens/kpo_hws/hw2/internal/app"

// exportResult saves the result for the student and mode at the supplied path.
func (r *REPL) exportResult(args string) (string, error) {
	studentID, rest := splitWord(args)
	mode, path := splitWord(rest)
	if path == "" {
		return "", errInvalidArguments
	}
	if err := r.gradebook.ExportResult(studentID, app.ResultMode(mode), path); err != nil {
		return "", err
	}
	return "Итог сохранён в " + path + ".", nil
}
