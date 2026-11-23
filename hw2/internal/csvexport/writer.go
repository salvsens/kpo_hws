// Package csvexport writes the result of a student to a CSV file.
package csvexport

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/salvsens/kpo_hws/hw2/internal/app"
)

const (
	// hundredthsPerPoint is the number of hundredths in one point.
	hundredthsPerPoint = 100
	// csvFilePermissions allow only the owner to read or write a new file.
	csvFilePermissions = 0o600
)

var _ app.ResultExporter = Writer{}

// Writer saves one result in a new CSV file.
type Writer struct{}

// NewWriter creates a CSV writer for gradebook results.
func NewWriter() Writer {
	return Writer{}
}

// WriteResult saves the result in a new file at the supplied path.
// It returns an error if the file already exists.
func (Writer) WriteResult(path string, result app.Result) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, csvFilePermissions)
	if err != nil {
		return fmt.Errorf("create CSV: %w", err)
	}

	writeErr := writeCSV(file, result)
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return fmt.Errorf("write CSV %q: %w", path, errors.Join(err, os.Remove(path)))
	}
	return nil
}

// writeCSV writes the header and one row.
func writeCSV(out io.Writer, result app.Result) error {
	score := result.ScoreHundredths
	records := [][]string{
		{"student_id", "student_name", "mode", "result"},
		{
			result.Student.ID,
			result.Student.Name,
			string(result.Mode),
			fmt.Sprintf("%d.%02d", score/hundredthsPerPoint, score%hundredthsPerPoint),
		},
	}
	return csv.NewWriter(out).WriteAll(records)
}
