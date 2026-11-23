package cli_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
	"github.com/salvsens/kpo_hws/hw2/internal/cli"
	"github.com/salvsens/kpo_hws/hw2/internal/memory"
)

func TestREPLGradebookScenario(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`add-student Анна Смирнова
add-criterion 10 40 Практика
add-criterion 20 60 Итоговый проект
students
criteria
set-grade S-001 C-002 15
set-grade S-001 C-001 8
grades S-001
result S-001 final
edit-grade S-001 C-001 12
grades S-001
result S-001 final
edit-grade S-001 C-002 20
result S-001 final
cancel-grade S-001 C-002
result S-001 current
result S-001 final
set-grade S-001 C-002 20
result S-001 final
exit
add-student После выхода`)
	repl := newTestREPL()
	var out bytes.Buffer

	if err := repl.Run(input, &out); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	assertOutputOrder(t, out.String(), []string{
		"Студент добавлен (S-001): Анна Смирнова.",
		"Критерий добавлен (C-001): Практика, максимум 10, вес 40%.",
		"Критерий добавлен (C-002): Итоговый проект, максимум 20, вес 60%.",
		"S-001  Анна Смирнова\n",
		"C-001  Практика  максимум 10  вес 40%\nC-002  Итоговый проект  максимум 20  вес 60%\n",
		"Оценка выставлена (Итоговый проект): 15 из 20.",
		"Оценка выставлена (Практика): 8 из 10.",
		"C-001  Практика  8 из 10  вес 40%\nC-002  Итоговый проект  15 из 20  вес 60%\n",
		"Финальный итог (S-001): 7.70 из 10.",
		"Ошибка. Балл 12 превышает максимум критерия «Практика». Введите балл от 0 до 10.",
		"C-001  Практика  8 из 10  вес 40%",
		"Финальный итог (S-001): 7.70 из 10.",
		"Оценка исправлена (Итоговый проект): 20 из 20.",
		"Финальный итог (S-001): 9.20 из 10.",
		"Оценка отменена (Итоговый проект).",
		"Текущий итог (S-001): 8.00 из 10.",
		"Финальный итог недоступен. Выставьте оценку по критерию «Итоговый проект» (C-002).",
		"Оценка выставлена (Итоговый проект): 20 из 20.",
		"Финальный итог (S-001): 9.20 из 10.",
		"Работа завершена.\n",
	})
	if strings.Contains(out.String(), "После выхода") {
		t.Errorf("Run() processed input after exit:\n%s", out.String())
	}
}

func TestREPLContinuesAfterCommandErrors(t *testing.T) {
	t.Parallel()

	input := strings.NewReader(`students
criteria
students лишний аргумент
несуществующая-команда
exit лишний аргумент
add-student
add-student	Иван Петров
grades S-001
result S-001 current
result S-001 final
result S-001 unknown
grades S-999
add-criterion 10 12.5 Неверный вес
add-criterion 10 0 Нулевой вес
add-criterion 0 100 Нулевой максимум
add-criterion 10.5 100 Тест с дробным максимумом
add-criterion 10 1 Лишний вес
set-grade S-001 C-999 8
set-grade S-001 C-001 8.251
set-grade S-001 C-001 8,25
set-grade S-001 C-001 -0.01
set-grade S-001 C-001 8.25
set-grade S-001 C-001 8
cancel-grade S-001 C-001
edit-grade S-001 C-001 8
help`)
	var out bytes.Buffer

	if err := newTestREPL().Run(input, &out); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	assertOutputOrder(t, out.String(), []string{
		"Студентов пока нет.",
		"Рубрикатор пуст.",
		"Используйте команду: students",
		"Неизвестная команда",
		"Используйте команду: exit",
		"Имя студента пустое.",
		"Студент добавлен (S-001): Иван Петров.",
		"Оценок пока нет.",
		"Текущий итог недоступен.",
		"Добавьте критерии, чтобы сумма весов стала 100%.",
		"Укажите current или final.",
		"Студент не найден.",
		"вес \"12.5\" задан неверно.",
		"Вес критерия вне диапазона.",
		"Максимум критерия должен быть больше нуля.",
		"Критерий добавлен (C-001): Тест с дробным максимумом, максимум 10.5, вес 100%.",
		"Сумма весов превысит 100%.",
		"Критерий не найден.",
		"число \"8.251\" задано неверно.",
		"число \"8,25\" задано неверно.",
		"Балл -0.01 меньше нуля",
		"Оценка выставлена (Тест с дробным максимумом): 8.25 из 10.5.",
		"Для изменения используйте edit-grade.",
		"Оценка отменена (Тест с дробным максимумом).",
		"Оценка не найдена.",
		"Доступные команды:",
		"add-student <имя>",
		"result <студент> <current|final>",
		"exit\n",
	})
}

func TestREPLReturnsStreamErrors(t *testing.T) {
	t.Parallel()

	streamErr := errors.New("stream failed")
	cases := []struct {
		writeFailure bool
		name         string
	}{
		{
			name: "read error",
		},
		{
			name:         "write error",
			writeFailure: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repl := newTestREPL()
			if tc.writeFailure {
				if err := repl.Run(strings.NewReader("help"), failedWriter{err: streamErr}); !errors.Is(err, streamErr) {
					t.Errorf("Run() error = %v, want %v", err, streamErr)
				}
				return
			}
			var out bytes.Buffer
			if err := repl.Run(iotest.ErrReader(streamErr), &out); !errors.Is(err, streamErr) {
				t.Errorf("Run() error = %v, want %v", err, streamErr)
			}
		})
	}
}

type failedWriter struct {
	err error
}

func (w failedWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}

func newTestREPL() *cli.REPL {
	return cli.New(app.NewGradebook(memory.NewRepository()))
}

func assertOutputOrder(t *testing.T, out string, parts []string) {
	t.Helper()

	remaining := out
	for _, part := range parts {
		_, rest, found := strings.Cut(remaining, part)
		if !found {
			t.Fatalf("output missing %q in expected order:\n%s", part, out)
		}
		remaining = rest
	}
}
