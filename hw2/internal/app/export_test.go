package app_test

import (
	"errors"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
	"github.com/salvsens/kpo_hws/hw2/internal/domain"
	"github.com/salvsens/kpo_hws/hw2/internal/memory"
)

func TestGradebookExportResult(t *testing.T) {
	t.Parallel()

	exporter := &recordingExporter{}
	gradebook := app.NewGradebook(memory.NewRepository(), exporter)
	student := addStudent(t, gradebook, "Анна Смирнова")
	criterion := addCriterion(t, gradebook, "Практика", 1000, 100)
	const path = "result.csv"

	if err := gradebook.ExportResult(student.ID, app.FinalMode, path); !errors.Is(err, domain.ErrGradeNotFound) {
		t.Errorf("ExportResult() error = %v, want %v", err, domain.ErrGradeNotFound)
	}
	if exporter.calls != 0 {
		t.Errorf("WriteResult() calls = %d, want 0 before a result is available", exporter.calls)
	}

	setGrade(t, gradebook, student.ID, criterion, 800)
	want, err := gradebook.CalculateResult(student.ID, app.FinalMode)
	if err != nil {
		t.Fatalf("CalculateResult() error = %v, want nil", err)
	}
	if err := gradebook.ExportResult(student.ID, app.FinalMode, path); err != nil {
		t.Fatalf("ExportResult() error = %v, want nil", err)
	}

	if exporter.calls != 1 || exporter.path != path || exporter.result != want {
		t.Errorf("exporter = %+v, want one call with (%q, %v)", exporter, path, want)
	}
	exporter.err = errors.New("write failed")
	if err := gradebook.ExportResult(student.ID, app.FinalMode, path); !errors.Is(err, exporter.err) {
		t.Errorf("ExportResult() error = %v, want %v", err, exporter.err)
	}
}

type recordingExporter struct {
	path   string
	result app.Result
	err    error
	calls  int
}

func (e *recordingExporter) WriteResult(path string, result app.Result) error {
	e.path = path
	e.result = result
	e.calls++
	return e.err
}
