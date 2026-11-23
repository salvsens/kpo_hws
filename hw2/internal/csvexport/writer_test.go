package csvexport_test

import (
	"encoding/csv"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
	"github.com/salvsens/kpo_hws/hw2/internal/csvexport"
)

func TestWriteResult(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		student    string
		mode       app.ResultMode
		hundredths int64
		want       string
	}{
		{
			name:       "final with trailing zero",
			student:    "Анна Смирнова",
			mode:       app.FinalMode,
			hundredths: 770,
			want:       "7.70",
		},
		{
			name:       "current with commas and quotes",
			student:    "Анна, группа \"ПОС\"",
			mode:       app.CurrentMode,
			hundredths: 8,
			want:       "0.08",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "итог с пробелами.csv")
			const id = "S-001"
			result := app.Result{
				Student: app.StudentSummary{
					ID:   id,
					Name: tc.student,
				},
				Mode:            tc.mode,
				ScoreHundredths: tc.hundredths,
			}

			if err := csvexport.NewWriter().WriteResult(path, result); err != nil {
				t.Fatalf("WriteResult(%q) error = %v, want nil", path, err)
			}

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
				{id, tc.student, string(tc.mode), tc.want},
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("CSV records = %q, want %q", got, want)
			}
		})
	}
}

func TestWriteResultKeepsExistingFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "result.csv")
	const original = "existing report\n"
	if err := os.WriteFile(path, []byte(original), 0o600); err != nil {
		t.Fatalf("create existing file: %v", err)
	}

	err := csvexport.NewWriter().WriteResult(path, app.Result{})

	if !errors.Is(err, fs.ErrExist) {
		t.Errorf("WriteResult(%q) error = %v, want %v", path, err, fs.ErrExist)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read existing file: %v", err)
	}
	if got := string(data); got != original {
		t.Errorf("file contents = %q, want %q", got, original)
	}
}
