package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/salvsens/kpo_hws/hw2/internal/domain"
)

// demoStep gives the command and the expected error for one demo step.
type demoStep struct {
	line    string
	wantErr error
}

// runDemo runs the sample commands on an empty journal and exports its final result.
func (r *REPL) runDemo(path string) (string, error) {
	if path == "" {
		return "", errInvalidArguments
	}

	demo := New(r.newDemo(), r.newDemo)
	lines := []string{"Демонстрация ПОС. Отдельный журнал."}
	for _, step := range demoScript(path) {
		lines = append(lines, "> "+step.line)
		response := demo.execute(step.line)
		if !errors.Is(response.err, step.wantErr) {
			if response.err != nil {
				return strings.Join(lines, "\n") + "\n", response.err
			}
			return strings.Join(lines, "\n") + "\n", fmt.Errorf("демо ожидало отказ для команды %q", step.line)
		}
		lines = append(lines, response.text)
	}
	return strings.Join(lines, "\n"), nil
}

// demoScript gives the commands from the specification.
func demoScript(path string) []demoStep {
	// finalResultCommand shows the final result of the demo student.
	const finalResultCommand = "result S-001 final"

	return []demoStep{
		{line: "add-student Анна Смирнова"},
		{line: "add-criterion 10 40 Практика"},
		{line: "add-criterion 20 60 Проект"},
		{line: studentsCommand},
		{line: criteriaCommand},
		{line: "set-grade S-001 C-001 8"},
		{line: "set-grade S-001 C-002 15"},
		{line: "grades S-001"},
		{line: finalResultCommand},
		{line: "edit-grade S-001 C-001 12", wantErr: domain.ErrScoreOutOfRange},
		{line: "grades S-001"},
		{line: finalResultCommand},
		{line: "edit-grade S-001 C-002 20"},
		{line: finalResultCommand},
		{line: "cancel-grade S-001 C-002"},
		{line: "result S-001 current"},
		{line: finalResultCommand, wantErr: domain.ErrGradeNotFound},
		{line: "set-grade S-001 C-002 20"},
		{line: finalResultCommand},
		{line: "export S-001 final " + path},
	}
}
