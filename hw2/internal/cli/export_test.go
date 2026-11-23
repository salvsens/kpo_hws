package cli_test

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestREPLExport(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	unavailable := filepath.Join(dir, "итог без оценки.csv")
	path := filepath.Join(dir, "текущий итог.csv")
	input := strings.Join([]string{
		"add-student Анна Смирнова",
		"add-criterion 10 100 Практика",
		"export S-001 final " + unavailable,
		"set-grade S-001 C-001 8",
		"result S-001 current",
		"export S-001 current " + path,
		"edit-grade S-001 C-001 9",
		"export S-001 current " + path,
		"exit",
	}, "\n")
	var out bytes.Buffer

	if err := newTestREPL().Run(strings.NewReader(input), &out); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	assertOutputOrder(t, out.String(), []string{
		"Финальный итог недоступен.",
		"Текущий итог (S-001): 8.00 из 10.",
		"Итог сохранён в " + path + ".",
		"Файл уже существует. Выберите новый путь.",
		"Работа завершена.",
	})
	if _, err := os.Stat(unavailable); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Stat(%q) error = %v, want %v", unavailable, err, fs.ErrNotExist)
	}
	assertCSV(t, path, []string{"S-001", "Анна Смирнова", "current", "8.00"})
}

func TestREPLDemoUsesSeparateJournal(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	first := filepath.Join(dir, "первое демо.csv")
	second := filepath.Join(dir, "второе демо.csv")
	input := strings.Join([]string{
		"add-student Иван Петров",
		"add-criterion 10 100 Контрольная",
		"set-grade S-001 C-001 5",
		"demo " + first,
		"demo " + second,
		"students",
		"criteria",
		"grades S-001",
	}, "\n")
	var out bytes.Buffer

	if err := newTestREPL().Run(strings.NewReader(input), &out); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	assertOutputOrder(t, out.String(), []string{
		"Студент добавлен (S-001): Иван Петров.",
		"Демонстрация ПОС. Отдельный журнал.",
		"Студент добавлен (S-001): Анна Смирнова.",
		"Финальный итог (S-001): 7.70 из 10.",
		"Балл 12 превышает максимум критерия «Практика».",
		"Финальный итог (S-001): 7.70 из 10.",
		"Финальный итог (S-001): 9.20 из 10.",
		"Текущий итог (S-001): 8.00 из 10.",
		"Финальный итог недоступен.",
		"Итог сохранён в " + first + ".",
		"Демонстрация ПОС. Отдельный журнал.",
		"Студент добавлен (S-001): Анна Смирнова.",
		"Итог сохранён в " + second + ".",
		"S-001  Иван Петров\nC-001  Контрольная  максимум 10  вес 100%\nC-001  Контрольная  5 из 10  вес 100%\n",
	})
	want := []string{"S-001", "Анна Смирнова", "final", "9.20"}
	assertCSV(t, first, want)
	assertCSV(t, second, want)
}

func TestREPLDemoReportsExportFailure(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "missing", "result.csv")
	input := strings.NewReader("demo " + path + "\nexit\n")
	var out bytes.Buffer

	if err := newTestREPL().Run(input, &out); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}

	assertOutputOrder(t, out.String(), []string{
		"> export S-001 final " + path,
		"Ошибка. Папка для CSV не найдена. Проверьте путь к файлу.",
		"Работа завершена.",
	})
	if strings.Contains(out.String(), "Итог сохранён") {
		t.Errorf("Run() reported a successful export:\n%s", out.String())
	}
}

func assertCSV(t *testing.T, path string, row []string) {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open CSV: %v", err)
	}
	got, readErr := csv.NewReader(file).ReadAll()
	if err := errors.Join(readErr, file.Close()); err != nil {
		t.Fatalf("read CSV: %v", err)
	}
	want := [][]string{
		{"student_id", "student_name", "mode", "result"},
		row,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CSV records = %q, want %q", got, want)
	}
}
